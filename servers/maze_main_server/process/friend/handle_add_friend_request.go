package friend

import (
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/Friend"
	"maze_game_server/services/friendservice"
	"maze_game_server/usecase/online"
)

func (f *FriendComponent) OnAddFriendRequest_10552_10553(s *session.Session, req *Friend.AddFriendRequestRQ) (err error) {
	userId := uint64(s.UID())
	logger := log.Clone("OnAddFriendRequest", userId, 0)
	res := &Friend.AddFriendRequestRS{}

	logger.InfoWF("OnAddFriendRequest start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnAddFriendRequest end", zap.Any("res", res))
	}()

	toID := req.GetUserId()
	if userId == toID {
		logger.ErrorWF("OnAddFriendRequest userId args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("不能加自己为好友")
		return err
	}

	codeErr := friendservice.GlobalFriendService.AddFriendRequest(logger, userId, toID)
	if codeErr != nil {
		logger.ErrorWF("OnAddFriendRequest FriendRequest failed ", zap.Error(err), zap.Uint64("userId", userId), zap.Uint64("toID", toID))
		res.ErrInfo = codeErr.ToInfo()
		return codeErr
	}

	pushMsg := &Friend.AddFriendRequestID{
		FromUserId: proto.Uint64(userId),
	}
	// 通知对方请求加好友
	online.Push(logger, toID, 10554, pushMsg)

	return nil
}
