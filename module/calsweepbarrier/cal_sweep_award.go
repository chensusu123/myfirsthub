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
	"maze_game_server/config/GMazeConfigV8Cfg"
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

// 获取扫荡关卡的奖励 equipItem 装备奖励 ohterItem 物品奖励 expItem 经验奖励
func GetSweepBarrierAward(logger fklog.FKLogI, uid uint64, barrierId int32) (equipItem map[int32]int32, ohterItem []*MazeCommon.MazeItem, expItem *MazeCommon.MazeItem, err error) {
	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, uid)
	if err != nil {
		logger.ErrorWF("GetSweepBarrierAward GetUserInfoV2 fail", zap.Error(err))
		return
	}

	barrierCfg := GMazeBarriesV8Cfg.Get(barrierId)
	if barrierCfg == nil {
		logger.ErrorWF("GetSweepBarrierAward get barrier cfg fail", zap.Any("barrierId", barrierId))
		err = errors.New("barrier cfg nil")
		return
	}

	//获取关卡所有怪物集合
	foeCountMap := make(map[int32]int32)
	allFoe := GMazeBrushFoeV8Cfg.GetAll()
	if len(allFoe) == 0 {
		logger.ErrorWF("GetSweepBarrierAward barrier foe empty", zap.Any("barrierId", barrierId))
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
		logger.ErrorWF("GetSweepBarrierAward GetMoneyExtraAdditionEquip fail", zap.Error(err))
		return
	}
	expExtra, err := mazecommonvalue.GetExpExtraAdditionEquip(logger, uid)
	if err != nil {
		logger.ErrorWF("GetSweepBarrierAward GetExpExtraAdditionEquip fail", zap.Error(err))
		return
	}

	//计算所有怪物可获得的经验和钱和装备
	foeExpMap := make(map[int32]int64)
	foeMoneyMap := make(map[int32]int64)
	foeEquipPointMap := make(map[int32]int64)
	var addItem1, addItem2 int32
	var addExp, addMoney, addEquipPoint int64
	for foeId, num := range foeCountMap {
		foeCfg := GMazeFoeV8Cfg.Get(foeId)
		if foeCfg == nil {
			logger.ErrorWF("GetSweepBarrierAward get foe cfg fail", zap.Any("foeId", foeId))
			err = errors.New("foe cfg nil")
			return
		}
		foeExpMap[foeId] = int64(foeCfg.Drop_exp_num)
		foeMoneyMap[foeId] = int64(foeCfg.Drop_coin_num)
		foeEquipPointMap[foeId] = int64(foeCfg.Drop_equip_score_num)
		addExp += (int64(foeCfg.Drop_exp_num) + expExtra) * int64(num)
		addMoney += (int64(foeCfg.Drop_coin_num) + moneyExtra) * int64(num)
		addEquipPoint += int64(foeCfg.Drop_equip_score_num) * int64(num)
		addItem1 += foeCfg.Drop_item1_score_num * num
		addItem2 += foeCfg.Drop_item2_score_num * num
	}

	// 计算加成 由于没有武力值 暂时没有额外加成
	// 获取关卡宝箱掉落奖励 包含装备和物品
	addItems, equipMap, err := mazebarrier.GetBarrierPassAward(logger, barrierId)
	if err != nil {
		logger.ErrorWF("GetSweepBarrierAward GetBarrierPassAward fail", zap.Any("barrierId", barrierId))
		return
	}

	//取存储的装备分 加上扫荡新增的分数 计算掉落的装备
	calLv := equipdropservice.GlobalEquipDropService.GetMazeBarrierLv(int32(userInfo.Level), barrierId)
	shopCfg := GMazeShopV8Cfg.Get(calLv)
	if shopCfg == nil {
		logger.ErrorWF("GetSweepBarrierAward get shop cfg fail", zap.Any("calLv", calLv))
		err = errors.New("shop cfg nil")
		return
	}

	dropInfo, err := equipdropmodel.NewEquipSpecialDropModel(logger, uid)
	if err != nil {
		logger.ErrorWF("GetSweepBarrierAward GetEquipSpecialDropModel fail", zap.Uint64("uid", uid))
		return
	}

	//根据装备积分额外增加装备
	newTotal := dropInfo.EquipPoints + int32(addEquipPoint)
	dropInfo.EquipPoints = newTotal % barrierCfg.Need_equip_score
	equipNum := newTotal / barrierCfg.Need_equip_score

	// fmt.Println("newTotal", newTotal)
	// fmt.Println("equipNum", equipNum)

	// 掉落道具一
	item1Num := addItem1 / barrierCfg.Need_item1_score
	if item1Num > 0 {
		// id从通用配置获取
		itemConfig := GMazeConfigV8Cfg.Get(901)
		for _, v := range itemConfig.Value_map {
			addItems[int32(v)] += int64(item1Num) * int64(barrierCfg.Item1_nums_per_pile)
		}
	}

	// 掉落道具二
	item2Num := addItem2 / barrierCfg.Need_item2_score
	if item2Num > 0 {
		// id从通用配置获取
		itemConfig := GMazeConfigV8Cfg.Get(902)
		for _, v := range itemConfig.Value_map {
			addItems[int32(v)] += int64(item2Num) * int64(barrierCfg.Item2_nums_per_pile)
		}
	}

	addEquipMap, err := equipdropservice.GlobalEquipDropService.GetNewEquip(logger, uid, calLv, barrierId, equipNum)
	if err != nil {
		logger.ErrorWF("GetSweepBarrierAward GetNewEquip fail", zap.Error(err), zap.Any("barrier", barrierId), zap.Any("calLv", calLv))
		return
	}

	// logger.InfoWF("equipMap", zap.Any("equipMap", equipMap))
	for k, v := range addEquipMap {
		equipMap[k] += v
	}
	// logger.InfoWF("addEquipMap", zap.Any("addEquipMap", addEquipMap))
	// logger.InfoWF("equipMap", zap.Any("equipMap", equipMap))
	equipItem = equipMap

	// 经验奖励
	if addExp > 0 {
		expItem = &MazeCommon.MazeItem{ItemId: proto.Int32(constdef.MazeCommonItemExp), Count: proto.Int64(addExp)}
	}

	addItems[constdef.MazeCommonItemCoin] += addMoney
	if len(addItems) > 0 {
		awardItems := itemutil.Map2Common(addItems)
		ohterItem = append(ohterItem, awardItems...)
	}

	return
}

// 发送扫荡奖励
func CalUserSweepBarrierAward(logger fklog.FKLogI, uid uint64, barrierId int32, header *Common.PacketHeader) (awardItem []*MazeCommon.MazeItem, rareItem []*MazeCommon.MazeItem, err error) {
	// 获取扫荡奖励
	equipItem, ohterItem, expItem, err := GetSweepBarrierAward(logger, uid, barrierId)
	if err != nil {
		logger.ErrorWF("CalUserSweepBarrierAward GetSweepBarrierAward fail", zap.Error(err))
		return
	}

	// 获取用户信息
	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, uid)
	if err != nil {
		logger.ErrorWF("CalUserSweepBarrierAward GetUserInfoV2 fail", zap.Error(err))
		return
	}

	barrierCfg := GMazeBarriesV8Cfg.Get(barrierId)
	if barrierCfg == nil {
		logger.ErrorWF("GetSweepBarrierAward get barrier cfg fail", zap.Any("barrierId", barrierId))
		err = errors.New("barrier cfg nil")
		return
	}

	//稀有材料集合
	rareMap := make(map[int32]struct{})
	for _, v := range barrierCfg.Rare_items_show {
		rareMap[v] = struct{}{}
	}

	// 发送经验
	if expItem != nil {
		addExp := expItem.GetCount()
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

	if expItem != nil {
		awardItem = append(awardItem, expItem)
	}

	tradeNo := gentradeno.GetTradeNum()

	// 发送物品
	if len(ohterItem) > 0 {
		//697	UN_CGK_COMMON_BILL_TYPE_697	迷宫扫荡
		errInfo := gentradeno.AddItemEx(logger, uid, 697, tradeNo, header, ohterItem...)
		if errInfo != nil {
			logger.ErrorWF("CalUserSweepBarrierAward AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("awardItems", ohterItem))
		}

		for _, v := range ohterItem {
			_, ok := rareMap[v.GetItemId()]
			if ok {
				rareItem = append(rareItem, v)
			} else {
				awardItem = append(awardItem, v)
			}
		}
	}

	// 发送装备
	if len(equipItem) > 0 {
		rs, err := addequip.AddEquipToBag(logger, uid, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_SWEEP_AWARD), tradeNo, equipItem)
		if err != nil {
			logger.ErrorWF("CalUserSweepBarrierAward addEquipToBag fail", zap.Error(err), zap.Any("optype", int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_BOX_AWARD)),
				zap.Any("tradeNo", tradeNo), zap.Any("addEquip", equipItem))
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
		Awards:  getAwards(rareItem, awardItem),
	}

	mazebarrieruserkafka.PushMazeBarrierUserRecord(logger, sweepRecord)

	return
}

func getAwards(awardMap []*MazeCommon.MazeItem, awardEquip []*MazeCommon.MazeItem) string {
	awardStr := make([]string, 0)

	for _, v := range awardMap {
		awardStr = append(awardStr, fmt.Sprintf("%d:%d", v.GetItemId(), v.GetCount()))
	}

	for _, v := range awardEquip {
		awardStr = append(awardStr, fmt.Sprintf("%d:%d", v.GetItemId(), v.GetCount()))
	}
	return strings.Join(awardStr, "_")
}
