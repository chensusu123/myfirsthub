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
	"google.golang.org/protobuf/proto"
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
		logger.ErrorWF("OnQueryMessages WrapUser error", zap.Error(err), zap.Any("req", req))
		return err
	}

	messages, err := p2pservice.Default.QueryMessages(s.Context(), logger, app.Maze, user, peerId, lastMsgID, 20)
	if err != nil {
		logger.ErrorWF("OnQueryMessages QueryMessages error", zap.Error(err), zap.Any("req", req))
		return err
	}
	res.PeerId = proto.Uint64(peerId)
	for _, message := range messages {
		res.MsgList = append(res.MsgList, PbMessage(message))
	}
	return nil
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
		logger.ErrorWF("OnSendMessage WrapUser error", zap.Error(err), zap.Any("req", req))
		return err
	}

	messageID, err := p2pservice.Default.SendMessage(s.Context(), logger, app.Maze, user, peerId, _type, content)
	if err != nil {
		logger.ErrorWF("OnSendMessage SendMessage error", zap.Error(err), zap.Any("req", req))
		return err
	}
	res.MsgId = proto.Uint64(messageID)
	res.PeerId = proto.Uint64(peerId)
	return nil
}
