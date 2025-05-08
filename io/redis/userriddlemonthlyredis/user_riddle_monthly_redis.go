package userriddlemonthlyredis

import (
	"context"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
)

var gRedis = &fkredis.FkRedis{}

// u:%llu:riddle:monthly
func init() {
	fkconfig.RegisterNameNode("UserRiddleMonthlyRedis", 21681, gRedis)
}

func GetMazeCardExpirationTime(logger fklog.FKLogI, userId uint64) (int64, error) {
	key := gRedis.GetKey(userId)
	expirationTime, err := redis.Int64(gRedis.Do(context.TODO(), "get", key))
	if err != nil && err != redis.ErrNil {
		logger.InfoWF("GetMazeCardExpirationTime get failed", zap.String("key", key), zap.Error(err))
		return 0, err
	}

	logger.InfoWF("GetMazeCardExpirationTime end", zap.String("key", key), zap.Int64("time", expirationTime))
	return expirationTime, nil
}
