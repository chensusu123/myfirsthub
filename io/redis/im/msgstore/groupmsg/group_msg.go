package groupmsg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	globalredis "maze_game_server/io/redis"
	"maze_game_server/io/redis/im/msgstore"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
)

type Message = msgstore.Message

// getKey 获取缓存操作key。
func getKey(args ...interface{}) string {
	return fmt.Sprintf("im:app:%d:group:%d:history", args...)
}

// QueryMessages 分页查询会话中的历史消息
func QueryMessages(logger fklog.FKLogI, appID int32, groupID int32, lastID uint64, limit int) (messages []Message, err error) {
	var (
		key = getKey(appID, groupID)
		ctx = context.Background()
	)
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("QueryMessages Client fail",
			zap.Error(err),
			zap.Any("key", key),
		)
		return nil, err
	}
	ret, err := cli.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{Max: strconv.FormatUint(lastID, 10), Count: 20}).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		} else {
			logger.ErrorWF("QueryMessages ZRevRangeByScore fail",
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
			logger.ErrorWF("QueryMessages Unmarshal fail", zap.Error(err), zap.String("value", value))
			return nil, err
		}
		messages = append(messages, message)
	}
	logger.DebugWF("QueryMessages success", zap.Any("key", key), zap.Any("messages", messages))
	return
}

// SaveMessage 在会话保存历史消息
func SaveMessage(logger fklog.FKLogI, appID int32, groupID int32, message Message) (err error) {
	var (
		key = getKey(appID, groupID)
		ctx = context.Background()
	)
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("SaveMessage Client fail",
			zap.Error(err),
			zap.Any("key", key),
		)
		return err
	}
	data, err := json.Marshal(message)
	if err != nil {
		logger.ErrorWF("SaveMessage Marshal fail",
			zap.Error(err),
			zap.Any("key", key),
		)
		return err
	}

	err = cli.ZAdd(ctx, key, redis.Z{Score: float64(message.MessageID), Member: data}).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		} else {
			logger.ErrorWF("SaveMessage ZAdd fail",
				zap.Error(err),
				zap.Any("key", key),
			)
			return err
		}
	}
	logger.DebugWF("SaveMessage success", zap.Any("key", key), zap.Any("message", message))
	return
}
