package friend

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/Friend"
	"maze_game_server/services/friendservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (f *FriendComponent) OnRemoveBlack_10715_10716(s *session.Session, req *Friend.RemoveBlackRQ) (err error) {
	userId := uint64(s.UID())
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &Friend.RemoveBlackRS{}
	res.Header = req.Header
	res.UserId = req.UserId

	logger.CtxInfo(ctx, "OnRemoveBlack start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnRemoveBlack end", zap.Any("res", res))
	}()

	toID := req.GetUserId()
	if userId == uint64(toID) {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("不能加自己为好友")
		return
	}

	err = friendservice.GlobalFriendService.RemoveBlacklist(ctx, userId, uint64(toID))
	if err != nil {
		logger.CtxError(ctx, "OnRemoveBlack RemoveBlacklist failed", zap.Error(err), zap.Uint64("userId", userId), zap.Int64("toID", toID))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return err
	}

	return nil
}
