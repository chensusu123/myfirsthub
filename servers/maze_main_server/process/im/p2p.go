package im

import (
	"maze_game_server/app"
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeIM"
	"maze_game_server/services/p2pservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

func (im *IM) OnQueryMessages_10643_10644(s *session.Session, req *MazeIM.QueryMessagesRQ) (err error) {
	defer fkprometheus.InfoPMT("OnQueryMessages")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeIM.QueryMessagesRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	logger.InfoWF("OnQueryMessages start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnQueryMessages end", zap.Any("res", res))
	}()

	var (
		userId    = uint64(s.UID())
		peerId    = req.GetPeerId()
		lastMsgID = req.GetLastMsgId()
	)

	user, err := app.WrapUser(userId, "")
	if err != nil {
		return err
	}

	messages, err := p2pservice.Default.QueryMessages(s.Context(), logger, app.Maze, user, peerId, lastMsgID, 20)
	_ = messages
	return
}

func (im *IM) OnSendMessage_10645_10646(s *session.Session, req *MazeIM.SendMessageRQ) (err error) {
	defer fkprometheus.InfoPMT("OnSendMessage")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeIM.SendMessageRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	logger.InfoWF("OnSendMessage start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnSendMessage end", zap.Any("res", res))
	}()

	var (
		userId  = uint64(s.UID())
		peerId  = req.GetPeerId()
		_type   = req.GetType()
		content = req.GetContent()
	)

	user, err := app.WrapUser(userId, "")
	if err != nil {
		return err
	}

	messageID, err := p2pservice.Default.SendMessage(s.Context(), logger, app.Maze, user, peerId, _type, content)
	_ = messageID
	return
}
