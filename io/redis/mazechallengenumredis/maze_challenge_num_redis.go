package mazechallengenumredis

import (
	"context"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
)

var gRedis = &fkredis.FkRedis{}

func init() {
	fkconfig.RegisterNameNode("mazechallengenumredis", 21737, gRedis)
}

func GetUserChallengeNum(logger fklog.FKLogI, userId uint64, dateTime int) (num int, err error) {
	key := gRedis.GetKey(userId, dateTime)
	num, err = redis.Int(gRedis.Do(context.TODO(), "get", key))
	if err == redis.ErrNil {
		err = nil
		logger.InfoWF("GetUserChallengeNum nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.ErrorWF("GetUserChallengeNum fail", zap.String("key", key), zap.Error(err))
		return
	}

	return
}

func AddUserChallengeNum(logger fklog.FKLogI, userId uint64, dateTime int, count int32) (err error) {
	key := gRedis.GetKey(userId, dateTime)
	newCount, err := redis.Int(gRedis.Do(context.TODO(), "incrby", key, count))
	if err != nil {
		logger.ErrorWF("AddUserChallengeNum incr fail", zap.String("key", key), zap.Error(err))
		return
	}
	logger.InfoWF("AddUserChallengeNum succ", zap.Any("key", key), zap.Any("count", count), zap.Any("newCount", newCount))
	return
}

func GMDel(logger fklog.FKLogI, userId uint64, dateTime int) (err error) {
	key := gRedis.GetKey(userId, dateTime)
	_, err = redis.Int(gRedis.Do(context.TODO(), "del", key))
	if err != nil {
		logger.ErrorWF("GMDel del fail", zap.String("key", key), zap.Error(err))
		return
	}
	logger.InfoWF("GMDel succ", zap.Any("key", key))
	return
}
