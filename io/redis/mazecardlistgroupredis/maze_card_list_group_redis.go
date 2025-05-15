package mazecardlistgroupredis

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"go.uber.org/zap"
)

var gRedis = &fkredis.FkRedis{}

// maze:card:group:%d:list
func init() {
	fkconfig.RegisterNameNode("MazeCardListRedis", 21683, gRedis)
}

func SetMazeCard(logger fklog.FKLogI, userId uint64, expirationTime int64) error {
	key := fmt.Sprintf("maze:card:group:%d:list", fkconfig.EnvVal.GroupID)
	_, err := gRedis.Do(context.TODO(), "zadd", key, expirationTime, userId)
	if err != nil {
		logger.ErrorWF("SetMazeCard zadd failed", zap.String("key", key), zap.Uint64("userId", userId),
			zap.Int64("time", expirationTime), zap.Error(err))
		return err
	}

	logger.InfoWF("SetMazeCard end", zap.Uint64("userId", userId), zap.Int64("time", expirationTime))
	return nil
}

func GetMazeCard(logger fklog.FKLogI, userId uint64) (int64, error) {
	key := fmt.Sprintf("maze:card:group:%d:list", fkconfig.EnvVal.GroupID)
	expirationTime, err := redis.Int64(gRedis.Do(context.TODO(), "zscore", key, userId))
	if err != nil && err != redis.ErrNil {
		logger.ErrorWF("GetMazeCard zscore failed", zap.Uint64("userId", userId), zap.Error(err))
		return 0, err
	}

	logger.InfoWF("GetMazeCard end", zap.Uint64("userId", userId), zap.Int64("time", expirationTime))
	return 0, nil
}

func BatchDelMazeCard(logger fklog.FKLogI, userList []int64) error {
	if len(userList) == 0 {
		return nil
	}

	param := make([]interface{}, 0, len(userList)+2)
	param = append(param, fmt.Sprintf("maze:card:group:%d:list", fkconfig.EnvVal.GroupID))
	for _, userId := range userList {
		param = append(param, userId)
	}

	_, err := gRedis.Do(context.TODO(), "zrem", param...)
	if err != nil {
		logger.ErrorWF("BatchDelMazeCard zrem failed", zap.Int64s("userList", userList), zap.Error(err))
		return err
	}

	logger.InfoWF("BatchDelMazeCard end", zap.Int64s("userList", userList))
	return nil
}

func DelMazeCard(logger fklog.FKLogI, userId uint64) error {
	key := fmt.Sprintf("maze:card:group:%d:list", fkconfig.EnvVal.GroupID)
	_, err := gRedis.Do(context.TODO(), "zrem", key, userId)
	if err != nil {
		logger.ErrorWF("DelMazeCard zrem failed", zap.Uint64("userId", userId), zap.Error(err))
		return err
	}

	logger.InfoWF("DelMazeCard end", zap.Uint64("userId", userId))
	return nil
}

func GetMazeCardExpirationList(logger fklog.FKLogI) ([]int64, error) {
	key := fmt.Sprintf("maze:card:group:%d:list", fkconfig.EnvVal.GroupID)
	userList, err := redis.Int64s(gRedis.Do(context.TODO(), "zrangebyscore", key, "-inf", time.Now().Unix()))
	if err != nil {
		logger.ErrorWF("GetMazeCardExpirationList zrangebyscore failed", zap.Error(err))
		return nil, err
	}

	logger.InfoWF("GetMazeCardExpirationList end", zap.Int64s("userList", userList))
	return userList, nil
}
