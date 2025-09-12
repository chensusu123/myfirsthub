package familyservice

import (
	"context"
	"errors"
	"maze_game_server/model/familymodel"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 检查家族申请是否合理
func (r *service) CheckApplyFamily(ctx context.Context, familyID int32, applyUser *familymodel.FamilyMember) (bool, error) {
	logger := fklog.ContextAppLogger(ctx)
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(ctx, familyID)
	if err != nil {
		logger.CtxError(ctx, "CheckApplyFamily LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return false, err
	}
	// 检查家族人数
	if err = familyInfoModel.CheckFamilyMemberCount(ctx); err != nil {
		logger.CtxError(ctx, "CheckApplyFamily CheckFamilyMemberCount err",
			zap.Int32("familyID", familyID), zap.Uint64("userID", applyUser.UserID))
		return false, errors.New("家族人数已满")
	}
	// 检查用户是否在家族
	if familyInfoModel.CheckUserInFamily(ctx, applyUser.UserID) {
		logger.CtxError(ctx, "CheckApplyFamily CheckUserInFamily err",
			zap.Int32("familyID", familyID), zap.Uint64("userID", applyUser.UserID))
		return false, errors.New("用户已在家族中")
	}
	return true, nil
}

// 设置玩家所在家族
func (r *service) SetUserFamily(ctx context.Context, userID uint64, familyID int32) error {
	userFamilyModel := familymodel.NewUserFamilyModel(ctx, userID)
	return userFamilyModel.SetUserFamily(ctx, familyID)
}

// 获取玩家所在家族
func (r *service) GetUserFamily(ctx context.Context, userID uint64) (int32, error) {
	userFamilyModel := familymodel.NewUserFamilyModel(ctx, userID)
	return userFamilyModel.GetUserFamily(ctx)
}

func (r *service) DeleteUserFamily(ctx context.Context, userID uint64) error {
	userFamilyModel := familymodel.NewUserFamilyModel(ctx, userID)
	return userFamilyModel.DelUserFamily(ctx)
}

// 设置玩家上次退出家族时间
func (r *service) SetUserLastLeaveFamilyTime(ctx context.Context, userID uint64) error {
	userFamilyModel := familymodel.NewUserFamilyModel(ctx, userID)
	return userFamilyModel.SetUserLastLeaveFamilyTime(ctx)
}

// 获取玩家上次退出家族时间
func (r *service) GetUserLastLeaveFamilyTime(ctx context.Context, userID uint64) (int64, error) {
	userFamilyModel := familymodel.NewUserFamilyModel(ctx, userID)
	return userFamilyModel.GetUserLastLeaveFamilyTime(ctx)
}

// 检测玩家距离上次退出家族时间是否超过一天
func (r *service) CheckUserLastLeaveFamilyTime(ctx context.Context, userID uint64) (bool, error) {
	logger := fklog.ContextAppLogger(ctx)
	userFamilyModel := familymodel.NewUserFamilyModel(ctx, userID)
	lastLeaveFamilyTime, err := userFamilyModel.GetUserLastLeaveFamilyTime(ctx)
	if err != nil {
		logger.CtxError(ctx, "CheckUserLastLeaveFamilyTime GetUserLastLeaveFamilyTime err",
			zap.Uint64("userID", userID), zap.Error(err))
		return false, err
	}

	if lastLeaveFamilyTime == 0 {
		return true, nil
	}

	if time.Now().Unix()-lastLeaveFamilyTime > 24*60*60 {
		return true, nil
	}
	return false, nil
}

// 检查玩家是否有家族
func (r *service) CheckUserHaveFamily(ctx context.Context, userID uint64) (bool, error) {
	logger := fklog.ContextAppLogger(ctx)
	userFamilyModel := familymodel.NewUserFamilyModel(ctx, userID)
	familyID, err := userFamilyModel.GetUserFamily(ctx)
	if err != nil {
		logger.CtxError(ctx, "CheckUserHaveFamily GetUserFamily err",
			zap.Uint64("userID", userID), zap.Error(err))
		return false, err
	}
	return familyID > 0, nil
}

// 检查用户是否在家族的请求列表中
func (r *service) CheckUserInApplyList(ctx context.Context, familyID int32, userID uint64) (bool, error) {
	logger := fklog.ContextAppLogger(ctx)
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(ctx, familyID)
	if err != nil {
		logger.CtxError(ctx, "CheckUserInApplyList LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Uint64("userID", userID), zap.Error(err))
		return false, err
	}
	return familyInfoModel.CheckUserInApplyList(ctx, userID), nil
}

// 检查用户是否在家族中
func (r *service) CheckUserInFamily(ctx context.Context, familyID int32, userID uint64) (bool, error) {
	logger := fklog.ContextAppLogger(ctx)
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(ctx, familyID)
	if err != nil {
		logger.CtxError(ctx, "CheckUserInFamily LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Uint64("userID", userID), zap.Error(err))
		return false, err
	}
	return familyInfoModel.CheckUserInFamily(ctx, userID), nil
}
