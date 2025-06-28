package familyredis

import (
	"context"
	"fmt"
	globalredis "maze_game_server/io/redis"

	"github.com/redis/go-redis/v9"
)

func getFamilyInfoRedisKey(familyId int32) string {
	return fmt.Sprintf("family:info:%d", familyId)
}

// SetFamilyInfo 设置家族信息
func SetFamilyInfo(familyId int32, data []byte) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Set(context.TODO(), db.MakeSectionKey(getFamilyInfoRedisKey(familyId)), data, 0).Err()
}

// GetFamilyInfo 获取家族信息
func GetFamilyInfo(familyId int32) ([]byte, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	ret, err := db.Get(context.TODO(), db.MakeSectionKey(getFamilyInfoRedisKey(familyId))).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

// BatchGetFamilyInfo 批量获取家族信息
func BatchGetFamilyInfo(familyIds []int32) ([][]byte, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	args := make([]string, len(familyIds))
	for i, familyID := range familyIds {
		args[i] = db.MakeSectionKey(getFamilyInfoRedisKey(familyID))
	}
	result, err := db.MGet(context.TODO(), args...).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	ret := make([][]byte, len(result))
	for _, v := range result {
		if v != nil {
			if b, ok := v.(string); ok {
				ret = append(ret, []byte(b))
			}
		}
	}
	return ret, err
}

// DelFamilyInfo 删除家族信息
func DelFamilyInfo(familyID int32) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	return db.Del(context.TODO(), db.MakeSectionKey(getFamilyInfoRedisKey(familyID))).Err()
}
