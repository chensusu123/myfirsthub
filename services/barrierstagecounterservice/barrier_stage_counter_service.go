package barrierstagecounterservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

type BarrierStageCounterService interface {
	GetBarrierStageCounter(logger fklog.FKLogI, userId uint64, barrierId int32) (killMonsterNum int32, totalDamage, totalExp int64, guidList []int64, err error)
	AddKillMonsterNum(logger fklog.FKLogI, userId uint64, barrierId, stageId, monsterId, addVal int32, monsterGuid int64) (killMonsterNum int32, guidList []int64, err error)
	AddDamage(logger fklog.FKLogI, userId uint64, barrierId, stageId int32, addVal int64) (totalDamage int64, err error)
	DelBarrierStageCounter(logger fklog.FKLogI, userId uint64, barrierId, stageId int32) error
}

var GlobalBarrierStageCounterService BarrierStageCounterService

func init() {
	GlobalBarrierStageCounterService = newBarrierStageCounterService()
}

type service struct {
}

func newBarrierStageCounterService() BarrierStageCounterService {
	return &service{}
}
