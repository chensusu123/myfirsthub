package userdisableredis

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
	_ = fkconfig.RegisterNameNode("userdisableredis", 17545, gRedis)
}

func getKey(userId uint64) string {
	return fmt.Sprintf("uid:%d", fmt.Sprintf("%d", userId))
}

// 获取用户信息是否封禁
func IsUserDisable(ctx context.Context, userId uint64) bool {
	logger := fklog.ContextAppLogger(ctx)
	r, err := redis.Bool(gRedis.Do(context.TODO(), "EXISTS", getKey(userId)))
	if err != nil {
		if err != redis.ErrNil {
			logger.CtxError(ctx, "IsUserDisable exists error", zap.Uint64("userId", userId), zap.Error(err))
		}
		return false
	}
	return r
}
