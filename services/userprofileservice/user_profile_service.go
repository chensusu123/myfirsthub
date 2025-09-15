package userprofileservice

import (
	"context"
	"maze_game_server/model/userprofilemodel"
)

var GlobalUserProfileService UserProfileService

type UserProfileService interface {
	// 批量获取用户资料
	GetBatchUserProfile(ctx context.Context, users []uint64) (profiles []*userprofilemodel.UserProfileModel, err error)
	// 获取用户资料
	GetUserProfile(ctx context.Context, userID uint64) (profile *userprofilemodel.UserProfileModel, err error)
	// 修改用户资料
	AlterUserProfile(ctx context.Context, userId uint64, alterProfile *userprofilemodel.UserProfileModel) (err error)
	//用于当前用户查询其他人的详细信息
	GetUserDetailInfo(ctx context.Context, userID uint64, peerID uint64) (detail *userprofilemodel.UserDetailModel, err error)
}

type service struct {
}

func newUserProfileService() UserProfileService {
	return &service{}
}

func init() {
	GlobalUserProfileService = newUserProfileService()
}
