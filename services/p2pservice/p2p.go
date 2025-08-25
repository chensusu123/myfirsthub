package p2pservice

import (
	"context"
	"maze_game_server/app"
	"maze_game_server/io/redis/im/msgstore/p2pmsg"
	"maze_game_server/lib/idgenerator"
	"time"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
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
	QueryMessages(ctx context.Context, logger fklog.FKLogI, a app.App, user app.User, peerID uint64, lastID uint64, limit int) (messages []app.Message, err error)

	// SendMessage 向指定用户发送私聊消息
	//
	// 参数:
	//	- a: 应用
	//	- userID: 发送用户
	//	- peerID: 接收用户
	//	- _type: 消息类型
	//	- content: 消息内容
	SendMessage(ctx context.Context, logger fklog.FKLogI, a app.App, user app.User, peerID uint64, _type int32, content string) (messageID uint64, err error)

	// RemoveMessage 删除私聊中的指定消息(只删除自己这边的私聊记录)
	//
	// 参数:
	//	- a: 应用
	//	- userID: 发送用户
	//	- peerID: 接收用户
	//	- messageID: 消息ID
	RemoveMessage(ctx context.Context, logger fklog.FKLogI, a app.App, user app.User, peerID uint64, messageID uint64) (err error)
}

var (
	Default P2PService = &p2p{}
)

type p2p struct {
}

// QueryMessages implements P2PService.
func (p *p2p) QueryMessages(ctx context.Context, logger fklog.FKLogI, a app.App, user app.User, peerID uint64, lastID uint64, limit int) (messages []app.Message, err error) {
	if lastID <= 0 {
		lastID = idgenerator.MaxMessageID
	}
	if limit <= 0 {
		limit = 20
	}
	return p2pmsg.QueryMessages(logger, a.ID(), user.UserID(), peerID, lastID, limit)
}

// SendMessage implements P2PService.
func (p *p2p) SendMessage(ctx context.Context, logger fklog.FKLogI, a app.App, user app.User, peerID uint64, _type int32, content string) (messageID uint64, err error) {
	message := app.Message{}
	messageID, err = idgenerator.NextID()
	if err != nil {
		return 0, err
	}
	message.MessageID = messageID
	message.UserID = user.UserID()
	message.CreateTime = time.Now().Unix()
	message.Type = _type
	message.Content = content
	err = p2pmsg.SaveMessage(logger, a.ID(), user.UserID(), peerID, message)
	if err == nil {
		err = p2pmsg.SaveMessage(logger, a.ID(), peerID, user.UserID(), message)
	}
	return
}

// RemoveMessage implements P2PService.
func (p *p2p) RemoveMessage(ctx context.Context, logger fklog.FKLogI, a app.App, user app.User, peerID uint64, messageID uint64) (err error) {
	return
}
