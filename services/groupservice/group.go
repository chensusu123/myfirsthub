package groupservice

import (
	"context"
	"maze_game_server/app"
	grouppkg "maze_game_server/io/redis/im/group"
	"maze_game_server/io/redis/im/msgstore/groupmsg"
	"maze_game_server/lib/idgenerator"
	"time"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
)

type GroupService interface {
	// QueryMessages 查询群组聊天消息
	//
	// 参数:
	//	- a: 应用
	//	- groupID: 组ID
	//	- lastID: 客户端的最后一条消息ID
	// 	- limit: 读取时限制读取条数
	QueryMessages(ctx context.Context, logger fklog.FKLogI, a app.App, groupID int32, lastID uint64, limit int) (messages []app.Message, err error)

	// SendMessage 向群组中发送消息
	//
	// 参数:
	//	- a: 应用
	//	- groupID: 组ID
	// 	- sender: 消息发送者
	//	- _type: 消息类型
	//	- content: 消息内容
	SendMessage(ctx context.Context, logger fklog.FKLogI, a app.App, groupID int32, sender uint64, _type int32, content string) (messageID uint64, err error)

	// GetGroupInfo 获取聊天组信息
	//
	// 参数:
	//	- a: 应用
	//	- groupID: 组ID
	GetGroupInfo(ctx context.Context, logger fklog.FKLogI, a app.App, groupID int32) (g *app.Group, err error)

	// CreateGroup 创建聊天组
	//
	// 参数:
	//	- a: 应用
	//	- user: 群主用户
	//	- memberIDs: 成员ID
	CreateGroup(ctx context.Context, logger fklog.FKLogI, a app.App, user app.User, memberIDs []uint64) (g *app.Group, err error)

	// InviteMember 群组邀请成员
	//
	// 参数:
	//	- a: 应用
	//	- groupID: 组ID
	//	- memberID: 成员ID
	InviteMember(ctx context.Context, logger fklog.FKLogI, a app.App, groupID int32, memberID uint64) (err error)
}

var (
	Default GroupService = &group{}
)

type group struct {
}

// QueryMessages implements GroupService.
func (g *group) QueryMessages(ctx context.Context, logger fklog.FKLogI, a app.App, groupID int32, lastID uint64, limit int) (messages []app.Message, err error) {
	if lastID <= 0 {
		lastID = idgenerator.MaxMessageID
	}
	if limit <= 0 {
		limit = 20
	}
	return groupmsg.QueryMessages(logger, a.ID(), groupID, lastID, limit)
}

// SendMessage implements GroupService.
func (g *group) SendMessage(ctx context.Context, logger fklog.FKLogI, a app.App, groupID int32, sender uint64, _type int32, content string) (messageID uint64, err error) {
	message := app.Message{}
	messageID, err = idgenerator.NextID()
	if err != nil {
		return 0, err
	}
	message.MessageID = messageID
	message.UserID = sender
	message.CreateTime = time.Now().Unix()
	message.Type = _type
	message.Content = content
	err = groupmsg.SaveMessage(logger, a.ID(), groupID, message)
	return
}

// GetGroupInfo implements GroupService.
func (*group) GetGroupInfo(ctx context.Context, logger fklog.FKLogI, a app.App, groupID int32) (g *app.Group, err error) {
	return grouppkg.GetGroupInfo(logger, a.ID(), groupID)
}

// CreateGroup implements GroupService.
func (*group) CreateGroup(ctx context.Context, logger fklog.FKLogI, a app.App, user app.User, memberIDs []uint64) (g *app.Group, err error) {
	// TODO 群组ID生成器
	groupID := int32(time.Now().Unix())
	return grouppkg.CreateGroup(logger, a.ID(), user.UserID(), groupID, memberIDs)
}

// InviteMember implements GroupService.
func (g *group) InviteMember(ctx context.Context, logger fklog.FKLogI, a app.App, groupID int32, memberID uint64) (err error) {
	return grouppkg.InviteMember(logger, a.ID(), groupID, memberID)
}
