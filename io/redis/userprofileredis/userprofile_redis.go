// @Author pangchenyang 2025/6/9 22:02:00
// @Desc: 
package userprofileredis

import (
	"context"
	"encoding/json"
	"fmt"

	"maze_game_server/pb/common/UserProfile"

	"github.com/go-redis/redis/v8"
	"sync"
	"maze_game_server/usecase/redisconfig"
)

var db *redis.Client
var once sync.Once

const (
	userProfileCacheKey = "user:profile:%d" // 永不过期 todo 需要优化
)

// OpUserProfile 用户资料
type OpUserProfile struct {
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

// NewOpUserProfile 构造函数
func NewOpUserProfile(rdb *redis.Client) *OpUserProfile {
	return &OpUserProfile{rdb: rdb}
}

// SetProfile 设置用户资料
func SetProfile(userID uint64, profile *UserProfile.UserProfile) error {
	ctx := context.Background()
	key := GetKey(userID)

	jsonData, err := json.Marshal(profile) // todo 先用json存，后期考虑加密的话用proto打包
	if err != nil {
		return err
	}

	return db.Set(ctx, key, jsonData, 0).Err()
}

// GetProfile 获取用户资料
func GetProfile(userID uint64) (*UserProfile.UserProfile, error) {
	ctx := context.Background()
	key := GetKey(userID)

	data, err := db.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 资料不存在不返回错误
		}
		return nil, err
	}

	var profile UserProfile.UserProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// DeleteProfile 删除用户资料 // todo 后续监听用户换服、移民、注销等事件
func DeleteProfile(userID uint64) error {
	ctx := context.Background()
	key := GetKey(userID)
	return db.Del(ctx, key).Err()
}

// BatchGetProfile 批量获取用户资料
func BatchGetProfile(userIDs []uint64) ([]*UserProfile.UserProfile, error) {
	ctx := context.Background()

	keys := make([]string, len(userIDs))
	for i, userID := range userIDs {
		keys[i] = GetKey(userID)
	}

	values, err := db.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	result := make([]*UserProfile.UserProfile, len(userIDs))

	for _, value := range values {
		if value == nil {
			continue
		}

		var profile *UserProfile.UserProfile
		if err := json.Unmarshal([]byte(value.(string)), &profile); err != nil {
			return nil, err
		}
		result = append(result, profile)
	}

	return result, nil
}
