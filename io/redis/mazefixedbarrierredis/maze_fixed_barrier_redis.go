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
	fkconfig.RegisterNameNode("mazefixedbarrierredis", 0, gRedis)
}

func GetUserFixedBarrierID(logger fklog.FKLogI, userId uint64) (barrierId int32, err error) {

	key := fmt.Sprintf("maze:u:%d:fixed:barrier", userId)

	res, err := redis.Int(gRedis.Do(context.TODO(), "get", key))
	if err == redis.ErrNil {
		err = nil
		logger.InfoWF("GetUserFixedBarrierID nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.ErrorWF("GetUserFixedBarrierID get fail", zap.String("key", key), zap.Error(err))
		return
	}

	barrierId = int32(res)

	logger.InfoWF("GetUserFixedBarrierID succ", zap.String("key", key), zap.Any("barrierId", barrierId))
	return
}

func SetUserFixedBarrierID(logger fklog.FKLogI, userId uint64, barrierId int32) (err error) {
	key := fmt.Sprintf("maze:u:%d:fixed:barrier", userId)
	_, err = gRedis.Do(context.TODO(), "set", key, barrierId)
	if err != nil {
		logger.ErrorWF("SetUserFixedBarrierID set fail", zap.String("key", key), zap.Error(err))
		return
	}
	logger.InfoWF("SetUserFixedBarrierID succ", zap.String("key", key), zap.Any("barrierId", barrierId))
	return
}

func DelUserFixedBarrierID(logger fklog.FKLogI, userId uint64) (err error) {
	key := fmt.Sprintf("maze:u:%d:fixed:barrier", userId)
	_, err = gRedis.Do(context.TODO(), "del", key)
	if err != nil {
		logger.ErrorWF("DelUserFixedBarrierID set fail", zap.String("key", key), zap.Error(err))
		return
	}
	logger.InfoWF("DelUserFixedBarrierID succ", zap.String("key", key))
	return
}
