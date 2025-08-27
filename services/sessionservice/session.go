package sessionservice

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"maze_game_server/app"
	"maze_game_server/io/redis/im/msgstore/p2pmsg"
	sessionpkg "maze_game_server/io/redis/im/session"

	"maze_game_server/pb/common/MazeIM"

	"google.golang.org/protobuf/proto"
)

type SessionService interface {
	// QueryRecentSessions 查询用户最近的会话记录
	//
	// 参数:
	//	- a: 应用
	// 	- user: 用户标识
	QueryRecentSessions(ctx context.Context, a app.App, user app.User) (sessions map[string]app.Session, err error)

	// CreateNormalSession 创建普通私聊会话
	//
	// 参数:
	//	- a: 应用
	// 	- user: 用户标识
	//	- peerID: 对方ID
	CreateNormalSession(ctx context.Context, a app.App, user app.User, peerID uint64) (err error)

	// CreateGroupSession 创建群聊会话
	//
	// 参数:
	//	- a: 应用
	// 	- user: 用户标识
	//	- groupID: 群组ID
	CreateGroupSession(ctx context.Context, a app.App, user app.User, groupID int32) (err error)

	// RemoveSession 删除会话记录
	//
	// 参数:
	//	- a: 应用
	// 	- user: 用户标识
	//	- sessionID: 会话ID
	RemoveSession(ctx context.Context, a app.App, user app.User, sessionID string) (err error)
	// GetMessageInfo 获取会话消息信息
	//
	// 参数:
	//	- a: 应用
	// 	- user: 用户标识
	//	- sessions: 会话列表
	GetMessageInfo(ctx context.Context, a app.App, user app.User, sessions map[string]app.Session) (messageInfo []*MazeIM.Session, err error)
}

var (
	Default SessionService = &session{}
)

type session struct {
}

// QueryRecentSessions implements SessionService.
func (s *session) QueryRecentSessions(ctx context.Context, a app.App, user app.User) (sessions map[string]app.Session, err error) {
	return sessionpkg.QuerySessions(ctx, a.ID(), user.UserID())
}

// CreateNormalSession implements SessionService.
func (s *session) CreateNormalSession(ctx context.Context, a app.App, user app.User, peerID uint64) (err error) {
	sessionID := s.NormalSessionID(peerID)
	return sessionpkg.AddP2PSession(ctx, a.ID(), user.UserID(), sessionID, peerID)
}

// CreateGroupSession implements SessionService.
func (s *session) CreateGroupSession(ctx context.Context, a app.App, user app.User, groupID int32) (err error) {
	sessionID := s.GroupSessionID(groupID)
	return sessionpkg.AddGroupSession(ctx, a.ID(), user.UserID(), sessionID, groupID)
}

// RemoveSession implements SessionService.
func (s *session) RemoveSession(ctx context.Context, a app.App, user app.User, sessionID string) (err error) {
	return sessionpkg.RemoveSession(ctx, a.ID(), user.UserID(), sessionID)
}

func (s *session) NormalSessionID(userID uint64) string {
	return sum([]byte(fmt.Sprintf("p2p:%d", userID)))
}

func (s *session) GroupSessionID(groupID int32) string {
	return sum([]byte(fmt.Sprintf("group:%d", groupID)))
}

// sha1 returns the SHA-1 checksum of the data.
func sum(data []byte) string {
	sum := sha1.Sum(data)
	return hex.EncodeToString(sum[:])
}

func (s *session) GetMessageInfo(ctx context.Context, a app.App, user app.User, sessions map[string]app.Session) (messageInfo []*MazeIM.Session, err error) {
	if len(sessions) == 0 {
		return nil, nil
	}
	for _, session := range sessions {
		peerID := session.PeerID
		if peerID > 0 {
			p2pmsg, err := p2pmsg.QueryMessages(ctx, a.ID(), user.UserID(), peerID, uint64(0), 10)
			if err != nil {
				return nil, err
			}
			if len(p2pmsg) > 0 {
				messageInfo = append(messageInfo, &MazeIM.Session{
					SessionId:  proto.String(session.ID),
					Type:       proto.Int32(1),
					GroupId:    proto.Int32(0),
					CreateTime: proto.Int64(session.CreateTime),
					Recent:     PbSessionMessage(p2pmsg),
				})
			}
		} else {
			groupID := session.GroupID
			if groupID > 0 {
				groupmsg, err := p2pmsg.QueryMessages(ctx, a.ID(), user.UserID(), 0, uint64(0), 10)
				if err != nil {
					return nil, err
				}
				if len(groupmsg) > 0 {
					messageInfo = append(messageInfo, &MazeIM.Session{
						SessionId:  proto.String(session.ID),
						Type:       proto.Int32(2),
						GroupId:    proto.Int32(groupID),
						CreateTime: proto.Int64(session.CreateTime),
						Recent:     PbSessionMessage(groupmsg),
					})
				}
			}
		}
	}
	return messageInfo, nil
}
func PbSessionMessage(messages []p2pmsg.Message) []*MazeIM.Message {
	pbMessages := make([]*MazeIM.Message, 0, len(messages))
	for _, message := range messages {
		pbMessages = append(pbMessages, &MazeIM.Message{
			MsgId:   proto.Uint64(message.MessageID),
			Type:    proto.Int32(message.Type),
			Content: []byte(message.Content),
			Sender:  proto.Uint64(message.UserID),
		})
	}
	return pbMessages
}
