package equipbaggm

import (
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeEquipAffixRandPoolV8Cfg"
	"maze_game_server/config/GMazeEquipAffixRollTypeV8Cfg"
	"maze_game_server/io/redis/mazebagequipredis"
	"maze_game_server/pb/server/MazeEquipCache"
	"maze_game_server/servers/maze_main_server/process/equip"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func fixAllEquipAttrLimit(logger fklog.FKLogI, userId uint64) error {
	equipInfoMap, err := mazebagequipredis.GetAllEquipInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("fixAllEquipAttrLimit GetAllEquipInfo fail", zap.Error(err))
		return err
	}
	chgEquipList := make([]*MazeEquipCache.MazeEquipInfoDb, 0)
	for _, equipInfo := range equipInfoMap {
		newEquipInfo, chgEquip, err := fixEquipAttr(logger, userId, equipInfo)
		if err != nil {
			logger.ErrorWF("fixAllEquipAttrLimit fixEquipAttr fail", zap.Error(err))
			return err
		}
		if chgEquip {
			chgEquipList = append(chgEquipList, newEquipInfo)
		}
	}
	if len(chgEquipList) > 0 {
		mazebagequipredis.BatchSaveEquipInfo(logger, userId, chgEquipList)
	}
	return err
}

func fixEquipAttr(logger fklog.FKLogI, userId uint64, equipInfo *MazeEquipCache.MazeEquipInfoDb) (*MazeEquipCache.MazeEquipInfoDb, bool, error) {
	chgEquip := false
	for _, attrInfo := range equipInfo.BaseAttrs {
		cfg := GMazeEquipAffixRandPoolV8Cfg.Get(attrInfo.GetAttrGroup())
		if cfg == nil {
			logger.ErrorWF("FixEquipAttr GMazeEquipAffixRandPoolV8Cfg fail",
				zap.Any("attrGroupId", attrInfo.GetAttrGroup()))
			return nil, false, errors.New("属性配置不存在")
		}
		rollTypeCfg := GMazeEquipAffixRollTypeV8Cfg.Get(cfg.Roll_type)
		if rollTypeCfg == nil {
			logger.ErrorWF("FixEquipAttr GMazeEquipAffixRandPoolV8Cfg err",
				zap.Int32("rollId", cfg.Roll_type))
			return nil, false, errors.New("配置不存在")
		}
		calcAddCount, calcRangeCount := equip.CalcAttrRealRoll(logger, cfg.Add_attr_min, cfg.Add_attr_max, cfg.Show_attr_min, cfg.Show_attr_max, int64(attrInfo.GetRandWeight()), rollTypeCfg.Round_value)
		for index := 0; index <= len(attrInfo.ShowAttrList); index++ {
			showAttrInfo := attrInfo.ShowAttrList[index]
			attrValue := equip.CalcAttrValByRoll(cfg.Show_attr_min[showAttrInfo.GetAttrId()], cfg.Show_attr_max[showAttrInfo.GetAttrId()], rollTypeCfg.Round_value, calcAddCount, calcRangeCount)
			if showAttrInfo.GetAttrValue() != attrValue {
				showAttrInfo.AttrValue = proto.Int64(attrValue)
				chgEquip = true
			}
		}
		for index := 0; index <= len(attrInfo.RealAttrList); index++ {
			realAttrInfo := attrInfo.RealAttrList[index]
			attrValue := equip.CalcAttrValByRoll(cfg.Add_attr_min[realAttrInfo.GetAttrId()], cfg.Add_attr_max[realAttrInfo.GetAttrId()], rollTypeCfg.Round_value, calcAddCount, calcRangeCount)
			if realAttrInfo.GetAttrValue() != attrValue {
				realAttrInfo.AttrValue = proto.Int64(attrValue)
				chgEquip = true
			}
		}
	}

	return equipInfo, chgEquip, nil
}

// func resetInsAllEquip(logger fklog.FKLogI, userId uint64) error {
// 	equipInfoMap, err := mazebagequipredis.GetAllEquipInfo(logger, userId)
// 	if err != nil {
// 		logger.ErrorWF("fixAllEquipAttrLimit GetAllEquipInfo fail", zap.Error(err))
// 		return err
// 	}
// 	chgEquipList := make([]*MazeEquipCache.MazeEquipInfoDb, 0)
// 	for _, equipInfo := range equipInfoMap {
// 		newEquipInfo, err := resetInsEquip(logger, userId, equipInfo)
// 		if err != nil {
// 			logger.ErrorWF("fixAllEquipAttrLimit fixEquipAttr fail", zap.Error(err))
// 			return err
// 		}
// 		chgEquipList = append(chgEquipList, newEquipInfo)
// 	}
// 	if len(chgEquipList) > 0 {
// 		mazebagequipredis.BatchSaveEquipInfo(logger, userId, chgEquipList)
// 	}
// 	return err
// }

// func resetInsEquip(logger fklog.FKLogI, userId uint64, equipInfo *MazeEquipCache.MazeEquipInfoDb) (*MazeEquipCache.MazeEquipInfoDb, error) {
// 	// 装备实例化开始
// 	up := &process.DEIUWParam{}
// 	up.ChargeStage = 1
// 	up.OpType = 0
// 	up.TradeNum = 0
// 	// up.BarrierMap = make(map[int32]int64)
// 	up.PoolLimitMap = process.GetPoolLimitMap(logger, mazeequipconfigv8.GetEquipGoodAttrIds())
// 	up.PoolLimitBase2Map = process.GetPoolLimitMap(logger, mazeequipconfigv8.GetEquipBadAttrIds())
// 	equipSvr := &MazeEquipSvr.SvrEquipInfo{
// 		EquipId: proto.Int32(equipInfo.GetEquipId()),
// 	}
// 	ep := process.NewDEIEWParam(equipSvr)
// 	ep.Guid = equipInfo.GetEquipGuid()
// 	ep.Score = equipInfo.GetEquipScore()
// 	r, e := process.DoInsEquip(logger, userId, up, ep)
// 	if e != nil {
// 		logger.ErrorWF("ResetEquipAttr DoInsEquip fail", zap.Error(e))
// 		return nil, e
// 	}
// 	equipInfo.BaseAttrs = r.EquipInfo.BaseAttrs
// 	equipInfo.ModAttrs = r.EquipInfo.ModAttrs
// 	equipInfo.SpecialAttrs = r.EquipInfo.SpecialAttrs
// 	return equipInfo, nil
// }
