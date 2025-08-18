package game

import (
	"fmt"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/addequip"
	"maze_game_server/common/function/gentradeno"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/common/function/packtopb/equiptoitem"
	"maze_game_server/config/GMazeConfigV8Cfg"
	"maze_game_server/io/kafka/mazebarrieruserkafka"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/io/redis/mazebarriereventredis"
	"maze_game_server/io/redis/mazeuserbarrierredis"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazecommonvalue"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/pb/server/MazeEquipSvr"
	"maze_game_server/servers/maze_main_server/process/game/events"
	"maze_game_server/services/awardservice"
	"maze_game_server/services/barrierstagecounterservice"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (g *Game) OnMazeBarrierDeathRQ_10449_10450(s *session.Session, req *MazeGame.BarrierDeathRQ) (err error) {
	defer fkprometheus.InfoPMT("OnMazeBarrierDeathRQ")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeGame.BarrierDeathRS{}

	logger.InfoWF("OnMazeBarrierDeathRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnMazeBarrierDeathRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.BarrierId = req.BarrierId

	userId := uint64(s.UID())

	if req.GetBarrierId() <= 0 {
		logger.ErrorWF("OnMazeBarrierDeathRQ req barrier invalid", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡id未设置")
		return
	}
	if req.GetFoeExp() < 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("经验设置错误")
		return
	}

	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierDeathRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	if req.GetBarrierId() != userInfo.Barrier {
		logger.ErrorWF("OnMazeBarrierDeathRQ req barrier lt pass barrier", zap.Any("req", req), zap.Int32("save", userInfo.Barrier))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("请求的关卡id和存储的不一致")
		return
	}

	userBarrier, err := mazeuserbarrierredis.GetUserBarrierInfo(logger, userId, req.GetBarrierId())
	if err != nil {
		logger.ErrorWF("OnMazeBarrierDeathRQ GetUserBarrierInfo fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	// 计算出失败的奖励
	realItem, showItem, realEquip, showEquip, showExp, err := awardservice.GlobalAwardService.GetBarrierDeathAward(logger, userId, req.GetBarrierId())
	if err != nil {
		logger.ErrorWF("OnMazeBarrierDeathRQ GetBarrierDeathAward fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	logger.InfoWF("OnMazeBarrierDeathRQ GetBarrierDeathAward", zap.Any("realItem", realItem), zap.Any("showItem", showItem), zap.Any("realEquip", realEquip), zap.Any("showEquip", showEquip), zap.Any("showExp", showExp))

	var nowExp int64
	showExp -= int64(req.GetFoeExp())
	if showExp > 0 {
		tmpSum := int64(showExp) * int64(GMazeConfigV8Cfg.Get(911).Value_int)
		nowExp = tmpSum/10000 + int64(req.GetFoeExp())
	} else {
		nowExp = int64(req.GetFoeExp())
	}

	//更新等级经验
	oldLevel := userInfo.Level
	oldExp := userInfo.TotalExp
	err = userInfo.AddExp(nowExp)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierDeathRQ addExp fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	newLevel := userInfo.Level
	err = mazeuserinfo.SetUserInfoV2(logger, userId, userInfo)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierDeathRQ SetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	mazecommonvalue.HandleUserLevelExpChg(logger, userId, userInfo.Level, userInfo.Exp, req.GetHeader().GetSession())

	res.BarrierAward = append(res.BarrierAward, &MazeCommon.MazeItem{
		ItemId: proto.Int32(constdef.MazeCommonItemExp),
		Count:  proto.Int64(nowExp),
	})

	defer func() {
		if oldLevel != newLevel {
			levelRecord := &mazeuserlevelkafka.MazeUserLevelRecord{
				UserId:      userId,
				OldLevel:    int32(oldLevel),
				OldTotalExp: oldExp,
				NewLevel:    int32(newLevel),
				NewTotalExp: int32(userInfo.TotalExp),
			}
			mazeuserlevelkafka.PushMazeLevelRecord(logger, levelRecord)
		}
	}()

	// 死亡之后是否需要清空复活次数
	userBarrier.BarrierId = proto.Int32(req.GetBarrierId())
	userBarrier.RebornCount = proto.Int32(0)
	userBarrier.BarrierStatus = proto.Int32(1)
	userBarrier.EndTime = proto.Int64(time.Now().Unix())
	err = mazeuserbarrierredis.SetUserBarrierInfo(logger, userId, 0, userBarrier)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierDeathRQ SetUserBarrierInfo fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	err = mazebarriereventredis.LeaveBarrier(logger, userId, req.GetBarrierId(), false)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierPassRQ LeaveBarrier fail", zap.Error(err))
	}

	// 触发离开关卡事件
	events.OnLeaveBarrier(logger, userId, 0, time.Now().UnixMilli(), &MazeGame.BattleEventLeaveBarrier{BarrierId: proto.Int32(req.GetBarrierId()), Result: MazeGame.BarrierResult_DEATH.Enum()})

	// 发送道具和装备奖励
	tradeNo := gentradeno.GetTradeNum()
	otherItem := make([]*MazeCommon.MazeItem, 0)
	if len(realItem) > 0 {
		awardItems := itemutil.Map2Common(realItem)
		otherItem = append(otherItem, awardItems...)
	}

	if len(otherItem) > 0 {
		//697	UN_CGK_COMMON_BILL_TYPE_697	迷宫扫荡
		errInfo := gentradeno.AddItemEx(logger, uint64(s.UID()), 697, tradeNo, req.GetHeader(), otherItem...)
		if errInfo != nil {
			logger.ErrorWF("CalUserSweepBarrierAward AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("otherItem", otherItem))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}
		// res.BarrierAward = append(res.BarrierAward, otherItem...)
	}

	// 发送装备
	if len(realEquip) > 0 {
		_, err := addequip.AddEquipToBag(logger, uint64(s.UID()), int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_SWEEP_AWARD), tradeNo, realEquip)
		if err != nil {
			logger.ErrorWF("CalUserSweepBarrierAward addEquipToBag fail", zap.Error(err), zap.Any("optype", int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_BOX_AWARD)),
				zap.Any("tradeNo", tradeNo), zap.Any("addEquip", realEquip))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return err
		}
	}

	// 展示获取的奖励
	if len(realItem) > 0 {
		res.BarrierAward = append(res.BarrierAward, itemutil.Map2Common(realItem)...)
	}

	if len(realEquip) > 0 {
		for equipId, count := range realEquip {
			for i := 0; i < int(count); i++ {
				itemEquip, err := equiptoitem.PackMazeEquipInfoSvrToItem(equipId)
				if err != nil {
					logger.ErrorWF("CalUserSweepBarrierAward PackMazeEquipInfoSvrToItem fail", zap.Error(err), zap.Any("equipId", equipId))
					continue
				}
				res.BarrierAward = append(res.BarrierAward, itemEquip)
			}
		}
	}

	logger.InfoWF("OnMazeBarrierDeathRQ showAward", zap.Any("realItem", realItem), zap.Any("realEquip", realEquip))
	// logger.InfoWF("OnMazeBarrierDeathRQ addItems", zap.Any("addItems", addItems), zap.Any("equipItem", equipItem), zap.Any("expCount", expCount), zap.Any("nowExp", nowExp))

	killMonsterNum, totalDamage, _, _, err := barrierstagecounterservice.GlobalBarrierStageCounterService.GetBarrierStageCounter(logger, userId, req.GetBarrierId())
	if err != nil {
		logger.ErrorWF("OnMazeBarrierPassRQ GetBarrierAreaRecord fail", zap.Error(err))
		return err
	}
	res.TotalDamage = proto.Int64(totalDamage)
	res.KillMonsterNum = proto.Int32(killMonsterNum)

	passRecord := &mazebarrieruserkafka.MazeBarrierUserGameRecord{
		UserId:  userId,
		Barrier: req.GetBarrierId(),
		GameRet: mazebarrieruserkafka.GameRetDeath,
		Awards:  getmapAwards(realItem, realEquip),
	}

	mazebarrieruserkafka.PushMazeBarrierUserRecord(logger, passRecord)

	return nil
}

func GetDeathPunish(logger fklog.FKLogI, userId uint64) (subDeathPer int64, lostMin, lostMax int64, err error) {
	// forceVal, err := mazecalcattrredis.GetMazeForce(logger, userId)
	// if err != nil {
	// 	logger.ErrorWF("GetDeathPunish GetMazeForce fail", zap.Error(err))
	// 	return
	// }

	// allDeathCfg := GMazeKongfuDeathV8Cfg.GetAll()
	// if len(allDeathCfg) > 0 {
	// 	for _, cfg := range allDeathCfg {
	// 		if forceVal >= cfg.Kongfu_min && forceVal <= cfg.Kongfu_max {
	// 			return int64(cfg.Death_less_pro), cfg.Death_less_min, cfg.Death_less_max, nil
	// 		}
	// 	}
	// }

	// err = errors.New("not found death cfg")
	return
}

func getmapAwards(awardMap map[int32]int64, awardEquip map[int32]int32) string {
	awardStr := make([]string, 0)

	for k, v := range awardMap {
		awardStr = append(awardStr, fmt.Sprintf("%d:%d", k, v))
	}

	for k, v := range awardEquip {
		awardStr = append(awardStr, fmt.Sprintf("%d:%d", k, v))
	}
	return strings.Join(awardStr, "_")
}
