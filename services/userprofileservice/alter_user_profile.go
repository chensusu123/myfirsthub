package userprofileservice

import (
	"context"
	"fmt"
	"maze_game_server/model/userprofilemodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// AlterUserProfile
func (s *service) AlterUserProfile(ctx context.Context, userId uint64, alterProfile *userprofilemodel.UserProfileModel) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	if alterProfile == nil {
		return fmt.Errorf("AlterUserProfile incorrent, userID:%v, alterProfile:%v", userId, alterProfile)
	}

	model, err := userprofilemodel.LoadUserProfileModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "AlterUserProfile LoadUserProfileModel error", zap.Any("userId", userId),
			zap.Error(err))
		return
	}

	// check
	if model != nil && model.UserID != userId {
		logger.CtxError(ctx, "AlterUserProfile check userId error", zap.Any("userId", userId),
			zap.Error(err))
		return
	}

	// alter
	model.FillModelData(alterProfile)
	model.Save(ctx, userId)

	return
}
