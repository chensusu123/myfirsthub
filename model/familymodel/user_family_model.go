package familymodel

import (
	"context"
	"time"

	"maze_game_server/io/redis/familyredis"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

type UserFamilyModel struct {
	UserID uint64 `json:"user_id,omitempty"`
}

func NewUserFamilyModel(ctx context.Context, userID uint64) *UserFamilyModel {
	return &UserFamilyModel{
		UserID: userID,
	}
}

// 获取用户所在家族
func (r *UserFamilyModel) GetUserFamily(ctx context.Context) (int32, error) {
	logger := fklog.ContextAppLogger(ctx)
	familyID, err := familyredis.GetUserFamilyID(ctx, r.UserID)
	if err != nil {
		logger.CtxError(ctx, "GetUserFamily GetUserFamilyID err",
			zap.Uint64("userID", r.UserID), zap.Error(err))
		return 0, err
	}
	return familyID, nil
}

// 设置用户上次退出家族时间
func (r *UserFamilyModel) SetUserLastLeaveFamilyTime(ctx context.Context) error {
	logger := fklog.ContextAppLogger(ctx)
	lastLeaveFamilyTime := time.Now().Unix()
	err := familyredis.SetUserLastLeaveFamilyTime(ctx, r.UserID, lastLeaveFamilyTime)
	if err != nil {
		logger.CtxError(ctx, "SetUserLastLeaveFamilyTime SetUserLastLeaveFamilyTime err",
			zap.Uint64("userID", r.UserID), zap.Int64("lastLeaveFamilyTime", lastLeaveFamilyTime), zap.Error(err))
		return err
	}
	return nil
}

// 获取用户上次退出家族时间
func (r *UserFamilyModel) GetUserLastLeaveFamilyTime(ctx context.Context) (int64, error) {
	logger := fklog.ContextAppLogger(ctx)
	lastLeaveFamilyTime, err := familyredis.GetUserLastLeaveFamilyTime(ctx, r.UserID)
	if err != nil {
		logger.CtxError(ctx, "GetUserLastLeaveFamilyTime GetUserLastLeaveFamilyTime err",
			zap.Uint64("userID", r.UserID), zap.Error(err))
		return 0, err
	}
	return lastLeaveFamilyTime, nil
}

// 设置玩家对应的家族
func (r *UserFamilyModel) SetUserFamily(ctx context.Context, familyID int32) error {
	logger := fklog.ContextAppLogger(ctx)
	err := familyredis.SetUserFamilyID(ctx, r.UserID, familyID)
	if err != nil {
		logger.CtxError(ctx, "SetUserFamilyID SetUserFamilyID err",
			zap.Uint64("userID", r.UserID), zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	return nil
}

// 删除玩家对应的家族
func (r *UserFamilyModel) DelUserFamily(ctx context.Context) error {
	logger := fklog.ContextAppLogger(ctx)
	err := familyredis.DelUserFamilyID(ctx, r.UserID)
	if err != nil {
		logger.CtxError(ctx, "DelUserFamilyID DelUserFamilyID err",
			zap.Uint64("userID", r.UserID), zap.Error(err))
		return err
	}
	return nil
}
