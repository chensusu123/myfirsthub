package familyredis

import (
	"context"
	"fmt"
	"strconv"

	globalredis "maze_game_server/io/redis"

	"github.com/redis/go-redis/v9"
)

func getUserFamilyIDRedisKey(userID uint64) string {
	return fmt.Sprintf("user:family:%d", userID)
}

func getUserLastLeaveFamilyTimeKey(userID uint64) string {
	return fmt.Sprintf("user:last:leave:family:%d", userID)
}

// 设置用户对应家族ID
func SetUserFamilyID(ctx context.Context, userID uint64, familyID int32) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Set(ctx, db.MakeSectionKey(getUserFamilyIDRedisKey(userID)), familyID, 0).Err()
}

// 获取用户对应的家族ID
func GetUserFamilyID(ctx context.Context, userID uint64) (int32, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return 0, err
	}
	ret, err := db.Get(ctx, db.MakeSectionKey(getUserFamilyIDRedisKey(userID))).Result()
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

func DelUserFamilyID(ctx context.Context, userID uint64) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Del(ctx, db.MakeSectionKey(getUserFamilyIDRedisKey(userID))).Err()
}

// 设置用户上次离开家族时间
func SetUserLastLeaveFamilyTime(ctx context.Context, userID uint64, lastLeaveFamilyTime int64) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Set(ctx, db.MakeSectionKey(getUserLastLeaveFamilyTimeKey(userID)), lastLeaveFamilyTime, 0).Err()
}

// 获取用户上次离开家族时间
func GetUserLastLeaveFamilyTime(ctx context.Context, userID uint64) (int64, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return 0, err
	}
	ret, err := db.Get(ctx, db.MakeSectionKey(getUserLastLeaveFamilyTimeKey(userID))).Result()
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
