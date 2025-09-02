package allianceservice

import (
	"context"
	"maze_game_server/model/alliancemodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) QueryAllianceInfo(ctx context.Context, allianceID int32) (*alliancemodel.AllianceInfoModel, error) {
	allianceID = 1
	logger := fklog.ContextAppLogger(ctx)
	allianceInfoModel, err := alliancemodel.LoadAllianceInfoModel(ctx, allianceID)
	if err != nil {
		logger.CtxError(ctx, "QueryAllianceInfo LoadAllianceInfoModel err",
			zap.Int32("allianceID", allianceID), zap.Error(err))
		return nil, err
	}
	return allianceInfoModel, nil
}

func (s *service) QueryAllianceList(ctx context.Context) (*alliancemodel.AllianceListModel, error) {
	logger := fklog.ContextAppLogger(ctx)
	allianceListModel, err := alliancemodel.LoadAllianceListModel(ctx)
	if err != nil {
		logger.CtxError(ctx, "QueryAllianceList LoadAllianceListModel err", zap.Error(err))
		return nil, err
	}
	return allianceListModel, nil
}
