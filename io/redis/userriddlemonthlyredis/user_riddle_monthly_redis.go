package userriddlemonthlyredis

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
)

var gRedis = &fkredis.FkRedis{}

// u:%llu:riddle:monthly
func init() {
	fkconfig.RegisterNameNode("UserRiddleMonthlyRedis", 21681, gRedis)
}

func GetMazeCardExpirationTime(ctx context.Context, userId uint64) (int64, error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("u:%d:riddle:monthly", userId)
	expirationTime, err := redis.Int64(gRedis.Do(ctx, "get", key))
	if err != nil && err != redis.ErrNil {
		logger.CtxInfo(ctx, "GetMazeCardExpirationTime get failed", zap.String("key", key), zap.Error(err))
		return 0, err
	}

	logger.CtxInfo(ctx, "GetMazeCardExpirationTime end", zap.String("key", key), zap.Int64("time", expirationTime))
	return expirationTime, nil
}
