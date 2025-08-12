package barrierarearecordredis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/redis"
)

// getKey 获取缓存操作key。
func gePassAreaKey(userId uint64, barrierId int32) string {
	return fmt.Sprintf("area:u:%d:barrier:record:%d", userId, barrierId)
}

func GetBarrierAreaNumRecord(logger fklog.FKLogI, userId uint64, barrierId int32) ([]byte, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	key := db.MakeSectionKey(gePassAreaKey(userId, barrierId))

	bytes, err := db.Get(context.TODO(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		logger.ErrorWF("GetBarrierAreaNumRecord GET err", zap.String("key", key), zap.Error(err))
		return nil, err
	}

	return bytes, nil
}

func SetBarrierAreaNumRecord(logger fklog.FKLogI, userId uint64, barrierId int32, data []byte) (err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(gePassAreaKey(userId, barrierId))

	err = db.Set(context.TODO(), key, data, 0).Err()
	if err != nil {
		logger.ErrorWF("SetBarrierAreaNumRecord err", zap.String("key", key), zap.Any("passArea", string(data)),
			zap.Error(err))
		return err
	}
	logger.InfoWF("SetBarrierAreaNumRecord end", zap.String("key", key), zap.Any("passArea", string(data)))
	return nil
}

func DelBarrierAreaNumRecord(logger fklog.FKLogI, userId uint64, barrierId int32) (err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(gePassAreaKey(userId, barrierId))

	err = db.Del(context.TODO(), key).Err()
	if err != nil {
		logger.ErrorWF("DelBarrierAreaNumRecord err", zap.String("key", key), zap.Error(err))
		return err
	}
	logger.InfoWF("DelBarrierAreaNumRecord end", zap.String("key", key))
	return nil
}
