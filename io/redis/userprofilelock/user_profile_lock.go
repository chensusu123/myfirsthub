// @Author pangchenyang 2025/6/10 20:00:00
// @Desc: 
package userprofilelock

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"maze_game_server/io/redis/redisconfig"
	"context"
)

var (
	redisClient *redis.Client
	lockExpire  = 10 * time.Second // 锁的默认过期时间
)

const (
	profileLockKey = "maze:profile:init:lock:%d"
)

func GetKey(userID uint64) string {
	return fmt.Sprintf(profileLockKey, userID)
}

func init() {
	redisClient = redisconfig.GetClient("ProfileLockRedis")
}

// Lock 获取初始化锁
func Lock(userID uint64) (bool, error) {
	ctx := context.Background()
	key := GetKey(userID)
	return redisClient.SetNX(ctx, key, "1", lockExpire).Result()
}

// Unlock 释放初始化锁
func Unlock(userID uint64) error {
	ctx := context.Background()
	key := GetKey(userID)
	return redisClient.Del(ctx, key).Err()
}
