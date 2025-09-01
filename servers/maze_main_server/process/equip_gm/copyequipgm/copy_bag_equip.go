package copyequipgm

import (
	"context"
	"maze_game_server/io/redis/mazebagequipredis"
	"maze_game_server/io/redis/mazeequipgetnumredis"
	"maze_game_server/io/redis/mazeequipguidredis"
	"maze_game_server/pb/server/MazeEquipCache"
	"maze_game_server/servers/maze_main_server/process/equip_gm/copyequipgm/copyinterface"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func CopyBagEquipData(ctx context.Context, srcUserId uint64, dstUsers []uint64, param copyinterface.CopyParam) error {
	logger := fklog.ContextAppLogger(ctx)
	var err error
	var curUser uint64
	defer func() {
		if err != nil {
			logger.CtxError(ctx, "CopyBagEquipData fail", zap.Error(err),
				zap.Uint64("src", srcUserId),
				zap.Uint64("dst", curUser))
		} else {
			logger.CtxInfo(ctx, "CopyBagEquipData succ",
				zap.Uint64("src", srcUserId),
				zap.Int("dstLen", len(dstUsers)))
		}
	}()
	allEquips, err := mazebagequipredis.GetAllEquipInfo(ctx, srcUserId)
	if err != nil {
		logger.CtxError(ctx, "CopyBagEquipData GetAllEquipInfo err", zap.Error(err))
		return err
	}
	equipScoreMap, err := mazeequipgetnumredis.GetAllEquipGetNum(ctx, srcUserId)
	if err != nil {
		logger.CtxError(ctx, "CopyBagEquipData GetAllEquipGetNum error", zap.Error(err))
		return err
	}

	equipMaxGuid, err := mazeequipguidredis.GetNewGuid(ctx, srcUserId, 1)
	if err != nil {
		logger.CtxError(ctx, "CopyBagEquipData GetNewGuid error", zap.Error(err))
		return err
	}

	for _, dstId := range dstUsers {
		// 拷贝装备数据
		curUser = dstId
		err = mazebagequipredis.GMDelEquip(ctx, dstId)
		if err != nil {
			logger.CtxError(ctx, "CopyBagEquipData GMDelEquip err", zap.Error(err))
			return err
		}
		equipList := make([]*MazeEquipCache.MazeEquipInfoDb, 0)
		for _, equipInfo := range allEquips {
			equipList = append(equipList, equipInfo)
		}
		err = mazeequipguidredis.SetEquipGuidGm(ctx, dstId, equipMaxGuid)
		if err != nil {
			logger.CtxError(ctx, "CopyBagEquipData SetEquipGuidGm error", zap.Error(err), zap.Any("equipMaxGuid", equipMaxGuid))
			return err
		}
		err = mazebagequipredis.BatchSaveEquipInfo(ctx, dstId, equipList)
		if err != nil {
			logger.CtxError(ctx, "CopyBagEquipData BatchSaveEquipInfo error", zap.Error(err), zap.Any("equipList", equipList))
			return err
		}
		err = mazeequipgetnumredis.BatchSetEquipGetNum(ctx, dstId, equipScoreMap)
		if err != nil {
			logger.CtxError(ctx, "CopyBagEquipData BatchSetEquipGetNum error", zap.Error(err), zap.Any("equipScoreMap", equipScoreMap))
			return err
		}
	}
	return nil
}
