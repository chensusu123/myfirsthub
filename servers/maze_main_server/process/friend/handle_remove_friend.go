package friend

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/Friend"
	"maze_game_server/services/friendservice"
	"maze_game_server/services/userprofileservice"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (f *FriendComponent) OnDeleteFriend_10706_10707(s *session.Session, req *Friend.DeleteFriendRQ) (err error) {
	userId := s.UID()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &Friend.DeleteFriendRS{}

	logger.CtxInfo(ctx, "OnDeleteFriend start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnDeleteFriend end", zap.Any("res", res))
	}()

	toID := req.GetUserId()

	toInfo, inInfo, err := friendservice.GlobalFriendService.RemoveFriend(ctx, uint64(userId), uint64(toID))
	if err != nil {
		logger.CtxError(ctx, "OnDeleteFriend RemoveFriend failed", zap.Error(err), zap.Int64("userId", userId), zap.Int64("toID", toID))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return
	}

	delUserInfo, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, uint64(toID))
	if err != nil {
		logger.CtxError(ctx, "OnDeleteFriend GetUserProfile Fail",
			zap.Any("toID", toID),
			zap.Error(err),
		)
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return
	}

	nowUserInfo, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, uint64(userId))
	if err != nil {
		logger.CtxError(ctx, "OnDeleteFriend GetUserProfile Fail",
			zap.Any("userID", userId),
			zap.Error(err),
		)
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return
	}

	// 推给自己好友列表变化包
	tofriendListChange := &Friend.FriendListChangeID{
		DelFriendList: []*Friend.FriendInfo{
			{
				UserInfo: &Friend.User{
					UserId:     proto.Int64(toID),
					UserName:   proto.String(delUserInfo.NickName),
					UserGender: proto.Int32(delUserInfo.Sex),
					AvaterUrl:  proto.String(delUserInfo.Avatar),
				},
				FriendType: proto.Int32(toInfo.FriendType),
				AddTime:    proto.Int64(toInfo.CreateAt),
			},
		},
	}
	online.ClusterPush(ctx, uint64(userId), 10708, tofriendListChange)

	// 推给对方好友列表变化包
	infriendListChange := &Friend.FriendListChangeID{
		DelFriendList: []*Friend.FriendInfo{
			{
				UserInfo: &Friend.User{
					UserId:     proto.Int64(userId),
					UserName:   proto.String(nowUserInfo.NickName),
					UserGender: proto.Int32(nowUserInfo.Sex),
					AvaterUrl:  proto.String(nowUserInfo.Avatar),
				},
				FriendType: proto.Int32(inInfo.FriendType),
				AddTime:    proto.Int64(inInfo.CreateAt),
			},
		},
	}
	online.ClusterPush(ctx, uint64(toID), 10708, infriendListChange)

	return nil
}
