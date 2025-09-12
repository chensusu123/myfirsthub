package allianceredis

import (
	"context"
	"fmt"
	"strconv"

	globalredis "maze_game_server/io/redis"

	"github.com/redis/go-redis/v9"
)

func getFamilyToAllianceRedisKey(familyID int32) string {
	return fmt.Sprintf("family:to:alliance:%d", familyID)
}

func GetFamilyToAlliance(ctx context.Context, familyID int32) (int32, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return 0, err
	}
	value, err := db.Get(ctx, db.MakeSectionKey(getFamilyToAllianceRedisKey(familyID))).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	allianceID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(allianceID), nil
}

func SetFamilyToAlliance(ctx context.Context, familyID int32, allianceID int32) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Set(ctx, db.MakeSectionKey(getFamilyToAllianceRedisKey(familyID)), allianceID, 0).Err()
}

func DelFamilyToAlliance(ctx context.Context, familyID int32) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Del(ctx, db.MakeSectionKey(getFamilyToAllianceRedisKey(familyID))).Err()
}
