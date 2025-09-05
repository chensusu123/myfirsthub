package tempbuffservice

import (
	"context"
	"maze_game_server/model/buffenergy"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) AddTmpBuffEnergy(ctx context.Context, userID uint64, barrierID, areaID, areaIndex int32, energyCount int32) error {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "AddTmpBuffEnergy start",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Int32("areaID", areaID),
		zap.Int32("areaIndex", areaIndex),
	)

	defer func() {
		logger.CtxInfo(ctx, "AddTmpBuffEnergy end",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("areaID", areaID),
			zap.Int32("areaIndex", areaIndex),
		)
	}()

	areaData, err := buffenergy.NewBuffEnergy(ctx, userID, barrierID, areaID, areaIndex)
	if err != nil {
		logger.CtxError(ctx, "AddTmpBuffEnergy NewBuffEnergy fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("areaID", areaID),
			zap.Int32("areaIndex", areaIndex),
		)
		return err
	}

	areaData.NowEnergy += energyCount

	// todo 到当前等级满能量了 推一个三选一列表
	return err
}
