package tasktimer

import (
	"context"
	"errors"

	"maze_game_server/pb/server/SeaTaskSvr"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

func init() {
}

var gOnTimeout TimeoutFunc

func RegOnTimeoutFunc(f TimeoutFunc) {
	gOnTimeout = f
}

func SetSeaTask(ctx context.Context, uid, session uint64, op SeaTaskSvr.TaskNotifyType, info *SeaTaskSvr.TaskInfo) (err error) {
	if gOnTimeout == nil {
		return errors.New("gOnTimeout is nil")
	}
	logger := fklog.ContextAppLogger(ctx)
	return GTaskTimerBusiness.SetSeaTask(logger, uid, session, op, info)
}

type TimeoutFunc func(ctx context.Context, shardingID uint64, req SeaTaskSvr.TaskExpireNotifyRQ) (res SeaTaskSvr.TaskExpireNotifyRS, err error)
