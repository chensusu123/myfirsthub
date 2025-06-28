package allianceredis

import (
	"context"
	"fmt"
	globalredis "maze_game_server/io/redis"

	"github.com/redis/go-redis/v9"
)

func getAllianceInfoRedisKey(allianceID int32) string {
	return fmt.Sprintf("alliance:info:%d", allianceID)
}

func getAllianceIDRedisKey() string {
	return "alliance:id"
}

func SetAllianceInfo(allianceID int32, data []byte) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Set(context.TODO(), db.MakeSectionKey(getAllianceInfoRedisKey(allianceID)), data, 0).Err()
}

func GetAllianceInfo(allianceID int32) ([]byte, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	res, err := db.Get(context.TODO(), db.MakeSectionKey(getAllianceInfoRedisKey(allianceID))).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

func BatchGetAllianceInfo(allianceIDs []int32) ([][]byte, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	args := make([]string, len(allianceIDs))
	for i, allianceID := range allianceIDs {
		args[i] = db.MakeSectionKey(getAllianceInfoRedisKey(allianceID))
	}
	result, err := db.MGet(context.TODO(), args...).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	ret := make([][]byte, len(result))
	for i, v := range result {
		if v == nil {
			ret[i] = nil
		} else {
			ret[i] = v.([]byte)
		}
	}
	return ret, err
}

func GetAllianceID() int32 {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return 0
	}
	value, err := db.Incr(context.TODO(), db.MakeSectionKey(getAllianceIDRedisKey())).Result()
	if err == redis.Nil {
		return 0
	}
	return int32(value)
}
