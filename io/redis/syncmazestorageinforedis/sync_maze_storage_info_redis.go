package syncmazestorageinforedis

import (
	"context"
	"fmt"
	"strconv"

	"maze_game_server/pb/common/MazeGame"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"google.golang.org/protobuf/proto"
)

const (
	StorageInfo_RoleItemData    = "roleItemData"    // 角色物品列表
	StorageInfo_PassLevel       = "passLevel"       // 角色通关值
	StorageInfo_RolePos         = "rolePos"         // 角色位置
	StorageInfo_StorageItemInfo = "storageItemInfo" // 机关列表
	StorageInfo_StageLevel      = "stageLevel"      // 阶段等级
	StorageInfo_MonsterAreaInfo = "monsterAreaInfo" // 已经打过的刷怪区域
)

var (
	gRedis   = &fkredis.FkRedis{}
	redisKey = "u:%d:stage:%d:storage:info"
)

func init() {
	fkconfig.RegisterNameNode("syncmazestorageinforedis", 21727, gRedis)
}

func DelSyncMazeStorageInfo(ctx context.Context, userId uint64, barrierId int32) error {
	key := fmt.Sprintf(redisKey, userId, barrierId)
	_, err := gRedis.Do(ctx, "DEL", key)
	return err
}

func DelSyncMazeStorageInfoByKey(ctx context.Context, userId uint64, barrierId int32, field string) error {
	key := fmt.Sprintf(redisKey, userId, barrierId)
	_, err := gRedis.Do(ctx, "HDEL", key, field)
	return err
}

func SaveSyncMazeStorageInfo(ctx context.Context, userId uint64, barrierId int32, field string, data interface{}) error {
	key := fmt.Sprintf(redisKey, userId, barrierId)
	_, err := gRedis.Do(ctx, "HSET", key, field, data)
	return err
}

func GetSyncMazeStorageInfo(ctx context.Context, userId uint64, barrierId int32) (info *MazeGame.MazeStorageInfo, err error) {
	key := fmt.Sprintf(redisKey, userId, barrierId)
	res, err := redis.StringMap(gRedis.Do(ctx, "hgetall", key))
	if err != nil {
		return
	}
	if len(res) == 0 {
		return
	}

	info = &MazeGame.MazeStorageInfo{}
	for k, v := range res {
		switch k {
		case StorageInfo_RoleItemData:
			info.RoleItemData = proto.String(v)
		case StorageInfo_PassLevel:
			a, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				continue
			}
			info.PassLevel = proto.Int64(a)
		case StorageInfo_RolePos:
			info.RolePos = proto.String(v)
		case StorageInfo_StorageItemInfo:
			info.StorageItemInfo = proto.String(v)
		case StorageInfo_StageLevel:
			a, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				continue
			}
			info.StageLevel = proto.Int32(int32(a))
		case StorageInfo_MonsterAreaInfo:
			info.MonsterAreaInfo = proto.String(v)
		}
	}

	return
}
