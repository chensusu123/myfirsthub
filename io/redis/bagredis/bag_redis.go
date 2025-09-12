package bagredis

import (
	"context"
	"fmt"
	globalredis "maze_game_server/io/redis"
	"strconv"
)

func getKey(userId uint64) string {
	return fmt.Sprintf("bag:u:%d", userId)
}

func Incr(ctx context.Context, userId uint64, itemId int32, count int64) (curCount int64, err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return 0, err
	}
	key := db.MakeSectionKey(getKey(userId))

	curCount, err = db.HIncrBy(ctx, key, strconv.FormatInt(int64(itemId), 10), count).Result()
	return
}

func BatchGet(ctx context.Context, userId uint64, itemIds []string) ([]interface{}, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	key := db.MakeSectionKey(getKey(userId))

	res, err := db.HMGet(ctx, key, itemIds...).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func BatchDel(ctx context.Context, userId uint64, itemIds []string) (err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(getKey(userId))

	_, err = db.HDel(ctx, key, itemIds...).Result()
	if err != nil {
		return err
	}
	return nil
}

func GetAll(ctx context.Context, userId uint64) (itemMap map[string]string, err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	key := db.MakeSectionKey(getKey(userId))

	mpString, err := db.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return mpString, nil
}

func Del(ctx context.Context, userId uint64) (err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(getKey(userId))

	err = db.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
