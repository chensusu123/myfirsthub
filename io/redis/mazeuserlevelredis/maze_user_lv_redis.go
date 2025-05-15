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

var USER_LEVEL = "level"
var USER_EXP = "exp"
var USER_TOTAL_EXP = "total_exp"
var USER_TYPE = "user_type"
var USER_BARRIER = "barrier"
var USER_PASS_BARRIER = "pass_barrier"
var USER_EQUIP_POINT = "equip_point"
var USER_ENERGY = "energy"
var USER_ENERGY_LAST_TIME = "energy_last_time"

var gRedis = &fkredis.FkRedis{}

func init() {
	fkconfig.RegisterNameNode("mazeuserlevelredis", 21613, gRedis)
}

func GetUserLevel(logger fklog.FKLogI, userId uint64) (level int64, err error) {
	key := fmt.Sprintf("maze:u:%d:level:exp", userId)
	level, err = redis.Int64(gRedis.Do(context.TODO(), "hget", key, USER_LEVEL))
	if err == redis.ErrNil {
		err = nil
		logger.InfoWF("GetUserLevel nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.ErrorWF("GetUserLevel fail", zap.String("key", key), zap.Error(err))
		return
	}

	return
}

func GetUserInfo(logger fklog.FKLogI, userId uint64) (userInfo map[string]int64, err error) {
	key := fmt.Sprintf("maze:u:%d:level:exp", userId)
	res, err := redis.StringMap(gRedis.Do(context.TODO(), "hgetall", key))
	if err == redis.ErrNil {
		err = nil
		logger.InfoWF("GetUserInfo nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.ErrorWF("GetUserInfo fail", zap.String("key", key), zap.Error(err))
		return
	}

	userInfo = make(map[string]int64)

	for k, v := range res {
		userInfo[k] = fkutil.ToInt64(v)
	}

	logger.InfoWF("GetUserInfo succ", zap.Any("userInfo", userInfo))
	return
}

func SetUserInfo(logger fklog.FKLogI, userId uint64, userInfo map[string]int64) (err error) {
	key := fmt.Sprintf("maze:u:%d:level:exp", userId)
	var args []interface{}
	args = append(args, key)
	for k, v := range userInfo {
		args = append(args, k)
		args = append(args, v)
	}
	if len(args) == 1 {
		logger.ErrorWF("SetUserInfo params empty", zap.Any("userInfo", userInfo))
		return
	}
	_, err = gRedis.Do(context.TODO(), "hmset", args...)
	if err != nil {
		logger.ErrorWF("SetUserInfo hmset fail", zap.String("key", key), zap.Any("args", args), zap.Error(err))
		return
	}
	logger.InfoWF("SetUserInfo succ", zap.String("key", key), zap.Any("userInfo", userInfo))
	return
}

func GMDel(logger fklog.FKLogI, userId uint64) (err error) {
	key := fmt.Sprintf("maze:u:%d:level:exp", userId)
	_, err = redis.Int64(gRedis.Do(context.TODO(), "del", key))
	if err != nil {
		logger.ErrorWF("GMDel fail", zap.String("key", key), zap.Error(err))
		return
	}

	logger.InfoWF("GMDel succ", zap.Any("key", key))
	return
}
