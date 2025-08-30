package mazechallengenumredis

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

func init() {
	fkconfig.RegisterNameNode("mazechallengenumredis", 21737, gRedis)
}

func GetUserChallengeNum(ctx context.Context, userId uint64, dateTime int) (num int, err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:u:%d:date:%d", userId, dateTime)
	num, err = redis.Int(gRedis.Do(context.TODO(), "get", key))
	if err == redis.ErrNil {
		err = nil
		logger.CtxInfo(ctx, "GetUserChallengeNum nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.CtxError(ctx, "GetUserChallengeNum fail", zap.String("key", key), zap.Error(err))
		return
	}

	return
}

func AddUserChallengeNum(ctx context.Context, userId uint64, dateTime int, count int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:u:%d:date:%d", userId, dateTime)
	newCount, err := redis.Int(gRedis.Do(context.TODO(), "incrby", key, count))
	if err != nil {
		logger.CtxError(ctx, "AddUserChallengeNum incr fail", zap.String("key", key), zap.Error(err))
		return
	}
	logger.CtxInfo(ctx, "AddUserChallengeNum succ", zap.Any("key", key), zap.Any("count", count), zap.Any("newCount", newCount))
	return
}

func GMDel(ctx context.Context, userId uint64, dateTime int) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:u:%d:date:%d", userId, dateTime)
	_, err = redis.Int(gRedis.Do(context.TODO(), "del", key))
	if err != nil {
		logger.CtxError(ctx, "GMDel del fail", zap.String("key", key), zap.Error(err))
		return
	}
	logger.CtxInfo(ctx, "GMDel succ", zap.Any("key", key))
	return
}
