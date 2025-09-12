package mazebarrieropstatusredis

import (
	"context"
	"fmt"
	"maze_game_server/lib/serialize"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
)

type BarrierOpStatus struct {
	Status map[string]int64 `json:"status,omitempty"`
}

var (
	cli = &fkredis.FkRedis{}
)

func init() {
	fkconfig.RegisterNameNode("mazebarrieropstatusredis", 21728, cli)
}

// getKey 获取缓存操作key。
func getKey(args ...interface{}) string {
	return fmt.Sprintf("maze:u:%d:barrier:%d:op:status", args...)
}

// IsTriggered 用来校验关卡中某个操作是否已经进行过，幂等校验。
func IsTriggered(ctx context.Context, userID uint64, barrierID int32, op string) (triggered bool, triggerFn func() (err error), err error) {
	logger := fklog.ContextAppLogger(ctx)
	var (
		key = getKey(userID, barrierID)
	)
	result, err := redis.Bytes(cli.Do(ctx, "GET", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			logger.CtxError(ctx, "IsTriggered GET fail",
				zap.Error(err),
				zap.Any("key", key),
				zap.Uint64("userID", userID),
				zap.Int32("barrierID", barrierID),
				zap.String("op", op),
			)
			return false, nil, err
		}
	}
	var status BarrierOpStatus
	if len(result) <= 0 {
		status.Status = make(map[string]int64)
	} else {
		err = serialize.Unmarshal(result, &status)
		if err != nil {
			logger.CtxError(ctx, "IsTriggered Unmarshal fail",
				zap.Error(err),
				zap.Any("key", key),
				zap.Uint64("userID", userID),
				zap.Int32("barrierID", barrierID),
				zap.String("op", op),
				zap.ByteString("result", result),
			)
			return false, nil, err
		}
	}

	_, triggered = status.Status[op]

	if !triggered {
		triggerFn = func() (err error) {
			status.Status[op] = time.Now().UnixMilli()

			data, err := serialize.Marshal(status)
			if err != nil {
				logger.CtxError(ctx, "IsTriggered Marshal fail",
					zap.Error(err),
					zap.Any("key", key),
					zap.Uint64("userID", userID),
					zap.Int32("barrierID", barrierID),
					zap.String("op", op),
					zap.Any("status", status),
				)
				return err
			}

			_, err = redis.Bytes(cli.Do(ctx, "SET", key, data))
			if err != nil {
				logger.CtxError(ctx, "IsTriggered SET fail",
					zap.Error(err),
					zap.Any("key", key),
					zap.Uint64("userID", userID),
					zap.Int32("barrierID", barrierID),
					zap.String("op", op),
				)
				return err
			}

			logger.CtxDebug(ctx, "IsTriggered trigger success",
				zap.Uint64("userID", userID),
				zap.Int32("barrierID", barrierID),
				zap.String("op", op),
			)
			return nil
		}
	} else {
		triggerFn = func() (err error) { return nil }
	}

	logger.CtxDebug(ctx, "IsTriggered success",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.String("op", op),
		zap.Bool("triggered", triggered),
	)
	return triggered, triggerFn, nil
}

// ClearOpStatus 清理关卡操作状态
func ClearOpStatus(ctx context.Context, userID uint64, barrierID int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	var (
		key = getKey(userID, barrierID)
	)
	_, err = cli.Do(ctx, "DEL", key)
	if err != nil {
		logger.CtxError(ctx, "ClearOpStatus DEL fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
		)
		return err
	}
	logger.CtxDebug(ctx, "ClearOpStatus success", zap.Any("key", key), zap.Uint64("userID", userID), zap.Int32("barrierID", barrierID))
	return nil
}
