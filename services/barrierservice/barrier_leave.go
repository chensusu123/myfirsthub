package barrierservice

import (
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/addequip"
	"maze_game_server/common/function/gentradeno"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/common/function/packtopb/equiptoitem"
	"maze_game_server/config/GMazeActionCountV8Cfg"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/config/GMazeConfigV8Cfg"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/io/redis/mazechallengenumredis"
	"maze_game_server/model/userbarriermodel"
	"maze_game_server/model/userinfomodel"
	"maze_game_server/module/mazebarrier"
	"maze_game_server/module/mazecommonvalue"
	"maze_game_server/pb/common/Common"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MessageType"
	"maze_game_server/pb/server/MazeEquipSvr"
	"maze_game_server/services/awardservice"
	"maze_game_server/services/barrierstagecounterservice"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// BarrierPass implements BarrierService.
func (b *barrier) BarrierPass(logger fklog.FKLogI, header *Common.PacketHeader, userID uint64, barrierID int32, foeExp int32) (
	killMonsterNum int32, totalDamage int64, awards, rareAwards []*MazeCommon.MazeItem, errinfo *MessageType.ErrorInfo) {
	cfg := GMazeBarriesV8Cfg.Get(barrierID)
	if cfg == nil {
		logger.ErrorWF("BarrierPass get barrier cfg fail", zap.Any("barrier", barrierID), zap.Any("foeExp", foeExp))
		return 0, 0, nil, nil, errors.COMMON_ERROR_TIPS.Wrap("关卡配置数据获取失败")
	}
	rareMap := make(map[int32]struct{})
	for _, v := range cfg.Rare_items_show {
		rareMap[v] = struct{}{}
	}
	rareMap = map[int32]struct{}{}

	userInfo, err := userinfomodel.NewUserInfoModel(logger, userID)
	if err != nil {
		logger.ErrorWF("BarrierPass GetUserInfoV2 fail", zap.Error(err), zap.Any("barrier", barrierID), zap.Any("foeExp", foeExp))
		return 0, 0, nil, nil, errors.MODULE_ERROR.ToInfo()
	}

	if userInfo.PassBarrier >= barrierID {
		logger.ErrorWF("BarrierPass req barrier lt pass barrier", zap.Any("barrierID", barrierID), zap.Int32("pass", userInfo.PassBarrier))
		return 0, 0, nil, nil, errors.COMMON_ERROR_TIPS.Wrap("该关卡已上报过通关")
	}

	if userInfo.Barrier != barrierID {
		logger.ErrorWF("BarrierPass userinfo barrier not match", zap.Any("barrierID", barrierID), zap.Int32("save", userInfo.Barrier))
		return 0, 0, nil, nil, errors.COMMON_ERROR_TIPS.Wrap("记录用户正在打的关卡与上报通关id不匹配")
	}

	userBarrier, err := userbarriermodel.NewUserBarrierModel(logger, userID)
	if err != nil {
		logger.ErrorWF("BarrierPass GetUserBarrierInfo fail", zap.Error(err), zap.Any("barrier", barrierID), zap.Any("foeExp", foeExp))
		return 0, 0, nil, nil, errors.MODULE_ERROR.ToInfo()
	}

	// userInfo.SetPassBarrier(barrierID)
	userInfo.PassBarrier = barrierID
	// userInfo.SetBarrier(cfg.Next_id)
	oldLevel := userInfo.Level
	oldExp := userInfo.TotalExp
	err = userInfo.AddExp(int64(foeExp))
	if err != nil {
		logger.ErrorWF("BarrierPass AddExp fail", zap.Error(err), zap.Any("barrier", barrierID), zap.Any("foeExp", foeExp))
		return 0, 0, nil, nil, errors.MODULE_ERROR.ToInfo()
	}
	newLevel := userInfo.Level
	err = userInfo.Save(logger)
	if err != nil {
		logger.ErrorWF("BarrierPass SetUserInfoV2 fail", zap.Error(err), zap.Any("barrier", barrierID), zap.Any("foeExp", foeExp))
		return 0, 0, nil, nil, errors.MODULE_ERROR.ToInfo()
	}
	mazecommonvalue.HandleUserLevelExpChg(logger, userID, userInfo.Level, userInfo.Exp, header.GetSession())

	defer func() {
		if oldLevel != newLevel {
			levelRecord := &mazeuserlevelkafka.MazeUserLevelRecord{
				UserId:      userID,
				OldLevel:    int32(oldLevel),
				OldTotalExp: oldExp,
				NewLevel:    int32(newLevel),
				NewTotalExp: int32(userInfo.TotalExp),
			}
			mazeuserlevelkafka.PushMazeLevelRecord(logger, levelRecord)
		}
	}()

	if foeExp > 0 {
		expItem := &MazeCommon.MazeItem{ItemId: proto.Int32(constdef.MazeCommonItemExp), Count: proto.Int64(int64(foeExp))}
		_, ok := rareMap[constdef.MazeCommonItemExp]
		if ok {
			rareAwards = append(rareAwards, expItem)
		} else {
			awards = append(awards, expItem)
		}
	}

	//成功通关需要 (1和2通过kafka清理)
	//1. 清除临时buff
	//2. 影响挂机生产
	//3. 清除客户端透传数据(通过切换关卡id 切换不同的key 目前没清)

	// 死亡之后是否需要清空复活次数
	userBarrier.RebornCount = 0
	userBarrier.BarrierStatus = 3
	userBarrier.EndTime = time.Now().Unix()
	err = userBarrier.Save(logger)
	if err != nil {
		logger.ErrorWF("BarrierPass SetUserBarrierInfo fail", zap.Error(err), zap.Any("barrier", barrierID), zap.Any("foeExp", foeExp))
		return 0, 0, nil, nil, errors.MODULE_ERROR.ToInfo()
	}

	now := time.Now()
	today := now.Year()*10000 + int(now.Month())*100 + now.Day()

	var awardBarrierNum int32
	awardBarrierNumCfg := GMazeActionCountV8Cfg.Get(102)
	if awardBarrierNumCfg != nil {
		awardBarrierNum = awardBarrierNumCfg.Day_count_v8
	}

	useNumToday, err2 := mazechallengenumredis.GetUserChallengeNum(logger, userID, today)
	if err2 != nil {
		logger.ErrorWF("BarrierPass GetUserChallengeNum fail", zap.Error(err2))
	} else if useNumToday > 0 {
		if awardBarrierNum > int32(useNumToday) {
			awardBarrierNum = int32(useNumToday)
		}
		err = mazechallengenumredis.AddUserChallengeNum(logger, userID, today, 0-awardBarrierNum)
		if err != nil {
			logger.ErrorWF("BarrierPass AddUserChallengeNum fail", zap.Error(err))
		}
	}

	awardMap, equipMap, err := mazebarrier.GetBarrierPassAwardWithFirst(logger, barrierID)
	if err != nil {
		logger.ErrorWF("BarrierPass GetBarrierPassAwardWithFirst fail", zap.Error(err), zap.Any("barrier", barrierID))
	} else {
		//696	UN_CGK_COMMON_BILL_TYPE_696	迷宫通关
		tradeNo := gentradeno.GetTradeNum()
		if len(awardMap) > 0 {
			awardItems := itemutil.Map2Common(awardMap)
			errInfo := gentradeno.AddItemEx(logger, userID, 696, tradeNo, header, awardItems...)
			if errInfo != nil {
				logger.ErrorWF("BarrierPass AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("awardItems", awardItems))
			}

			for _, item := range awardItems {
				_, ok := rareMap[item.GetItemId()]
				if ok {
					rareAwards = append(rareAwards, item)
				} else {
					awards = append(awards, item)
				}
			}
		}

		if len(equipMap) > 0 {
			//MAZE_EQUIP_PASS_AWARD = 9;//迷宫通关奖励 张登元
			rs, err := addequip.AddEquipToBag(logger, userID, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_PASS_AWARD), tradeNo, equipMap)
			if err != nil {
				logger.ErrorWF("BarrierPass addEquipToBag fail", zap.Error(err), zap.Any("optype", int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_BOX_AWARD)),
					zap.Any("tradeNo", tradeNo), zap.Any("addEquip", equipMap))
			}

			for _, equip := range rs.GetEquipList() {
				itemEquip, err := equiptoitem.PackEquipToItem(equip)
				if err != nil {
					logger.ErrorWF("BarrierPass PackEquipToItem fail", zap.Error(err), zap.Any("equip", equip))
					continue
				}
				_, ok := rareMap[itemEquip.GetItemId()]
				if ok {
					rareAwards = append(rareAwards, itemEquip)
				} else {
					awards = append(awards, itemEquip)
				}
			}
		}
	}

	killMonsterNum, totalDamage, _, _, err = barrierstagecounterservice.GlobalBarrierStageCounterService.GetBarrierStageCounter(logger, userID, barrierID)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierPassRQ GetBarrierAreaRecord fail", zap.Error(err))
		return 0, 0, nil, nil, errors.MODULE_ERROR.ToInfo()
	}

	return killMonsterNum, totalDamage, awards, rareAwards, nil
}

// BarrierDeath implements BarrierService.
func (b *barrier) BarrierDeath(logger fklog.FKLogI, header *Common.PacketHeader, userID uint64, barrierID int32, foeExp int32) (
	killMonsterNum int32, totalDamage int64, awards []*MazeCommon.MazeItem, errinfo *MessageType.ErrorInfo) {
	userInfo, err := userinfomodel.NewUserInfoModel(logger, userID)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierDeathRQ GetUserInfoV2 fail", zap.Error(err), zap.Any("barrier", barrierID), zap.Any("foeExp", foeExp))
		return 0, 0, nil, errors.MODULE_ERROR.ToInfo()
	}

	if barrierID != userInfo.Barrier {
		logger.ErrorWF("OnMazeBarrierDeathRQ req barrier lt pass barrier", zap.Any("barrierID", barrierID), zap.Int32("save", userInfo.Barrier))
		return 0, 0, nil, errors.COMMON_ERROR_TIPS.Wrap("请求的关卡id和存储的不一致")
	}

	userBarrier, err := userbarriermodel.NewUserBarrierModel(logger, userID)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierDeathRQ GetUserBarrierInfo fail", zap.Error(err), zap.Any("barrier", barrierID), zap.Any("foeExp", foeExp))
		return 0, 0, nil, errors.MODULE_ERROR.ToInfo()
	}

	// 计算出失败的奖励
	realItem, showItem, realEquip, showEquip, showExp, err := awardservice.GlobalAwardService.GetBarrierDeathAward(logger, userID, barrierID)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierDeathRQ GetBarrierDeathAward fail", zap.Error(err), zap.Any("barrier", barrierID), zap.Any("foeExp", foeExp))
		return 0, 0, nil, errors.MODULE_ERROR.ToInfo()
	}

	logger.InfoWF("OnMazeBarrierDeathRQ GetBarrierDeathAward", zap.Any("realItem", realItem), zap.Any("showItem", showItem), zap.Any("realEquip", realEquip), zap.Any("showEquip", showEquip), zap.Any("showExp", showExp))

	var nowExp int64
	showExp -= int64(foeExp)
	if showExp > 0 {
		tmpSum := int64(showExp) * int64(GMazeConfigV8Cfg.Get(911).Value_int)
		nowExp = tmpSum/10000 + int64(foeExp)
	} else {
		nowExp = int64(foeExp)
	}

	//更新等级经验
	oldLevel := userInfo.Level
	oldExp := userInfo.TotalExp
	err = userInfo.AddExp(nowExp)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierDeathRQ addExp fail", zap.Error(err), zap.Any("barrier", barrierID), zap.Any("foeExp", foeExp))
		return 0, 0, nil, errors.MODULE_ERROR.ToInfo()
	}
	newLevel := userInfo.Level
	err = userInfo.Save(logger)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierDeathRQ SetUserInfoV2 fail", zap.Error(err), zap.Any("barrier", barrierID), zap.Any("foeExp", foeExp))
		return 0, 0, nil, errors.MODULE_ERROR.ToInfo()
	}
	mazecommonvalue.HandleUserLevelExpChg(logger, userID, userInfo.Level, userInfo.Exp, header.GetSession())

	awards = append(awards, &MazeCommon.MazeItem{
		ItemId: proto.Int32(constdef.MazeCommonItemExp),
		Count:  proto.Int64(nowExp),
	})

	defer func() {
		if oldLevel != newLevel {
			levelRecord := &mazeuserlevelkafka.MazeUserLevelRecord{
				UserId:      userID,
				OldLevel:    int32(oldLevel),
				OldTotalExp: oldExp,
				NewLevel:    int32(newLevel),
				NewTotalExp: int32(userInfo.TotalExp),
			}
			mazeuserlevelkafka.PushMazeLevelRecord(logger, levelRecord)
		}
	}()

	// 死亡之后是否需要清空复活次数
	userBarrier.BarrierID = barrierID
	userBarrier.RebornCount = 0
	userBarrier.BarrierStatus = 1
	userBarrier.EndTime = time.Now().Unix()
	err = userBarrier.Save(logger)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierDeathRQ SetUserBarrierInfo fail", zap.Error(err), zap.Any("barrier", barrierID), zap.Any("foeExp", foeExp))
		return 0, 0, nil, errors.MODULE_ERROR.ToInfo()
	}

	// 发送道具和装备奖励
	// 发送道具和装备奖励
	tradeNo := gentradeno.GetTradeNum()
	otherItem := make([]*MazeCommon.MazeItem, 0)
	if len(realItem) > 0 {
		awardItems := itemutil.Map2Common(realItem)
		otherItem = append(otherItem, awardItems...)
	}

	if len(otherItem) > 0 {
		//697	UN_CGK_COMMON_BILL_TYPE_697	迷宫扫荡
		errInfo := gentradeno.AddItemEx(logger, userID, 697, tradeNo, header, otherItem...)
		if errInfo != nil {
			logger.ErrorWF("CalUserSweepBarrierAward AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("otherItem", otherItem))
			return 0, 0, nil, errors.MODULE_ERROR.ToInfo()
		}
		// res.BarrierAward = append(res.BarrierAward, otherItem...)
	}

	// 发送装备
	if len(realEquip) > 0 {
		_, err := addequip.AddEquipToBag(logger, userID, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_SWEEP_AWARD), tradeNo, realEquip)
		if err != nil {
			logger.ErrorWF("CalUserSweepBarrierAward addEquipToBag fail", zap.Error(err), zap.Any("optype", int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_BOX_AWARD)),
				zap.Any("tradeNo", tradeNo), zap.Any("addEquip", realEquip))
			return 0, 0, nil, errors.MODULE_ERROR.ToInfo()
		}
	}

	// 展示获取的奖励
	if len(realItem) > 0 {
		awards = append(awards, itemutil.Map2Common(realItem)...)
	}

	if len(realEquip) > 0 {
		for equipId, count := range realEquip {
			for i := 0; i < int(count); i++ {
				itemEquip, err := equiptoitem.PackMazeEquipInfoSvrToItem(equipId)
				if err != nil {
					logger.ErrorWF("CalUserSweepBarrierAward PackMazeEquipInfoSvrToItem fail", zap.Error(err), zap.Any("equipId", equipId))
					continue
				}
				awards = append(awards, itemEquip)
			}
		}
	}

	logger.InfoWF("OnMazeBarrierDeathRQ showAward", zap.Any("realItem", realItem), zap.Any("realEquip", realEquip))
	// logger.InfoWF("OnMazeBarrierDeathRQ addItems", zap.Any("addItems", addItems), zap.Any("equipItem", equipItem), zap.Any("expCount", expCount), zap.Any("nowExp", nowExp))

	killMonsterNum, totalDamage, _, _, err = barrierstagecounterservice.GlobalBarrierStageCounterService.GetBarrierStageCounter(logger, userID, barrierID)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierPassRQ GetBarrierAreaRecord fail", zap.Error(err))
		return 0, 0, nil, errors.MODULE_ERROR.ToInfo()
	}

	return killMonsterNum, totalDamage, awards, nil
}
