package barrierarearecordservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

type BarrierAreaRecordService interface {
	GetBarrierAreaRecord(logger fklog.FKLogI, userId uint64, stageId int32) (killMonsterNum int32, totalDamage int64, guidList []int64, err error)
	AddKillMonsterNum(logger fklog.FKLogI, userId uint64, stageId, areaId, areaIndex, monsterId, addVal int32, monsterGuid int64) (killMonsterNum int32, guidList []int64, err error)
	AddDamage(logger fklog.FKLogI, userId uint64, stageId, areaId, areaIndex int32, addVal int64) (totalDamage int64, err error)
	DelBarrierAreaRecord(logger fklog.FKLogI, userId uint64, barrierId, stageId int32) error
}

var GlobalBarrierAreaRecordService BarrierAreaRecordService

func init() {
	GlobalBarrierAreaRecordService = newBarrierAreaRecordService()
}

type service struct {
}

func newBarrierAreaRecordService() BarrierAreaRecordService {
	return &service{}
}
