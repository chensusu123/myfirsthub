package userprofileservice

import (
	"context"
	"maze_game_server/model/userprofilemodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// GetUserProfile 查询用户信息
func (s *service) GetUserProfile(ctx context.Context, userID uint64) (profile *userprofilemodel.UserProfileModel, err error) {
	logger := fklog.ContextAppLogger(ctx)
	model, err := userprofilemodel.LoadUserProfileModel(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "GetUserProfile LoadUserProfileModel error", zap.Any("userID", userID),
			zap.Error(err))
		return
	}
	return model, nil
}
