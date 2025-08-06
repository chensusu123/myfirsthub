package mazeequipspecialdropredis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	globalredis "maze_game_server/io/redis"
)

func getRedisKey(userId uint64) string {
	return fmt.Sprintf("maze:equip:special:drop:u:%d", userId)
}

func GetMazeEquipSpecialDropInfo(logger fklog.FKLogI, userId uint64) ([]byte, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	key := db.MakeSectionKey(getRedisKey(userId))
	bytes, err := db.Get(context.TODO(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		logger.ErrorWF("GetMazeEquipSpecialDropInfo GET", zap.String("key", key), zap.Error(err))
		return nil, err
	}
	logger.InfoWF("GetMazeEquipSpecialDropInfo success", zap.String("key", key))
	return bytes, nil
}

func SetMazeEquipSpecialDropInfo(logger fklog.FKLogI, userId uint64, bytes []byte) (err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(getRedisKey(userId))
	_, err = db.Set(context.TODO(), key, bytes, 0).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		logger.ErrorWF("SetMazeEquipSpecialDropInfo GET", zap.String("key", key), zap.Error(err))
		return err
	}
	logger.InfoWF("SetMazeEquipSpecialDropInfo success", zap.Any("bytes", bytes), zap.String("key", key))
	return
}

// gm删除
func GMDel(logger fklog.FKLogI, userId uint64) (err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(getRedisKey(userId))
	err = db.Del(context.TODO(), key).Err()
	if err != nil {
		logger.ErrorWF("SetMazeEquipSpecialDropInfo GMDel fail", zap.String("key", key), zap.Error(err))
		return
	}
	return
}
