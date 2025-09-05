package mazefixedbarrierredis

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
)

// maze:u:%d:fixed:barrier
var gRedis = &fkredis.FkRedis{}

func init() {
	fkconfig.RegisterNameNode("mazefixedbarrierredis", 21689, gRedis)
}

func GetUserFixedBarrierID(ctx context.Context, userId uint64) (barrierId int32, err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:u:%d:fixed:barrier", userId)

	res, err := redis.Int(gRedis.Do(ctx, "get", key))
	if err == redis.ErrNil {
		err = nil
		logger.CtxInfo(ctx, "GetUserFixedBarrierID nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.CtxError(ctx, "GetUserFixedBarrierID get fail", zap.String("key", key), zap.Error(err))
		return
	}

	barrierId = int32(res)

	logger.CtxInfo(ctx, "GetUserFixedBarrierID succ", zap.String("key", key), zap.Any("barrierId", barrierId))
	return
}

func SetUserFixedBarrierID(ctx context.Context, userId uint64, barrierId int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:u:%d:fixed:barrier", userId)
	_, err = gRedis.Do(ctx, "set", key, barrierId)
	if err != nil {
		logger.CtxError(ctx, "SetUserFixedBarrierID set fail", zap.String("key", key), zap.Error(err))
		return
	}
	logger.CtxInfo(ctx, "SetUserFixedBarrierID succ", zap.String("key", key), zap.Any("barrierId", barrierId))
	return
}

func DelUserFixedBarrierID(ctx context.Context, userId uint64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:u:%d:fixed:barrier", userId)
	_, err = gRedis.Do(ctx, "del", key)
	if err != nil {
		logger.CtxError(ctx, "DelUserFixedBarrierID set fail", zap.String("key", key), zap.Error(err))
		return
	}
	logger.CtxInfo(ctx, "DelUserFixedBarrierID succ", zap.String("key", key))
	return
}
