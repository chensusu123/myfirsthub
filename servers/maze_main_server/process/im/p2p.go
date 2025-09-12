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

func (im *IM) OnMessageList_10663_10664(s *session.Session, req *MazeIM.MessageListRQ) (err error) {
	defer fkprometheus.InfoPMT("OnMessageList")()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)

	res := &MazeIM.MessageListRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	logger.CtxInfo(ctx, "OnMessageList start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnMessageList end", zap.Any("res", res))
	}()

	var (
		userId    = s.UID()
		peerId    = req.GetPeerId()
		lastMsgID = req.GetLastMsgId()
		newest    = req.GetNewest()
	)

	// 检查用户和对端是否存在
	user, errInfo := p2pservice.Default.CheckUserAndPeer(ctx, userId, peerId)
	if errInfo != "" {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(errInfo)
		logger.CtxError(ctx, "OnMessageList CheckUserAndPeer error", zap.Error(err), zap.Any("req", req))
		return err
	}

	messages, err := p2pservice.Default.QueryMessages(ctx, app.Maze, user, peerId, lastMsgID, newest)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("获取消息失败")
		logger.CtxError(ctx, "OnMessageList QueryMessages error", zap.Error(err), zap.Any("req", req))
		return err
	}
	res.PeerId = proto.Int64(peerId)
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
		userId  = s.UID()
		peerId  = req.GetPeerId()
		_type   = req.GetType()
		content = req.GetContent()
	)

	// 检查用户和对端是否存在
	user, errInfo := p2pservice.Default.CheckUserAndPeer(ctx, userId, peerId)
	if errInfo != "" {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(errInfo)
		logger.CtxError(ctx, "OnQueryMessages CheckUserAndPeer error", zap.Error(err), zap.Any("req", req))
		return err
	}

	messageID, err := p2pservice.Default.SendMessage(ctx, app.Maze, user, peerId, _type, content)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("发送消息失败")
		logger.CtxError(ctx, "OnSendMessage SendMessage error", zap.Error(err), zap.Any("req", req))
		return err
	}
	res.MsgId = proto.Uint64(messageID)
	res.PeerId = proto.Int64(peerId)
	res.Type = proto.Int32(_type)
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
		userId = s.UID()
		peerId = req.GetPeerId()
		msgID  = req.GetMsgId()
	)

	// 检查用户和对端是否存在
	user, errInfo := p2pservice.Default.CheckUserAndPeer(ctx, userId, peerId)
	if errInfo != "" {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(errInfo)
		logger.CtxError(ctx, "OnQueryMessages CheckUserAndPeer error", zap.Error(err), zap.Any("req", req))
		return err
	}

	err = p2pservice.Default.ReadMessage(ctx, app.Maze, user, peerId, msgID)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("已读消息失败")
		logger.CtxError(ctx, "OnReadMessage SendMessage error", zap.Error(err), zap.Any("req", req))
		return err
	}
	return nil
}
