package awardservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

type AwardService interface {
	GetBarrierDeathAward(logger fklog.FKLogI, userId uint64, barrier int32) (awardMap map[int32]int64, equipMap map[int32]int32, expItem int64, err error)
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
