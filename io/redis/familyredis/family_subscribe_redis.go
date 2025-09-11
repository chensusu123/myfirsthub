package familyredis

import (
	"context"
	"encoding/json"
	"fmt"
	globalredis "maze_game_server/io/redis"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/database/nanoredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// getKey 获取缓存操作key。
func getKey(db nanoredis.NanoRedisClient, args ...interface{}) string {
	return db.MakeSectionKey(fmt.Sprintf("subscribe:group:family:u:%d", args...))
}

func SaveSubscribe(ctx context.Context, userID uint64, groupIDs []int64) error {
	logger := fklog.ContextAppLogger(ctx)

	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "family SaveSubscribe Client fail",
			zap.Error(err),
		)
		return err
	}
	key := getKey(cli, userID)
	groupIDsStr, err := json.Marshal(groupIDs)
	if err != nil {
		logger.CtxError(ctx, "json.Marshal groupIDs fail",
			zap.Error(err),
		)
		return err
	}
	//hash field userid
	return cli.Set(ctx, key, string(groupIDsStr), 0).Err()

}

func GetSubscribe(ctx context.Context, userID uint64) ([]int64, error) {
	logger := fklog.ContextAppLogger(ctx)

	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "family GetSubscribe Client fail",
			zap.Error(err),
		)
		return nil, err
	}
	key := getKey(cli, userID)
	//hash field userid
	group, err := cli.Get(ctx, key).Result()

	if err != nil {
		return nil, err
	}
	groupIDs := []int64{}
	err = json.Unmarshal([]byte(group), &groupIDs)
	if err != nil {
		logger.CtxError(ctx, "json.Unmarshal groupIDs fail",
			zap.Error(err),
		)
		return nil, err
	}
	return groupIDs, nil
}

func DelSubscribe(ctx context.Context, userID uint64) error {
	logger := fklog.ContextAppLogger(ctx)

	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "family DelSubscribe Client fail",
			zap.Error(err),
		)
		return err
	}
	key := getKey(cli, userID)
	return cli.Del(ctx, key).Err()
}
