package mazebarriermoneyredis

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
)

var gRedis = &fkredis.FkRedis{}

func init() {
	fkconfig.RegisterNameNode("mazebarriermoneyredis", 21685, gRedis)
}

// 查所有
func GetAllMoney(logger fklog.FKLogI, userId uint64) (moneyMap map[int32]int64, err error) {
	moneyMap = make(map[int32]int64)
	key := fmt.Sprintf("maze:barrier:money:u:%d", userId)
	res, err := redis.Int64s(gRedis.Do(context.TODO(), "hgetall", key))
	if err == redis.ErrNil {
		err = nil
		logger.InfoWF("GetAllMoney nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.ErrorWF("GetAllMoney fail", zap.String("key", key), zap.Error(err))
		return
	}

	sLen := len(res)
	for i := 0; i < sLen; i += 2 {
		moneyMap[int32(res[i])] = res[i+1]
	}

	logger.InfoWF("GetAllMoney succ", zap.Any("moneyMap", moneyMap))
	return
}

// 指定查
func GetMoneyCount(logger fklog.FKLogI, userId uint64, moneyId int32) (moneyCount int64, err error) {
	key := fmt.Sprintf("maze:barrier:money:u:%d", userId)
	moneyCount, err = redis.Int64(gRedis.Do(context.TODO(), "hget", key, moneyId))
	if err == redis.ErrNil {
		err = nil
		logger.InfoWF("GetMoneyCount nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.ErrorWF("GetMoneyCount cfgId fail", zap.Error(err))
		return
	}

	logger.InfoWF("GetMoneyCount succ", zap.Any("moneyCount", moneyCount))
	return
}

func SetMoney(logger fklog.FKLogI, userId uint64, moneyId int32, moneyCount int64) (err error) {
	key := fmt.Sprintf("maze:barrier:money:u:%d", userId)
	_, err = gRedis.Do(context.TODO(), "hset", key, moneyId, moneyCount)
	if err != nil {
		logger.ErrorWF("SetMoney hset fail", zap.Error(err))
		return
	}

	logger.InfoWF("SetMoney succ", zap.String("key", key), zap.Any("moneyId", moneyId), zap.Any("moneyCount", moneyCount))
	return
}

func HMSetMoney(logger fklog.FKLogI, userId uint64, moneyMap map[int32]int64) (err error) {
	if len(moneyMap) <= 0 {
		return
	}
	key := fmt.Sprintf("maze:barrier:money:u:%d", userId)
	var args []interface{}
	args = append(args, key)
	for k, v := range moneyMap {
		args = append(args, k)
		args = append(args, v)
	}
	_, err = gRedis.Do(context.TODO(), "hmset", args...)
	if err != nil {
		logger.ErrorWF("HMSetMoney hmset fail", zap.String("key", key), zap.Any("args", args), zap.Error(err))
		return
	}

	return nil
}

// 加钱
func AddMoney(logger fklog.FKLogI, userId uint64, moneyId int32, moneyCount int64) (newCount int64, err error) {
	key := fmt.Sprintf("maze:barrier:money:u:%d", userId)
	newCount, err = redis.Int64(gRedis.Do(context.TODO(), "hincrby", key, moneyId, moneyCount))
	if err != nil {
		logger.ErrorWF("AddMoney hincrby fail", zap.Error(err))
		return
	}

	logger.InfoWF("AddMoney succ", zap.String("key", key), zap.Any("moneyId", moneyId), zap.Any("moneyCount", moneyCount))
	return
}

// 扣钱
func SubMoney(logger fklog.FKLogI, userId uint64, moneyId int32, moneyCount int64) (err error) {
	key := fmt.Sprintf("maze:barrier:money:u:%d", userId)
	_, err = gRedis.Do(context.TODO(), "hincrby", key, moneyId, 0-moneyCount)
	if err != nil {
		logger.ErrorWF("SubMoney hincrby fail", zap.Error(err))
		return
	}

	logger.InfoWF("SubMoney succ", zap.String("key", key), zap.Any("moneyId", moneyId), zap.Any("moneyCount", moneyCount))
	return
}

// 删除所有数据
func GMDel(logger fklog.FKLogI, userId uint64) (err error) {
	key := fmt.Sprintf("maze:barrier:money:u:%d", userId)

	_, err = gRedis.Do(context.TODO(), "DEL", key)
	if err != nil {
		logger.ErrorWF("GMDel fail", zap.Error(err), zap.String("key", key))
		return err
	}
	logger.InfoWF("GMDel succ ", zap.String("key", key))
	return nil
}

func GMSet(logger fklog.FKLogI, userId uint64, moneyId int32, moneyCount int64) (err error) {
	key := fmt.Sprintf("maze:barrier:money:u:%d", userId)

	_, err = gRedis.Do(context.TODO(), "HSET", key, moneyId, moneyCount)
	if err != nil {
		logger.ErrorWF("GMSet fail", zap.Error(err), zap.String("key", key))
		return err
	}
	logger.InfoWF("GMSet succ ", zap.String("key", key), zap.Any("moneyId", moneyId), zap.Any("moneyCount", moneyCount))
	return nil
}
