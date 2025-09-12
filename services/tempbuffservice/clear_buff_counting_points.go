package tempbuffservice

import (
	"context"
	"maze_game_server/excel/dollmappuzzlenewcfgex"
	"maze_game_server/model/buffenergy"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) ClearBuffCountingPoints(ctx context.Context, userID uint64, barrierID int32, stageID int32) error {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "ClearBuffCountingPoints start",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Int32("stageID", stageID),
	)

	defer func() {
		logger.CtxInfo(ctx, "ClearBuffCountingPoints end",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("stageID", stageID),
		)
	}()

	unPassAreas := dollmappuzzlenewcfgex.GetUnPassAreaInfos(barrierID, stageID)

	for _, unpassArea := range unPassAreas {
		areaData, err := buffenergy.NewBuffEnergy(ctx, userID, barrierID, unpassArea.AreaId, unpassArea.AreaIndex)
		if err != nil {
			logger.CtxError(ctx, "ClearBuffCountingPoints NewBuffEnergy Fail",
				zap.Uint64("userID", userID),
				zap.Int32("barrierID", barrierID),
				zap.Int32("stageID", stageID),
				zap.Int32("areaID", unpassArea.AreaId),
				zap.Int32("areaIndex", unpassArea.AreaIndex),
			)
			return err
		}

		err = areaData.Del(ctx, userID, barrierID, unpassArea.AreaId, unpassArea.AreaIndex)
		if err != nil {
			logger.CtxError(ctx, "ClearBuffCountingPoints Del Fail",
				zap.Uint64("userID", userID),
				zap.Int32("barrierID", barrierID),
				zap.Int32("stageID", stageID),
				zap.Int32("areaID", unpassArea.AreaId),
				zap.Int32("areaIndex", unpassArea.AreaIndex),
			)
			return err
		}
	}

	return nil
}
