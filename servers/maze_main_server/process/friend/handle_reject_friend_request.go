package friend

import (
	"go.uber.org/zap"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/Friend"
	"maze_game_server/services/friendservice"
)

func (f *FriendComponent) OnRejectFriendRequest_10558_10559(s *session.Session, req *Friend.RejectFriendRequestRQ) (err error) {
	userId := uint64(s.UID())
	logger := log.Clone("OnRejectFriendRequest", userId, 0)
	res := &Friend.RejectFriendRequestRS{}

	logger.InfoWF("OnRejectFriendRequest start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnRejectFriendRequest end", zap.Any("res", res))
	}()

	toID := req.GetUserId()

	codeErr := friendservice.GlobalFriendService.RejectFriendRequest(logger, userId, toID)
	if codeErr != nil {
		logger.ErrorWF("OnRejectFriendRequest RejectFriendRequest failed", zap.Error(err), zap.Uint64("userId", userId), zap.Uint64("toID", toID))
		res.ErrInfo = codeErr.ToInfo()
		return codeErr
	}

	return nil
}
