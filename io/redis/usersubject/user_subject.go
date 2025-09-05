package usersubject

import (
	"context"
	"fmt"
	"time"

	globalredis "maze_game_server/io/redis"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/database/nanoredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// getKey 获取缓存操作key。
func getKey(db nanoredis.NanoRedisClient, userID int64) string {
	return db.MakeSectionKey(fmt.Sprintf("user:subject:%d", userID))
}

// Add
func Add(ctx context.Context, userID int64, subject string) error {
	logger := fklog.ContextAppLogger(ctx)
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "Subject Add Client fail",
			zap.Int64("userID", userID),
			zap.String("subject", subject),
			zap.Error(err),
		)
		return err
	}
	key := getKey(cli, userID)
	_, err = cli.HSet(ctx, key, subject, time.Now().Unix()).Result()
	return err
}

func Del(ctx context.Context, userID int64, subject string) error {
	logger := fklog.ContextAppLogger(ctx)
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "Subject Del Client fail",
			zap.Int64("userID", userID),
			zap.String("subject", subject),
			zap.Error(err),
		)
		return err
	}
	key := getKey(cli, userID)
	_, err = cli.HDel(ctx, key, subject).Result()
	return err
}
