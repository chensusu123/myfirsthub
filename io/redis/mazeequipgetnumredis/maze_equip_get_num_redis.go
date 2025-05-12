package mazeequipgetnumredis

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/plate/freetk/fkutil"
	"go.uber.org/zap"
)

var gRedis = &fkredis.FkRedis{}

func init() {
	// 21643 maze:equip:get:num:%llu 迷宫游戏装备获取次数信息
	fkconfig.RegisterNameNode("mazeequipgetnumredis", 21643, gRedis)
}

func GMDel(logger fklog.FKLogI, userId uint64) error {
	key := fmt.Sprintf("maze:equip:get:num:%d", userId)
	_, err := redis.Int64(gRedis.Do(context.TODO(), "del", key))
	if err != nil {
		logger.ErrorWF("GMDel failed with", zap.Error(err), zap.Any("key", key))
		return err
	}
	logger.InfoWF("mazeequipgetnumredis GMDel succ", zap.Any("key", key))
	return nil
}

func GetEquipGetNum(logger fklog.FKLogI, userId uint64, equipId int32) (int64, error) {
	key := fmt.Sprintf("maze:equip:get:num:%d", userId)
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

func GetEquipGetNumInc(logger fklog.FKLogI, userId uint64, equipId int32, addition int32) (int64, error) {
	key := fmt.Sprintf("maze:equip:get:num:%d", userId)
	ret, err := redis.Int64(gRedis.Do(context.TODO(), "HINCRBY", key, equipId, addition))
	if err != nil {
		logger.ErrorWF("GetEquipGetNumInc get new guid failed with", zap.Error(err), zap.Int32("equipId", equipId), zap.Int32("add", addition))
		return 0, err
	}
	logger.InfoWF("GetEquipGetNumInc ret new guid with", zap.Int32("equipId", equipId), zap.Int32("add", addition), zap.Int64("score", ret))
	return ret, nil
}

// 保存装备的次数分值
func SetEquipGetNum(logger fklog.FKLogI, userId uint64, equipId, score int32) error {
	key := fmt.Sprintf("maze:equip:get:num:%d", userId)
	_, err := gRedis.Do(context.TODO(), "hset", key, equipId, score)
	if err != nil {
		logger.ErrorWF("SetEquipGetNum set score failed with", zap.Error(err),
			zap.Int32("equipId", equipId),
			zap.Int32("score", score))
		return err
	}
	logger.InfoWF("SetEquipGetNum set score with",
		zap.Int32("equipId", equipId),
		zap.Int32("score", score))
	return err
}

func BatchSetEquipGetNum(logger fklog.FKLogI, userId uint64, equipScoreMap map[int32]int32) (err error) {
	key := fmt.Sprintf("maze:equip:get:num:%d", userId)

	args := make([]interface{}, 0, len(equipScoreMap)*2+1)
	args = append(args, key)
	for k, v := range equipScoreMap {
		args = append(args, k)
		args = append(args, v)
	}
	if len(args) <= 1 {
		return nil
	}
	_, err = gRedis.Do(context.TODO(), "HMSET", args...)
	if err != nil {
		logger.ErrorWF("BatchSetEquipGetNum fail", zap.Error(err), zap.Any("equipScoreMap", equipScoreMap), zap.String("key", key))
		return err
	}
	logger.InfoWF("BatchSetEquipGetNum succ", zap.Any("equipScoreMap", equipScoreMap), zap.String("key", key))
	return
}

func GetBatchEquipGetNum(logger fklog.FKLogI, userId uint64, equipIds []int32) (equipScoreMap map[int32]int32, err error) {
	equipScoreMap = make(map[int32]int32)
	key := fmt.Sprintf("maze:equip:get:num:%d", userId)
	args := make([]interface{}, 0, len(equipIds)+1)
	args = append(args, key)
	for _, field := range equipIds {
		args = append(args, field)
	}
	resp, err := redis.Int64s(gRedis.Do(context.TODO(), "HMGET", args...))
	if err != nil {
		if err == redis.ErrNil {
			logger.InfoWF("GetBatchEquipGetNum call HMGET empty", zap.Int32s("equipIds", equipIds))
			return equipScoreMap, nil
		}
		logger.ErrorWF("GetBatchEquipGetNum call HMGET failed", zap.Int32s("equipIds", equipIds), zap.Any("err", err))
		return equipScoreMap, err
	}
	if len(resp) != len(equipIds) {
		logger.WarnWF("GetBatchEquipGetNum count not match", zap.String("key", key), zap.Int32s("fields", equipIds), zap.Any("resp", resp))
		return
	}
	equipScoreMap = make(map[int32]int32)
	for i, field := range equipIds {
		equipScoreMap[field] = int32(resp[i])
	}
	logger.InfoWF("GetBatchEquipGetNum success", zap.String("key", key), zap.Any("equipScoreMap", equipScoreMap))
	return equipScoreMap, nil
}

func GetAllEquipGetNum(logger fklog.FKLogI, userId uint64) (map[int32]int32, error) {
	equipScoreMap := make(map[int32]int32)
	key := fmt.Sprintf("maze:equip:get:num:%d", userId)
	ret, err := redis.Int64Map(gRedis.Do(context.TODO(), "HGETALL", key))
	if err != nil {
		logger.ErrorWF("GetAllEquipGetNum redis op failed with ",
			zap.Error(err),
			zap.String("key", key))
		return equipScoreMap, err
	}
	for k, v := range ret {
		field := fkutil.ToInt32(k)
		degree := int32(v)
		equipScoreMap[field] = degree
	}

	logger.InfoWF("GetAllEquipGetNum success",
		zap.Any("equipScoreMap", equipScoreMap),
		zap.String("key", key))
	return equipScoreMap, nil
}
