package awardservice

import (
	"context"
)

type AwardService interface {
	// 获取扫荡奖励 经验外部处理 realItemMap 失败新增物品奖励 showItemMap 总物品奖励 realEquipMap 失败新增装备奖励 showEquipMap 总装备奖励 expCount 扫荡能获取到的经验
	GetBarrierDeathAward(ctx context.Context, userId uint64, barrier int32) (realItemMap map[int32]int64, showItemMap map[int32]int64,
		realEquipMap map[int32]int32, showEquipMap map[int32]int32, expCount int64, err error)
}

var GlobalAwardService AwardService

func init() {
	GlobalAwardService = newAwardService()
}

type service struct {
}

func newAwardService() AwardService {
	return &service{}
}
