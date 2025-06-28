package familyservice

import (
	"errors"
	"maze_game_server/model/familymodel"
	"time"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
)

// 检查家族申请是否合理
func (r *service) CheckApplyFamily(logger fklog.FKLogI, familyID int32, applyUser *familymodel.FamilyMember) (bool, error) {
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(logger, familyID)
	if err != nil {
		logger.ErrorWF("CheckApplyFamily LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return false, err
	}
	// 检查家族人数
	if err = familyInfoModel.CheckFamilyMemberCount(logger); err != nil {
		logger.ErrorWF("CheckApplyFamily CheckFamilyMemberCount err",
			zap.Int32("familyID", familyID), zap.Uint64("userID", applyUser.UserID))
		return false, errors.New("家族人数已满")
	}
	// 检查用户是否在家族
	if familyInfoModel.CheckUserInFamily(logger, applyUser.UserID) {
		logger.ErrorWF("CheckApplyFamily CheckUserInFamily err",
			zap.Int32("familyID", familyID), zap.Uint64("userID", applyUser.UserID))
		return false, errors.New("用户已在家族中")
	}
	return true, nil
}

// 设置玩家所在家族
func (r *service) SetUserFamily(logger fklog.FKLogI, userID uint64, familyID int32) error {
	userFamilyModel := familymodel.NewUserFamilyModel(logger, userID)
	return userFamilyModel.SetUserFamily(logger, familyID)
}

// 获取玩家所在家族
func (r *service) GetUserFamily(logger fklog.FKLogI, userID uint64) (int32, error) {
	userFamilyModel := familymodel.NewUserFamilyModel(logger, userID)
	return userFamilyModel.GetUserFamily(logger)
}

func (r *service) DeleteUserFamily(logger fklog.FKLogI, userID uint64) error {
	userFamilyModel := familymodel.NewUserFamilyModel(logger, userID)
	return userFamilyModel.DelUserFamily(logger)
}

// 设置玩家上次退出家族时间
func (r *service) SetUserLastLeaveFamilyTime(logger fklog.FKLogI, userID uint64) error {
	userFamilyModel := familymodel.NewUserFamilyModel(logger, userID)
	return userFamilyModel.SetUserLastLeaveFamilyTime(logger)
}

// 获取玩家上次退出家族时间
func (r *service) GetUserLastLeaveFamilyTime(logger fklog.FKLogI, userID uint64) (int64, error) {
	userFamilyModel := familymodel.NewUserFamilyModel(logger, userID)
	return userFamilyModel.GetUserLastLeaveFamilyTime(logger)
}

// 检测玩家距离上次退出家族时间是否超过一天
func (r *service) CheckUserLastLeaveFamilyTime(logger fklog.FKLogI, userID uint64) (bool, error) {
	userFamilyModel := familymodel.NewUserFamilyModel(logger, userID)
	lastLeaveFamilyTime, err := userFamilyModel.GetUserLastLeaveFamilyTime(logger)
	if err != nil {
		logger.ErrorWF("CheckUserLastLeaveFamilyTime GetUserLastLeaveFamilyTime err",
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
func (r *service) CheckUserHaveFamily(logger fklog.FKLogI, userID uint64) (bool, error) {
	userFamilyModel := familymodel.NewUserFamilyModel(logger, userID)
	familyID, err := userFamilyModel.GetUserFamily(logger)
	if err != nil {
		logger.ErrorWF("CheckUserHaveFamily GetUserFamily err",
			zap.Uint64("userID", userID), zap.Error(err))
		return false, err
	}
	return familyID > 0, nil
}

// 检查用户是否在家族的请求列表中
func (r *service) CheckUserInApplyList(logger fklog.FKLogI, familyID int32, userID uint64) (bool, error) {
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(logger, familyID)
	if err != nil {
		logger.ErrorWF("CheckUserInApplyList LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Uint64("userID", userID), zap.Error(err))
		return false, err
	}
	return familyInfoModel.CheckUserInApplyList(logger, userID), nil
}

// 检查用户是否在家族中
func (r *service) CheckUserInFamily(logger fklog.FKLogI, familyID int32, userID uint64) (bool, error) {
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(logger, familyID)
	if err != nil {
		logger.ErrorWF("CheckUserInFamily LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Uint64("userID", userID), zap.Error(err))
		return false, err
	}
	return familyInfoModel.CheckUserInFamily(logger, userID), nil
}
