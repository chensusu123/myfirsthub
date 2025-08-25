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
	GroupID    int32  `json:"group_id,omitempty"`
	OldestID   uint64 `json:"oldest_id,omitempty"`
}

// getKey 获取缓存操作key。
func getKey(args ...interface{}) string {
	return fmt.Sprintf("im:app:%d:u:%d:session:list", args...)
}

// QuerySessions
func QuerySessions(logger fklog.FKLogI, appID int32, userID uint64) (sessions map[string]Session, err error) {
	var (
		key = getKey(appID, userID)
		ctx = context.Background()
	)
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("QuerySessions Client fail",
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
			logger.ErrorWF("QuerySessions HGETALL fail",
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
			logger.ErrorWF("QuerySessions Unmarshal fail",
				zap.Error(err),
				zap.Any("key", key),
			)
			return nil, err
		}
		sessions[k] = session
	}
	logger.DebugWF("QuerySessions success", zap.Any("key", key), zap.Any("sessions", sessions))
	return
}

// AddP2PSession 创建私聊会话
func AddP2PSession(logger fklog.FKLogI, appID int32, userID uint64, sessionID string, peerID uint64) (err error) {
	var (
		key = getKey(appID, userID)
		ctx = context.Background()
	)
	session := Session{
		PeerID:     peerID,
		CreateTime: time.Now().Unix(),
	}
	value, err := json.Marshal(session)
	if err != nil {
		logger.ErrorWF("AddP2PSession Marshal fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return err
	}
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("AddP2PSession Client fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return err
	}
	err = cli.HSet(ctx, key, sessionID, value).Err()
	if err != nil {
		logger.ErrorWF("AddP2PSession HSET fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return err
	}
	logger.DebugWF("AddP2PSession success", zap.Any("key", key), zap.String("sessionID", sessionID), zap.Any("session", session))
	return
}

// AddGroupSession 创建群聊会话
func AddGroupSession(logger fklog.FKLogI, appID int32, userID uint64, sessionID string, groupID int32) (err error) {
	var (
		key = getKey(appID, userID)
		ctx = context.Background()
	)
	session := Session{
		GroupID:    groupID,
		CreateTime: time.Now().Unix(),
	}
	value, err := json.Marshal(session)
	if err != nil {
		logger.ErrorWF("AddGroupSession Marshal fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return err
	}
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("AddGroupSession Client fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return err
	}
	err = cli.HSet(ctx, key, groupID, value).Err()
	if err != nil {
		logger.ErrorWF("AddGroupSession HSET fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return err
	}
	logger.DebugWF("AddGroupSession success", zap.Any("key", key), zap.String("sessionID", sessionID), zap.Any("session", session))
	return
}

// RemoveSession 移除私聊会话
func RemoveSession(logger fklog.FKLogI, appID int32, userID uint64, sessionID string) (err error) {
	var (
		key = getKey(appID, userID)
		ctx = context.Background()
	)
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("RemoveSession Client fail",
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
			logger.ErrorWF("RemoveSession HDEL fail",
				zap.Error(err),
				zap.Any("key", key),
				zap.String("sessionID", sessionID),
			)
			return err
		}
	}
	logger.DebugWF("RemoveSession success", zap.Any("key", key), zap.String("sessionID", sessionID))
	return
}
