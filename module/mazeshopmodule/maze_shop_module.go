package mazeshopmodule

import (
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeShopEquipListV8Cfg"
	"maze_game_server/config/GMazeShopV8Cfg"
	"maze_game_server/io/redis/mazeshopseqredis"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

//func GetMaxAreaId(barrierId, maxAreaId int32) int32 {
//	if maxAreaId == 0 {
//		maxAreaId = 10000
//		if barrierId > 0 {
//			maxAreaId = maxAreaId * barrierId
//		}
//	} else {
//		areaCfg := GDollMazeBrushAreaV8Cfg.Get(maxAreaId)
//		if areaCfg == nil {
//			return maxAreaId
//		}
//		if barrierId > areaCfg.Barries {
//			maxAreaId = barrierId * 10000
//		}
//	}
//	return maxAreaId
//}

//func GetMazeShopInfo(logger fklog.FKLogI, userId uint64, mazeLevel int32, areaId int32) (mazeShopInfo *mazeshopseqredis.MazeShopInfo, err error) {
//	newLevel := GetMazeShopLv(mazeLevel, areaId)
//	mazeShopInfo, err = mazeshopseqredis.GetMazeShopInfo(logger, userId, newLevel)
//	if err != nil {
//		logger.ErrorWF("GetMazeShopInfo error!", zap.Error(err), zap.Any("newLevel", newLevel))
//		return
//	}
//	cfg := GMazeShopV8Cfg.Get(newLevel)
//	if cfg == nil {
//		logger.ErrorWF("GetMazeShopInfo GMazeShopV8Cfg fail",
//			zap.Int32("newLevel", newLevel))
//		return mazeShopInfo, errors.New("获取配置失败")
//	}
//	if mazeShopInfo == nil {
//		mazeShopInfo, err = HandleMazeShopSeqInit(logger, cfg)
//		if err != nil || mazeShopInfo == nil {
//			logger.ErrorWF("GetMazeShopInfo HandleMazeShopSeqInit fail",
//				zap.Int32("newLevel", newLevel),
//				zap.Error(err))
//			return mazeShopInfo, errors.New("获取配置失败")
//		}
//	}
//	return mazeShopInfo, nil
//}

//func HandleMazeShopSeqInit(logger fklog.FKLogI, cfg *GMazeShopV8Cfg.MazeShopV8ConfigRow) (mazeShopInfo *mazeshopseqredis.MazeShopInfo, err error) {
//	mazeShopInfo = &mazeshopseqredis.MazeShopInfo{}
//	if cfg == nil {
//		logger.ErrorWF("HandleMazeShopSeqInit GMazeShopV8Cfg err")
//		return nil, errors.New("配置不存在")
//	}
//	seqWeight := make(map[int32]int32, 0)
//	for _, seqId := range cfg.Buy_list {
//		seqWeight[seqId] = 1000
//	}
//	seqId := randfuncs.RandByWeightV2(logger, seqWeight, false)
//	if seqId <= 0 {
//		logger.ErrorWF("HandleMazeShopSeqInit GMazeShopV8Cfg err",
//			zap.Any("seqWeight", seqWeight))
//		return nil, errors.New("配置不存在")
//	}
//	mazeShopInfo.SeqId = seqId
//	mazeShopInfo.BackId = cfg.Buy_list_base2
//	mazeShopInfo.ShopSlotNum = make(map[int32]int32)
//	return mazeShopInfo, nil
//}

func GetMazeSeqEquipId(logger fklog.FKLogI, userId uint64, mazeLevel int32, areaId int32, mazeShopInfo *mazeshopseqredis.MazeShopInfo, addCount, assignPos int32, equipMaxMap map[int32]int64) (equipMap map[int32]int32, err error) {
	newLevel := GetMazeShopLv(mazeLevel, areaId)
	equipMap = make(map[int32]int32)
	shopCfg := GMazeShopV8Cfg.Get(newLevel)
	if shopCfg == nil {
		logger.ErrorWF("GetMazeSeqEquipId GMazeShopV8Cfg fail",
			zap.Int32("newLevel", newLevel))
		return equipMap, errors.New("区域id配置不存在")
	}
	// attrMaxCount := 0
	// for k, v := range shopCfg.Equip_value_max {
	// 	if equipMaxMap[k] >= v {
	// 		attrMaxCount++
	// 	}
	// }
	// if attrMaxCount >= len(shopCfg.Equip_value_max) {
	// 	logger.ErrorWF("GetMazeSeqEquipId attribute max",
	// 		zap.Any("attrMaxCount", attrMaxCount),
	// 		zap.Any("equipMaxMap", equipMaxMap),
	// 		zap.Any("EquipValueMax", shopCfg.Equip_value_max))
	// 	return equipMap, errors.New("所有部位属性全部达到最大值")
	// }
	// if assignPos > 0 {
	// 	if equipMaxMap[assignPos] >= shopCfg.Equip_value_max[assignPos] {
	// 		logger.ErrorWF("GetMazeSeqEquipId attribute max",
	// 			zap.Any("equipMaxMap", equipMaxMap),
	// 			zap.Any("assignPos", assignPos))
	// 		return equipMap, errors.New("该部位属性已达到最大值")
	// 	}
	// }

	var seqCfg *GMazeShopEquipListV8Cfg.MazeShopEquipListV8ConfigRow
	var backCfg *GMazeShopEquipListV8Cfg.MazeShopEquipListV8ConfigRow
	for _, cfg := range GMazeShopEquipListV8Cfg.GetAll() {
		if (cfg.List_id == mazeShopInfo.SeqId) && (cfg.Pos_id == 0) {
			seqCfg = cfg
		}
		if (cfg.List_id == mazeShopInfo.BackId) && (cfg.Pos_id == 0) {
			backCfg = cfg
		}
	}
	if (seqCfg == nil) || (backCfg == nil) {
		logger.ErrorWF("GetMazeShopSeq GMazeShopEquipListV8Cfg err",
			zap.Any("mazeShopInfo", mazeShopInfo))
		return equipMap, errors.New("配置不存在")
	}
	for i := int32(0); i < addCount; i++ {
		var equipId int32
		if mazeShopInfo.BackIndex > 0 || (int(mazeShopInfo.CurSeqIndex) >= len(seqCfg.Equip_id)) {
			if int(mazeShopInfo.BackIndex) >= len(backCfg.Equip_id) {
				mazeShopInfo.BackIndex = 0
			}
			mazeShopInfo.BackIndex++
			mazeShopInfo.TotalCount++
			equipId = backCfg.Equip_id[mazeShopInfo.BackIndex-1]
		} else {
			mazeShopInfo.CurSeqIndex++
			mazeShopInfo.TotalCount++
			equipId = seqCfg.Equip_id[mazeShopInfo.CurSeqIndex-1]
		}
		convEquipId, err := ConvEquipId(logger, equipId, assignPos, equipMaxMap, shopCfg.Equip_value_max)
		if equipId != convEquipId {
			logger.InfoWF("GetMazeSeqEquipId conv equipId",
				zap.Int32("equipId", equipId),
				zap.Int32("convEquipId", convEquipId),
				zap.Any("assignPos", assignPos),
				zap.Any("equipMaxMap", equipMaxMap),
				zap.Any("equipValueMax", shopCfg.Equip_value_max))
		}
		if err != nil || convEquipId == 0 {
			logger.ErrorWF("GetMazeShopSeq ConvEquipId err",
				zap.Error(err),
				zap.Any("equipMaxMap", equipMaxMap),
				zap.Any("equipId", equipId))
			return equipMap, errors.New("配置不存在")
		}
		equipMap[convEquipId] += 1
	}
	logger.InfoWF("GetMazeSeqEquipId end",
		zap.Int32("newLevel", newLevel),
		zap.Int32("mazeLevel", mazeLevel),
		zap.Int32("areaId", areaId),
		zap.Any("mazeShopInfo", mazeShopInfo),
		zap.Any("addCount", addCount),
		zap.Any("assignPos", assignPos),
		zap.Any("equipMaxMap", equipMaxMap),
		zap.Any("equipMap", equipMap),
		zap.Any("equipValueMax", shopCfg.Equip_value_max))
	return equipMap, nil
}

func ConvEquipId(logger fklog.FKLogI, equipId int32, assignPos int32, equipMaxMap, equipValueMax map[int32]int64) (int32, error) {
	convEquipId := equipId
	if assignPos > 0 {
		if assignPos == 1 {
			convEquipId = 427*100000 + (equipId/10000)%10*10000 + assignPos*1000 + equipId%1000
		} else {
			convEquipId = 426*100000 + (equipId/10000)%10*10000 + assignPos*1000 + equipId%1000
		}
	} else {
		// equipPos := equipId / 1000 % 10
		// if equipMaxMap[equipPos] >= equipValueMax[equipPos] {
		// 	posWeight := make(map[int32]int32, 0)
		// 	for k, v := range equipValueMax {
		// 		if equipMaxMap[k] >= v {
		// 			continue
		// 		}
		// 		posWeight[k] = 1000
		// 	}
		// 	posId := randfuncs.RandByWeightV2(logger, posWeight, false)
		// 	if posId <= 0 {
		// 		logger.ErrorWF("ConvEquipId RandByWeightV2 err",
		// 			zap.Any("posWeight", posWeight))
		// 		return 0, errors.New("配置不存在")
		// 	}
		// 	if posId == 1 {
		// 		convEquipId = 427*100000 + (equipId/10000)%10*10000 + posId*1000 + equipId%1000
		// 	} else {
		// 		convEquipId = 426*100000 + (equipId/10000)%10*10000 + posId*1000 + equipId%1000
		// 	}
		// }
	}

	return convEquipId, nil
}

func GetEquipMaxMap(logger fklog.FKLogI, userId uint64, mazeLevel int32, areaId int32) (equipMaxMap map[int32]int64, isFull bool, err error) {
	// newLevel := GetMazeShopLv(mazeLevel, areaId)
	// shopCfg := GMazeShopV8Cfg.Get(newLevel)
	// if shopCfg == nil {
	// 	logger.ErrorWF("GetMazeSeqEquipId GMazeShopV8Cfg fail",
	// 		zap.Int32("newLevel", newLevel))
	// 	return equipMaxMap, false, errors.New("区域id配置不存在")
	// }
	equipMaxMap = make(map[int32]int64)
	// // 获取身上的装备信息
	// assembleInfoMap, err := dollassemblesuitredis.GetAllDollAssembleSuit(logger, userId)
	// if err != nil {
	// 	logger.ErrorWF("GetEquipMaxMap GetAllDollAssembleSuit error", zap.Error(err))
	// 	return equipMaxMap, false, err
	// }
	// equipGuids := make([]int64, 0)
	// for _, equipList := range assembleInfoMap {
	// 	for _, v := range equipList {
	// 		if v.GetEquipGuid() > 0 {
	// 			equipGuids = append(equipGuids, v.GetEquipGuid())
	// 		}
	// 	}
	// }
	// if len(equipGuids) <= 0 {
	// 	return equipMaxMap, false, nil
	// }
	// equipMap, err := mazebagequipredis.GetBatchEquipInfo(logger, userId, equipGuids...)
	// if err != nil {
	// 	logger.ErrorWF("GetEquipMaxMap GetBatchEquipInfo error", zap.Error(err))
	// 	return equipMaxMap, false, err
	// }
	// attrMaxCount := 0
	// for _, equipInfo := range equipMap {
	// 	equipPos := equipInfo.GetEquipId() / 1000 % 10
	// 	if _, ok := shopCfg.Equip_value_max[equipPos]; !ok {
	// 		continue
	// 	}
	// 	for _, attrInfoDb := range equipInfo.BaseAttrs {
	// 		if attrInfoDb.GetIndex() != 1 {
	// 			continue
	// 		}
	// 		for _, showAttrInfo := range attrInfoDb.ShowAttrList {
	// 			if showAttrInfo.GetAttrValue() >= shopCfg.Equip_value_max[equipPos] {
	// 				attrMaxCount++
	// 			}
	// 			equipMaxMap[equipPos] = showAttrInfo.GetAttrValue()
	// 			break
	// 		}
	// 	}
	// }
	// if attrMaxCount >= len(shopCfg.Equip_value_max) {
	// 	isFull = true
	// }
	// logger.InfoWF("GetEquipMaxMap end",
	// 	zap.Int32("newLevel", newLevel),
	// 	zap.Int32("mazeLevel", mazeLevel),
	// 	zap.Int32("areaId", areaId),
	// 	zap.Any("equipMaxMap", equipMaxMap),
	// 	zap.Any("attrMaxCount", attrMaxCount),
	// 	zap.Any("equipValueMax", shopCfg.Equip_value_max))
	return equipMaxMap, isFull, nil
}

func GetMazeShopLv(level int32, areaId int32) int32 {
	// cfg := GMazeShopAreaV8Cfg.Get(areaId)
	// if cfg == nil{
	// 	return level
	// }
	// if level < cfg.Shop_level_min{
	// 	return cfg.Shop_level_min
	// }
	// if level > cfg.Shop_level_max{
	// 	return cfg.Shop_level_max
	// }
	return level
}
