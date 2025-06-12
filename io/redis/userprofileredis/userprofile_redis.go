// @Author pangchenyang 2025/6/9 22:02:00
// @Desc: 
package userprofileredis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"maze_game_server/pb/common/UserProfile"

	"github.com/go-redis/redis/v8"
	"sync"
	"maze_game_server/usecase/redisconfig"
)

var db *redis.Client
var once sync.Once

const (
	userProfileCacheKey = "user:profile:%d"
	cacheExpiration     = 24 * time.Hour * 30 // 缓存30天
)

// UserProfileCache 用户资料缓存操作结构体
type UserProfileCache struct {
	rdb *redis.Client
}

func GetKey(userID uint64) string {
	once.Do(func() {
		var err error
		db, err = redisconfig.GetRedisService().GetClient()
		if err != nil {
			panic(err)
		}
	})
	return fmt.Sprintf(userProfileCacheKey, userID)
}

// NewUserProfileCache 构造函数
func NewUserProfileCache(rdb *redis.Client) *UserProfileCache {
	return &UserProfileCache{rdb: rdb}
}

// SetCache 设置用户资料缓存
func SetCache(userID uint64, profile *UserProfile.UserProfile) error {
	ctx := context.Background()
	key := GetKey(userID)

	jsonData, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	return db.Set(ctx, key, jsonData, cacheExpiration).Err()
}

// GetCache 获取用户资料缓存
func GetCache(userID uint64) (*UserProfile.UserProfile, error) {
	ctx := context.Background()
	key := GetKey(userID)

	data, err := db.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存不存在不返回错误
		}
		return nil, err
	}

	var profile UserProfile.UserProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// DeleteCache 删除用户资料缓存
func DeleteCache(userID uint64) error {
	ctx := context.Background()
	key := GetKey(userID)
	return db.Del(ctx, key).Err()
}

// BatchGetCache 批量获取用户资料缓存
func BatchGetCache(userIDs []uint64) (map[uint64]*UserProfile.UserProfile, error) {
	ctx := context.Background()
	pipe := db.Pipeline()

	cmds := make([]*redis.StringCmd, len(userIDs))
	for i, userID := range userIDs {
		key := GetKey(userID)
		cmds[i] = pipe.Get(ctx, key)
	}

	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		return nil, err
	}

	result := make(map[uint64]*UserProfile.UserProfile)
	for i, cmd := range cmds {
		data, err := cmd.Bytes()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			return nil, err
		}

		var profile UserProfile.UserProfile
		if err := json.Unmarshal(data, &profile); err != nil {
			return nil, err
		}
		result[userIDs[i]] = &profile
	}

	return result, nil
}
