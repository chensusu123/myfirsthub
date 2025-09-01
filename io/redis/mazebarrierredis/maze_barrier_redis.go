package mazebarrierredis

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

var CURR_BARRIER_FIELD = "curr_barrier"
var CURR_BARRIER_AREA_FIELD = "curr_area"
var HIGH_BARRIER_AREA_FIELD = "high_area"

var gRedis = &fkredis.FkRedis{}

func init() {
	fkconfig.RegisterNameNode("mazebarrierredis", 21688, gRedis)
}

// 查询当前关卡id
func GetCurrBarrier(ctx context.Context, userId uint64) (barrierId int32, err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:u:%d:barrier:%d", userId)

	res, err := redis.Int(gRedis.Do(context.TODO(), "hget", key, CURR_BARRIER_FIELD))
	if err == redis.ErrNil {
		err = nil
		logger.CtxInfo(ctx, "GetCurrBarrier nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.CtxError(ctx, "GetCurrBarrier hget fail", zap.String("key", key), zap.Error(err))
		return
	}

	barrierId = int32(res)

	logger.CtxInfo(ctx, "GetCurrBarrier succ", zap.String("key", key), zap.Any("barrierId", barrierId))
	return
}

func GetBarrierAndArea(ctx context.Context, userId uint64) (barrierId int32, areaId int32, highArea int32, err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:u:%d:barrier:%d", userId)
	res, err := redis.StringMap(gRedis.Do(context.TODO(), "hgetall", key))
	if err == redis.ErrNil {
		err = nil
		logger.CtxInfo(ctx, "GetCurrBarrierAndArea nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.CtxError(ctx, "GetCurrBarrierAndArea fail", zap.String("key", key), zap.Error(err))
		return
	}

	for k, v := range res {
		if k == CURR_BARRIER_FIELD {
			barrierId = fkutil.ToInt32(v)
		} else if k == CURR_BARRIER_AREA_FIELD {
			areaId = fkutil.ToInt32(v)
		} else if k == HIGH_BARRIER_AREA_FIELD {
			highArea = fkutil.ToInt32(v)
		}
	}

	logger.CtxInfo(ctx, "GetCurrBarrierAndArea succ", zap.Any("barrierId", barrierId), zap.Any("areaId", areaId), zap.Any("highArea", highArea))
	return
}

func SetBarrier(ctx context.Context, userId uint64, barrierId int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	setMap := make(map[string]int32)
	setMap[CURR_BARRIER_FIELD] = barrierId

	err = BatchSet(ctx, userId, setMap)
	if err != nil {
		logger.CtxError(ctx, "SetBarrier hmset fail", zap.Error(err), zap.Any("setMap", setMap))
		return
	}
	logger.CtxInfo(ctx, "SetBarrier succ", zap.Any("setMap", setMap))
	return
}

func BatchSet(ctx context.Context, userId uint64, setMap map[string]int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:u:%d:barrier:%d", userId)
	var args []interface{}
	args = append(args, key)
	for k, v := range setMap {
		args = append(args, k)
		args = append(args, v)
	}
	_, err = gRedis.Do(context.TODO(), "hmset", args...)
	if err != nil {
		logger.CtxError(ctx, "BatchSet hmset fail", zap.String("key", key), zap.Any("args", args), zap.Error(err))
		return
	}

	return nil
}

// 删除
func GMDel(ctx context.Context, userId uint64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:u:%d:barrier:%d", userId)

	_, err = gRedis.Do(context.TODO(), "DEL", key)
	if err != nil {
		logger.CtxError(ctx, "GMDel fail", zap.Error(err), zap.String("key", key))
		return err
	}
	logger.CtxInfo(ctx, "GMDel succ ", zap.String("key", key))
	return nil
}

func GMHDEL(ctx context.Context, userId uint64, barrierId int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:u:%d:barrier:%d", userId)

	_, err = gRedis.Do(context.TODO(), "HDEL", key, barrierId)
	if err != nil {
		logger.CtxError(ctx, "GMHDEL fail", zap.Error(err), zap.String("key", key), zap.Any("barrierId", barrierId))
		return err
	}
	logger.CtxInfo(ctx, "GMHDEL succ ", zap.String("key", key), zap.Any("barrierId", barrierId))
	return nil
}
