package barrieritemservice

import (
	"context"
	"fmt"
	"maze_game_server/excel/mazebarriesv8config"
	"maze_game_server/model/barrieritemsmodel"
	"maze_game_server/servers/maze_main_server/process/item"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) AddEquipScore(ctx context.Context, userID uint64, barrierID int32, score int32, monsterGuid int64, monsterPos string) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "AddEquipScore Start",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Int32("score", score),
		zap.Int64("monsterGuid", monsterGuid),
		zap.String("monsterPos", monsterPos),
	)

	defer func() {
		logger.CtxInfo(ctx, "AddEquipScore End",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("score", score),
			zap.Int64("monsterGuid", monsterGuid),
			zap.String("monsterPos", monsterPos),
		)
	}()

	data, err := barrieritemsmodel.NewBarrierItems(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "AddEquipScore NewBarrierItems Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("score", score),
			zap.Int64("monsterGuid", monsterGuid),
			zap.String("monsterPos", monsterPos),
			zap.Any("data", data),
			zap.Error(err),
		)
		return err
	}

	barrierCfg := mazebarriesv8config.GetStageConfig(ctx, barrierID)
	if barrierCfg.Need_equip_score == 0 {
		logger.CtxError(ctx, "AddEquipScore barrierCfg.Need_equip_score Equal zero",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("score", score),
			zap.Int64("monsterGuid", monsterGuid),
			zap.String("monsterPos", monsterPos),
			zap.Any("data", data),
		)
		return nil
	}

	nowScore := data.EquipScore + score
	fmt.Println(nowScore)

	equipNum := nowScore / barrierCfg.Need_equip_score
	// if equipNum
	data.EquipScore = nowScore % barrierCfg.Need_equip_score
	equips, err := s.FallOffEquip(ctx, userID, barrierID, equipNum)
	if err != nil {
		logger.CtxError(ctx, "AddEquipScore DropEquip Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("score", score),
			zap.Int64("monsterGuid", monsterGuid),
			zap.String("monsterPos", monsterPos),
			zap.Any("data", data),
			zap.Error(err),
		)
		return err
	}

	for _, equip := range equips {
		if equip.Count > 0 {
			data.Equips[uint64(equip.ItemId)] += uint64(equip.Count)
		}
	}

	err = data.Save(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "AddEquipScore Save Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("score", score),
			zap.Int64("monsterGuid", monsterGuid),
			zap.String("monsterPos", monsterPos),
			zap.Any("data", data),
			zap.Error(err),
		)
		return
	}

	// 推包
	if len(equips) > 0 {
		err = item.OnSendItemsPack(ctx, userID, nil, equips, monsterGuid, monsterPos)
		if err != nil {
			logger.CtxWarn(ctx, "AddEquipScore OnSendItemsPack Fail",
				zap.Uint64("userID", userID),
				zap.Int32("barrierID", barrierID),
				zap.Int32("score", score),
				zap.Int64("monsterGuid", monsterGuid),
				zap.String("monsterPos", monsterPos),
				zap.Any("equips", equips),
				zap.Error(err),
			)
			return
		}
	}

	return err
}
