package groupmsg

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

// getKey 获取缓存操作key。
func getKey(db nanoredis.NanoRedisClient, args ...interface{}) string {
	return db.MakeSectionKey(fmt.Sprintf("im:app:%d:group:%d:history", args...))
}

// QueryMessages 分页查询会话中的历史消息
func QueryMessages(ctx context.Context, appID int32, groupID int64, lastID uint64, limit int) (messages []Message, err error) {
	logger := fklog.ContextAppLogger(ctx)
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "QueryMessages Client fail",
			zap.Error(err),
		)
		return nil, err
	}
	key := getKey(cli, appID, groupID)
	ret, err := cli.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{Max: strconv.FormatUint(exchangeTextId(lastID), 10), Count: 20}).Result()
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
func SaveMessage(ctx context.Context, appID int32, groupID int64, message Message) (err error) {
	logger := fklog.ContextAppLogger(ctx)

	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "SaveMessage Client fail",
			zap.Error(err),
		)
		return err
	}
	key := getKey(cli, appID, groupID)
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

func exchangeTextId(textId uint64) uint64 {
	t1 := (textId >> 12) & 0xffffffff00000
	t2 := (textId & 0xfffff)
	return t1 | t2
}
