package im

import (
	"maze_game_server/app"
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeIM"
	"maze_game_server/services/sessionservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

func (im *IM) OnQueryRecentSessions_10652_10653(s *session.Session, req *MazeIM.QueryRecentSessionsRQ) (err error) {
	defer fkprometheus.InfoPMT("OnQueryMessages")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeIM.QueryRecentSessionsRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	logger.InfoWF("OnQueryMessages start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnQueryMessages end", zap.Any("res", res))
	}()

	var (
		userId = uint64(s.UID())
	)

	user, err := app.WrapUser(userId, "")
	if err != nil {
		return err
	}

	messages, err := sessionservice.Default.QueryRecentSessions(s.Context(), logger, app.Maze, user)
	_ = messages
	return
}

func (im *IM) OnRemoveSession_10654_10655(s *session.Session, req *MazeIM.RemoveSessionRQ) (err error) {
	defer fkprometheus.InfoPMT("OnRemoveSession")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeIM.RemoveSessionRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	logger.InfoWF("OnRemoveSession start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnRemoveSession end", zap.Any("res", res))
	}()

	var (
		userId    = uint64(s.UID())
		sessionID = req.GetSessionId()
	)

	user, err := app.WrapUser(userId, "")
	if err != nil {
		return err
	}

	err = sessionservice.Default.RemoveSession(s.Context(), logger, app.Maze, user, sessionID)
	_ = err
	return
}
