package im

import (
	"maze_game_server/app"
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeIM"
	"maze_game_server/services/groupservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (im *IM) OnQueryGroupMessages_10647_10648(s *session.Session, req *MazeIM.QueryGroupMessagesRQ) (err error) {
	defer fkprometheus.InfoPMT("OnQueryGroupMessages")()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)

	res := &MazeIM.QueryGroupMessagesRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	logger.CtxInfo(ctx, "OnQueryGroupMessages start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnQueryGroupMessages end", zap.Any("res", res))
	}()

	var (
		groupId   = req.GetGroupId()
		lastMsgID = req.GetLastMsgId()
	)

	messages, err := groupservice.Default.QueryMessages(ctx, app.Maze, groupId, lastMsgID, 20)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("获取消息失败")
		logger.CtxError(ctx, "OnQueryGroupMessages QueryMessages error", zap.Error(err))
		return
	}
	for _, message := range messages {
		res.MsgList = append(res.MsgList, PbMessage(message))
	}
	return
}

func (im *IM) OnSendGroupMessage_10649_10650(s *session.Session, req *MazeIM.SendGroupMessageRQ) (err error) {
	defer fkprometheus.InfoPMT("OnSendGroupMessage")()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)

	res := &MazeIM.SendGroupMessageRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	logger.CtxInfo(ctx, "OnSendGroupMessage start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnSendGroupMessage end", zap.Any("res", res))
	}()

	var (
		userId  = uint64(s.UID())
		groupId = req.GetGroupId()
		_type   = req.GetType()
		content = req.GetContent()
	)

	messageID, err := groupservice.Default.SendMessage(ctx, app.Maze, groupId, userId, _type, content)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.Wrap("发送群聊消息失败")
		logger.CtxInfo(ctx, "OnSendGroupMessage SendMessage error", zap.Error(err))
		return
	}
	res.GroupId = proto.Int32(groupId)
	res.MsgId = proto.Uint64(messageID)
	return
}

// 订阅群组消息 groupids填写这些群组的id
func (im *IM) SubscribeMessages_10660_10661(s *session.Session, req *MazeIM.SubMessageRQ) (err error) {
	defer fkprometheus.InfoPMT("SubscribeMessages")()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)

	res := &MazeIM.SubMessageRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	logger.CtxInfo(ctx, "SubscribeMessages start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "SubscribeMessages end", zap.Any("res", res))
	}()

	// var (
	// 	userId   = uint64(s.UID())
	// 	groupIds = req.GetGroupIds()
	// )

	// user, err := app.WrapUser(userId, "")
	// if err != nil {
	// 	res.ErrInfo = errors.MODULE_ERROR.Wrap("获取用户信息失败")
	// 	logger.CtxError(ctx, "SubscribeMessages WrapUser error", zap.Error(err))
	// 	return err
	// }

	// err = groupservice.Default.SubscribeMessages(ctx, app.Maze, user, groupIds)
	// if err != nil {
	// 	res.ErrInfo = errors.MODULE_ERROR.Wrap("订阅消息失败")
	// 	logger.CtxError(ctx, "SubscribeMessages SubscribeMessages error", zap.Error(err))
	// 	return err
	// }
	return
}
