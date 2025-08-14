package mazemailredis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/redis"
)

// getKey 获取缓存操作key
func getMailKey(userId uint64) string {
	return fmt.Sprintf("maze:mail:u:%d", userId)
}

func GetMail(logger fklog.FKLogI, userId uint64) ([]byte, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	key := db.MakeSectionKey(getMailKey(userId))

	bytes, err := db.Get(context.TODO(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		logger.ErrorWF("GetMail GET err", zap.String("key", key), zap.Error(err))
		return nil, err
	}
	return bytes, nil
}

func SetMail(logger fklog.FKLogI, userId uint64, data []byte) (err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(getMailKey(userId))

	err = db.Set(context.TODO(), key, data, 0).Err()
	if err != nil {
		logger.ErrorWF("SetMail Set err", zap.String("key", key), zap.Any("mail", string(data)),
			zap.Error(err))
		return err
	}
	logger.InfoWF("SetMail Set end", zap.String("key", key), zap.Any("mail", string(data)))
	return nil
}

func DelMail(logger fklog.FKLogI, userId uint64, barrierId int32) (err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(getMailKey(userId))

	err = db.Del(context.TODO(), key).Err()
	if err != nil {
		logger.ErrorWF("DelMail Del err", zap.String("key", key), zap.Error(err))
		return err
	}
	logger.InfoWF("DelMail Del end", zap.String("key", key))
	return nil
}
