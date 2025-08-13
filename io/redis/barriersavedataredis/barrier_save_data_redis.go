package barriersavedataredis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	globalredis "maze_game_server/io/redis"
)

func getKey(userId uint64, barrier int32) string {
	return fmt.Sprintf("save:u:%d:barrier:%d", userId, barrier)
}

func SetBarrierSaveData(logger fklog.FKLogI, userId uint64, barrier int32, bytes []byte) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(getKey(userId, barrier))

	err = db.Set(context.TODO(), key, bytes, 0).Err()
	if err != nil {
		logger.ErrorWF("SetBarrierSaveData", zap.String("key", key), zap.Int32("barrier", barrier), zap.Any("bytes", string(bytes)),
			zap.Error(err))
		return err
	}
	logger.InfoWF("SetBarrierSaveData success", zap.String("key", key), zap.Int32("barrier", barrier))
	return nil
}

func GetBarrierSaveData(logger fklog.FKLogI, userId uint64, barrier int32) ([]byte, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	key := db.MakeSectionKey(getKey(userId, barrier))

	bytes, err := db.Get(context.TODO(), key).Bytes()

	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		logger.ErrorWF("GetBarrierSaveData GET", zap.String("key", key), zap.Int32("barrier", barrier), zap.Error(err))
		return nil, err
	}

	logger.InfoWF("GetBarrierSaveData success", zap.String("key", key), zap.Int32("barrier", barrier))
	return bytes, nil
}

func DelBarrierSaveData(logger fklog.FKLogI, userId uint64, barrier int32) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(getKey(userId, barrier))

	err = db.Del(context.TODO(), key).Err()
	if err != nil {
		logger.ErrorWF("DelBarrierSaveData DEL", zap.String("key", key), zap.Error(err), zap.Int32("barrier", barrier))
		return err
	}
	logger.InfoWF("DelBarrierSaveData success", zap.String("key", key), zap.Int32("barrier", barrier))
	return nil
}
