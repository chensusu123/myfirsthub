package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	globalredis "maze_game_server/io/redis"
	"maze_game_server/io/redis/im/msgstore"
	"maze_game_server/io/redis/im/msgstore/p2pmsg"
	"time"

	"github.com/redis/go-redis/v9"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/database/nanoredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

type Session struct {
	ID          string             `json:"id,omitempty"`
	CreateTime  int64              `json:"create_time,omitempty"`
	PeerID      int64              `json:"peer_id,omitempty"`
	GroupID     int64              `json:"group_id,omitempty"`
	MessageTime int64              `json:"message_time,omitempty"`
	UnreadCount int64              `json:"unread_count,omitempty"`
	Recent      []msgstore.Message `json:"recent,omitempty"`
}

// getKey 获取缓存操作key。
func getKey(db nanoredis.NanoRedisClient, args ...interface{}) string {
	return db.MakeSectionKey(fmt.Sprintf("im:app:%d:u:%d:session:list", args...))
}

// QuerySessions
func QuerySessions(ctx context.Context, appID int32, userID uint64) (sessions map[string]Session, err error) {
	logger := fklog.ContextAppLogger(ctx)

	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "QuerySessions Client fail",
			zap.Error(err),
		)
		return nil, err
	}
	key := getKey(cli, appID, userID)
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
func AddP2PSession(ctx context.Context, appID int32, userID uint64, sessionID string, peerID int64, messageTime int64, isReceiver bool) (session *Session, err error) {
	logger := fklog.ContextAppLogger(ctx)

	initUnreadCount := int64(0)
	if isReceiver {
		initUnreadCount = 1
	}
	session = &Session{
		PeerID:      peerID,
		CreateTime:  time.Now().Unix(),
		MessageTime: messageTime,
		UnreadCount: initUnreadCount,
	}
	value, err := json.Marshal(session)
	if err != nil {
		logger.CtxError(ctx, "AddP2PSession Marshal fail",
			zap.Error(err),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return nil, err
	}
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "AddP2PSession Client fail",
			zap.Error(err),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return nil, err
	}
	key := getKey(cli, appID, userID)
	err = cli.HSet(ctx, key, sessionID, value).Err()
	if err != nil {
		logger.CtxError(ctx, "AddP2PSession HSET fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return nil, err
	}
	logger.CtxInfo(ctx, "AddP2PSession success", zap.Any("key", key), zap.String("sessionID", sessionID), zap.Any("session", session))
	return
}

// UpdateSession 更新私聊会话
func UpdateSession(ctx context.Context, appID int32, userID uint64, sessionID string, session *Session) (err error) {
	logger := fklog.ContextAppLogger(ctx)

	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "UpdateSession Client fail",
			zap.Error(err),
			zap.String("sessionID", sessionID),
		)
		return err
	}

	value, err := json.Marshal(session)
	if err != nil {
		logger.CtxError(ctx, "AddP2PSession Marshal fail",
			zap.Error(err),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return err
	}
	key := getKey(cli, appID, userID)
	err = cli.HSet(ctx, key, sessionID, value).Err()
	if err != nil {
		logger.CtxError(ctx, "UpdateSession HSET fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return err
	}
	logger.CtxInfo(ctx, "UpdateSession success", zap.Any("key", key), zap.String("sessionID", sessionID), zap.Any("session", session))
	return
}

// AddGroupSession 创建群聊会话
func AddGroupSession(ctx context.Context, appID int32, userID uint64, sessionID string, groupID int64) (err error) {
	logger := fklog.ContextAppLogger(ctx)

	session := Session{
		GroupID:    groupID,
		CreateTime: time.Now().Unix(),
	}
	value, err := json.Marshal(session)
	if err != nil {
		logger.CtxError(ctx, "AddGroupSession Marshal fail",
			zap.Error(err),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return err
	}
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "AddGroupSession Client fail",
			zap.Error(err),
			zap.String("sessionID", sessionID),
			zap.Any("session", session),
		)
		return err
	}
	key := getKey(cli, appID, userID)
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

	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "RemoveSession Client fail",
			zap.Error(err),
			zap.String("sessionID", sessionID),
		)
		return err
	}
	key := getKey(cli, appID, userID)
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

func GetNormalSession(ctx context.Context, appID int32, userID uint64, sessionID string) (session *Session, err error) {
	logger := fklog.ContextAppLogger(ctx)

	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "GetNormalSession Client fail",
			zap.Error(err),
			zap.String("sessionID", sessionID),
		)
		return session, err
	}
	key := getKey(cli, appID, userID)
	value, err := cli.HGet(ctx, key, sessionID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
			return session, nil
		} else {
			logger.CtxError(ctx, "GetNormalSession HGET fail",
				zap.Error(err),
				zap.Any("key", key),
				zap.String("sessionID", sessionID),
			)
			return session, err
		}
	}
	err = json.Unmarshal([]byte(value), &session)
	if err != nil {
		logger.CtxError(ctx, "GetNormalSession Unmarshal fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.String("value", value),
		)
		return session, err
	}
	logger.CtxInfo(ctx, "GetNormalSession success", zap.Any("key", key), zap.String("sessionID", sessionID), zap.Any("session", session))
	return session, nil
}

// GetMessageRecent 获取最近一条消息
func GetMessageRecent(ctx context.Context, appID int32, userID uint64, peerID int64) (message msgstore.Message, err error) {
	logger := fklog.ContextAppLogger(ctx)

	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "GetMessageRecent Client fail",
			zap.Error(err),
			zap.Int64("peerID", peerID),
		)
		return message, err
	}
	key := p2pmsg.GetKey(cli, appID, userID, peerID)
	//获取最新1条消息
	value, err := cli.ZRevRangeWithScores(ctx, key, 0, 0).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
			return message, nil
		} else {
			logger.CtxError(ctx, "GetMessageRecent LRANGE fail",
				zap.Error(err),
				zap.Any("key", key),
				zap.Int64("peerID", peerID),
			)
			return message, err
		}
	}
	// err = json.Unmarshal([]byte(value[0].Member), &message)
	err = json.Unmarshal([]byte(value[0].Member.(string)), &message)
	if err != nil {
		logger.CtxError(ctx, "GetMessageRecent Unmarshal fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.Int64("peerID", peerID),
			zap.Any("message", message),
		)
		return message, err
	}
	logger.CtxInfo(ctx, "GetMessageRecent success", zap.Any("key", key), zap.Int64("peerID", peerID), zap.Any("messages", message))
	return message, nil
}

// 设置未读数为0
func SetRemoveUnreadCount(ctx context.Context, appID int32, userID int64, sessionID string) (err error) {
	logger := fklog.ContextAppLogger(ctx)

	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "SetRemoveUnreadCount Client fail",
			zap.Error(err),
		)
		return err
	}
	key := getKey(cli, appID, userID)
	value, err := cli.HGet(ctx, key, sessionID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		} else {
			logger.CtxError(ctx, "SetRemoveUnreadCount HGET fail",
				zap.Error(err),
				zap.String("sessionID", sessionID),
			)
			return err
		}
	}

	var session Session
	err = json.Unmarshal([]byte(value), &session)
	if err != nil {
		logger.CtxError(ctx, "SetRemoveUnreadCount Unmarshal fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.String("value", value),
		)
		return err
	}
	session.UnreadCount = 0

	jsonValue, err := json.Marshal(session)
	if err != nil {
		logger.CtxError(ctx, "SetRemoveUnreadCount Marshal fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.String("value", value),
		)
		return err
	}
	err = cli.HSet(ctx, key, sessionID, jsonValue).Err()
	if err != nil {
		logger.CtxError(ctx, "SetRemoveUnreadCount HSET fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.String("sessionID", sessionID),
			zap.String("value", value),
		)
		return err
	}

	logger.CtxInfo(ctx, "SetRemoveUnreadCount success", zap.Any("key", key), zap.Any("jsonValue", jsonValue))
	return nil
}
