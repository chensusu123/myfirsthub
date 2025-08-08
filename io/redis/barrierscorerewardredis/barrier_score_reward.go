package barrierscorerewardredis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	globalredis "maze_game_server/io/redis"
)

func getKey(userId uint64, barrier int32) string {
	return fmt.Sprintf("score:reward:u:%d:barrier:%d", userId, barrier)
}

func SetBarrierScoreReward(logger fklog.FKLogI, userId uint64, barrier int32, bytes []byte) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(getKey(userId, barrier))

	err = db.Set(context.TODO(), key, bytes, 0).Err()
	if err != nil {
		logger.ErrorWF("SetBarrierScoreReward", zap.String("key", key), zap.Int32("barrier", barrier), zap.Any("bytes", bytes),
			zap.Error(err))
		return err
	}
	logger.InfoWF("SetBarrierScoreReward success", zap.String("key", key), zap.Int32("barrier", barrier))
	return nil
}

func GetBarrierScoreReward(logger fklog.FKLogI, userId uint64, barrier int32) ([]byte, error) {
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
		logger.ErrorWF("GetBarrierScoreReward GET", zap.String("key", key), zap.Int32("barrier", barrier), zap.Error(err))
		return nil, err
	}

	logger.InfoWF("GetBarrierScoreReward success", zap.String("key", key), zap.Int32("barrier", barrier))
	return bytes, nil
}

func DelBarrierScoreReward(logger fklog.FKLogI, userId uint64, barrier int32) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(getKey(userId, barrier))

	err = db.Del(context.TODO(), key).Err()
	if err != nil {
		logger.ErrorWF("DelBarrierScoreReward DEL", zap.String("key", key), zap.Error(err), zap.Int32("barrier", barrier))
		return err
	}
	logger.InfoWF("DelBarrierScoreReward success", zap.String("key", key), zap.Int32("barrier", barrier))
	return nil
}
