package mazemailredis

import (
	"context"
	"fmt"
	globalredis "maze_game_server/io/redis"

	"github.com/redis/go-redis/v9"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// getKey 获取缓存操作key
func getMailKey(userId uint64) string {
	return fmt.Sprintf("maze:mail:u:%d", userId)
}

func GetMail(ctx context.Context, userId uint64) ([]byte, error) {
	logger := fklog.ContextAppLogger(ctx)
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	key := db.MakeSectionKey(getMailKey(userId))

	bytes, err := db.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		logger.CtxError(ctx, "GetMail GET err", zap.String("key", key), zap.Error(err))
		return nil, err
	}
	return bytes, nil
}

func SetMail(ctx context.Context, userId uint64, data []byte) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(getMailKey(userId))

	err = db.Set(ctx, key, data, 0).Err()
	if err != nil {
		logger.CtxError(ctx, "SetMail Set err", zap.String("key", key), zap.Any("mail", string(data)),
			zap.Error(err))
		return err
	}
	logger.CtxInfo(ctx, "SetMail Set end", zap.String("key", key), zap.Any("mail", string(data)))
	return nil
}

func DelMail(ctx context.Context, userId uint64, barrierId int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(getMailKey(userId))

	err = db.Del(ctx, key).Err()
	if err != nil {
		logger.CtxError(ctx, "DelMail Del err", zap.String("key", key), zap.Error(err))
		return err
	}
	logger.CtxInfo(ctx, "DelMail Del end", zap.String("key", key))
	return nil
}
