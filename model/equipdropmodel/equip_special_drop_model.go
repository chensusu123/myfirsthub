package equipdropmodel

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/redis/mazeequipspecialdropredis"
	"maze_game_server/lib/serialize"
)

type EquipSpecialDropModel struct {
	OrderId          int32 `json:"order_id"`           // MazeEquDropV8ConfigRow表序列id
	SpecialDropIndex int32 `json:"special_drop_index"` // 已经掉落的特殊掉落下标
	EquipPoints      int32 `json:"equip_points"`       // 装备分
}

func NewEquipSpecialDropModel(logger fklog.FKLogI, userID uint64) (*EquipSpecialDropModel, error) {
	info := &EquipSpecialDropModel{}
	if err := info.load(logger, userID); err != nil {
		return nil, err
	}
	return info, nil
}

func (info *EquipSpecialDropModel) load(logger fklog.FKLogI, userID uint64) (err error) {
	bytes, err := mazeequipspecialdropredis.GetMazeEquipSpecialDropInfo(logger, userID)
	if err != nil {
		return err
	}
	if bytes == nil {
		return nil
	}
	err = serialize.Unmarshal(bytes, info)
	if err != nil {
		logger.ErrorWF("EquipSpecialDropModel load Unmarshal failed", zap.Error(err), zap.Uint64("userID", userID))
		return err
	}
	return
}

func (info *EquipSpecialDropModel) Save(logger fklog.FKLogI, userID uint64) (err error) {
	bytes, err := serialize.Marshal(info)
	if err != nil {
		logger.ErrorWF("EquipSpecialDropModel save Marshal failed", zap.Error(err), zap.Uint64("userID", userID))
		return err
	}
	return mazeequipspecialdropredis.SetMazeEquipSpecialDropInfo(logger, userID, bytes)
}

func (info *EquipSpecialDropModel) Del(logger fklog.FKLogI, userID uint64) (err error) {
	return mazeequipspecialdropredis.GMDel(logger, userID)
}
