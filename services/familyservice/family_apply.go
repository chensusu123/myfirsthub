package familyservice

import (
	"maze_game_server/model/familymodel"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
)

// 申请加入家族
func (r *service) ApplyFamily(logger fklog.FKLogI, familyID int32, applyUser familymodel.FamilyMember) (*familymodel.FamilyInfoModel, error) {
	// 加入申请列表
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(logger, familyID)
	if err != nil {
		logger.ErrorWF("ApplyFamily LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	familyInfoModel.AddApplyUser(logger, applyUser)
	// 保存
	err = familyInfoModel.Save(logger, familyID)
	if err != nil {
		logger.ErrorWF("ApplyFamily familyInfoModel.Save err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	return familyInfoModel, nil
}

// AgreeApplyFamily 同意加入家族
func (r *service) AgreeApplyFamily(logger fklog.FKLogI, familyID int32, userID uint64,
	applyUser familymodel.FamilyMember) (*familymodel.FamilyInfoModel, error) {
	// 加入申请列表
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(logger, familyID)
	if err != nil {
		logger.ErrorWF("AgreeApplyFamily LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	// 添加用户
	familyInfoModel.AddMember(logger, applyUser)
	// 移除申请用户
	familyInfoModel.RemApplyUser(logger, applyUser)
	// 保存
	err = familyInfoModel.Save(logger, familyID)
	if err != nil {
		logger.ErrorWF("AgreeApplyFamily familyInfoModel.Save err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	return familyInfoModel, nil
}

// RefuseApplyFamily 拒绝加入家族
func (r *service) RefuseApplyFamily(logger fklog.FKLogI, familyID int32, userID uint64,
	applyUser familymodel.FamilyMember) (*familymodel.FamilyInfoModel, error) {
	// 加入申请列表
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(logger, familyID)
	if err != nil {
		logger.ErrorWF("AgreeApplyFamily LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	// 移除申请用户
	familyInfoModel.RemApplyUser(logger, applyUser)
	// 保存
	err = familyInfoModel.Save(logger, familyID)
	if err != nil {
		logger.ErrorWF("AgreeApplyFamily familyInfoModel.Save err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	return familyInfoModel, nil
}
