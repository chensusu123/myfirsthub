package barrierguiditemredis

import (
	"context"
	"fmt"
	"strconv"

	globalredis "maze_game_server/io/redis"
)

func getKey(userID uint64, barrierID int32) string {
	return fmt.Sprintf("drop:guid:user:%d:barrier:%d", userID, barrierID)
}

func GetNowGuid(ctx context.Context, userID uint64, barrierID int32) (int64, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return 0, err
	}

	tmp, err := db.Get(ctx, db.MakeSectionKey(getKey(userID, barrierID))).Result()
	if err != nil {
		return 0, err
	}

	guid, err := strconv.ParseInt(tmp, 10, 64)
	if err != nil {
		return 0, err
	}
	return guid, nil
}

func IncrNowGuid(ctx context.Context, userID uint64, barrierID int32) (int64, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return 0, err
	}

	return db.Incr(ctx, db.MakeSectionKey(getKey(userID, barrierID))).Result()
}
