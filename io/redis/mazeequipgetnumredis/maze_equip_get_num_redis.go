package mazeequipgetnumredis

import (
	"context"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
)

var gRedis = &fkredis.FkRedis{}

func init() {
	// 21643 maze:equip:get:num:%llu 迷宫游戏装备获取次数信息
	fkconfig.RegisterNameNode("mazeequipgetnumredis", 21643, gRedis)
}

func GMDel(logger fklog.FKLogI, userId uint64) error {
	key := gRedis.GetKey(userId)
	_, err := redis.Int64(gRedis.Do(context.TODO(), "del", key))
	if err != nil {
		logger.ErrorWF("GMDel failed with", zap.Error(err), zap.Any("key", key))
		return err
	}
	logger.InfoWF("mazeequipgetnumredis GMDel succ", zap.Any("key", key))
	return nil
}

func GetEquipGetNum(logger fklog.FKLogI, userId uint64, equipId int32) (int64, error) {
	key := gRedis.GetKey(userId)
	ret, err := redis.Int64(gRedis.Do(context.TODO(), "hget", key, equipId))
	if err == redis.ErrNil {
		err = nil
		return 0, err
	}
	if err != nil {
		logger.ErrorWF("GetEquipGetNum get score failed with", zap.Error(err), zap.Int32("equipId", equipId))
		return 0, err
	}
	logger.InfoWF("GetEquipGetNum get score with", zap.Int32("equipId", equipId), zap.Int64("score", ret))
	return ret, nil
}
