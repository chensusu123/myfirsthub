package mazecardlistgroupredis

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"go.uber.org/zap"
)

var gRedis = &fkredis.FkRedis{}

// maze:card:group:%d:list
func init() {
	fkconfig.RegisterNameNode("MazeCardListRedis", 21683, gRedis)
}

func getKey() string {
	sectionID := appconfig.GlobalConfig().Global.SectionID
	return fmt.Sprintf("maze:card:group:%s:list", sectionID)
}

func SetMazeCard(ctx context.Context, userId uint64, expirationTime int64) error {
	logger := fklog.ContextAppLogger(ctx)
	key := getKey()
	_, err := gRedis.Do(ctx, "zadd", key, expirationTime, userId)
	if err != nil {
		logger.CtxError(ctx, "SetMazeCard zadd failed", zap.String("key", key), zap.Uint64("userId", userId),
			zap.Int64("time", expirationTime), zap.Error(err))
		return err
	}

	logger.CtxInfo(ctx, "SetMazeCard end", zap.Uint64("userId", userId), zap.Int64("time", expirationTime))
	return nil
}

func GetMazeCard(ctx context.Context, userId uint64) (int64, error) {
	logger := fklog.ContextAppLogger(ctx)
	key := getKey()
	expirationTime, err := redis.Int64(gRedis.Do(ctx, "zscore", key, userId))
	if err != nil && err != redis.ErrNil {
		logger.CtxError(ctx, "GetMazeCard zscore failed", zap.Uint64("userId", userId), zap.Error(err))
		return 0, err
	}

	logger.CtxInfo(ctx, "GetMazeCard end", zap.Uint64("userId", userId), zap.Int64("time", expirationTime))
	return 0, nil
}

func BatchDelMazeCard(ctx context.Context, userList []int64) error {
	if len(userList) == 0 {
		return nil
	}
	logger := fklog.ContextAppLogger(ctx)
	param := make([]interface{}, 0, len(userList)+2)
	param = append(param, getKey())
	for _, userId := range userList {
		param = append(param, userId)
	}

	_, err := gRedis.Do(ctx, "zrem", param...)
	if err != nil {
		logger.CtxError(ctx, "BatchDelMazeCard zrem failed", zap.Int64s("userList", userList), zap.Error(err))
		return err
	}

	logger.CtxInfo(ctx, "BatchDelMazeCard end", zap.Int64s("userList", userList))
	return nil
}

func DelMazeCard(ctx context.Context, userId uint64) error {
	logger := fklog.ContextAppLogger(ctx)
	key := getKey()
	_, err := gRedis.Do(ctx, "zrem", key, userId)
	if err != nil {
		logger.CtxError(ctx, "DelMazeCard zrem failed", zap.Uint64("userId", userId), zap.Error(err))
		return err
	}

	logger.CtxInfo(ctx, "DelMazeCard end", zap.Uint64("userId", userId))
	return nil
}

func GetMazeCardExpirationList(ctx context.Context) ([]int64, error) {
	logger := fklog.ContextAppLogger(ctx)
	key := getKey()
	userList, err := redis.Int64s(gRedis.Do(ctx, "zrangebyscore", key, "-inf", time.Now().Unix()))
	if err != nil {
		logger.CtxError(ctx, "GetMazeCardExpirationList zrangebyscore failed", zap.Error(err))
		return nil, err
	}

	logger.CtxInfo(ctx, "GetMazeCardExpirationList end", zap.Int64s("userList", userList))
	return userList, nil
}
