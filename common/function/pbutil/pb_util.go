package pbutil

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipCache"
	"gitlab.ifreetalk.com/maze/maze_game_server/config/GMazeEquipSuiteNameV8Cfg"
	"gitlab.ifreetalk.com/maze/maze_game_server/config/GMazeItemsV8Cfg"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeequipconfigv8"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeequiptyperesv8"
	"go.uber.org/zap"
)

func GetDollEquipName(equipInfo *MazeEquipCache.MazeEquipInfoDb, equipType int32) (equipResId int32, equipName string) {
	equipId := equipInfo.GetEquipId()
	if equipInfo.GetEquipSubType() > 0 {
		equipType = equipInfo.GetEquipSubType()
	}
	equipTypeResCfg := mazeequiptyperesv8.GetEquipTypeResCfg(equipId, 0, 0)
	if equipTypeResCfg != nil {
		equipName = equipTypeResCfg.Name
		equipResId = equipTypeResCfg.Order
	}
	if equipInfo.GetSuitId() > 0 {
		equipSuiteNameCfg := GMazeEquipSuiteNameV8Cfg.Get(equipId)
		if equipSuiteNameCfg != nil {
			equipName = equipSuiteNameCfg.Suite_equip_name[equipInfo.GetSuitId()]
		}
	}
	if equipName == "" {
		itemCfg := GMazeItemsV8Cfg.Get(equipId)
		if itemCfg != nil {
			equipName = itemCfg.Prop_name
		}
	}
	return equipResId, equipName
}

func GetDollEquipNameEx(equipInfo *MazeEquipCache.MazeEquipInfoDb, equipType int32) (equipResId int32, equipName string, mazeModel int32, icon string, iconAtlas string) {
	equipId := equipInfo.GetEquipId()
	if equipInfo.GetEquipSubType() > 0 {
		equipType = equipInfo.GetEquipSubType()
	}
	equipTypeResCfg := mazeequiptyperesv8.GetEquipTypeResCfg(equipId, 0, 0)
	if equipTypeResCfg != nil {
		equipName = equipTypeResCfg.Name
		equipResId = equipTypeResCfg.Order
		mazeModel = equipTypeResCfg.Maze_model
		icon = equipTypeResCfg.Icon
		iconAtlas = equipTypeResCfg.IconAtlas
	}
	if equipInfo.GetSuitId() > 0 {
		equipSuiteNameCfg := GMazeEquipSuiteNameV8Cfg.Get(equipId)
		if equipSuiteNameCfg != nil {
			equipName = equipSuiteNameCfg.Suite_equip_name[equipInfo.GetSuitId()]
		}
	}
	if equipName == "" {
		itemCfg := GMazeItemsV8Cfg.Get(equipId)
		if itemCfg != nil {
			equipName = itemCfg.Prop_name
		}
	}
	return
}

func GetEquipResId(logger fklog.FKLogI, equipInfo *MazeEquipCache.MazeEquipInfoDb) (equipResId int32) {
	equipId := equipInfo.GetEquipId()
	if equipInfo.GetEquipSubType() == 0 {
		return
	}

	var hurtId int32
	for _, attrInfo := range equipInfo.GetBaseAttrs() {
		for _, realAttr := range attrInfo.RealAttrList {
			if mazeequipconfigv8.GetSuitEffectType(realAttr.GetAttrId()) > 0 {
				hurtId = realAttr.GetAttrId()
				break
			}
		}
	}
	if hurtId == 0 {
		logger.WarnWF("GetEquipResId cannot find hurtId",
			zap.Int32("cfgId", equipInfo.GetEquipId()),
			zap.Int64("guid", equipInfo.GetEquipGuid()),
			zap.Any("baseAttrs", equipInfo.GetBaseAttrs()))
		return
	}
	equipTypeResCfg := mazeequiptyperesv8.GetEquipTypeResCfg(equipId, equipInfo.GetEquipSubType(), hurtId)
	if equipTypeResCfg != nil {
		equipResId = equipTypeResCfg.Order
	} else {
		logger.WarnWF("GetEquipResId cannot find resID",
			zap.Int32("cfgId", equipInfo.GetEquipId()),
			zap.Int64("guid", equipInfo.GetEquipGuid()),
			zap.Any("subType", equipInfo.GetEquipSubType()),
			zap.Int32("hurtId", hurtId))
	}
	return equipResId
}

// func ConvertIdentifyEquipDb(logger fklog.FKLogI, equipInfo *MazeEquipCache.MazeEquipInfoDb) (*MazeEquipCache.MazeEquipInfoDb, error) {
//	return equipInfo, nil
// }

// // 批量转换接口
// func BatchConvertIdentifyEquipDb(logger fklog.FKLogI, originalEquips map[int64]*MazeEquipCache.MazeEquipInfoDb) (realEquipMap map[int64]*MazeEquipCache.MazeEquipInfoDb, err error) {
//	realEquipMap = make(map[int64]*MazeEquipCache.MazeEquipInfoDb)
//	for k, v := range originalEquips {
//		r, e := ConvertIdentifyEquipDb(logger, v)
//		if e != nil {
//			err = e
//			return nil, err
//		}
//		realEquipMap[k] = r
//	}
//	return realEquipMap, err
// }
