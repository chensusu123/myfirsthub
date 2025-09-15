package userprofileservice

import (
	"context"
	"maze_game_server/model/userprofilemodel"
	"maze_game_server/services/friendservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

const (
	SHOWID_MAN   = 368000001
	SHOWID_WOMAN = 368000002
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

// 用户详情 （当前用户查询其他人的数据）
func (s *service) GetUserDetailInfo(ctx context.Context, userID, peerID uint64) (userDetail *userprofilemodel.UserDetailModel, err error) {
	logger := fklog.ContextAppLogger(ctx)
	profile, err := s.GetUserProfile(ctx, peerID)
	if err != nil {
		return
	}
	//这里是查询用户信息
	userDetail = &userprofilemodel.UserDetailModel{
		UserID:    int64(peerID),
		NickName:  profile.NickName,
		Sex:       profile.Sex,
		IconToken: profile.Avatar,
	}
	if userDetail.Sex == 1 {
		userDetail.ShowID = SHOWID_MAN
	} else {
		userDetail.ShowID = SHOWID_WOMAN
	}

	//判断是否是好友
	isFriend, err := friendservice.GlobalFriendService.CheckFriend(ctx, uint64(userID), uint64(peerID))
	if err != nil {
		logger.CtxError(ctx, "AddBlacklist CheckFriend error", zap.Error(err))
		return
	}
	userDetail.IsFriend = isFriend

	//判断是否是黑名单
	isBlack, err := friendservice.GlobalFriendService.IsBlacklist(ctx, uint64(userID), uint64(peerID))
	if err != nil {
		logger.CtxError(ctx, "AddBlacklist IsBlacklist error", zap.Error(err))
	}
	userDetail.IsBlack = isBlack
	return
}
