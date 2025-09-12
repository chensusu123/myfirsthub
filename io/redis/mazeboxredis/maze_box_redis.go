package mazeboxredis

import (
	"context"
	"fmt"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
)

var (
	cli = &fkredis.FkRedis{}
)

func init() {
	fkconfig.RegisterNameNode("mazeboxredis", 21642, cli)
}

// getKey 获取缓存操作key。
func getKey(args ...interface{}) string {
	return fmt.Sprintf("u:%d:barrier:%d:opened:box", args...)
}

// IsOpenedBox
func IsOpenedBox(ctx context.Context, userID uint64, barrierID, boxID int32) (openTime int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	var (
		key = getKey(userID, barrierID)
	)
	openTime, err = redis.Int64(cli.Do(ctx, "ZSCORE", key, boxID))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			logger.CtxError(ctx, "IsOpenedBox ZSCORE fail",
				zap.Error(err),
				zap.Any("key", key),
			)
			return 0, err
		}
	}
	logger.CtxDebug(ctx, "IsOpenedBox success", zap.Uint64("userID", userID), zap.Int32("barrierID", barrierID), zap.Int32("boxID", boxID), zap.Int64("openTime", openTime))
	return
}

// SetOpenBoxTime
func SetOpenBoxTime(ctx context.Context, userID uint64, barrierID, boxID int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	var (
		key = getKey(userID, barrierID)
	)
	_, err = cli.Do(ctx, "ZADD", key, time.Now().Unix(), boxID)
	if err != nil {
		logger.CtxError(ctx, "SetOpenBoxTime ZADD fail",
			zap.Error(err),
			zap.Any("key", key),
		)
		return err
	}
	logger.CtxDebug(ctx, "SetOpenBoxTime success", zap.Uint64("userID", userID), zap.Int32("barrierID", barrierID), zap.Int32("boxID", boxID))
	return
}

// ClearOpenBoxTime
func ClearOpenBoxTime(ctx context.Context, userID uint64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	var (
		keys = make([]interface{}, 0)
	)
	// 查询所有关卡ID
	for _, row := range GMazeBarriesV8Cfg.GetAll() {
		keys = append(keys, getKey(userID, row.Order))
	}
	_, err = cli.Do(ctx, "DEL", keys...)
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			logger.CtxError(ctx, "ClearOpenBoxTime DEL fail",
				zap.Error(err),
				zap.Any("keys", keys),
			)
			return err
		}
	}
	logger.CtxDebug(ctx, "ClearOpenBoxTime success", zap.Uint64("userID", userID), zap.Any("keys", keys))
	return
}
