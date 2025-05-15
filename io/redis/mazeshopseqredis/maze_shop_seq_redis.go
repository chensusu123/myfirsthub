package mazeshopseqredis

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
)

var gRedis = &fkredis.FkRedis{}

func init() {
	// maze:shop:seq:%d
	fkconfig.RegisterNameNode("mazeshopseqredis", 21579, gRedis)
}

// 人偶商城信息
type MazeShopInfo struct {
	SeqId       int32           `json:"seq_id"`        // 序列id
	BackId      int32           `json:"back_id"`       //备用id
	CurSeqIndex int32           `json:"cur_index"`     //当前位置
	BackIndex   int32           `json:"back_index"`    //备用位置
	TotalCount  int32           `json:"total_count"`   //累计数量
	ShopSlotNum map[int32]int32 `json:"shop_slot_num"` //槽位购买数量
	EquipPoints int32           `json:"equip_points"`  // 装备分
}

func GetMazeShopInfo(logger fklog.FKLogI, userId uint64, level int32) (mazeShopInfo *MazeShopInfo, err error) {
	key := fmt.Sprintf("maze:shop:seq:%d", userId)
	res, err := redis.Bytes(gRedis.Do(context.TODO(), "hget", key, level))
	if err == redis.ErrNil {
		err = nil
		return nil, err
	}
	if err != nil {
		logger.ErrorWF("GetMazeShopSeqInfo get count failed with", zap.Error(err), zap.Int32("level", level))
		return nil, err
	}
	mazeShopInfo = &MazeShopInfo{}
	err = json.Unmarshal(res, mazeShopInfo)
	if err != nil {
		logger.ErrorWF("GetMazeShopSeqInfo Unmarshal fail", zap.Error(err), zap.String("key", key))
		return nil, err
	}
	logger.InfoWF("GetMazeShopSeqInfo succ", zap.Any("level", level), zap.Any("mazeShopInfo", mazeShopInfo), zap.String("key", key))
	return mazeShopInfo, nil
}

func SetMazeShopInfo(logger fklog.FKLogI, userId uint64, level int32, mazeShopInfo *MazeShopInfo) (err error) {
	key := fmt.Sprintf("maze:shop:seq:%d", userId)
	data, err := json.Marshal(mazeShopInfo)
	if err != nil {
		logger.ErrorWF("SetRideAttrInfo marshal fail", zap.Error(err))
		return
	}

	_, err = gRedis.Do(context.TODO(), "hset", key, level, data)
	if err != nil {
		logger.ErrorWF("SetMazeShopInfo fail", zap.Error(err), zap.Any("mazeShopInfo", mazeShopInfo), zap.String("key", key))
		return err
	}
	logger.InfoWF("SetMazeShopInfo succ", zap.Any("level", level), zap.Any("mazeShopInfo", mazeShopInfo), zap.String("key", key))
	return
}

// gm删除
func GMDel(logger fklog.FKLogI, userId uint64) (err error) {
	key := fmt.Sprintf("maze:shop:seq:%d", userId)
	_, err = redis.Int(gRedis.Do(context.TODO(), "DEL", key))
	if err != nil {
		logger.ErrorWF("GMDel fail", zap.String("key", key), zap.Error(err))
		return
	}
	return
}
