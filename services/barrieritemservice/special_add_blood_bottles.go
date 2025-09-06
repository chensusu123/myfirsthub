package barrieritemservice

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/model/barrieritemsmodel"
	"maze_game_server/services/itemservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) SpecialAddBloodBottles(ctx context.Context, userID uint64, barrierID int32) error {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "SpecialAddBloodBottles Start",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
	)

	defer func() {
		logger.CtxInfo(ctx, "SpecialAddBloodBottles End",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
		)
	}()

	data, err := barrieritemsmodel.NewBarrierItems(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "SpecialAddBloodBottles NewBarrierItems Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Error(err),
		)
		return err
	}

	logger.CtxInfo(ctx, "SpecialAddBloodBottles GetData Successful",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Any("nowdata", data),
	)

	// bloodMap := mazeconfigv8config.GetMazeConfig(ctx, constdef.MazeCfgId951)

	var dropItem []*itemservice.ItemInfo
	data.Items[constdef.BloodBottleID] += 1
	dropItem = append(dropItem, &itemservice.ItemInfo{
		ItemId: constdef.BloodBottleID,
	})

	err = data.Save(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "SpecialAddBloodBottles Save",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Any("nowdata", data),
			zap.Any("adddata", dropItem),
		)
		return err
	}

	logger.CtxInfo(ctx, "SpecialAddBloodBottles Add Successful",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Any("nowdata", data),
		zap.Any("adddata", dropItem),
	)

	return err
}
