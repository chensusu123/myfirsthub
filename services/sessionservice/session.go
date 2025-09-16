package sessionservice

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"maze_game_server/app"
	"maze_game_server/io/redis/im/msgstore"
	"maze_game_server/io/redis/im/msgstore/p2pmsg"
	sessionpkg "maze_game_server/io/redis/im/session"
	"maze_game_server/usecase/online"

	"maze_game_server/pb/common/MazeIM"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	SessionChangeID = 10693
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
	CreateNormalSession(ctx context.Context, a app.App, user app.User, peerID int64, messageTime int64, isReceiver bool) (err error)

	// CreateGroupSession 创建群聊会话
	//
	// 参数:
	//	- a: 应用
	// 	- user: 用户标识
	//	- groupID: 群组ID
	CreateGroupSession(ctx context.Context, a app.App, user app.User, groupID int64) (err error)

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
	// UpdateNormalSession 更新私聊会话
	//
	// 参数:
	//	- a: 应用
	// 	- user: 用户标识
	//	- peerID: 对方ID
	//	- session: 会话信息
	UpdateNormalSession(ctx context.Context, a app.App, user app.User, peerID int64, session *sessionpkg.Session) error
	// SaveNormalSession 保存私人会话
	//
	// 参数:
	//	- a: 应用
	// 	- user: 用户标识
	//	- peerID: 对方ID
	//	- messageTime: 消息时间
	SaveNormalSession(ctx context.Context, a app.App, user app.User, peerID int64, messageTime int64, isReceiver bool) error

	// NotifyNormalSession 通知 私聊会话
	//
	// 参数:
	//	- a: 应用
	// 	- notifyUser: 通知用户
	//	- session: 会话信息
	NotifyNormalSession(ctx context.Context, a app.App, notifyUser int64, session *sessionpkg.Session) error

	// NotifyRemoveSession 通知 删除会话
	//
	// 参数:
	//	- a: 应用
	// 	- notifyUser: 通知用户
	//	- session: 会话信息
	NotifyRemoveSession(ctx context.Context, a app.App, notifyUser int64, session *sessionpkg.Session) error

	// GroupSessionID 私人会话ID
	//
	// 参数:
	//	- userID: 用户ID
	NormalSessionID(userID int64) string
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
func (s *session) CreateNormalSession(ctx context.Context, a app.App, user app.User, peerID int64, messageTime int64, isReceiver bool) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	sessionID := s.NormalSessionID(peerID)
	sessionInfo, err := sessionpkg.AddP2PSession(ctx, a.ID(), user.UserID(), sessionID, peerID, messageTime, isReceiver)
	if err != nil {
		logger.CtxError(ctx, "AddP2PSession error", zap.Error(err))
		return err
	}
	//新建私聊会话 通知双方
	recent, err := sessionpkg.GetMessageRecent(ctx, a.ID(), user.UserID(), peerID)
	if err != nil {
		logger.CtxError(ctx, "GetMessageRecent error", zap.Error(err))
		return err
	}
	sessionInfo.Recent = []msgstore.Message{recent}
	//通知对方
	err = s.NotifyNormalSession(ctx, a, int64(peerID), sessionInfo)
	if err != nil {
		logger.CtxError(ctx, "NotifyNormalSession error", zap.Error(err))
		return err
	}
	//通知自己
	err = s.NotifyNormalSession(ctx, a, int64(user.UserID()), sessionInfo)
	if err != nil {
		logger.CtxError(ctx, "NotifyNormalSession error", zap.Error(err))
		return err
	}
	return nil
}

func (s *session) UpdateNormalSession(ctx context.Context, a app.App, user app.User, peerID int64, session *sessionpkg.Session) error {
	return sessionpkg.UpdateSession(ctx, a.ID(), user.UserID(), s.NormalSessionID(peerID), session)
}

// SaveNormalSession 保存私人会话
func (s *session) SaveNormalSession(ctx context.Context, a app.App, user app.User, peerID int64, messageTime int64, isReceiver bool) error {
	logger := fklog.ContextAppLogger(ctx)
	sessionID := s.NormalSessionID(peerID)
	//是否存在当前聊天对象perrID的session记录
	session, err := sessionpkg.GetNormalSession(ctx, a.ID(), user.UserID(), sessionID)
	if err != nil {
		logger.CtxError(ctx, "GetNormalSession error", zap.Error(err))
		return err
	}
	//不存在则创建
	if session == nil {
		err = s.CreateNormalSession(ctx, a, user, peerID, messageTime, isReceiver)
		if err != nil {
			logger.CtxError(ctx, "CreateNormalSession error", zap.Error(err))
			return err
		}
	} else {
		//存在则更新
		if isReceiver {
			session.UnreadCount += 1
		}
		session.MessageTime = messageTime
		err = s.UpdateNormalSession(ctx, a, user, peerID, session)
		if err != nil {
			logger.CtxError(ctx, "UpdateNormalSession error", zap.Error(err))
			return err
		}
	}
	return nil
}

// CreateGroupSession implements SessionService.
func (s *session) CreateGroupSession(ctx context.Context, a app.App, user app.User, groupID int64) (err error) {
	sessionID := s.GroupSessionID(groupID)
	return sessionpkg.AddGroupSession(ctx, a.ID(), user.UserID(), sessionID, groupID)
}

// RemoveSession implements SessionService.
func (s *session) RemoveSession(ctx context.Context, a app.App, user app.User, sessionID string) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	//先取出来，等会通知自己
	session, err := sessionpkg.GetNormalSession(ctx, a.ID(), user.UserID(), sessionID)
	if err != nil {
		logger.CtxError(ctx, "GetSession error", zap.Error(err))
		return err
	}
	err = sessionpkg.RemoveSession(ctx, a.ID(), user.UserID(), sessionID)
	if err != nil {
		logger.CtxError(ctx, "RemoveSession error", zap.Error(err))
		return err
	}
	//给自己发送删除通知
	err = s.NotifyRemoveSession(ctx, a, int64(user.UserID()), session)
	if err != nil {
		logger.CtxError(ctx, "NotifyRemoveSession error", zap.Error(err))
		return err
	}
	return nil
}

func (s *session) NormalSessionID(userID int64) string {
	return sum([]byte(fmt.Sprintf("p2p:%d", userID)))
}

func (s *session) GroupSessionID(groupID int64) string {
	return sum([]byte(fmt.Sprintf("group:%d", groupID)))
}

// sha1 returns the SHA-1 checksum of the data.
func sum(data []byte) string {
	sum := sha1.Sum(data)
	return hex.EncodeToString(sum[:])
}

// 获取message信息
func (s *session) GetMessageInfo(ctx context.Context, a app.App, user app.User, sessions map[string]app.Session) (messageInfo []*MazeIM.Session, err error) {
	if len(sessions) == 0 {
		return nil, nil
	}
	for _, session := range sessions {
		peerID := session.PeerInfo.GetUserId()
		if peerID > 0 {
			p2pmsg, err := p2pmsg.QueryMessages(ctx, a.ID(), int64(user.UserID()), peerID, uint64(0), true)
			if err != nil {
				return nil, err
			}
			if len(p2pmsg) > 0 {
				messageInfo = append(messageInfo, &MazeIM.Session{
					SessionId:  proto.String(session.ID),
					CreateTime: proto.Int64(session.CreateTime),
					Recent:     PbSessionMessage(p2pmsg),
				})
			}
		} else {
			groupID := session.GroupID
			if groupID > 0 {
				groupmsg, err := p2pmsg.QueryMessages(ctx, a.ID(), int64(user.UserID()), 0, uint64(0), true)
				if err != nil {
					return nil, err
				}
				if len(groupmsg) > 0 {
					messageInfo = append(messageInfo, &MazeIM.Session{
						SessionId:  proto.String(session.ID),
						CreateTime: proto.Int64(session.CreateTime),
						Recent:     PbSessionMessage(groupmsg),
					})
				}
			}
		}
	}
	return messageInfo, nil
}

// 通知 私聊会话
func (s *session) NotifyNormalSession(ctx context.Context, a app.App, notifyUser int64, session *sessionpkg.Session) error {
	logger := fklog.ContextAppLogger(ctx)
	// 推送消息给集群
	notifyMessage := &MazeIM.SessionChangeID{
		AddSessionList: pbSession(session),
	}
	logger.CtxInfo(ctx, "NotifyNormalSession start", zap.Int64("peerId", notifyUser), zap.Any("NotifyNormalSession", notifyMessage))
	err := online.ClusterPush(ctx, uint64(notifyUser), SessionChangeID, notifyMessage)
	if err != nil {
		logger.CtxError(ctx, "NotifyNormalSession error", zap.Error(err), zap.Any("NotifyNormalSession", notifyMessage))
	}
	return err
}

// 通知 删除会话
func (s *session) NotifyRemoveSession(ctx context.Context, a app.App, notifyUser int64, session *sessionpkg.Session) error {
	logger := fklog.ContextAppLogger(ctx)
	// 推送消息给集群
	notifyMessage := &MazeIM.SessionChangeID{
		DelSessionList: pbSession(session),
	}
	logger.CtxInfo(ctx, "NotifyRemoveSession start", zap.Int64("peerId", notifyUser), zap.Any("NotifyRemoveSession", notifyMessage))
	err := online.ClusterPush(ctx, uint64(notifyUser), SessionChangeID, notifyMessage)
	if err != nil {
		logger.CtxError(ctx, "NotifyRemoveSession error", zap.Error(err), zap.Any("NotifyRemoveSession", notifyMessage))
	}
	return err
}

// PbSessionMessage 转换为pb的message
func PbSessionMessage(messages []p2pmsg.Message) []*MazeIM.Message {
	pbMessages := make([]*MazeIM.Message, 0, len(messages))
	for _, message := range messages {
		pbMessages = append(pbMessages, &MazeIM.Message{
			MsgId:   proto.Uint64(message.MessageID),
			Type:    proto.Int32(message.Type),
			Content: []byte(message.Content),
			Sender:  proto.Int64(message.UserID),
		})
	}
	return pbMessages
}

// pbSession 转换为pb的session
func pbSession(session *sessionpkg.Session) []*MazeIM.Session {
	pbSessions := make([]*MazeIM.Session, 0, 1)
	pbSessions = append(pbSessions, &MazeIM.Session{
		SessionId:  proto.String(session.ID),
		CreateTime: proto.Int64(session.CreateTime),
		Recent:     PbSessionMessage(session.Recent),
	})
	return pbSessions
}
