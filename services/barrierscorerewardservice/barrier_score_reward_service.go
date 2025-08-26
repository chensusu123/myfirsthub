package barrierscorerewardservice

import (
	"context"
)

type BarrierScoreRewardService interface {
	// 保存积分装备奖励信息
	SaveBarrierScoreRewardEquip(ctx context.Context, userId uint64, barrier int32, equipList map[int32]int32) (err error)
	// 保存积分道具奖励信息
	SaveBarrierScoreRewardItem(ctx context.Context, userId uint64, barrier int32, itemList map[int32]int64) (err error)
	// 保存积分装备和道具信息
	SaveBarrierScoreReward(ctx context.Context, userId uint64, barrier int32, equipList map[int32]int32, itemList map[int32]int64) (err error)
	// 获取当前积分装备和道具信息
	GetBarrierScoreReward(ctx context.Context, userId uint64, barrier int32) (equipList map[int32]int32, itemList map[int32]int64, err error)
	// 删除关卡积分奖励记录
	DelBarrierScoreRewardItem(ctx context.Context, userId uint64, barrier int32) (err error)
}

var GlobalScoreRewardService BarrierScoreRewardService

func init() {
	GlobalScoreRewardService = newBarrierScoreRewardService()
}

type service struct {
}

func newBarrierScoreRewardService() BarrierScoreRewardService {
	return &service{}
}
