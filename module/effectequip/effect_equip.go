/*
 * @Author: majian
 * @Date: 2024-12-26 20:50:12
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-07 17:40:33
 */
package effectequip

import (
	"errors"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/redis/mazebagequipredis"
	"maze_game_server/pb/server/MazeEquipCache"
)

var EquipNoExist = errors.New("equip not exist")

func GetEffectEquipInfo(logger fklog.FKLogI, userId uint64, equipGuid int64) (effectEquip *MazeEquipCache.MazeEquipInfoDb, err error) {
	effectEquip, err = mazebagequipredis.GetEquipInfo(logger, userId, equipGuid)
	if err != nil {
		logger.ErrorWF("GetEffectEquipInfo GetEquipInfo fail", zap.Error(err),
			zap.Int64("guid", equipGuid))
		return nil, err
	}

	if effectEquip == nil || effectEquip.GetEquipGuid() <= 0 {
		logger.ErrorWF("GetEffectEquipInfo equip not exist", zap.Int64("equipGuid", equipGuid))
		return nil, EquipNoExist
	}

	// effectEquip, err = pbutil.ConvertIdentifyEquipDb(logger, equipInfo)
	// if err != nil {
	//	logger.ErrorWF("GetEffectEquipInfo ConvertIdentifyEquipDb fail",
	//		zap.Error(err),
	//		zap.Any("equipInfo", equipInfo))
	//
	//	return nil, err
	// }
	logger.DebugWF("GetEffectEquipInfo succ",
		zap.Any("effectEquip", effectEquip))
	return effectEquip, nil
}

func BatchGetEffectEquipInfo(logger fklog.FKLogI, userId uint64, equipGuids ...int64) (effectEquipMap map[int64]*MazeEquipCache.MazeEquipInfoDb, err error) {
	// 从背包查询装备信息
	effectEquipMap, err = mazebagequipredis.GetBatchEquipInfo(logger, userId, equipGuids...)
	if err != nil {
		logger.ErrorWF("BatchGetEffectEquipInfo GetBatchEquipInfo fail", zap.Error(err),
			zap.Any("equipGuids", equipGuids))
		return
	}
	// effectEquipMap, err = pbutil.BatchConvertIdentifyEquipDb(logger, equipDatails)
	// if err != nil {
	//	logger.ErrorWF("BatchGetEffectEquipInfo BatchConvertIdentifyEquipDb fail", zap.Error(err),
	//		zap.Any("equipDatails", equipDatails))
	//	return
	// }
	logger.DebugWF("BatchGetEffectEquipInfo succ",
		zap.Any("effectEquip", effectEquipMap))
	return effectEquipMap, nil
}

func GetAllEffectEquipInfo(logger fklog.FKLogI, userId uint64) (effectEquipMap map[int64]*MazeEquipCache.MazeEquipInfoDb, err error) {
	effectEquipMap, err = mazebagequipredis.GetAllEquipInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("GetAllEffectEquipInfo GetAllEquipInfo error", zap.Error(err))
		return nil, err
	}
	// effectEquipMap, err = pbutil.BatchConvertIdentifyEquipDb(logger, equipMap)
	// if err != nil {
	//	logger.ErrorWF("GetAllEffectEquipInfo BatchConvertIdentifyEquipDb fail", zap.Error(err),
	//		zap.Any("equipDatails", equipMap))
	//	return
	// }
	logger.DebugWF("BatchGetEffectEquipInfo succ",
		zap.Any("effectEquip", effectEquipMap))
	return effectEquipMap, nil
}

// 获取实例化装备
// func GetInstanceEffectEquipInfo(logger fklog.FKLogI, userId uint64, equipGuid int64) (effectEquip *MazeEquipCache.MazeEquipInfoDb, err error) {
// 	effectEquip, err = dollequipinstanceredis.GetEquipInstance(logger, userId, equipGuid)
// 	if err != nil {
// 		logger.ErrorWF("GetInstanceEffectEquipInfo GetEquipInfo fail", zap.Error(err),
// 			zap.Int64("guid", equipGuid))
// 		return nil, err
// 	}

// 	if effectEquip == nil || effectEquip.GetEquipGuid() <= 0 {
// 		logger.ErrorWF("GetInstanceEffectEquipInfo equip not exist", zap.Int64("equipGuid", equipGuid))
// 		return nil, EquipNoExist
// 	}

// 	//effectEquip, err = pbutil.ConvertIdentifyEquipDb(logger, equipInfo)
// 	//if err != nil {
// 	//	logger.ErrorWF("GetInstanceEffectEquipInfo ConvertIdentifyEquipDb fail",
// 	//		zap.Error(err),
// 	//		zap.Any("equipInfo", equipInfo))
// 	//
// 	//	return nil, err
// 	//}
// 	logger.DebugWF("GetInstanceEffectEquipInfo succ",
// 		zap.Any("effectEquip", effectEquip))
// 	return effectEquip, nil
// }
