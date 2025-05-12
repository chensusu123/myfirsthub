package mazebarrierredis

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/plate/freetk/fkutil"
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
func GetCurrBarrier(logger fklog.FKLogI, userId uint64) (barrierId int32, err error) {
	key := fmt.Sprintf("maze:u:%d:barrier:%d", userId)

	res, err := redis.Int(gRedis.Do(context.TODO(), "hget", key, CURR_BARRIER_FIELD))
	if err == redis.ErrNil {
		err = nil
		logger.InfoWF("GetCurrBarrier nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.ErrorWF("GetCurrBarrier hget fail", zap.String("key", key), zap.Error(err))
		return
	}

	barrierId = int32(res)

	logger.InfoWF("GetCurrBarrier succ", zap.String("key", key), zap.Any("barrierId", barrierId))
	return
}

func GetBarrierAndArea(logger fklog.FKLogI, userId uint64) (barrierId int32, areaId int32, highArea int32, err error) {
	key := fmt.Sprintf("maze:u:%d:barrier:%d", userId)
	res, err := redis.StringMap(gRedis.Do(context.TODO(), "hgetall", key))
	if err == redis.ErrNil {
		err = nil
		logger.InfoWF("GetCurrBarrierAndArea nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.ErrorWF("GetCurrBarrierAndArea fail", zap.String("key", key), zap.Error(err))
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

	logger.InfoWF("GetCurrBarrierAndArea succ", zap.Any("barrierId", barrierId), zap.Any("areaId", areaId), zap.Any("highArea", highArea))
	return
}

func SetBarrier(logger fklog.FKLogI, userId uint64, barrierId int32) (err error) {
	setMap := make(map[string]int32)
	setMap[CURR_BARRIER_FIELD] = barrierId

	err = BatchSet(logger, userId, setMap)
	if err != nil {
		logger.ErrorWF("SetBarrier hmset fail", zap.Error(err), zap.Any("setMap", setMap))
		return
	}
	logger.InfoWF("SetBarrier succ", zap.Any("setMap", setMap))
	return
}

func BatchSet(logger fklog.FKLogI, userId uint64, setMap map[string]int32) (err error) {
	key := fmt.Sprintf("maze:u:%d:barrier:%d", userId)
	var args []interface{}
	args = append(args, key)
	for k, v := range setMap {
		args = append(args, k)
		args = append(args, v)
	}
	_, err = gRedis.Do(context.TODO(), "hmset", args...)
	if err != nil {
		logger.ErrorWF("BatchSet hmset fail", zap.String("key", key), zap.Any("args", args), zap.Error(err))
		return
	}

	return nil
}

// 删除
func GMDel(logger fklog.FKLogI, userId uint64) (err error) {
	key := fmt.Sprintf("maze:u:%d:barrier:%d", userId)

	_, err = gRedis.Do(context.TODO(), "DEL", key)
	if err != nil {
		logger.ErrorWF("GMDel fail", zap.Error(err), zap.String("key", key))
		return err
	}
	logger.InfoWF("GMDel succ ", zap.String("key", key))
	return nil
}

func GMHDEL(logger fklog.FKLogI, userId uint64, barrierId int32) (err error) {
	key := fmt.Sprintf("maze:u:%d:barrier:%d", userId)

	_, err = gRedis.Do(context.TODO(), "HDEL", key, barrierId)
	if err != nil {
		logger.ErrorWF("GMHDEL fail", zap.Error(err), zap.String("key", key), zap.Any("barrierId", barrierId))
		return err
	}
	logger.InfoWF("GMHDEL succ ", zap.String("key", key), zap.Any("barrierId", barrierId))
	return nil
}
