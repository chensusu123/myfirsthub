package useridredis

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
	_ = fkconfig.RegisterNameNode("useridredis", 17545, gRedis)
}

var gRegionID = uint64(1)

const maxUserID = 10000000

func SetRegionID(regionID uint64) {
	gRegionID = regionID
}

func getKey(userId uint64) string {
	return fmt.Sprintf("uid:generate", userId)
}

// 获取用户信息是否封禁
func Generate(logger fklog.FKLogI) uint64 {
	groupID := fkconfig.EnvVal.GroupID
	if groupID == 5 {
		gRegionID = uint64(1)
	} else {
		gRegionID = uint64(groupID)
	}

	key := "uid:generate"
	userId, err := redis.Int64(gRedis.Do(context.TODO(), "INCR", key))
	if err != nil {
		logger.ErrorWF("Generate error", zap.Error(err))
		return 0
	}
	if userId <= maxUserID {
		begin := gRegionID * maxUserID
		gRedis.Do(context.TODO(), "SET", key, begin)
	}
	userId, err = redis.Int64(gRedis.Do(context.TODO(), "INCR", key))
	if err != nil {
		logger.ErrorWF("Generate error", zap.Error(err))
		return 0
	}
	return uint64(userId)
}
