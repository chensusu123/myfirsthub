package barriersavedataservice

import (
	"context"
	"maze_game_server/model/barriersavedatamodel"
)

// BarrierSaveDataService 关卡存档service
type BarrierSaveDataService interface {
	// 保存关卡存档
	SaveBarrierData(ctx context.Context, userId uint64, barrier, stageId, rescueValue, bossPower int32, bossProgress float32, rescueItems []*barriersavedatamodel.RescueItemInfo) error
	// 获取关卡存档
	GetBarrierSaveData(ctx context.Context, userId uint64, barrier int32) (*barriersavedatamodel.BarrierSaveDataModel, error)
	// 删除关卡存档
	DelBarrierSaveData(ctx context.Context, userId uint64, barrier int32) error
	// 获取通关值
	GetPassValue(ctx context.Context, userId uint64, barrier int32) (int64, error)
}

var GlobalBarrierSaveDataService BarrierSaveDataService

func init() {
	GlobalBarrierSaveDataService = newBarrierSaveDataService()
}

type service struct {
}

func newBarrierSaveDataService() BarrierSaveDataService {
	return &service{}
}
