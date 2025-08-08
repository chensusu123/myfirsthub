package barrierscorerewardservice

import "gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"

type BarrierScoreRewardService interface {
	// 保存积分装备奖励信息
	SaveBarrierScoreRewardEquip(logger fklog.FKLogI, userId uint64, barrier int32, equipList map[int32]int32) (err error)
	// 保存积分道具奖励信息
	SaveBarrierScoreRewardItem(logger fklog.FKLogI, userId uint64, barrier int32, itemList map[int32]int64) (err error)
	// 保存积分装备和道具信息
	SaveBarrierScoreReward(logger fklog.FKLogI, userId uint64, barrier int32, equipList map[int32]int32, itemList map[int32]int64) (err error)
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
