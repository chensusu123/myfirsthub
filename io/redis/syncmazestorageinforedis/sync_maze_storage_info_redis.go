package syncmazestorageinforedis

import (
	"context"
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"google.golang.org/protobuf/proto"
	"maze_game_server/pb/common/MazeGame"
	"strconv"
)

const (
	StorageInfo_RoleItemData    = "roleItemData"    //角色物品列表
	StorageInfo_PassLevel       = "passLevel"       //角色通关值
	StorageInfo_RolePos         = "rolePos"         //角色位置
	StorageInfo_StorageItemInfo = "storageItemInfo" //机关列表
)

var gRedis = &fkredis.FkRedis{}
var redisKey = "u:%d:stage:%d:storage:info"

func init() {
	fkconfig.RegisterNameNode("syncmazestorageinforedis", 21727, gRedis)
}

func DelSyncMazeStorageInfo(userId uint64, barrierId int32) error {
	key := fmt.Sprintf(redisKey, userId, barrierId)
	_, err := gRedis.Do(context.TODO(), "DEL", key)
	return err
}

func DelSyncMazeStorageInfoByKey(userId uint64, barrierId int32, field string) error {
	key := fmt.Sprintf(redisKey, userId, barrierId)
	_, err := gRedis.Do(context.TODO(), "HDEL", key, field)
	return err
}

func SaveSyncMazeStorageInfo(userId uint64, barrierId int32, field string, data interface{}) error {
	key := fmt.Sprintf(redisKey, userId, barrierId)
	_, err := gRedis.Do(context.TODO(), "HSET", key, field, data)
	return err
}

func GetSyncMazeStorageInfo(userId uint64, barrierId int32) (info *MazeGame.MazeStorageInfo, err error) {
	key := fmt.Sprintf(redisKey, userId, barrierId)
	res, err := redis.StringMap(gRedis.Do(context.TODO(), "hgetall", key))
	if err != nil {
		return
	}
	if len(res) == 0 {
		return
	}

	info = &MazeGame.MazeStorageInfo{}
	for k, v := range res {
		if k == StorageInfo_RoleItemData {
			info.RoleItemData = proto.String(v)
		}
		if k == StorageInfo_PassLevel {
			a, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				continue
			}
			info.PassLevel = proto.Int64(a)
		}
		if k == StorageInfo_RolePos {
			info.RolePos = proto.String(v)
		}
		if k == StorageInfo_StorageItemInfo {
			info.StorageItemInfo = proto.String(v)
		}
	}

	return
}
