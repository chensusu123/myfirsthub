package barrierenergyservice

import (
	"context"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

type BarrierEnergyService interface {
	GetBarrierEnergy(ctx context.Context, userId uint64) (curEnergy int32, nextTime int64, err error)
	AddEnergy(ctx context.Context, userId uint64, addVal int32) (curEnergy int32, nextTime int64, err error)
	SubEnergy(ctx context.Context, userId uint64, subVal int32) (int32, error)
	SendEnergyChgPack(logger fklog.FKLogI, userId uint64, curEnergy int32, nextRecoverTime int64) error
	GetEnergyMaxValue() int32
	GetEnergyRecoverCfg() int64
	GetEnergyItemCfg() (int32, int32)
	ResetEnergy(ctx context.Context, userId uint64) error
	PushEnergyRecord(ctx context.Context, userId uint64, oldEnergy, newEnergy, opType int32, lastTime int64)
}

var GlobalBarrierEnergyService BarrierEnergyService

func init() {
	GlobalBarrierEnergyService = newBarrierEnergyService()
}

type service struct{}

func newBarrierEnergyService() BarrierEnergyService {
	return &service{}
}
