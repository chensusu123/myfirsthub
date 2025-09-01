package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	globalredis "maze_game_server/io/redis"
	"time"

	"github.com/redis/go-redis/v9"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

type Session struct {
	ID         string `json:"id,omitempty"`
	CreateTime int64  `json:"create_time,omitempty"`
	PeerID     uint64 `json:"peer_id,omitempty"`
	GroupID    int64  `json:"group_id,omitempty"`
	OldestID   uint64 `json:"oldest_id,omitempty"`
}

// getKey 获取缓存操作key。
func getKey(args ...interface{}) string {
	return fmt.Sprintf("im:app:%d:u:%d:session:list", args...)
}

// QuerySessions
func QuerySessions(ctx context.Context, appID int32, userID uint64) (sessions map[string]Session, err error) {
	logger := fklog.ContextAppLogger(ctx)
	var (
		key = getKey(appID, userID)
	)
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "QuerySessions Client fail",
			zap.Error(err),
			zap.Any("key", key),
		)
		return nil, err
	}
	ret, err := cli.HGetAll(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		} else {
			logger.CtxError(ctx, "QuerySessions HGETALL fail",
				zap.Error(err),
				zap.Any("key", key),
			)
			return nil, err
		}
	}
	sessions = make(map[string]Session)
	for k, v := range ret {
		var session Session
		err = json.Unmarshal([]byte(v), &session)
		if err != nil {
			logger.CtxError(ctx, "QuerySessions Unmarshal fail",
				zap.Error(err),
				zap.Any("key", key),
			)
			return nil, err
		}
		sessions[k] = session
	}
	logger.CtxInfo(ctx, "QuerySessions success", zap.Any("key", key), zap.Any("sessions", sessions))
	return
}

// AddP2PSession 创建私聊会话
func AddP2PSession(ctx context.Context, appID int32, userID uint64, sessionID string, peerID uint64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	var (
		key = getKey(appID, userID)
	)
	session := Session{
		PeerID:     peerID,
		CreateTime: time.Now().Unix(),
	}
	value, err := json.Marshal(session)
	if err != nil {
		logger.CtxError(ctx, "AddP2PSession Marshal fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return err
	}
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "AddP2PSession Client fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return err
	}
	err = cli.HSet(ctx, key, sessionID, value).Err()
	if err != nil {
		logger.CtxError(ctx, "AddP2PSession HSET fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return err
	}
	logger.CtxInfo(ctx, "AddP2PSession success", zap.Any("key", key), zap.String("sessionID", sessionID), zap.Any("session", session))
	return
}

// AddGroupSession 创建群聊会话
func AddGroupSession(ctx context.Context, appID int32, userID uint64, sessionID string, groupID int64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	var (
		key = getKey(appID, userID)
	)
	session := Session{
		GroupID:    groupID,
		CreateTime: time.Now().Unix(),
	}
	value, err := json.Marshal(session)
	if err != nil {
		logger.CtxError(ctx, "AddGroupSession Marshal fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return err
	}
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "AddGroupSession Client fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return err
	}
	err = cli.HSet(ctx, key, groupID, value).Err()
	if err != nil {
		logger.CtxError(ctx, "AddGroupSession HSET fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return err
	}
	logger.CtxInfo(ctx, "AddGroupSession success", zap.Any("key", key), zap.String("sessionID", sessionID), zap.Any("session", session))
	return
}

// RemoveSession 移除私聊会话
func RemoveSession(ctx context.Context, appID int32, userID uint64, sessionID string) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	var (
		key = getKey(appID, userID)
	)
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "RemoveSession Client fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
		)
		return err
	}
	err = cli.HDel(ctx, key, sessionID).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		} else {
			logger.CtxError(ctx, "RemoveSession HDEL fail",
				zap.Error(err),
				zap.Any("key", key),
				zap.String("sessionID", sessionID),
			)
			return err
		}
	}
	logger.CtxInfo(ctx, "RemoveSession success", zap.Any("key", key), zap.String("sessionID", sessionID))
	return
}
