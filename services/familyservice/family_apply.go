package familyservice

import (
	"context"
	"maze_game_server/model/familymodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 申请加入家族
func (r *service) ApplyFamily(ctx context.Context, familyID int32, applyUser familymodel.FamilyMember) (*familymodel.FamilyInfoModel, error) {
	logger := fklog.ContextAppLogger(ctx)
	// 加入申请列表
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(ctx, familyID)
	if err != nil {
		logger.CtxError(ctx, "ApplyFamily LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	familyInfoModel.AddApplyUser(ctx, applyUser)
	// 保存
	err = familyInfoModel.Save(ctx, familyID)
	if err != nil {
		logger.CtxError(ctx, "ApplyFamily familyInfoModel.Save err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	return familyInfoModel, nil
}

// AgreeApplyFamily 同意加入家族
func (r *service) AgreeApplyFamily(ctx context.Context, familyID int32, userID uint64,
	applyUser familymodel.FamilyMember) (*familymodel.FamilyInfoModel, error) {
	logger := fklog.ContextAppLogger(ctx)
	// 加入申请列表
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(ctx, familyID)
	if err != nil {
		logger.CtxError(ctx, "AgreeApplyFamily LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	// 添加用户
	familyInfoModel.AddMember(ctx, applyUser)
	// 移除申请用户
	familyInfoModel.RemApplyUser(ctx, applyUser)
	// 保存
	err = familyInfoModel.Save(ctx, familyID)
	if err != nil {
		logger.CtxError(ctx, "AgreeApplyFamily familyInfoModel.Save err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	return familyInfoModel, nil
}

// RefuseApplyFamily 拒绝加入家族
func (r *service) RefuseApplyFamily(ctx context.Context, familyID int32, userID uint64,
	applyUser familymodel.FamilyMember) (*familymodel.FamilyInfoModel, error) {
	logger := fklog.ContextAppLogger(ctx)
	// 加入申请列表
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(ctx, familyID)
	if err != nil {
		logger.CtxError(ctx, "AgreeApplyFamily LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	// 移除申请用户
	familyInfoModel.RemApplyUser(ctx, applyUser)
	// 保存
	err = familyInfoModel.Save(ctx, familyID)
	if err != nil {
		logger.ErrorWF("AgreeApplyFamily familyInfoModel.Save err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	return familyInfoModel, nil
}
