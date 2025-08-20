package useridredis

import (
	"context"
	"strconv"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"go.uber.org/zap"
)

var gRedis = &fkredis.FkRedis{}

func init() {
	_ = fkconfig.RegisterNameNode("useridredis", 17545, gRedis)
}

var gRegionID = 0

const maxUserID = 10000000

// 获取用户信息是否封禁
func Generate(ctx context.Context, logger fklog.FKLogI) uint64 {
	if gRegionID == 0 {
		tmpVal, err := strconv.ParseUint(appconfig.GlobalConfig().Global.SectionID, 10, 32)
		if err != nil {
			gRegionID = 1
		} else {
			gRegionID = int(tmpVal)
		}
	}

	key := "uid:generate"
	userId, err := redis.Int64(gRedis.Do(ctx, "INCR", key))
	if err != nil {
		logger.ErrorWF("Generate error", zap.Error(err))
		return 0
	}
	if userId <= maxUserID {
		begin := gRegionID * maxUserID
		gRedis.Do(ctx, "SET", key, begin)
	}
	userId, err = redis.Int64(gRedis.Do(ctx, "INCR", key))
	if err != nil {
		logger.ErrorWF("Generate error", zap.Error(err))
		return 0
	}
	return uint64(userId)
}
