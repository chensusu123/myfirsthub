package familyredis

import (
	"context"

	globalredis "maze_game_server/io/redis"

	"github.com/redis/go-redis/v9"
)

// getFamilyListKey 获取家族列表key
func getFamilyListKey() string {
	return "familylist"
}

// GetFamilyIDs 查询家族列表信息
func GetFamilyList(ctx context.Context) ([]byte, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	ret, err := db.Get(ctx, db.MakeSectionKey(getFamilyListKey())).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

// SetFamilyIDs 保存家族列表信息
func SetFamilyIDs(ctx context.Context, data []byte) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Set(ctx, db.MakeSectionKey(getFamilyListKey()), data, 0).Err()
}

// DelFamilyInfo 删除家族列表信息
func DelFamilyList(ctx context.Context) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Del(ctx, db.MakeSectionKey(getFamilyListKey())).Err()
}
