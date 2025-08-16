package barriersavedatamodel

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/redis/barriersavedataredis"
	"maze_game_server/lib/serialize"
)

// 关卡存档数据
type BarrierSaveDataModel struct {
	StageId      int32   `json:"stage_id,omitempty"`
	RescueValue  int32   `json:"rescue_value,omitempty"`
	BossPower    int32   `json:"boss_power,omitempty"`
	BossProgress float32 `json:"boss_progress,omitempty"`
}

func NewBarrierSaveDataModel(logger fklog.FKLogI, userID uint64, barrierId int32, isLoad bool) (*BarrierSaveDataModel, error) {
	model := &BarrierSaveDataModel{}
	if isLoad {
		if err := model.load(logger, userID, barrierId); err != nil {
			return nil, err
		}
	}
	return model, nil
}

func (b *BarrierSaveDataModel) load(logger fklog.FKLogI, userID uint64, barrierId int32) (err error) {
	bytes, err := barriersavedataredis.GetBarrierSaveData(logger, userID, barrierId)
	if err != nil {
		return err
	}
	if bytes == nil {
		return nil
	}
	err = serialize.Unmarshal(bytes, b)
	if err != nil {
		logger.ErrorWF("BarrierSaveDataModel load Unmarshal err", zap.Error(err), zap.Int32("barrierId", barrierId))
		return err
	}
	return
}

func (b *BarrierSaveDataModel) Save(logger fklog.FKLogI, userID uint64, barrierId int32) (err error) {
	bytes, err := serialize.Marshal(b)
	if err != nil {
		logger.ErrorWF("BarrierSaveDataModel save Marshal err", zap.Error(err), zap.Int32("barrierId", barrierId))
		return err
	}
	return barriersavedataredis.SetBarrierSaveData(logger, userID, barrierId, bytes)
}

func (b *BarrierSaveDataModel) Del(logger fklog.FKLogI, userID uint64, barrierId int32) (err error) {
	return barriersavedataredis.DelBarrierSaveData(logger, userID, barrierId)
}
