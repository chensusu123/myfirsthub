package im

import (
	"maze_game_server/app"
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeIM"
	"maze_game_server/services/groupservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

func (im *IM) OnQueryGroupMessages_10647_10648(s *session.Session, req *MazeIM.QueryGroupMessagesRQ) (err error) {
	defer fkprometheus.InfoPMT("OnQueryGroupMessages")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeIM.QueryGroupMessagesRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	logger.InfoWF("OnQueryGroupMessages start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnQueryGroupMessages end", zap.Any("res", res))
	}()

	var (
		userId    = uint64(s.UID())
		groupId   = req.GetGroupId()
		lastMsgID = req.GetLastMsgId()
	)

	_ = userId

	messages, err := groupservice.Default.QueryMessages(s.Context(), logger, app.Maze, groupId, lastMsgID, 20)
	_ = messages
	return
}

func (im *IM) OnSendGroupMessage_10649_10650(s *session.Session, req *MazeIM.SendGroupMessageRQ) (err error) {
	defer fkprometheus.InfoPMT("OnSendGroupMessage")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeIM.SendGroupMessageRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	logger.InfoWF("OnSendGroupMessage start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnSendGroupMessage end", zap.Any("res", res))
	}()

	var (
		userId  = uint64(s.UID())
		groupId = req.GetGroupId()
		_type   = req.GetType()
		content = req.GetContent()
	)

	messageID, err := groupservice.Default.SendMessage(s.Context(), logger, app.Maze, groupId, userId, _type, content)
	_ = messageID
	return
}
