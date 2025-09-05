// @Author pangchenyang 2025/6/9 22:02:00
// @Desc:
package userprofileredis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/database/nanoredis"
)

var GlobalUserProfileRedis *UserProfileRedis

type UserProfileRedis struct {
	*nanoredis.NanoRedis
}

func NewRedisDemo(serviceName string, name string) *UserProfileRedis {
	GlobalUserProfileRedis.NanoRedis = nanoredis.NewNanoRedis(serviceName, name)
	return GlobalUserProfileRedis
}

func (r *UserProfileRedis) getKey(userID uint64) string {
	return fmt.Sprintf("user:profile:%d", userID)
}

// SetProfile 设置用户资料
func (r *UserProfileRedis) SetProfile(ctx context.Context, userID uint64, data []byte) error {
	db, err := r.GetDB()
	if err != nil {
		return err
	}
	return db.Set(ctx, r.getKey(userID), data, 0).Err()
}

// GetProfile 获取用户资料
func (r *UserProfileRedis) GetProfile(ctx context.Context, userID uint64) ([]byte, error) {
	db, err := r.GetDB()
	if err != nil {
		return nil, err
	}
	ret, err := db.Get(ctx, r.getKey(userID)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return ret, err
}

// BatchGetProfile 批量获取用户资料
func (r *UserProfileRedis) BatchGetProfile(ctx context.Context, userIDs []uint64) ([][]byte, error) {
	db, err := r.GetDB()
	if err != nil {
		return nil, err
	}
	args := make([]string, len(userIDs))
	for i, userID := range userIDs {
		args[i] = r.getKey(userID)
	}
	result, err := db.MGet(ctx, args...).Result()
	if err != redis.Nil {
		return nil, err
	}
	ret := make([][]byte, len(result))
	for _, v := range result {
		if v != nil {
			if b, ok := v.(string); ok {
				result = append(result, []byte(b))
			}
		}
	}
	return ret, err
}

// DeleteProfile 删除用户资料
func (r *UserProfileRedis) DelProfile(ctx context.Context, userID uint64) error {
	db, err := r.GetDB()
	if err != nil {
		return err
	}
	return db.Del(ctx, r.getKey(userID)).Err()
}
