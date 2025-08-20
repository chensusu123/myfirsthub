package barriersavedatamodel

import (
	"context"
	"fmt"
	"maze_game_server/io"
)

func getKey(userId uint64, barrier int32) string {
	return fmt.Sprintf("save:u:%d:barrier:%d", userId, barrier)
}

type RescueItemInfo struct {
	RescueItemId int32 `json:"rescue_item_id,omitempty"`
	MapConfigId  int32 `json:"map_config_id,omitempty"`
}

// 关卡存档数据
type BarrierSaveDataModel struct {
	StageId      int32             `json:"stage_id,omitempty"`
	RescueValue  int32             `json:"rescue_value,omitempty"`
	BossPower    int32             `json:"boss_power,omitempty"`
	BossProgress float32           `json:"boss_progress,omitempty"`
	RescueItems  []*RescueItemInfo `json:"rescue_items,omitempty"`
}

func NewBarrierSaveDataModel(ctx context.Context, userID uint64, barrierId int32, isLoad bool) (*BarrierSaveDataModel, error) {
	model := &BarrierSaveDataModel{}
	if isLoad {
		if err := model.load(ctx, userID, barrierId); err != nil {
			return nil, err
		}
	}
	return model, nil
}

func (b *BarrierSaveDataModel) load(ctx context.Context, userID uint64, barrierId int32) (err error) {
	return io.LoadSvrData(ctx, getKey(userID, barrierId), b)
}

func (b *BarrierSaveDataModel) Save(ctx context.Context, userID uint64, barrierId int32) (err error) {
	return io.SaveSvrData(ctx, getKey(userID, barrierId), b)
}

func (b *BarrierSaveDataModel) Del(ctx context.Context, userID uint64, barrierId int32) (err error) {
	return io.DeleteSvrData(ctx, getKey(userID, barrierId))
}
