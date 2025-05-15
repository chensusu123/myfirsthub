package tasktimer

import (
	"errors"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/protodef/SeaTaskSvr"
)

func init() {
}

var gOnTimeout TimeoutFunc

func RegOnTimeoutFunc(f TimeoutFunc) {
	gOnTimeout = f
}

func SetSeaTask(logger fklog.FKLogI, uid, session uint64, op SeaTaskSvr.TaskNotifyType, info *SeaTaskSvr.TaskInfo) (err error) {
	if gOnTimeout == nil {
		return errors.New("gOnTimeout is nil")
	}
	return GTaskTimerBusiness.SetSeaTask(logger, uid, session, op, info)
}

type TimeoutFunc func(logger fklog.FKLogI, shardingID uint64, req SeaTaskSvr.TaskExpireNotifyRQ) (res SeaTaskSvr.TaskExpireNotifyRS, err error)
