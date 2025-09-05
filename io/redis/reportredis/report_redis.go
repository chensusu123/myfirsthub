package reportredis

import (
	"context"
	"fmt"
	globalredis "maze_game_server/io/redis"

	"github.com/go-redis/redis"
)

func getKey(userID uint64) string {
	return fmt.Sprintf("report:entertime:%d", userID)
}

// 获取玩家进入关卡的时间
func GetEnterTime(ctx context.Context, userID uint64) (enterTime uint64, err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return 0, err
	}

	key := db.MakeSectionKey(getKey(userID))
	enterTime, err = db.Get(ctx, key).Uint64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	return
}

func SetEnterTime(ctx context.Context, userID uint64, enterTime uint64) (err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}

	key := db.MakeSectionKey(getKey(userID))
	return db.Set(ctx, key, enterTime, 0).Err()
}
