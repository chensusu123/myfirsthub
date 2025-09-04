package p2pservice

import (
	"context"
	"fmt"
	"maze_game_server/app"
	"maze_game_server/io/redis/im/msgstore/p2pmsg"
	"maze_game_server/lib/idgenerator"
	"maze_game_server/services/sessionservice"
	"time"

	"maze_game_server/pb/common/MazeIM"

	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type P2PService interface {
	// QueryMessages 查询群组聊天消息
	//
	// 参数:
	//	- a: 应用
	//	- userID: 发送用户
	//	- peerID: 接收用户
	//	- lastID: 客户端的最后一条消息ID
	// 	- limit: 读取时限制读取条数
	QueryMessages(ctx context.Context, a app.App, user app.User, peerID uint64, lastID uint64, limit int) (messages []app.Message, err error)

	// SendMessage 向指定用户发送私聊消息
	//
	// 参数:
	//	- a: 应用
	//	- userID: 发送用户
	//	- peerID: 接收用户
	//	- _type: 消息类型
	//	- content: 消息内容
	SendMessage(ctx context.Context, a app.App, user app.User, peerID uint64, _type int32, content []byte) (messageID uint64, err error)

	// ReadMessage 标记私聊中的指定消息已读
	//
	// 参数:
	//	- a: 应用
	//	- userID: 发送用户
	//	- peerID: 接收用户
	//	- messageID: 消息ID
	ReadMessage(ctx context.Context, a app.App, user app.User, peerID uint64, messageID uint64) (err error)

	// RemoveMessage 删除私聊中的指定消息(只删除自己这边的私聊记录)
	//
	// 参数:
	//	- a: 应用
	//	- userID: 发送用户
	//	- peerID: 接收用户
	//	- messageID: 消息ID
	RemoveMessage(ctx context.Context, a app.App, user app.User, peerID uint64, messageID uint64) (err error)

	// checkUserAndPeer 检查用户和接收者是否合法
	//
	// 参数:
	//	- ctx: 上下文
	//	- userId: 发送用户
	//	- peerId: 接收用户
	// 返回值:
	//	- user: 用户
	//	- errInfo: 错误信息(string)
	CheckUserAndPeer(ctx context.Context, userId uint64, peerId uint64) (user app.User, errInfo string)
}

var (
	Default P2PService = &p2p{}
)

type p2p struct {
}

var GlobalP2PService P2PService

func init() {
	GlobalP2PService = newP2PService()
}

func newP2PService() P2PService {
	return &p2p{}
}

// QueryMessages implements P2PService.
func (p *p2p) QueryMessages(ctx context.Context, a app.App, user app.User, peerID uint64, lastID uint64, limit int) (messages []app.Message, err error) {
	if lastID <= 0 {
		lastID = idgenerator.MaxMessageID
	}
	if limit <= 0 {
		limit = 20
	}
	return p2pmsg.QueryMessages(ctx, a.ID(), user.UserID(), peerID, lastID, limit)
}

// SendMessage implements P2PService.
func (p *p2p) SendMessage(ctx context.Context, a app.App, user app.User, peerID uint64, _type int32, content []byte) (messageID uint64, err error) {
	message := app.Message{}
	messageID = idgenerator.MessageID(time.Now().UnixMilli(), idgenerator.GenerateSeed())
	// if err != nil {
	// 	return 0, err
	// }
	message.MessageID = messageID
	message.UserID = user.UserID()
	message.CreateTime = time.Now().Unix()
	message.Type = _type
	message.Content = content
	logger := fklog.ContextAppLogger(ctx)
	//发送者的记录
	err = p2pmsg.SaveMessage(ctx, a.ID(), user.UserID(), peerID, message)
	if err != nil {
		logger.CtxError(ctx, "SaveMessage sender error", zap.Error(err))
	}
	//接收者的记录
	err = p2pmsg.SaveMessage(ctx, a.ID(), peerID, user.UserID(), message)
	if err != nil {
		logger.CtxError(ctx, "SaveMessage peer error", zap.Error(err))
	}
	// 通知接收者
	err = p.notifyMessage(ctx, user.UserID(), peerID, messageID, _type, 10651, content)
	if err != nil {
		logger.CtxError(ctx, "notifyMessage error", zap.Error(err))
	}
	//创建发送者会话
	sessionservice.Default.CreateNormalSession(ctx, a, user, peerID)
	//创建接收者会话
	peerUser, err := app.WrapUser(peerID, "")
	if err != nil {
		logger.CtxError(ctx, "WrapUser error", zap.Error(err))
		return 0, err
	}
	sessionservice.Default.CreateNormalSession(ctx, a, peerUser, user.UserID())
	return messageID, nil
}

// ReadMessage implements P2PService.
func (p *p2p) ReadMessage(ctx context.Context, a app.App, user app.User, peerID uint64, messageID uint64) (err error) {
	//这里peerid是此消息发送者的id user.UserID()是接收者的id
	err = p2pmsg.ReadMessage(ctx, a.ID(), peerID, user.UserID(), messageID)
	if err != nil {
		fmt.Println("ReadMessage error:", err)
		return err
	}
	//通知接收者
	err = p.notifyMessage(ctx, user.UserID(), peerID, messageID, 1, 10656, []byte(""))
	if err != nil {
		fmt.Println("notifyMessage error:", err)
		return err
	}
	return nil
}

// RemoveMessage implements P2PService.
func (p *p2p) RemoveMessage(ctx context.Context, a app.App, user app.User, peerID uint64, messageID uint64) (err error) {
	return p2pmsg.RemoveMessage(ctx, a.ID(), user.UserID(), peerID, messageID)

}

// notifyMessage 通知接收者
func (p *p2p) notifyMessage(ctx context.Context, userId uint64, peerID uint64, messageID uint64, _type int32, packId uint16, content []byte) error {
	logger := fklog.ContextAppLogger(ctx)
	// 推送消息给集群
	notifyMessage := &MazeIM.MessageNotificationID{
		UserId: proto.Uint64(peerID),
		Message: &MazeIM.Message{
			MsgId:      proto.Uint64(messageID),
			Type:       proto.Int32(_type),
			Content:    []byte(content),
			Sender:     proto.Uint64(userId),
			CreateTime: proto.Int64(time.Now().Unix()),
		},
	}
	logger.CtxInfo(ctx, "notifyMessage start", zap.Uint64("peerId", peerID), zap.Any("notifyMessage", notifyMessage))
	err := online.ClusterPush(ctx, peerID, packId, notifyMessage)
	if err != nil {
		logger.CtxError(ctx, "notifyMessage error", zap.Error(err), zap.Any("notifyMessage", notifyMessage))
	}
	return err
}

// CheckUserAndPeer 检查用户和接收者是否合法
func (p *p2p) CheckUserAndPeer(ctx context.Context, userId uint64, peerId uint64) (user app.User, errInfo string) {
	logger := fklog.ContextAppLogger(ctx)
	if userId == peerId {
		errInfo = "用户id不能和接收者id相同"
		logger.CtxError(ctx, "OnSendMessage user and peerId error", zap.Any("err", errInfo), zap.Any("user", userId), zap.Any("peerId", peerId))
		return nil, errInfo
	}
	if userId <= 0 {
		errInfo = "用户id不合法"
		logger.CtxError(ctx, "OnSendMessage user error", zap.Any("err", errInfo), zap.Any("user", userId))
		return nil, errInfo
	}
	//判断用户是否存在
	user, _ = app.WrapUser(userId, "")
	if user == nil {
		errInfo = "用户不存在"
		logger.CtxError(ctx, "OnSendMessage user error", zap.Any("err", errInfo), zap.Any("user", userId))
		return nil, errInfo
	}
	//判断peerId是否合法
	if peerId <= 0 {
		errInfo = "接收者id不合法"
		logger.CtxError(ctx, "OnSendMessage peerId error", zap.Any("err", errInfo), zap.Any("peerId", peerId))
		return nil, errInfo
	}
	peerUser, _ := app.WrapUser(peerId, "")
	if peerUser == nil {
		errInfo = "接收者不存在"
		logger.CtxError(ctx, "OnSendMessage peerUser error", zap.Any("err", errInfo), zap.Any("peerId", peerId))
		return nil, errInfo
	}
	//TODO 是否黑名单、禁言、拒接聊天等判断
	return user, ""
}
