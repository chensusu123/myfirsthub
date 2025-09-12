package barrierscorerewardmodel

import (
	"context"
	"fmt"
	"maze_game_server/io"
)

// 关卡积分掉落装备记录
func getKey(userId uint64, barrier int32) string {
	return fmt.Sprintf("score:reward:u:%d:barrier:%d", userId, barrier)
}

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

func NewBarrierScoreRewardModel(ctx context.Context, userID uint64, barrier int32) (*BarrierScoreRewardModel, error) {
	model := &BarrierScoreRewardModel{
		ItemList:  make([]*ItemRewardInfo, 0),
		EquipList: make([]*EquipRewardInfo, 0),
	}
	if err := model.load(ctx, userID, barrier); err != nil {
		return nil, err
	}
	return model, nil
}

func (b *BarrierScoreRewardModel) load(ctx context.Context, userID uint64, barrierId int32) (err error) {
	return io.LoadSvrData(ctx, getKey(userID, barrierId), b)
}

func (b *BarrierScoreRewardModel) Save(ctx context.Context, userID uint64, barrierId int32) (err error) {
	return io.SaveSvrData(ctx, getKey(userID, barrierId), b)
}

func (b *BarrierScoreRewardModel) Del(ctx context.Context, userID uint64, barrierId int32) (err error) {
	return io.DeleteSvrData(ctx, getKey(userID, barrierId))
}
