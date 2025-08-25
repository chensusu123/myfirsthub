package sessionservice

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"maze_game_server/app"
	sessionpkg "maze_game_server/io/redis/im/session"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

type SessionService interface {
	// QueryRecentSessions 查询用户最近的会话记录
	//
	// 参数:
	//	- a: 应用
	// 	- user: 用户标识
	QueryRecentSessions(ctx context.Context, logger fklog.FKLogI, a app.App, user app.User) (sessions map[string]app.Session, err error)

	// CreateNormalSession 创建普通私聊会话
	//
	// 参数:
	//	- a: 应用
	// 	- user: 用户标识
	//	- peerID: 对方ID
	CreateNormalSession(ctx context.Context, logger fklog.FKLogI, a app.App, user app.User, peerID uint64) (err error)

	// CreateGroupSession 创建群聊会话
	//
	// 参数:
	//	- a: 应用
	// 	- user: 用户标识
	//	- groupID: 群组ID
	CreateGroupSession(ctx context.Context, logger fklog.FKLogI, a app.App, user app.User, groupID int32) (err error)

	// RemoveSession 删除会话记录
	//
	// 参数:
	//	- a: 应用
	// 	- user: 用户标识
	//	- sessionID: 会话ID
	RemoveSession(ctx context.Context, logger fklog.FKLogI, a app.App, user app.User, sessionID string) (err error)
}

var (
	Default SessionService = &session{}
)

type session struct {
}

// QueryRecentSessions implements SessionService.
func (s *session) QueryRecentSessions(ctx context.Context, logger fklog.FKLogI, a app.App, user app.User) (sessions map[string]app.Session, err error) {
	return sessionpkg.QuerySessions(logger, a.ID(), user.UserID())
}

// CreateNormalSession implements SessionService.
func (s *session) CreateNormalSession(ctx context.Context, logger fklog.FKLogI, a app.App, user app.User, peerID uint64) (err error) {
	sessionID := s.NormalSessionID(peerID)
	return sessionpkg.AddP2PSession(logger, a.ID(), user.UserID(), sessionID, peerID)
}

// CreateGroupSession implements SessionService.
func (s *session) CreateGroupSession(ctx context.Context, logger fklog.FKLogI, a app.App, user app.User, groupID int32) (err error) {
	sessionID := s.GroupSessionID(groupID)
	return sessionpkg.AddGroupSession(logger, a.ID(), user.UserID(), sessionID, groupID)
}

// RemoveSession implements SessionService.
func (s *session) RemoveSession(ctx context.Context, logger fklog.FKLogI, a app.App, user app.User, sessionID string) (err error) {
	return sessionpkg.RemoveSession(logger, a.ID(), user.UserID(), sessionID)
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
