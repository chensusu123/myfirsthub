package familyredis

import (
	"context"
	"fmt"
	globalredis "maze_game_server/io/redis"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func getUserFamilyIDRedisKey(userID uint64) string {
	return fmt.Sprintf("user:family:%d", userID)
}

func getUserLastLeaveFamilyTimeKey(userID uint64) string {
	return fmt.Sprintf("user:last:leave:family:%d", userID)
}

// 设置用户对应家族ID
func SetUserFamilyID(userID uint64, familyID int32) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Set(context.TODO(), db.MakeSectionKey(getUserFamilyIDRedisKey(userID)), familyID, 0).Err()
}

// 获取用户对应的家族ID
func GetUserFamilyID(userID uint64) (int32, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return 0, err
	}
	ret, err := db.Get(context.TODO(), db.MakeSectionKey(getUserFamilyIDRedisKey(userID))).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	familyID, err := strconv.ParseInt(ret, 10, 64)
	if err != nil {
		return 0, err
	}
	return int32(familyID), nil
}

func DelUserFamilyID(userID uint64) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Del(context.TODO(), db.MakeSectionKey(getUserFamilyIDRedisKey(userID))).Err()
}

// 设置用户上次离开家族时间
func SetUserLastLeaveFamilyTime(userID uint64, lastLeaveFamilyTime int64) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Set(context.TODO(), db.MakeSectionKey(getUserLastLeaveFamilyTimeKey(userID)), lastLeaveFamilyTime, 0).Err()
}

// 获取用户上次离开家族时间
func GetUserLastLeaveFamilyTime(userID uint64) (int64, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return 0, err
	}
	ret, err := db.Get(context.TODO(), db.MakeSectionKey(getUserLastLeaveFamilyTimeKey(userID))).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	lastLeaveFamilyTime, err := strconv.ParseInt(ret, 10, 64)
	if err != nil {
		return 0, err
	}
	return lastLeaveFamilyTime, nil
}
