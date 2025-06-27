package userprofileservice

import (
	"fmt"
	"maze_game_server/model/userprofilemodel"
	"maze_game_server/pb/common/UserProfile"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
)

var GlobalUserProfileService = &UserProfileService{}

type UserProfileService struct {
}

// GetUserProfile 查询多个用户信息
func (u *UserProfileService) GetBatchUserProfile(logger fklog.FKLogI, users []uint64) (profiles []*UserProfile.UserProfile, err error) {
	if len(users) <= 0 {
		return
	}
	models, err := userprofilemodel.LoadBatchUserProfileModel(logger, users)
	if err != nil {
		logger.ErrorWF("GetBatchUserProfile LoadBatchUserProfileModel error", zap.Any("users", users),
			zap.Error(err))
		return
	}
	for _, v := range models {
		profiles = append(profiles, v.ModelDataToPb())
	}
	return
}

// GetUserProfile 查询用户信息
func (u *UserProfileService) GetUserProfile(logger fklog.FKLogI, userID uint64) (profile *UserProfile.UserProfile, err error) {
	model, err := userprofilemodel.LoadUserProfileModel(logger, userID)
	if err != nil {
		logger.ErrorWF("GetUserProfile LoadUserProfileModel error", zap.Any("userID", userID),
			zap.Error(err))
		return
	}
	return model.ModelDataToPb(), nil
}

// AlterUserProfile
func (u *UserProfileService) AlterUserProfile(logger fklog.FKLogI, userId uint64, alterProfile *UserProfile.UserProfile) (err error) {
	if alterProfile == nil {
		return fmt.Errorf("AlterUserProfile incorrent, userID:%v, alterProfile:%v", userId, alterProfile)
	}

	model, err := userprofilemodel.LoadUserProfileModel(logger, userId)
	if err != nil {
		logger.ErrorWF("AlterUserProfile LoadUserProfileModel error", zap.Any("userId", userId),
			zap.Error(err))
		return
	}

	// check
	if model != nil && model.UserID != userId {
		logger.ErrorWF("AlterUserProfile check userId error", zap.Any("userId", userId),
			zap.Error(err))
		return
	}

	// alter
	model.FillModelData(alterProfile)
	model.Save(logger, userId)

	return
}
