package friend

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/friendmodel"
	"maze_game_server/pb/common/Friend"
	"maze_game_server/services/friendservice"
	"maze_game_server/services/userprofileservice"
	"maze_game_server/usecase/online"

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
	res.Header = req.Header
	res.ReplyResult = req.ReplyResult

	logger.CtxInfo(ctx, "OnReplyFriendApply start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnReplyFriendApply end", zap.Any("res", res))
	}()

	toID := req.GetUserId()

	nowUserProfile, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnReplyFriendApply GetUserProfile fail",
			zap.Any("userID", userId),
			zap.Error(err),
		)
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return err
	}

	var skip bool
	handlerUserList := make([]*friendmodel.ReceiveFriendRequestInfo, 0)
	changeUserList := make([]*friendmodel.ReceiveFriendRequestInfo, 0)
	switch req.GetReplyResult() {
	case int32(Friend.REPLY_FRIEND_APPLY_RESULT_AGREE):
		handlerUserList, changeUserList, skip, err = friendservice.GlobalFriendService.AgreeFriendApply(ctx, userId, toID)
		if err != nil {
			if !skip {
				logger.CtxError(ctx, "OnReplyFriendApply AgreeFriendApply failed", zap.Error(err), zap.Any("toID", toID))
				res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
			} else {
				res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
			}
		}
	case int32(Friend.REPLY_FRIEND_APPLY_RESULT_REFUSE):
		handlerUserList, changeUserList, skip, err = friendservice.GlobalFriendService.RefuseFriendApply(ctx, userId, toID)
		if err != nil {
			if !skip {
				logger.CtxError(ctx, "OnReplyFriendApply AgreeFriendApply failed", zap.Error(err), zap.Any("toID", toID))
				res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
			} else {
				res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
			}
		}
	}

	logger.CtxInfo(ctx, "OnReplyFriendApply handler", zap.Any("handlerUserList", handlerUserList), zap.Any("changeUserList", changeUserList))
	// 推收到请求列表变化包
	toapplyID := &Friend.FriendApplyID{}

	for _, handlerUser := range handlerUserList {
		touserProfile, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, handlerUser.FromUserId)
		if err != nil {
			logger.CtxError(ctx, "OnReplyFriendApply GetUserProfile Fail",
				zap.Any("handlerUser.FromUserId", handlerUser.FromUserId),
				zap.Error(err),
			)
			continue
		}

		toapplyID.DelReceiveInfo = append(toapplyID.DelReceiveInfo, &Friend.ReceiveInfo{
			UserInfo: &Friend.User{
				UserId:     proto.Int64(int64(handlerUser.FromUserId)),
				UserName:   proto.String(touserProfile.NickName),
				UserGender: proto.Int32(touserProfile.Sex),
				AvaterUrl:  proto.String(touserProfile.Avatar),
			},
			ReceiveTime: proto.Int64(handlerUser.CreateAt),
			ExpireTime:  proto.Int64(handlerUser.CreateAt + int64(friendmodel.ExpireTime)),
		})

		res.UserId = append(res.UserId, int64(handlerUser.FromUserId))

		// 推发送请求变化
		sendRqMsg := &Friend.FriendApplyResultID{
			UserInfo: &Friend.User{
				UserId:     proto.Int64(int64(userId)),
				UserName:   proto.String(nowUserProfile.NickName),
				UserGender: proto.Int32(nowUserProfile.Sex),
				AvaterUrl:  proto.String(nowUserProfile.Avatar),
			},
			CreateTime:  proto.Int64(handlerUser.CreateAt),
			ReplyResult: proto.Int32(req.GetReplyResult()),
			From:        proto.Int32(handlerUser.From),
		}

		online.ClusterPush(ctx, handlerUser.FromUserId, 10699, sendRqMsg)
	}

	if len(handlerUserList) != 0 {
		// 推好友请求变化包 删除我的收到好友请求列表
		online.ClusterPush(ctx, userId, 10696, toapplyID)
	}

	for _, changeUser := range changeUserList {
		touserProfile, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, changeUser.FromUserId)
		if err != nil {
			logger.CtxError(ctx, "OnReplyFriendApply GetUserProfile Fail",
				zap.Any("changeUser.FromUserId", changeUser.FromUserId),
				zap.Error(err),
			)
			continue
		}
		if req.GetReplyResult() == int32(Friend.REPLY_FRIEND_APPLY_RESULT_AGREE) {
			//推好友列表变化包
			pushMsg := &Friend.FriendListChangeID{
				AddFriendList: []*Friend.FriendInfo{
					{
						UserInfo: &Friend.User{
							UserId:     proto.Int64(int64(userId)),
							UserName:   proto.String(nowUserProfile.NickName),
							UserGender: proto.Int32(nowUserProfile.Sex),
							AvaterUrl:  proto.String(nowUserProfile.Avatar),
						},
						FriendType: proto.Int32(changeUser.From),
						AddTime:    proto.Int64(changeUser.CreateAt),
					},
				},
			}
			// 告诉对方好友列表变化
			online.ClusterPush(ctx, changeUser.FromUserId, 10708, pushMsg)

			// 告诉自己好友列表变化
			handlerUserMsg := &Friend.FriendListChangeID{
				AddFriendList: []*Friend.FriendInfo{
					{
						UserInfo: &Friend.User{
							UserId:     proto.Int64(int64(changeUser.FromUserId)),
							UserName:   proto.String(touserProfile.NickName),
							UserGender: proto.Int32(touserProfile.Sex),
							AvaterUrl:  proto.String(touserProfile.Avatar),
						},
						FriendType: proto.Int32(changeUser.From),
						AddTime:    proto.Int64(changeUser.CreateAt),
					},
				},
			}
			online.ClusterPush(ctx, userId, 10708, handlerUserMsg)
		}
	}

	return nil
}
