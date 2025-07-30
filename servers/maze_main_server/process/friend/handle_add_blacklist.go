package friend

import (
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/Friend"
	"maze_game_server/services/friendservice"
)

func (f *FriendComponent) OnAddBlacklist_10564_10565(s *session.Session, req *Friend.AddBlacklistRQ) (err error) {
	userId := uint64(s.UID())
	logger := log.Clone("OnAddBlacklist", userId, 0)
	res := &Friend.AddBlacklistRS{}

	logger.InfoWF("OnAddBlacklist start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnAddBlacklist end", zap.Any("res", res))
	}()

	toID := req.GetUserId()
	if userId == toID {
		logger.ErrorWF("OnAddBlacklist userId args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("不能加自己为黑名单")
		return err
	}

	codeErr := friendservice.GlobalFriendService.AddBlacklist(logger, userId, toID)
	if codeErr != nil {
		logger.ErrorWF("OnAddBlacklist AddBlacklist failed", zap.Error(err), zap.Uint64("toID", toID))
		res.ErrInfo = codeErr.ToInfo()
		return codeErr
	}

	return nil
}
