package allianceservice

import (
	"context"
	"maze_game_server/model/alliancemodel"
	"maze_game_server/model/familymodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) QueryUserAlliance(ctx context.Context, userID uint64) (allianceID int32, err error) {
	logger := fklog.ContextAppLogger(ctx)
	userFamilyModel := familymodel.NewUserFamilyModel(ctx, userID)
	familyID, err := userFamilyModel.GetUserFamily(ctx)
	if err != nil {
		logger.CtxError(ctx, "QueryUserAlliance GetUserFamily err",
			zap.Uint64("userID", userID), zap.Error(err))
		return 0, err
	}

	familyToAllianceModel := alliancemodel.NewFamilyToAllianceModel(ctx, familyID)
	allianceID, err = familyToAllianceModel.GetUserAlliance(ctx)
	if err != nil {
		logger.CtxError(ctx, "QueryUserAlliance QueryUserAlliance err",
			zap.Uint64("userID", userID), zap.Error(err))
		return 0, err
	}
	return allianceID, nil
}
