package friend

import (
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/Friend"
	"maze_game_server/services/friendservice"
	"maze_game_server/usecase/online"
)

func (f *FriendComponent) OnAcceptFriendRequest_10555_10556(s *session.Session, req *Friend.AcceptFriendRequestRQ) (err error) {
	userId := uint64(s.UID())
	logger := log.Clone("OnAcceptFriendRequest", userId, 0)
	res := &Friend.AcceptFriendRequestRS{}

	logger.InfoWF("OnAcceptFriendRequest start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnAcceptFriendRequest end", zap.Any("res", res))
	}()

	toID := req.GetUserId()

	codeErr := friendservice.GlobalFriendService.AcceptFriendRequest(logger, userId, toID)
	if codeErr != nil {
		logger.ErrorWF("OnAcceptFriendRequest AcceptFriendRequest failed", zap.Error(err), zap.Uint64("toID", toID))
		res.ErrInfo = codeErr.ToInfo()
		return codeErr
	}

	// 通知同意加好友
	pushMsg := &Friend.AcceptFriendRequestID{
		UserId: proto.Uint64(userId),
	}
	// 通知对方同意加好友
	online.Push(logger, toID, 10557, pushMsg)

	return nil
}
