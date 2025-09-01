// @Author: ZhaoXiming 2025/3/24 15:35
// @Desc: 迷宫装备合成存储

package mazeequipmixdb

import (
	"context"
	"errors"
	"fmt"

	"maze_game_server/common/equipmix"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
)

/*
类型 HASH
FIELD	lv、cfg、idx
VALUE
*/

var db = &fkredis.FkRedis{}

func init() {
	_ = fkconfig.RegisterNameNode("mazeequipmixdb", 21731, db)
}

func getKey(uid uint64) string {
	return fmt.Sprintf("maze:equip:mix:%d", uid)
}

const (
	field_lv  = "lv"
	field_cfg = "cfg"
	field_idx = "idx"
)

func GetEquipMixData(ctx context.Context, uid uint64) (data *equipmix.MixData, err error) {
	k := getKey(uid)

	res, err := redis.Int64s(db.Do(context.Background(), "HMGET", k, field_lv, field_cfg, field_idx))
	if errors.Is(err, redis.ErrNil) {
		err = nil
	}
	if err != nil {
		return
	}
	if len(res) != 3 {
		err = errors.New("wrong data len")
		return
	}

	data = &equipmix.MixData{
		Lv:  int32(res[0]),
		Cfg: int32(res[1]),
		Idx: int(res[2]),
	}

	return
}

func SetEquipMixData(ctx context.Context, uid uint64, data *equipmix.MixData) (err error) {
	k := getKey(uid)
	_, err = db.Do(context.Background(), "HMSET", k, field_lv, data.Lv, field_cfg, data.Cfg, field_idx, data.Idx)
	return
}
