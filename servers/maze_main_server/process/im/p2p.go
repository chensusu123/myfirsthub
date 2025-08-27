package im

import (
	"maze_game_server/app"
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeIM"
	"maze_game_server/services/p2pservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (im *IM) OnQueryMessages_10643_10644(s *session.Session, req *MazeIM.QueryMessagesRQ) (err error) {
	defer fkprometheus.InfoPMT("OnQueryMessages")()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)

	res := &MazeIM.QueryMessagesRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	logger.CtxInfo(ctx, "OnQueryMessages start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnQueryMessages end", zap.Any("res", res))
	}()

	var (
		userId    = uint64(s.UID())
		peerId    = req.GetPeerId()
		lastMsgID = req.GetLastMsgId()
	)

	user, err := app.WrapUser(userId, "")
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("获取用户信息失败")
		logger.CtxError(ctx, "OnQueryMessages WrapUser error", zap.Error(err), zap.Any("req", req))
		return err
	}

	messages, err := p2pservice.Default.QueryMessages(ctx, app.Maze, user, peerId, lastMsgID, 20)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("获取消息失败")
		logger.CtxError(ctx, "OnQueryMessages QueryMessages error", zap.Error(err), zap.Any("req", req))
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
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)

	res := &MazeIM.SendMessageRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	logger.CtxInfo(ctx, "OnSendMessage start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnSendMessage end", zap.Any("res", res))
	}()

	var (
		userId  = uint64(s.UID())
		peerId  = req.GetPeerId()
		_type   = req.GetType()
		content = req.GetContent()
	)

	user, err := app.WrapUser(userId, "")
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("获取用户信息失败")
		logger.CtxError(ctx, "OnSendMessage WrapUser error", zap.Error(err), zap.Any("req", req))
		return err
	}

	messageID, err := p2pservice.Default.SendMessage(ctx, app.Maze, user, peerId, _type, content)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("发送消息失败")
		logger.CtxError(ctx, "OnSendMessage SendMessage error", zap.Error(err), zap.Any("req", req))
		return err
	}
	res.MsgId = proto.Uint64(messageID)
	res.PeerId = proto.Uint64(peerId)
	return nil
}

func (im *IM) OnReadMessage_10656_10657(s *session.Session, req *MazeIM.ReadMessageRQ) (err error) {
	defer fkprometheus.InfoPMT("OnReadMessage")()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)

	res := &MazeIM.ReadMessageRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	logger.CtxInfo(ctx, "OnReadMessage start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnReadMessage end", zap.Any("res", res))
	}()

	var (
		userId = uint64(s.UID())
		peerId = req.GetPeerId()
		msgID  = req.GetMsgId()
	)

	user, err := app.WrapUser(userId, "")
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数错误")
		logger.CtxError(ctx, "OnReadMessage WrapUser error", zap.Error(err), zap.Any("req", req))
		return err
	}

	err = p2pservice.Default.ReadMessage(ctx, app.Maze, user, peerId, msgID)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		logger.CtxError(ctx, "OnReadMessage SendMessage error", zap.Error(err), zap.Any("req", req))
		return err
	}
	return nil
}
