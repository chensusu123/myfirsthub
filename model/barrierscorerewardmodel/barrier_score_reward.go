package barrierscorerewardmodel

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/redis/barrierscorerewardredis"
	"maze_game_server/lib/serialize"
)

// 关卡积分掉落装备记录

type ItemRewardInfo struct {
	ItemId int32 `json:"item_id,omitempty"`
	Count  int64 `json:"count,omitempty"`
}
type EquipRewardInfo struct {
	EquipId int32 `json:"equip_id,omitempty"`
	Count   int32 `json:"count,omitempty"`
}
type BarrierScoreRewardModel struct {
	ItemList  []*ItemRewardInfo  `json:"item_list,omitempty"`
	EquipList []*EquipRewardInfo `json:"equip_list,omitempty"`
}

func NewBarrierScoreRewardModel(logger fklog.FKLogI, userID uint64, barrier int32) (*BarrierScoreRewardModel, error) {
	model := &BarrierScoreRewardModel{}
	if err := model.load(logger, userID, barrier); err != nil {
		return nil, err
	}
	return model, nil
}

func (b *BarrierScoreRewardModel) load(logger fklog.FKLogI, userID uint64, barrier int32) (err error) {
	bytes, err := barrierscorerewardredis.GetBarrierScoreReward(logger, userID, barrier)
	if err != nil {
		return err
	}
	if bytes == nil {
		b.ItemList = make([]*ItemRewardInfo, 0)
		b.EquipList = make([]*EquipRewardInfo, 0)
		return nil
	}
	err = serialize.Unmarshal(bytes, b)
	if err != nil {
		logger.ErrorWF("BarrierScoreReward load Unmarshal failed", zap.Error(err), zap.Int32("barrier", barrier))
		return err
	}
	return
}

func (b *BarrierScoreRewardModel) Save(logger fklog.FKLogI, userID uint64, barrier int32) (err error) {
	bytes, err := serialize.Marshal(b)
	if err != nil {
		logger.ErrorWF("BarrierScoreReward save Marshal failed", zap.Error(err), zap.Int32("barrier", barrier))
		return err
	}
	return barrierscorerewardredis.SetBarrierScoreReward(logger, userID, barrier, bytes)
}

func (b *BarrierScoreRewardModel) Del(logger fklog.FKLogI, userID uint64, barrier int32) (err error) {
	return barrierscorerewardredis.DelBarrierScoreReward(logger, userID, barrier)
}
