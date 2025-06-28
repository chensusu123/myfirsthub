package allianceservice

import (
	"maze_game_server/model/alliancemodel"
	"maze_game_server/model/familymodel"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
)

func (s *service) QueryUserAlliance(logger fklog.FKLogI, userID uint64) (allianceID int32, err error) {
	userFamilyModel := familymodel.NewUserFamilyModel(logger, userID)
	familyID, err := userFamilyModel.GetUserFamily(logger)
	if err != nil {
		logger.ErrorWF("QueryUserAlliance GetUserFamily err",
			zap.Uint64("userID", userID), zap.Error(err))
		return 0, err
	}

	familyToAllianceModel := alliancemodel.NewFamilyToAllianceModel(logger, familyID)
	allianceID, err = familyToAllianceModel.GetUserAlliance(logger)
	if err != nil {
		logger.ErrorWF("QueryUserAlliance QueryUserAlliance err",
			zap.Uint64("userID", userID), zap.Error(err))
		return 0, err
	}
	return allianceID, nil
}
