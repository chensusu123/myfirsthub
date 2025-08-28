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
func GetFamilyList() ([]byte, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	ret, err := db.Get(context.TODO(), db.MakeSectionKey(getFamilyListKey())).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

// SetFamilyIDs 保存家族列表信息
func SetFamilyIDs(data []byte) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Set(context.TODO(), db.MakeSectionKey(getFamilyListKey()), data, 0).Err()
}

// DelFamilyInfo 删除家族列表信息
func DelFamilyList() error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Del(context.TODO(), db.MakeSectionKey(getFamilyListKey())).Err()
}
