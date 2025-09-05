package userprofileservice

import (
	"context"
	"maze_game_server/model/userprofilemodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// GetUserProfile 查询多个用户信息
func (s *service) GetBatchUserProfile(ctx context.Context, users []uint64) (profiles []*userprofilemodel.UserProfileModel, err error) {
	logger := fklog.ContextAppLogger(ctx)
	if len(users) <= 0 {
		return
	}
	models, err := userprofilemodel.LoadBatchUserProfileModel(ctx, users)
	if err != nil {
		logger.CtxError(ctx, "GetBatchUserProfile LoadBatchUserProfileModel error", zap.Any("users", users),
			zap.Error(err))
		return
	}
	for _, v := range models {
		profiles = append(profiles, v)
	}
	return
}
