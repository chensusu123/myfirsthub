package allianceservice

import (
	"context"
	"maze_game_server/model/alliancemodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) ChangeAllianceInfo(ctx context.Context, allianceID int32, allianceName string) error {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "ChangeAllianceInfo Start",
		zap.Any("allianceName", allianceName),
	)
	defer func() {
		logger.CtxInfo(ctx, "ChangeAllianceInfo End",
			zap.Any("allianceName", allianceName),
		)
	}()

	allianceInfoModel, err := alliancemodel.LoadAllianceInfoModel(ctx, allianceID)
	if err != nil {
		logger.CtxError(ctx, "QueryAllianceInfo LoadAllianceInfoModel err",
			zap.Int32("allianceID", allianceID), zap.Error(err))
		return err
	}

	allianceInfoModel.AllianceName = allianceName

	return allianceInfoModel.Save(ctx)
}
