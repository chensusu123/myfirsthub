package calsweepbarrier

import (
	"errors"
	"fmt"
	"maze_game_server/model/equipdropmodel"
	"maze_game_server/services/equipdropservice"
	"strings"

	"maze_game_server/common/constdef"
	"maze_game_server/common/function/addequip"
	"maze_game_server/common/function/gentradeno"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/common/function/packtopb/equiptoitem"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/config/GMazeBrushFoeV8Cfg"
	"maze_game_server/config/GMazeFoeV8Cfg"
	"maze_game_server/config/GMazeShopV8Cfg"
	"maze_game_server/io/kafka/mazebarrieruserkafka"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/module/mazebarrier"
	"maze_game_server/module/mazecommonvalue"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/Common"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/server/MazeEquipSvr"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func CalUserSweepBarrierAward(logger fklog.FKLogI, uid uint64, barrierId int32, header *Common.PacketHeader) (awardItem []*MazeCommon.MazeItem, rareItem []*MazeCommon.MazeItem, err error) {

	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, uid)
	if err != nil {
		logger.ErrorWF("CalUserSweepBarrierAward GetUserInfoV2 fail", zap.Error(err))
		return
	}

	//稀有材料集合
	rareMap := make(map[int32]struct{})
	barrierCfg := GMazeBarriesV8Cfg.Get(barrierId)
	if barrierCfg == nil {
		logger.ErrorWF("CalUserSweepBarrierAward get barrier cfg fail", zap.Any("barrierId", barrierId))
		err = errors.New("barrier cfg nil")
		return
	}
	for _, v := range barrierCfg.Rare_items_show {
		rareMap[v] = struct{}{}
	}
	//获取关卡所有怪物集合
	foeCountMap := make(map[int32]int32)
	allFoe := GMazeBrushFoeV8Cfg.GetAll()
	if len(allFoe) == 0 {
		logger.ErrorWF("CalUserSweepBarrierAward barrier foe empty", zap.Any("barrierId", barrierId))
		err = errors.New("foe empty")
		return
	}
	for _, cfg := range allFoe {
		if cfg.Barries_id != barrierId {
			continue
		}
		for _, foe := range cfg.Monsters_id {
			if foe > 0 {
				foeCountMap[foe] += 1
			}
		}
		for foe, num := range cfg.Monsterslist_ids_and_nums {
			if foe > 0 {
				foeCountMap[foe] += num
			}
		}
	}

	moneyExtra, err := mazecommonvalue.GetMoneyExtraAdditionEquip(logger, uid)
	if err != nil {
		logger.ErrorWF("CalUserSweepBarrierAward GetMoneyExtraAdditionEquip fail", zap.Error(err))
		return
	}
	expExtra, err := mazecommonvalue.GetExpExtraAdditionEquip(logger, uid)
	if err != nil {
		logger.ErrorWF("CalUserSweepBarrierAward GetExpExtraAdditionEquip fail", zap.Error(err))
		return
	}

	//计算所有怪物可获得的经验和钱和装备
	foeExpMap := make(map[int32]int64)
	foeMoneyMap := make(map[int32]int64)
	foeEquipPointMap := make(map[int32]int64)
	var addExp, addMoney, addEquipPoint int64
	for foeId, num := range foeCountMap {
		foeCfg := GMazeFoeV8Cfg.Get(foeId)
		if foeCfg == nil {
			logger.ErrorWF("CalUserSweepBarrierAward get foe cfg fail", zap.Any("foeId", foeId))
			err = errors.New("foe cfg nil")
			return
		}
		foeExpMap[foeId] = foeCfg.Drop_exp_num[int32(userInfo.Level)]
		foeMoneyMap[foeId] = foeCfg.Drop_coin_num[int32(userInfo.Level)]
		foeEquipPointMap[foeId] = foeCfg.Drop_equip_score_num[int32(userInfo.Level)]
		addExp += (foeCfg.Drop_exp_num[int32(userInfo.Level)] + expExtra) * int64(num)
		addMoney += (foeCfg.Drop_coin_num[int32(userInfo.Level)] + moneyExtra) * int64(num)
		addEquipPoint += foeCfg.Drop_equip_score_num[int32(userInfo.Level)] * int64(num)
	}

	// 计算加成 由于没有武力值 暂时没有额外加成
	addItems, equipMap, err := mazebarrier.GetBarrierPassAward(logger, barrierId)
	if err != nil {
		logger.ErrorWF("CalUserSweepBarrierAward GetBarrierPassAward fail", zap.Any("barrierId", barrierId))
		return
	}

	//取存储的装备分 加上扫荡新增的分数 计算掉落的装备
	calLv := equipdropservice.GlobalEquipDropService.GetMazeBarrierLv(int32(userInfo.Level), barrierId)
	//shopInfo, err := calequipsequence.GetMazeShopInfo(logger, uid, calLv, barrierId)
	//if err != nil {
	//	logger.ErrorWF("CalUserSweepBarrierAward GetMazeShopInfo fail", zap.Any("barrierId", barrierId))
	//	return
	//}
	shopCfg := GMazeShopV8Cfg.Get(calLv)
	if shopCfg == nil {
		logger.ErrorWF("CalUserSweepBarrierAward get shop cfg fail", zap.Any("calLv", calLv))
		err = errors.New("shop cfg nil")
		return
	}

	dropInfo, err := equipdropmodel.NewEquipSpecialDropModel(logger, uid)
	if err != nil {
		logger.ErrorWF("CalUserSweepBarrierAward GetEquipSpecialDropModel fail", zap.Uint64("uid", uid))
		return
	}

	//根据装备积分额外增加装备
	newTotal := dropInfo.EquipPoints + int32(addEquipPoint)
	dropInfo.EquipPoints = newTotal % barrierCfg.Need_equip_score
	equipNum := newTotal / barrierCfg.Need_equip_score

	//addEquipMap, err := calequipsequence.GetNewEquip(logger, uid, barrierId, calLv, equipNum)
	//if err != nil {
	//	logger.ErrorWF("CalUserSweepBarrierAward GetNewEquip fail", zap.Error(err), zap.Any("barrier", barrierId), zap.Any("calLv", calLv))
	//	return
	//}
	addEquipMap, err := equipdropservice.GlobalEquipDropService.GetNewEquip(logger, uid, calLv, barrierId, equipNum)
	if err != nil {
		logger.ErrorWF("CalUserSweepBarrierAward GetNewEquip fail", zap.Error(err), zap.Any("barrier", barrierId), zap.Any("calLv", calLv))
		return
	}

	for k, v := range addEquipMap {
		equipMap[k] += v
	}

	//加经验 加钱 加装备
	if addExp > 0 {
		oldLevel := userInfo.Level
		oldExp := userInfo.TotalExp
		err = userInfo.AddExp(int64(addExp))
		if err != nil {
			logger.ErrorWF("CalUserSweepBarrierAward addExp fail", zap.Any("addExp", addExp))
			return
		}
		newLevel := userInfo.Level
		err = mazeuserinfo.SetUserInfoV2(logger, uid, userInfo)
		if err != nil {
			logger.ErrorWF("CalUserSweepBarrierAward SetUserInfoV2 fail", zap.Error(err))
			return
		}

		mazecommonvalue.HandleUserLevelExpChg(logger, uid, userInfo.Level, userInfo.Exp, header.GetSession())
		defer func() {
			if oldLevel != newLevel {
				levelRecord := &mazeuserlevelkafka.MazeUserLevelRecord{
					UserId:      uid,
					OldLevel:    int32(oldLevel),
					OldTotalExp: oldExp,
					NewLevel:    int32(newLevel),
					NewTotalExp: int32(userInfo.TotalExp),
				}
				mazeuserlevelkafka.PushMazeLevelRecord(logger, levelRecord)
			}
		}()
	}

	awardItem = make([]*MazeCommon.MazeItem, 0)
	rareItem = make([]*MazeCommon.MazeItem, 0)

	if addExp > 0 {
		awardItem = append(awardItem, &MazeCommon.MazeItem{ItemId: proto.Int32(constdef.MazeCommonItemExp), Count: proto.Int64(addExp)})
	}

	tradeNo := gentradeno.GetTradeNum()
	addItems[constdef.MazeCommonItemCoin] += addMoney
	if len(addItems) > 0 {
		awardItems := itemutil.Map2Common(addItems)
		//697	UN_CGK_COMMON_BILL_TYPE_697	迷宫扫荡
		errInfo := gentradeno.AddItemEx(logger, uid, 697, tradeNo, header, awardItems...)
		if errInfo != nil {
			logger.ErrorWF("CalUserSweepBarrierAward AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("awardItems", awardItems))
		} else {

			// queryItems, errInfo2 := gentradeno.QueryItems(logger, uid, &MazeCommon.MazeItem{ItemId: proto.Int32(constdef.MazeCommonItemCoin)}, &MazeCommon.MazeItem{ItemId: proto.Int32(constdef.MazeCommonItemDiamond)})
			// if errInfo2 == nil {
			// 	commonList := make([]*mazecommonvalue.CommonValueStruct, 0)
			// 	for _, v := range queryItems {
			// 		if v.GetItemId() == constdef.MazeCommonItemCoin {
			// 			coinCommon := mazecommonvalue.MakeCommonValueList(logger, map[int32]int64{
			// 				int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_MONEY): v.GetCount()},
			// 				map[int32]int32{}, map[int32]string{})
			// 			commonList = append(commonList, coinCommon...)
			// 		} else if v.GetItemId() == constdef.MazeCommonItemDiamond {
			// 			diamondCommon := mazecommonvalue.MakeCommonValueList(logger, map[int32]int64{
			// 				int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_DIAMOND): v.GetCount()},
			// 				map[int32]int32{}, map[int32]string{})
			// 			commonList = append(commonList, diamondCommon...)
			// 		}
			// 	}
			// 	mazecommonvalue.SendCommonValueIdPack(logger, uid, commonList)
			// }
		}
		for _, v := range awardItems {
			_, ok := rareMap[v.GetItemId()]
			if ok {
				rareItem = append(rareItem, v)
			} else {
				awardItem = append(awardItem, v)
			}
		}
	}
	if len(equipMap) > 0 {
		rs, err := addequip.AddEquipToBag(logger, uid, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_SWEEP_AWARD), tradeNo, equipMap)
		if err != nil {
			logger.ErrorWF("CalUserSweepBarrierAward addEquipToBag fail", zap.Error(err), zap.Any("optype", int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_BOX_AWARD)),
				zap.Any("tradeNo", tradeNo), zap.Any("addEquip", equipMap))
		}

		for _, equip := range rs.GetEquipList() {
			itemEquip, err := equiptoitem.PackEquipToItem(equip)
			if err != nil {
				logger.ErrorWF("CalUserSweepBarrierAward PackEquipToItem fail", zap.Error(err), zap.Any("equip", equip))
				continue
			}
			_, ok := rareMap[itemEquip.GetItemId()]
			if ok {
				rareItem = append(rareItem, itemEquip)
			} else {
				awardItem = append(awardItem, itemEquip)
			}
		}
	}

	sweepRecord := &mazebarrieruserkafka.MazeBarrierUserGameRecord{
		UserId:  uid,
		Barrier: barrierId,
		GameRet: mazebarrieruserkafka.GameRetSweep,
		Awards:  getAwards(addItems, equipMap),
	}

	mazebarrieruserkafka.PushMazeBarrierUserRecord(logger, sweepRecord)

	return
}

func getAwards(awardMap map[int32]int64, awardEquip map[int32]int32) string {
	awardStr := make([]string, 0)

	for k, v := range awardMap {
		awardStr = append(awardStr, fmt.Sprintf("%d:%d", k, v))
	}

	for k, v := range awardEquip {
		awardStr = append(awardStr, fmt.Sprintf("%d:%d", k, v))
	}
	return strings.Join(awardStr, "_")
}
