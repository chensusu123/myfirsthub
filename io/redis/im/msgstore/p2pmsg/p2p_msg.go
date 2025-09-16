package p2pmsg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	globalredis "maze_game_server/io/redis"
	"maze_game_server/io/redis/im/msgstore"
	"strconv"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/database/nanoredis"

	"github.com/redis/go-redis/v9"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

type Message = msgstore.Message

const (
	limitMessageCount = 20
)

// getKey 获取缓存操作key。
func GetKey(db nanoredis.NanoRedisClient, args ...interface{}) string {
	return db.MakeSectionKey(fmt.Sprintf("im:app:%d:p2p:%d:%d:history", args...))
}

// QueryMessages 分页查询会话中的历史消息
func QueryMessages(ctx context.Context, appID int32, userID, peerID int64, lastID uint64, newest bool) (messages []Message, err error) {
	logger := fklog.ContextAppLogger(ctx)

	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "QueryMessages Client fail",
			zap.Error(err),
		)
		return nil, err
	}
	key := GetKey(cli, appID, userID, peerID)
	var ret []string
	if lastID == 0 {
		ret, err = cli.ZRange(ctx, key, 0, limitMessageCount-1).Result()
	} else {
		if newest {
			ret, err = cli.ZRangeByScore(ctx, key, &redis.ZRangeBy{Min: strconv.FormatUint(exchangeTextId(lastID), 10), Count: limitMessageCount}).Result()
		} else {
			ret, err = cli.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{Max: strconv.FormatUint(exchangeTextId(lastID), 10), Count: limitMessageCount}).Result()
		}
	}

	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		} else {
			logger.CtxError(ctx, "QueryMessages ZRevRangeByScore fail",
				zap.Error(err),
				zap.Any("key", key),
			)
			return nil, err
		}
	}
	for _, value := range ret {
		var message Message
		err = json.Unmarshal([]byte(value), &message)
		if err != nil {
			logger.CtxError(ctx, "QueryMessages Unmarshal fail", zap.Error(err), zap.String("value", value))
			return nil, err
		}
		messages = append(messages, message)
	}
	logger.CtxInfo(ctx, "QueryMessages success", zap.Any("key", key), zap.Any("messages", messages))
	return
}

// SaveMessage 在会话保存历史消息
func SaveMessage(ctx context.Context, appID int32, userID, peerID int64, message Message) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "SaveMessage Client fail",
			zap.Error(err),
		)
		return err
	}
	key := GetKey(cli, appID, userID, peerID)
	data, err := json.Marshal(message)
	if err != nil {
		logger.CtxError(ctx, "SaveMessage Marshal fail",
			zap.Error(err),
			zap.Any("key", key),
		)
		return err
	}

	err = cli.ZAdd(ctx, key, redis.Z{Score: float64(exchangeTextId(message.MessageID)), Member: data}).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		} else {
			logger.CtxError(ctx, "SaveMessage ZAdd fail",
				zap.Error(err),
				zap.Any("key", key),
			)
			return err
		}
	}
	logger.CtxInfo(ctx, "SaveMessage success", zap.Any("key", key), zap.Any("message", message))
	return
}

// ReadMessage 标记消息为已读
func ReadMessage(ctx context.Context, appID int32, userID, peerID int64, messageID uint64) (err error) {
	logger := fklog.ContextAppLogger(ctx)

	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "ReadMessage Client fail",
			zap.Error(err),
		)
		return err
	}
	key := GetKey(cli, appID, userID, peerID)
	//取出messageid对应的value
	member, err := cli.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: strconv.FormatUint(exchangeTextId(messageID), 10),
		Max: strconv.FormatUint(exchangeTextId(messageID), 10),
	}).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		} else {
			logger.CtxError(ctx, "ReadMessage getZValue fail",
				zap.Error(err),
				zap.Any("key", key),
				zap.Any("messageID", messageID),
			)
			return err
		}
	}
	if len(member) == 0 {
		return nil
	}
	message := Message{}
	// 反序列化消息
	err = json.Unmarshal([]byte(member[0]), &message)
	if err != nil {
		logger.CtxError(ctx, "ReadMessage Unmarshal fail", zap.Error(err), zap.String("value", member[0]))
		return err
	}
	if message.HasRead {
		// 已经标记为已读，不需要再次标记
		logger.CtxInfo(ctx, "ReadMessage already read", zap.Any("key", key), zap.Any("message", message))
		return nil
	}
	message.HasRead = true
	messageData, err := json.Marshal(message)
	if err != nil {
		logger.CtxError(ctx, "ReadMessage Marshal fail", zap.Error(err), zap.Any("message", message))
		return err
	}

	//TODO: 这里应该是先删除以前的消息，再添加新的消息状态,后续不用这种方式
	//先删除以前的消息
	err = cli.ZRem(ctx, key, member[0]).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		} else {
			logger.CtxError(ctx, "ReadMessage ZRem fail",
				zap.Error(err),
				zap.Any("key", key),
			)
			return err
		}
	}
	// 再添加新的消息状态
	err = cli.ZAdd(ctx, key, redis.Z{Score: float64(exchangeTextId(messageID)), Member: messageData}).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		} else {
			logger.CtxError(ctx, "ReadMessage ZAdd fail",
				zap.Error(err),
				zap.Any("key", key),
			)
			return err
		}
	}
	logger.CtxInfo(ctx, "ReadMessage success", zap.Any("key", key), zap.Any("message", message))
	return
}

// 删除消息
func RemoveMessage(ctx context.Context, appID int32, userID, peerID int64, messageID uint64) (err error) {
	logger := fklog.ContextAppLogger(ctx)

	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "RemoveMessage Client fail",
			zap.Error(err),
		)
		return err
	}
	key := GetKey(cli, appID, userID, peerID)
	//取出messageid对应的value
	member, err := cli.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: strconv.FormatUint(exchangeTextId(messageID), 10),
		Max: strconv.FormatUint(exchangeTextId(messageID), 10),
	}).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		} else {
			logger.CtxError(ctx, "RemoveMessage getZValue fail",
				zap.Error(err),
				zap.Any("key", key),
				zap.Any("messageID", messageID),
			)
			return err
		}
	}
	if len(member) == 0 {
		return nil
	}
	// 删除消息
	err = cli.ZRem(ctx, key, member[0]).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		} else {
			logger.CtxError(ctx, "RemoveMessage ZRem fail",
				zap.Error(err),
				zap.Any("key", key),
			)
			return err
		}
	}
	logger.CtxInfo(ctx, "RemoveMessage success", zap.Any("key", key), zap.Any("messageID", messageID))
	return
}

func exchangeTextId(textId uint64) uint64 {
	t1 := (textId >> 12) & 0xffffffff00000
	t2 := (textId & 0xfffff)
	return t1 | t2
}
