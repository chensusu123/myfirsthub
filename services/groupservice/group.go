package groupservice

import (
	"context"
	"maze_game_server/app"
	grouppkg "maze_game_server/io/redis/im/group"
	"maze_game_server/io/redis/im/msgstore/groupmsg"
	"maze_game_server/lib/idgenerator"
	"maze_game_server/services/sessionservice"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"

	"maze_game_server/pb/common/MazeIM"

	"maze_game_server/usecase/online"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type GroupService interface {
	// QueryMessages 查询群组聊天消息
	//
	// 参数:
	//	- a: 应用
	//	- groupID: 组ID
	//	- lastID: 客户端的最后一条消息ID
	// 	- limit: 读取时限制读取条数
	QueryMessages(ctx context.Context, a app.App, groupID int64, lastID uint64, limit int) (messages []app.Message, err error)

	// SendMessage 向群组中发送消息
	//
	// 参数:
	//	- a: 应用
	//	- groupID: 组ID
	// 	- sender: 消息发送者
	//	- _type: 消息类型
	//	- content: 消息内容
	SendMessage(ctx context.Context, a app.App, groupID int64, sender uint64, _type int32, content []byte) (messageID uint64, err error)

	// GetGroupInfo 获取聊天组信息
	//
	// 参数:
	//	- a: 应用
	//	- groupID: 组ID
	GetGroupInfo(ctx context.Context, a app.App, groupID int64) (g *app.Group, err error)

	// CreateGroup 创建聊天组
	//
	// 参数:
	//	- a: 应用
	//	- user: 群主用户
	//	- memberIDs: 成员ID
	CreateGroup(ctx context.Context, a app.App, userID uint64, memberIDs []uint64) (g *app.Group, err error)

	// InviteMember 群组邀请成员
	//
	// 参数:
	//	- a: 应用
	//	- groupID: 组ID
	//	- memberID: 成员ID
	InviteMember(ctx context.Context, a app.App, groupID int64, memberID uint64) (err error)
}

var (
	Default GroupService = &group{}
)

type group struct {
}

var GlobalGroupService = newGroupService()

func init() {
	GlobalGroupService = newGroupService()
}

func newGroupService() GroupService {
	return &group{}
}

// QueryMessages implements GroupService.
func (g *group) QueryMessages(ctx context.Context, a app.App, groupID int64, lastID uint64, limit int) (messages []app.Message, err error) {
	if lastID <= 0 {
		lastID = idgenerator.MaxMessageID
	}
	if limit <= 0 {
		limit = 20
	}
	return groupmsg.QueryMessages(ctx, a.ID(), groupID, lastID, limit)
}

// SendMessage implements GroupService.
func (g *group) SendMessage(ctx context.Context, a app.App, groupID int64, sender uint64, _type int32, content []byte) (messageID uint64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	message := app.Message{}
	messageID = idgenerator.MessageID(time.Now().UnixMilli(), idgenerator.GenerateSeed())
	// if err != nil {
	// 	return 0, err
	// }
	//TODO 需要判断这个组是否存在，通过家族或者联盟判断
	message.MessageID = messageID
	message.UserID = sender
	message.CreateTime = time.Now().Unix()
	message.Type = _type
	message.Content = content
	err = groupmsg.SaveMessage(ctx, a.ID(), groupID, message)
	if err != nil {
		logger.CtxError(ctx, "SendMessage error", zap.Error(err), zap.Uint64("sender", sender), zap.Int32("_type", _type), zap.ByteString("content", content))
		return 0, err
	}
	err = g.notifyGroupMessage(ctx, a, sender, messageID, _type, groupID, 10651, content)
	if err != nil {
		logger.CtxError(ctx, "notifyGroupMessage error", zap.Error(err), zap.Uint64("sender", sender), zap.Int32("_type", _type), zap.ByteString("content", content))
		return 0, err
	}
	//创建群组成员会话
	user, err := app.WrapUser(sender, "")
	if err != nil {
		logger.CtxError(ctx, "WrapUser error", zap.Error(err), zap.Uint64("sender", sender))
		return 0, err
	}
	sessionservice.Default.CreateGroupSession(ctx, a, user, groupID)
	return messageID, nil
}

// GetGroupInfo implements GroupService.
func (*group) GetGroupInfo(ctx context.Context, a app.App, groupID int64) (g *app.Group, err error) {
	return grouppkg.GetGroupInfo(ctx, a.ID(), groupID)
}

// CreateGroup implements GroupService.
func (*group) CreateGroup(ctx context.Context, a app.App, userID uint64, memberIDs []uint64) (g *app.Group, err error) {
	// TODO 群组ID生成器
	groupID, err := idgenerator.NextID()
	if err != nil {
		return nil, err
	}
	return grouppkg.CreateGroup(ctx, a.ID(), userID, groupID, memberIDs)
}

// InviteMember implements GroupService.
func (g *group) InviteMember(ctx context.Context, a app.App, groupID int64, memberID uint64) (err error) {
	return grouppkg.InviteMember(ctx, a.ID(), groupID, memberID)
}

func (g *group) notifyGroupMessage(ctx context.Context, a app.App, userId uint64, messageID uint64, _type int32, groupID int64, packId uint16, content []byte) error {
	logger := fklog.ContextAppLogger(ctx)
	groupInfo, err := g.GetGroupInfo(ctx, a, groupID)
	if err != nil {
		logger.CtxError(ctx, "notifyGroupMessage error", zap.Error(err), zap.Uint64("userId", userId), zap.Uint64("messageID", messageID), zap.Int32("_type", _type), zap.Int64("groupID", groupID), zap.ByteString("content", content))
		return err
	}

	for _, member := range groupInfo.Members {
		if member.UserID == userId || member.UserID <= 0 {
			continue
		}
		// 推送消息给集群
		notifyMessage := &MazeIM.GroupMessageNotificationID{
			GroupId: proto.Int64(groupID),
			Message: &MazeIM.Message{
				MsgId:      proto.Uint64(messageID),
				Type:       proto.Int32(_type),
				Content:    []byte(content),
				Sender:     proto.Uint64(userId),
				CreateTime: proto.Int64(time.Now().Unix()),
			},
		}
		logger.CtxInfo(ctx, "notifyMessage group", zap.Any("notifyMessage group", notifyMessage))
		err = online.ClusterPush(ctx, userId, packId, notifyMessage)
		if err != nil {
			logger.CtxError(ctx, "notifyMessage error", zap.Error(err), zap.Any("notifyMessage", notifyMessage))
		}
		//创建群组成员会话
		user, err := app.WrapUser(member.UserID, "")
		if err != nil {
			logger.CtxError(ctx, "WrapUser error", zap.Error(err), zap.Uint64("userId", member.UserID))
			continue
		}
		sessionservice.Default.CreateGroupSession(ctx, a, user, groupID)
	}
	return nil
}

func (g *group) SubscribeMessages(ctx context.Context, a app.App, user app.User, groupIds []int32) (err error) {
	// logger := fklog.ContextAppLogger(ctx)

	return
}
