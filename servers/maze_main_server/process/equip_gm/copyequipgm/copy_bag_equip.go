package copyequipgm

import (
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/protodef/MazeEquipCache"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebagequipredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeequipgetnumredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeequipguidredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip_gm/copyequipgm/copyinterface"
)

func CopyBagEquipData(logger fklog.FKLogI, srcUserId uint64, dstUsers []uint64, param copyinterface.CopyParam) error {
	var err error
	var curUser uint64
	defer func() {
		if err != nil {
			logger.ErrorWF("CopyBagEquipData fail", zap.Error(err),
				zap.Uint64("src", srcUserId),
				zap.Uint64("dst", curUser))
		} else {
			logger.InfoWF("CopyBagEquipData succ",
				zap.Uint64("src", srcUserId),
				zap.Int("dstLen", len(dstUsers)))
		}
	}()
	allEquips, err := mazebagequipredis.GetAllEquipInfo(logger, srcUserId)
	if err != nil {
		logger.ErrorWF("CopyBagEquipData GetAllEquipInfo err", zap.Error(err))
		return err
	}
	equipScoreMap, err := mazeequipgetnumredis.GetAllEquipGetNum(logger, srcUserId)
	if err != nil {
		logger.ErrorWF("CopyBagEquipData GetAllEquipGetNum error", zap.Error(err))
		return err
	}

	equipMaxGuid, err := mazeequipguidredis.GetNewGuid(logger, srcUserId, 1)
	if err != nil {
		logger.ErrorWF("CopyBagEquipData GetNewGuid error", zap.Error(err))
		return err
	}

	for _, dstId := range dstUsers {
		// 拷贝装备数据
		curUser = dstId
		err = mazebagequipredis.GMDelEquip(logger, dstId)
		if err != nil {
			logger.ErrorWF("CopyBagEquipData GMDelEquip err", zap.Error(err))
			return err
		}
		equipList := make([]*MazeEquipCache.MazeEquipInfoDb, 0)
		for _, equipInfo := range allEquips {
			equipList = append(equipList, equipInfo)
		}
		err = mazeequipguidredis.SetEquipGuidGm(logger, dstId, equipMaxGuid)
		if err != nil {
			logger.ErrorWF("CopyBagEquipData SetEquipGuidGm error", zap.Error(err), zap.Any("equipMaxGuid", equipMaxGuid))
			return err
		}
		err = mazebagequipredis.BatchSaveEquipInfo(logger, dstId, equipList)
		if err != nil {
			logger.ErrorWF("CopyBagEquipData BatchSaveEquipInfo error", zap.Error(err), zap.Any("equipList", equipList))
			return err
		}
		err = mazeequipgetnumredis.BatchSetEquipGetNum(logger, dstId, equipScoreMap)
		if err != nil {
			logger.ErrorWF("CopyBagEquipData BatchSetEquipGetNum error", zap.Error(err), zap.Any("equipScoreMap", equipScoreMap))
			return err
		}
	}
	return nil
}
