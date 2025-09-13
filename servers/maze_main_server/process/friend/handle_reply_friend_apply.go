package friend

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/friendmodel"
	"maze_game_server/pb/common/Friend"
	"maze_game_server/services/friendservice"
	"maze_game_server/services/userprofileservice"
	"maze_game_server/usecase/online"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 处理好友请求 同意/拒绝
func (f *FriendComponent) OnReplyFriendApply_10697_10698(s *session.Session, req *Friend.ReplyFriendApplyRQ) (err error) {
	userId := uint64(s.UID())
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &Friend.ReplyFriendApplyRS{}

	logger.CtxInfo(ctx, "OnReplyFriendApply start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnReplyFriendApply end", zap.Any("res", res))
	}()

	toID := req.GetUserId()

	handlerUserList := make([]*friendmodel.ReceiveFriendRequestInfo, 0)
	switch req.GetReplyResult() {
	case int32(Friend.REPLY_FRIEND_APPLY_RESULT_AGREE):
		handlerUserList, err = friendservice.GlobalFriendService.AgreeFriendApply(ctx, userId, toID)
		if err != nil {
			logger.CtxError(ctx, "OnReplyFriendApply AgreeFriendApply failed", zap.Error(err), zap.Any("toID", toID))
			res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
			return err
		}
	case int32(Friend.REPLY_FRIEND_APPLY_RESULT_REFUSE):
		handlerUserList, err = friendservice.GlobalFriendService.RefuseFriendApply(ctx, userId, toID)
		if err != nil {
			logger.CtxError(ctx, "OnReplyFriendApply AgreeFriendApply failed", zap.Error(err), zap.Any("toID", toID))
			res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
			return err
		}
	}

	nowUserProfile, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnReplyFriendApply GetUserProfile Fail",
			zap.Any("userID", userId),
			zap.Error(err),
		)
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return err
	}

	applyID := &Friend.FriendApplyID{}
	// 推请求处理包
	for _, handlerUser := range handlerUserList {
		userProfile, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, handlerUser.FromUserId)
		if err != nil {
			logger.CtxError(ctx, "OnReplyFriendApply GetUserProfile Fail",
				zap.Any("handlerUser.FromUserId", handlerUser.FromUserId),
				zap.Error(err),
			)
			continue
		}
		pushMsg := &Friend.FriendApplyResultID{
			UserInfo: &Friend.User{
				UserId:     proto.Int64(int64(handlerUser.FromUserId)),
				UserName:   proto.String(userProfile.NickName),
				UserGender: proto.Int32(userProfile.Sex),
				AvaterUrl:  proto.String(userProfile.Avatar),
			},
			CreateTime:  proto.Int64(handlerUser.CreateAt),
			ReplyResult: proto.Int32(req.GetReplyResult()),
			From:        proto.Int32(handlerUser.From),
		}

		if req.GetReplyResult() == int32(Friend.REPLY_FRIEND_APPLY_RESULT_AGREE) {

		} else {

		}

		applyID.DelReceiveInfo = append(applyID.DelReceiveInfo, &Friend.ReceiveInfo{
			UserInfo: &Friend.User{
				UserId:     proto.Int64(int64(handlerUser.FromUserId)),
				UserName:   proto.String(userProfile.NickName),
				UserGender: proto.Int32(userProfile.Sex),
				AvaterUrl:  proto.String(userProfile.Avatar),
			},
			ReceiveTime: proto.Int64(handlerUser.CreateAt),
			ExpireTime:  proto.Int64(handlerUser.CreateAt + int64(friendmodel.ExpireTime)),
		})

		// 推好友请求处理包
		online.ClusterPush(ctx, handlerUser.FromUserId, 10699, pushMsg)

		nowTime := time.Now()
		// 推给对方好友请求变化包
		toApplyID := &Friend.FriendApplyID{
			DelReceiveInfo: []*Friend.ReceiveInfo{
				{
					UserInfo: &Friend.User{
						UserId:     proto.Int64(s.UID()),
						UserName:   proto.String(nowUserProfile.NickName),
						UserGender: proto.Int32(nowUserProfile.Sex),
						AvaterUrl:  proto.String(nowUserProfile.Avatar),
					},
					ReceiveTime: proto.Int64(nowTime.UnixMilli()),
					ExpireTime:  proto.Int64(nowTime.UnixMilli() + int64(friendmodel.ExpireTime)),
				},
			},
		}
		online.ClusterPush(ctx, handlerUser.FromUserId, 10696, toApplyID)
	}

	// 推好友请求变化包
	online.ClusterPush(ctx, userId, 10696, applyID)

	return nil
}
