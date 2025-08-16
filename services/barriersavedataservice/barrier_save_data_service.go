package barriersavedataservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"maze_game_server/model/barriersavedatamodel"
)

// BarrierSaveDataService 关卡存档service
type BarrierSaveDataService interface {
	// 保存关卡存档
	SaveBarrierData(logger fklog.FKLogI, userId uint64, barrier, stageId, rescueValue, bossPower int32) error
	// 获取关卡存档
	GetBarrierSaveData(logger fklog.FKLogI, userId uint64, barrier int32) (*barriersavedatamodel.BarrierSaveDataModel, error)
	// 删除关卡存档
	DelBarrierSaveData(logger fklog.FKLogI, userId uint64, barrier int32) error
	// 获取通关值
	GetPassValue(logger fklog.FKLogI, userId uint64, barrier int32) (int64, error)
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
