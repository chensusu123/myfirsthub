package friend

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/Friend"
	"maze_game_server/services/friendservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (f *FriendComponent) OnAddBlack_10713_10714(s *session.Session, req *Friend.AddBlackRQ) (err error) {
	userId := s.UID()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &Friend.AddBlackRS{}
	res.Header = req.Header
	res.UserId = req.UserId

	logger.CtxInfo(ctx, "OnAddBlack start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnAddBlack end", zap.Any("res", res))
	}()

	toID := req.GetUserId()
	if userId == toID {
		logger.CtxWarn(ctx, "OnAddBlack userId args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("不能加自己为黑名单")
		return err
	}

	codeErr := friendservice.GlobalFriendService.AddBlacklist(ctx, uint64(userId), uint64(toID))
	if codeErr != nil {
		logger.CtxError(ctx, "OnAddBlack AddBlacklist failed", zap.Error(err), zap.Int64("toID", toID))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return err
	}

	return nil
}
