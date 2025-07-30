package friend

import (
	"go.uber.org/zap"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/Friend"
	"maze_game_server/services/friendservice"
)

func (f *FriendComponent) OnRemoveFriend_10562_10563(s *session.Session, req *Friend.RemoveFriendRQ) (err error) {
	userId := uint64(s.UID())
	logger := log.Clone("OnRemoveFriend", userId, 0)
	res := &Friend.RemoveFriendRS{}

	logger.InfoWF("OnRemoveFriend start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnRemoveFriend end", zap.Any("res", res))
	}()

	toID := req.GetUserId()

	codeErr := friendservice.GlobalFriendService.RemoveFriend(logger, userId, toID)
	if codeErr != nil {
		logger.ErrorWF("OnRemoveFriend RemoveFriend failed", zap.Error(err), zap.Uint64("userId", userId), zap.Uint64("toID", toID))
		res.ErrInfo = codeErr.ToInfo()
		return codeErr
	}

	return nil
}
