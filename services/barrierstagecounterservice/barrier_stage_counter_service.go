package barrierstagecounterservice

import (
	"context"
)

type BarrierStageCounterService interface {
	// 获取关卡记录数据
	GetBarrierStageCounter(ctx context.Context, userId uint64, barrierId int32) (killMonsterNum int32, totalDamage, totalExp int64, guidList []int64, err error)
	// 增加杀怪数量
	AddKillMonsterNum(ctx context.Context, userId uint64, barrierId, stageId, areaId, areaIndex, addVal, monsterId int32, monsterGuid int64, curHp int64,
		maxHp int64, monsterPos string) (killMonsterNum int32, guidList []int64, err error)
	// 增加杀怪伤害值
	AddDamage(ctx context.Context, userId uint64, barrierId, stageId int32, addVal int64) (totalDamage int64, err error)
	DelBarrierStageCounter(ctx context.Context, userId uint64, barrierId, stageId int32) error
	DelBarrierStageCounterOnPass(ctx context.Context, userId uint64, barrierId int32) error
	// 检查当前上报怪物是否合法 怪物是否重复 上报区域是否正确 只能读取检查 不能修改model
	CheckMonsterInvalid(ctx context.Context, userID uint64, barrierID int32, stageID int32, areaID int32, areaIndex int32, monsterGuid int64) (needCheck bool, err error)
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
