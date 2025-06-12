// @Author pangchenyang 2025/6/10 20:00:00
// @Desc: 
package userprofilelock

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"maze_game_server/io/redis/redisconfig"
	"context"
	"sync"
)

var (
	db         *redis.Client
	lockExpire = 10 * time.Second // 锁的默认过期时间
	once       sync.Once
)

const (
	profileLockKey = "maze:profile:init:lock:%d"
)

func GetKey(userID uint64) string {
	once.Do(func() {
		db = redisconfig.GetClient(redisconfig.DefaultRedisName)
	})
	return fmt.Sprintf(profileLockKey, userID)
}

// func init() {
// 	redisClient = redisconfig.GetClient("ProfileLockRedis")
// }

// Lock 获取初始化锁
func Lock(userID uint64) (bool, error) {
	ctx := context.Background()
	key := GetKey(userID)
	return db.SetNX(ctx, key, "1", lockExpire).Result()
}

// Unlock 释放初始化锁
func Unlock(userID uint64) error {
	ctx := context.Background()
	key := GetKey(userID)
	return db.Del(ctx, key).Err()
}
