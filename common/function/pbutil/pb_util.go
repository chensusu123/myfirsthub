package pbutil

import (
	"context"
	"maze_game_server/config/GMazeEquipSuiteNameV8Cfg"
	"maze_game_server/config/GMazeItemsV8Cfg"
	"maze_game_server/excel/mazeequipconfigv8"
	"maze_game_server/excel/mazeequiptyperesv8"
	"maze_game_server/pb/server/MazeEquipCache"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func GetDollEquipName(ctx context.Context, equipInfo *MazeEquipCache.MazeEquipInfoDb, equipType int32) (equipResId int32, equipName string) {
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
		equipSuiteNameCfg := GMazeEquipSuiteNameV8Cfg.GetWithCtx(ctx, equipId)
		if equipSuiteNameCfg != nil {
			equipName = equipSuiteNameCfg.Suite_equip_name[equipInfo.GetSuitId()]
		}
	}
	if equipName == "" && equipId > int32(0) {
		itemCfg := GMazeItemsV8Cfg.GetWithCtx(ctx, equipId)
		if itemCfg != nil {
			equipName = itemCfg.Prop_name
		}
	}

	logger := fklog.ContextAppLogger(ctx)
	logger.CtxWarn(ctx, "GetDollEquipName  equipName",
		zap.Int32("cfgId", equipInfo.GetEquipId()),
		zap.Int64("guid", equipInfo.GetEquipGuid()),
		zap.Int32("equipType", equipType))
	return equipResId, equipName
}

func GetDollEquipNameEx(ctx context.Context, equipInfo *MazeEquipCache.MazeEquipInfoDb, equipType int32) (equipResId int32, equipName string, mazeModel int32, icon string, iconAtlas string) {
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
		equipSuiteNameCfg := GMazeEquipSuiteNameV8Cfg.GetWithCtx(ctx, equipId)
		if equipSuiteNameCfg != nil {
			equipName = equipSuiteNameCfg.Suite_equip_name[equipInfo.GetSuitId()]
		}
	}
	if equipName == "" {
		itemCfg := GMazeItemsV8Cfg.GetWithCtx(ctx, equipId)
		if itemCfg != nil {
			equipName = itemCfg.Prop_name
		}
	}
	return
}

func GetEquipResId(ctx context.Context, equipInfo *MazeEquipCache.MazeEquipInfoDb) (equipResId int32) {
	logger := fklog.ContextAppLogger(ctx)
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
		logger.CtxWarn(ctx, "GetEquipResId cannot find hurtId",
			zap.Int32("cfgId", equipInfo.GetEquipId()),
			zap.Int64("guid", equipInfo.GetEquipGuid()),
			zap.Any("baseAttrs", equipInfo.GetBaseAttrs()))
		return
	}
	equipTypeResCfg := mazeequiptyperesv8.GetEquipTypeResCfg(equipId, equipInfo.GetEquipSubType(), hurtId)
	if equipTypeResCfg != nil {
		equipResId = equipTypeResCfg.Order
	} else {
		logger.CtxWarn(ctx, "GetEquipResId cannot find resID",
			zap.Int32("cfgId", equipInfo.GetEquipId()),
			zap.Int64("guid", equipInfo.GetEquipGuid()),
			zap.Any("subType", equipInfo.GetEquipSubType()),
			zap.Int32("hurtId", hurtId))
	}
	return equipResId
}

// func ConvertIdentifyEquipDb(ctx context.Context, equipInfo *MazeEquipCache.MazeEquipInfoDb) (*MazeEquipCache.MazeEquipInfoDb, error) {
//	return equipInfo, nil
// }

// // 批量转换接口
// func BatchConvertIdentifyEquipDb(ctx context.Context, originalEquips map[int64]*MazeEquipCache.MazeEquipInfoDb) (realEquipMap map[int64]*MazeEquipCache.MazeEquipInfoDb, err error) {
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
