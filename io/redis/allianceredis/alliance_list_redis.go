package allianceredis

import (
	"context"

	globalredis "maze_game_server/io/redis"

	"github.com/redis/go-redis/v9"
)

func getAllianceListRedisKey() string {
	return "alliance:list"
}

func GetAllianceList(ctx context.Context) ([]byte, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	res, err := db.Get(ctx, db.MakeSectionKey(getAllianceListRedisKey())).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

func SetAllianceList(ctx context.Context, data []byte) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Set(ctx, db.MakeSectionKey(getAllianceListRedisKey()), data, 0).Err()
}

func DelAllianceList(ctx context.Context) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Del(ctx, db.MakeSectionKey(getAllianceListRedisKey())).Err()
}
