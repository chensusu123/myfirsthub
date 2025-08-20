/*
@Author: xiaobo
@Date: 2025/3/24 11:46
@Description:
*/

package moneyredis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	globalredis "maze_game_server/io/redis"
	"strconv"
)

func getKey(uid uint64) string {
	return fmt.Sprintf("money:u:%d", uid)
}

func GetAll(ctx context.Context, userId uint64) (moneyMap map[string]string, err error) {
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

func Get(ctx context.Context, userId uint64, moneyId int32) (moneyCount int64, err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return 0, err
	}
	key := db.MakeSectionKey(getKey(userId))

	value, err := db.HGet(ctx, key, strconv.FormatInt(int64(moneyId), 10)).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	return value, nil
}

func Set(ctx context.Context, userId uint64, moneyId int32, moneyCount int64) (err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(getKey(userId))
	_, err = db.HSet(ctx, key, moneyId, moneyCount).Result()
	if err != nil {
		return err
	}
	return nil
}

func Incr(ctx context.Context, userId uint64, moneyId int32, moneyCount int64) (newCount int64, err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return 0, err
	}
	key := db.MakeSectionKey(getKey(userId))

	curCount, err := db.HIncrBy(ctx, key, strconv.FormatInt(int64(moneyId), 10), moneyCount).Result()
	return curCount, err
}

func Sub(ctx context.Context, userId uint64, moneyId int32, moneyCount int64) (err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(getKey(userId))

	_, err = db.HIncrBy(ctx, key, strconv.FormatInt(int64(moneyId), 10), 0-moneyCount).Result()
	return err
}

// 删除所有数据
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

func BatchSet(ctx context.Context, userId uint64, itemMap map[int32]int64) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(getKey(userId))
	param := make([]interface{}, 0)
	for k, v := range itemMap {
		param = append(param, k)
		param = append(param, v)
	}
	_, err = db.HMSet(ctx, key, param...).Result()
	if err != nil {
		return err
	}
	return nil
}
