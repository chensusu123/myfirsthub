package allianceredis

import (
	"context"
	globalredis "maze_game_server/io/redis"

	"github.com/redis/go-redis/v9"
)

func getAllianceListRedisKey() string {
	return "alliance:list"
}

func GetAllianceList() ([]byte, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	res, err := db.Get(context.TODO(), db.MakeSectionKey(getAllianceListRedisKey())).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

func SetAllianceList(data []byte) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Set(context.TODO(), db.MakeSectionKey(getAllianceListRedisKey()), data, 0).Err()
}

func DelAllianceList() error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Del(context.TODO(), db.MakeSectionKey(getAllianceListRedisKey())).Err()
}
