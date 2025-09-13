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

// 添加好友请求包
func (f *FriendComponent) OnFriendApply_10694_10695(s *session.Session, req *Friend.FriendApplyRQ) (err error) {
	userId := uint64(s.UID())
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &Friend.FriendApplyRS{}

	res.UserId = req.UserId
	res.From = req.From

	logger.CtxInfo(ctx, "OnFriendApply start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnFriendApply end", zap.Any("res", res))
	}()

	toID := uint64(req.GetUserId())
	if userId == toID {
		logger.CtxError(ctx, "OnFriendApply userId args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("不能加自己为好友")
		return err
	}

	applyTime, err := friendservice.GlobalFriendService.AddFriendRequest(ctx, userId, toID, req.GetFrom())
	if err != nil {
		logger.ErrorWF("OnFriendApply FriendRequest failed ", zap.Error(err), zap.Uint64("userId", userId), zap.Uint64("toID", toID))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return
	}

	userProfiel, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, toID)
	if err != nil {
		logger.CtxError(ctx, "OnFriendApply GetUserProfile Fail",
			zap.Uint64("userID", userId),
			zap.Uint64("toID", toID),
			zap.Error(err),
		)
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return
	}

	pushMsg := &Friend.FriendApplyID{
		AddReceiveInfo: []*Friend.ReceiveInfo{
			{
				UserInfo: &Friend.User{
					UserId:     proto.Int64(s.UID()),
					UserName:   proto.String(userProfiel.NickName),
					UserGender: proto.Int32(userProfiel.Sex),
					AvaterUrl:  proto.String(userProfiel.Avatar),
				},
				ReceiveTime: proto.Int64(applyTime),
				ExpireTime:  proto.Int64(applyTime + int64(friendmodel.ExpireTime)),
			},
		},
	}
	// 通知对方请求加好友
	online.ClusterPush(ctx, toID, 10696, pushMsg)

	return nil
}
