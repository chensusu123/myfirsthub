package allianceservice

import (
	"maze_game_server/model/alliancemodel"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
)

func (s *service) QueryAllianceInfo(logger fklog.FKLogI, allianceID int32) (*alliancemodel.AllianceInfoModel, error) {
	allianceInfoModel, err := alliancemodel.LoadAllianceInfoModel(logger, allianceID)
	if err != nil {
		logger.ErrorWF("QueryAllianceInfo LoadAllianceInfoModel err",
			zap.Int32("allianceID", allianceID), zap.Error(err))
		return nil, err
	}
	return allianceInfoModel, nil
}

func (s *service) QueryAllianceList(logger fklog.FKLogI) (*alliancemodel.AllianceListModel, error) {
	allianceListModel, err := alliancemodel.LoadAllianceListModel(logger)
	if err != nil {
		logger.ErrorWF("QueryAllianceList LoadAllianceListModel err", zap.Error(err))
		return nil, err
	}
	return allianceListModel, nil
}
