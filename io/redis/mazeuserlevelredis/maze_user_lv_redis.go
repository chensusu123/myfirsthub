package mazeuserlevelredis

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

var (
	USER_LEVEL            = "level"
	USER_EXP              = "exp"
	USER_TOTAL_EXP        = "total_exp"
	USER_TYPE             = "user_type"
	USER_BARRIER          = "barrier"
	USER_PASS_BARRIER     = "pass_barrier"
	USER_EQUIP_POINT      = "equip_point"
	USER_ENERGY           = "energy"
	USER_ENERGY_LAST_TIME = "energy_last_time"
)

var gRedis = &fkredis.FkRedis{}

func init() {
	fkconfig.RegisterNameNode("mazeuserlevelredis", 21613, gRedis)
}

func GetUserLevel(ctx context.Context, userId uint64) (level int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:u:%d:level:exp", userId)
	level, err = redis.Int64(gRedis.Do(ctx, "hget", key, USER_LEVEL))
	if err == redis.ErrNil {
		err = nil
		logger.CtxError(ctx, "GetUserLevel nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.CtxError(ctx, "GetUserLevel fail", zap.String("key", key), zap.Error(err))
		return
	}

	return
}

func GetUserInfo(ctx context.Context, userId uint64) (userInfo map[string]int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:u:%d:level:exp", userId)
	res, err := redis.StringMap(gRedis.Do(ctx, "hgetall", key))
	if err == redis.ErrNil {
		err = nil
		logger.CtxInfo(ctx, "GetUserInfo nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.CtxError(ctx, "GetUserInfo fail", zap.String("key", key), zap.Error(err))
		return
	}

	userInfo = make(map[string]int64)

	for k, v := range res {
		userInfo[k] = fkutil.ToInt64(v)
	}

	logger.CtxInfo(ctx, "GetUserInfo succ", zap.Any("userInfo", userInfo))
	return
}

func SetUserInfo(ctx context.Context, userId uint64, userInfo map[string]int64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:u:%d:level:exp", userId)
	var args []interface{}
	args = append(args, key)
	for k, v := range userInfo {
		args = append(args, k)
		args = append(args, v)
	}
	if len(args) == 1 {
		logger.CtxError(ctx, "SetUserInfo params empty", zap.Any("userInfo", userInfo))
		return
	}
	_, err = gRedis.Do(ctx, "hmset", args...)
	if err != nil {
		logger.CtxError(ctx, "SetUserInfo hmset fail", zap.String("key", key), zap.Any("args", args), zap.Error(err))
		return
	}
	logger.CtxInfo(ctx, "SetUserInfo succ", zap.String("key", key), zap.Any("userInfo", userInfo))
	return
}

func GMDel(ctx context.Context, userId uint64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:u:%d:level:exp", userId)
	_, err = redis.Int64(gRedis.Do(ctx, "del", key))
	if err != nil {
		logger.CtxError(ctx, "GMDel fail", zap.String("key", key), zap.Error(err))
		return
	}

	logger.CtxInfo(ctx, "GMDel succ", zap.Any("key", key))
	return
}
