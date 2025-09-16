package p2pservice

import (
	"context"
	"maze_game_server/app"
	"maze_game_server/io/redis/im/msgstore/p2pmsg"
	"maze_game_server/io/redis/im/session"
	"maze_game_server/lib/idgenerator"
	"maze_game_server/services/sessionservice"
	"time"

	"maze_game_server/pb/common/MazeIM"

	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	MessageNotificationID     = 10651 //消息通知ID
	MessageReadNotificationID = 10688 //消息已读通知ID
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
	QueryMessages(ctx context.Context, a app.App, user app.User, peerID int64, lastID uint64, newest bool) (messages []app.Message, err error)

	// SendMessage 向指定用户发送私聊消息
	//
	// 参数:
	//	- a: 应用
	//	- userID: 发送用户
	//	- peerID: 接收用户
	//	- _type: 消息类型
	//	- content: 消息内容
	SendMessage(ctx context.Context, a app.App, user app.User, peerID int64, _type int32, content []byte) (messageID uint64, err error)

	// ReadMessage 标记私聊中的指定消息已读
	//
	// 参数:
	//	- a: 应用
	//	- userID: 发送用户
	//	- peerID: 接收用户
	//	- messageID: 消息ID
	ReadMessage(ctx context.Context, a app.App, user app.User, peerID int64, messageID uint64) (err error)

	// RemoveMessage 删除私聊中的指定消息(只删除自己这边的私聊记录)
	//
	// 参数:
	//	- a: 应用
	//	- userID: 发送用户
	//	- peerID: 接收用户
	//	- messageID: 消息ID
	RemoveMessage(ctx context.Context, a app.App, user app.User, peerID int64, messageID uint64) (err error)

	// checkUserAndPeer 检查用户和接收者是否合法
	//
	// 参数:
	//	- ctx: 上下文
	//	- userId: 发送用户
	//	- peerId: 接收用户
	// 返回值:
	//	- user: 用户
	//	- errInfo: 错误信息(string)
	CheckUserAndPeer(ctx context.Context, userID int64, peerID int64) (user app.User, errInfo string)

	// MessageReadNotify 消息已读通知
	//
	// 参数:
	//	- ctx: 上下文
	//	- userId: 发送用户
	//	- peerId: 接收用户
	//	- messageID: 消息ID
	// 返回值:
	//	- error: 错误信息
	MessageReadNotify(ctx context.Context, userID int64, peerID int64, messageID uint64) error
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
func (p *p2p) QueryMessages(ctx context.Context, a app.App, user app.User, peerID int64, lastID uint64, newest bool) (messages []app.Message, err error) {
	if lastID <= 0 {
		lastID = idgenerator.MaxMessageID
	}
	return p2pmsg.QueryMessages(ctx, a.ID(), int64(user.UserID()), peerID, lastID, newest)
}

// SendMessage implements P2PService.
func (p *p2p) SendMessage(ctx context.Context, a app.App, user app.User, peerID int64, _type int32, content []byte) (messageID uint64, err error) {
	message := app.Message{}
	messageID = idgenerator.MessageID(time.Now().UnixMilli(), idgenerator.GenerateSeed())
	// if err != nil {
	// 	return 0, err
	// }
	message.MessageID = messageID
	message.UserID = int64(user.UserID())
	message.CreateTime = time.Now().Unix()
	message.Type = _type
	message.Content = content
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "SendMessage start", zap.Any("user", user), zap.Int64("peerId", peerID), zap.Any("message", message))
	//发送者的记录
	err = p2pmsg.SaveMessage(ctx, a.ID(), int64(user.UserID()), peerID, message)
	if err != nil {
		logger.CtxError(ctx, "SaveMessage sender error", zap.Error(err))
	}
	//接收者的记录
	err = p2pmsg.SaveMessage(ctx, a.ID(), peerID, int64(user.UserID()), message)
	if err != nil {
		logger.CtxError(ctx, "SaveMessage peer error", zap.Error(err))
	}
	//通知发送者
	err = p.notifyMessage(ctx, int64(user.UserID()), int64(user.UserID()), messageID, _type, MessageNotificationID, content)
	if err != nil {
		logger.CtxError(ctx, "notifyMessage sender error", zap.Error(err))
	}
	// 通知接收者
	err = p.notifyMessage(ctx, int64(user.UserID()), peerID, messageID, _type, MessageNotificationID, content)
	if err != nil {
		logger.CtxError(ctx, "notifyMessage peer error", zap.Error(err))
	}

	//保存发送者会话
	sessionservice.Default.SaveNormalSession(ctx, a, user, peerID, message.CreateTime, false)
	//保存接收者会话
	peerUser, err := app.WrapUser(uint64(peerID), "")
	if err != nil {
		logger.CtxError(ctx, "WrapUser error", zap.Error(err))
		return 0, err
	}
	sessionservice.Default.SaveNormalSession(ctx, a, peerUser, int64(user.UserID()), message.CreateTime, true)
	return messageID, nil
}

// ReadMessage implements P2PService.
func (p *p2p) ReadMessage(ctx context.Context, a app.App, user app.User, peerID int64, messageID uint64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	//这里peerid是此消息发送者的id user.UserID()是接收者的id
	err = p2pmsg.ReadMessage(ctx, a.ID(), peerID, int64(user.UserID()), messageID)
	if err != nil {
		logger.CtxError(ctx, "ReadMessage error", zap.Error(err), zap.Any("messageID", messageID))
		return err
	}
	//session把未读消息设置成0
	sessionID := sessionservice.Default.NormalSessionID(peerID)
	err = session.SetRemoveUnreadCount(ctx, a.ID(), int64(user.UserID()), sessionID)
	if err != nil {
		logger.CtxError(ctx, "SetRemoveUnreadCount error", zap.Error(err), zap.Any("messageID", messageID))
		return err
	}
	//通知接收者
	err = p.MessageReadNotify(ctx, int64(user.UserID()), peerID, messageID)
	if err != nil {
		logger.CtxError(ctx, "notifyReadMessage error", zap.Error(err), zap.Any("messageID", messageID))
		return err
	}
	return nil
}

// RemoveMessage implements P2PService.
func (p *p2p) RemoveMessage(ctx context.Context, a app.App, user app.User, peerID int64, messageID uint64) (err error) {
	return p2pmsg.RemoveMessage(ctx, a.ID(), int64(user.UserID()), peerID, messageID)

}

// notifyMessage 通知接收者
func (p *p2p) notifyMessage(ctx context.Context, userID int64, peerID int64, messageID uint64, _type int32, packId uint16, content []byte) error {
	logger := fklog.ContextAppLogger(ctx)
	userInfo, err := session.GetUserInfo(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "GetUserInfo error", zap.Error(err), zap.Any("userId", userID))
		return err
	}
	// 推送消息给集群
	notifyMessage := &MazeIM.MessageNotificationID{
		UserId: proto.Int64(peerID),
		Message: &MazeIM.Message{
			MsgId:      proto.Uint64(messageID),
			Type:       proto.Int32(_type),
			Content:    []byte(content),
			Sender:     proto.Int64(userID),
			CreateTime: proto.Int64(time.Now().Unix()),
		},
		UserInfo: userInfo,
	}
	logger.CtxInfo(ctx, "notifyMessage start", zap.Int64("peerId", peerID), zap.Any("notifyMessage", notifyMessage))
	err = online.ClusterPush(ctx, uint64(peerID), packId, notifyMessage)
	if err != nil {
		logger.CtxError(ctx, "notifyMessage error", zap.Error(err), zap.Any("notifyMessage", notifyMessage))
	}
	return err
}

// MessageReadNotify 消息已读通知
func (p *p2p) MessageReadNotify(ctx context.Context, userID int64, peerID int64, messageID uint64) error {
	logger := fklog.ContextAppLogger(ctx)
	// 推送消息给消息发送者
	notifyReadMessage := &MazeIM.MessageReadNotificationID{
		PeerId: proto.Int64(userID),
		MsgId:  proto.Uint64(messageID),
	}
	logger.CtxInfo(ctx, "MessageReadNotify start", zap.Int64("userId", userID), zap.Any("notifyMessage", notifyReadMessage))
	err := online.ClusterPush(ctx, uint64(peerID), MessageReadNotificationID, notifyReadMessage)
	if err != nil {
		logger.CtxError(ctx, "MessageReadNotify error", zap.Error(err), zap.Any("notifyMessage", notifyReadMessage))
		return err
	}
	return nil
}

// CheckUserAndPeer 检查用户和接收者是否合法
func (p *p2p) CheckUserAndPeer(ctx context.Context, userID int64, peerID int64) (user app.User, errInfo string) {
	logger := fklog.ContextAppLogger(ctx)
	if userID == peerID {
		errInfo = "用户id不能和接收者id相同"
		logger.CtxError(ctx, "OnSendMessage user and peerId error", zap.Any("err", errInfo), zap.Any("user", userID), zap.Any("peerId", peerID))
		return nil, errInfo
	}
	if userID <= 0 {
		errInfo = "用户id不合法"
		logger.CtxError(ctx, "OnSendMessage user error", zap.Any("err", errInfo), zap.Any("user", userID))
		return nil, errInfo
	}
	//判断用户是否存在
	user, _ = app.WrapUser(uint64(userID), "")
	if user == nil {
		errInfo = "用户不存在"
		logger.CtxError(ctx, "OnSendMessage user error", zap.Any("err", errInfo), zap.Any("user", userID))
		return nil, errInfo
	}
	//判断peerId是否合法
	if peerID <= 0 {
		errInfo = "接收者id不合法"
		logger.CtxError(ctx, "OnSendMessage peerId error", zap.Any("err", errInfo), zap.Any("peerId", peerID))
		return nil, errInfo
	}
	peerUser, _ := app.WrapUser(uint64(peerID), "")
	if peerUser == nil {
		errInfo = "接收者不存在"
		logger.CtxError(ctx, "OnSendMessage peerUser error", zap.Any("err", errInfo), zap.Any("peerId", peerID))
		return nil, errInfo
	}
	//TODO 是否黑名单、禁言、拒接聊天等判断
	return user, ""
}
