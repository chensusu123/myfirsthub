package allianceservice

import (
	"maze_game_server/model/alliancemodel"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
)

func (s *service) AddAlliance(logger fklog.FKLogI, allianceName string) error {
	allianceListModel, err := alliancemodel.LoadAllianceListModel(logger)
	if err != nil {
		logger.ErrorWF("AddAlliance LoadAllianceListModel err", zap.Error(err))
		return err
	}
	allianceID := allianceListModel.GetAllianceID()

	allianceInfoModel := alliancemodel.NewAllianceInfoModel(logger, allianceID, allianceName)
	err = allianceInfoModel.Save(logger)
	if err != nil {
		logger.ErrorWF("AddAlliance Save allianceInfoModel err", zap.Error(err))
		return err
	}

	err = allianceListModel.AddAlliance(logger, allianceID)
	if err != nil {
		logger.ErrorWF("AddAlliance AddAlliance err", zap.Error(err))
		return err
	}
	return nil
}
