package im

import (
	"maze_game_server/app"
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeIM"
	"maze_game_server/services/sessionservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

func (im *IM) OnQueryRecentSessions_10652_10653(s *session.Session, req *MazeIM.QueryRecentSessionsRQ) (err error) {
	defer fkprometheus.InfoPMT("OnQueryMessages")()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)

	res := &MazeIM.QueryRecentSessionsRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	logger.CtxInfo(ctx, "OnQueryMessages start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnQueryMessages end", zap.Any("res", res))
	}()

	var (
		userId = uint64(s.UID())
	)

	user, err := app.WrapUser(userId, "")
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("获取用户信息失败")
		logger.CtxError(ctx, "OnQueryMessages WrapUser error", zap.Error(err))
		return err
	}

	messages, err := sessionservice.Default.QueryRecentSessions(ctx, app.Maze, user)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("获取最近会话失败")
		logger.CtxError(ctx, "OnQueryMessages QueryRecentSessions error", zap.Error(err))
		return err
	}

	messagesList, err := sessionservice.Default.GetMessageInfo(ctx, app.Maze, user, messages)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("获取消息信息失败")
		logger.CtxError(ctx, "OnQueryMessages GetMessageInfo error", zap.Error(err))
		return err
	}
	res.SessionList = messagesList
	return
}

func (im *IM) OnRemoveSession_10654_10655(s *session.Session, req *MazeIM.RemoveSessionRQ) (err error) {
	defer fkprometheus.InfoPMT("OnRemoveSession")()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)

	res := &MazeIM.RemoveSessionRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	logger.CtxInfo(ctx, "OnRemoveSession start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnRemoveSession end", zap.Any("res", res))
	}()

	var (
		userId    = uint64(s.UID())
		sessionID = req.GetSessionId()
	)

	user, err := app.WrapUser(userId, "")
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("获取用户信息失败")
		logger.CtxError(ctx, "OnRemoveSession WrapUser error", zap.Error(err))
		return err
	}

	err = sessionservice.Default.RemoveSession(ctx, app.Maze, user, sessionID)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("删除会话失败")
		logger.CtxError(ctx, "OnRemoveSession RemoveSession error", zap.Error(err))
		return err
	}

	return
}
