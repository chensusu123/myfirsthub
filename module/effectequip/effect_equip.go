/*
 * @Author: majian
 * @Date: 2024-12-26 20:50:12
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-07 17:40:33
 */
package effectequip

import (
	"context"
	"errors"

	"maze_game_server/io/redis/mazebagequipredis"
	"maze_game_server/pb/server/MazeEquipCache"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

var EquipNoExist = errors.New("equip not exist")

func GetEffectEquipInfo(ctx context.Context, userId uint64, equipGuid int64) (effectEquip *MazeEquipCache.MazeEquipInfoDb, err error) {
	logger := fklog.ContextAppLogger(ctx)
	effectEquip, err = mazebagequipredis.GetEquipInfo(ctx, userId, equipGuid)
	if err != nil {
		logger.CtxError(ctx, "GetEffectEquipInfo GetEquipInfo fail", zap.Error(err),
			zap.Int64("guid", equipGuid))
		return nil, err
	}

	if effectEquip == nil || effectEquip.GetEquipGuid() <= 0 {
		logger.CtxError(ctx, "GetEffectEquipInfo equip not exist", zap.Int64("equipGuid", equipGuid))
		return nil, EquipNoExist
	}

	// effectEquip, err = pbutil.ConvertIdentifyEquipDb(logger, equipInfo)
	// if err != nil {
	//	logger.CtxError(ctx,"GetEffectEquipInfo ConvertIdentifyEquipDb fail",
	//		zap.Error(err),
	//		zap.Any("equipInfo", equipInfo))
	//
	//	return nil, err
	// }
	logger.CtxDebug(ctx, "GetEffectEquipInfo succ",
		zap.Any("effectEquip", effectEquip))
	return effectEquip, nil
}

func BatchGetEffectEquipInfo(ctx context.Context, userId uint64, equipGuids ...int64) (effectEquipMap map[int64]*MazeEquipCache.MazeEquipInfoDb, err error) {
	logger := fklog.ContextAppLogger(ctx)
	// 从背包查询装备信息
	effectEquipMap, err = mazebagequipredis.GetBatchEquipInfo(ctx, userId, equipGuids...)
	if err != nil {
		logger.CtxError(ctx, "BatchGetEffectEquipInfo GetBatchEquipInfo fail", zap.Error(err),
			zap.Any("equipGuids", equipGuids))
		return
	}
	// effectEquipMap, err = pbutil.BatchConvertIdentifyEquipDb(logger, equipDatails)
	// if err != nil {
	//	logger.CtxError(ctx,"BatchGetEffectEquipInfo BatchConvertIdentifyEquipDb fail", zap.Error(err),
	//		zap.Any("equipDatails", equipDatails))
	//	return
	// }
	logger.CtxDebug(ctx, "BatchGetEffectEquipInfo succ",
		zap.Any("effectEquip", effectEquipMap))
	return effectEquipMap, nil
}

func GetAllEffectEquipInfo(ctx context.Context, userId uint64) (effectEquipMap map[int64]*MazeEquipCache.MazeEquipInfoDb, err error) {
	logger := fklog.ContextAppLogger(ctx)
	effectEquipMap, err = mazebagequipredis.GetAllEquipInfo(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "GetAllEffectEquipInfo GetAllEquipInfo error", zap.Error(err))
		return nil, err
	}
	// effectEquipMap, err = pbutil.BatchConvertIdentifyEquipDb(logger, equipMap)
	// if err != nil {
	//	logger.CtxError(ctx,"GetAllEffectEquipInfo BatchConvertIdentifyEquipDb fail", zap.Error(err),
	//		zap.Any("equipDatails", equipMap))
	//	return
	// }
	logger.CtxDebug(ctx, "BatchGetEffectEquipInfo succ",
		zap.Any("effectEquip", effectEquipMap))
	return effectEquipMap, nil
}

// 获取实例化装备
// func GetInstanceEffectEquipInfo(ctx context.Context, userId uint64, equipGuid int64) (effectEquip *MazeEquipCache.MazeEquipInfoDb, err error) {
// 	effectEquip, err = dollequipinstanceredis.GetEquipInstance(logger, userId, equipGuid)
// 	if err != nil {
// 		logger.CtxError(ctx,"GetInstanceEffectEquipInfo GetEquipInfo fail", zap.Error(err),
// 			zap.Int64("guid", equipGuid))
// 		return nil, err
// 	}

// 	if effectEquip == nil || effectEquip.GetEquipGuid() <= 0 {
// 		logger.CtxError(ctx,"GetInstanceEffectEquipInfo equip not exist", zap.Int64("equipGuid", equipGuid))
// 		return nil, EquipNoExist
// 	}

// 	//effectEquip, err = pbutil.ConvertIdentifyEquipDb(logger, equipInfo)
// 	//if err != nil {
// 	//	logger.CtxError(ctx,"GetInstanceEffectEquipInfo ConvertIdentifyEquipDb fail",
// 	//		zap.Error(err),
// 	//		zap.Any("equipInfo", equipInfo))
// 	//
// 	//	return nil, err
// 	//}
// 	logger.CtxDebug(ctx,"GetInstanceEffectEquipInfo succ",
// 		zap.Any("effectEquip", effectEquip))
// 	return effectEquip, nil
// }
