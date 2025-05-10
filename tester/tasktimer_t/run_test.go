package tasktimer_t

import (
	"testing"
	"time"

	"gitlab.ifreetalk.com/maze/maze_game_server/usecase/tasktimer"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/protodef/SeaTaskSvr"
	"go.uber.org/zap"
)

func ProcessTimeOut(logger fklog.FKLogI, shardingID uint64, req SeaTaskSvr.TaskExpireNotifyRQ) (res SeaTaskSvr.TaskExpireNotifyRS, err error) {
	res.TaskInfo = &SeaTaskSvr.TaskInfo{UserId: req.TaskInfo.UserId}
	// res.ErrInfo = errors.NO_ERROR

	defer func() {
		logger.InfoWF("ProcessTimeOut end", zap.Any("res", res))
	}()
	logger.InfoWF("ProcessTimeOut with ", zap.Any("Msg", req))

	taskInfo := req.GetTaskInfo()
	if req.GetTaskInfo() == nil {
		logger.ErrorWF("ProcessTimeOut taskinfo nil", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("task info nil")
		return
	}
	if uint32(time.Now().Unix()) < taskInfo.GetTime() {
		logger.ErrorWF("ProcessTimeOut check time failed. time not touch,call later.",
			zap.Uint64("uid", taskInfo.GetUserId()),
			zap.Stringer("task", taskInfo), zap.Uint64("shardingId", shardingID),
		)
		res.ErrInfo = errors.ARGS_NOT_MATCH.Wrap("时间还没到")
		return
	}
	if taskInfo.GetType() != uint32(234) {
		logger.ErrorWF("ProcessTimeOut task typ not match", zap.Any("req", req), zap.Uint32("myType", 234))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("task typ not match")
		return
	}
	return
}

func TestDelTask(t *testing.T) {
	logger := gTestLogger.Clone("TestTask")
	logger.DebugWF("TestTask")
	tasktimer.RegOnTimeoutFunc(ProcessTimeOut)
	task := tasktimer.TaskTimerBusiness{}
	task.OnInit(logger, nil)
	task.OnStart(logger, nil)

	task.SetSeaTask(logger, 2, 3, SeaTaskSvr.TaskNotifyType_ENUM_TASK_NOTIFY_TYPE_ADD, &SeaTaskSvr.TaskInfo{
		UserId:  proto.Uint64(123),
		Time:    proto.Uint32(5),
		Type:    proto.Uint32(234),
		Context: []byte("sss"),
	})
	task.SetSeaTask(logger, 2, 3, SeaTaskSvr.TaskNotifyType_ENUM_TASK_NOTIFY_TYPE_DEL, &SeaTaskSvr.TaskInfo{
		UserId:  proto.Uint64(123),
		Time:    proto.Uint32(2),
		Type:    proto.Uint32(234),
		Context: []byte("sss"),
	})

	time.Sleep(10 * time.Second)
	task.OnStop(logger)
	task.OnFinish(logger)
}

func TestAddTask(t *testing.T) {
	logger := gTestLogger.Clone("TestTask")
	logger.DebugWF("TestTask")
	tasktimer.RegOnTimeoutFunc(ProcessTimeOut)
	task := tasktimer.TaskTimerBusiness{}
	task.OnInit(logger, nil)
	task.OnStart(logger, nil)

	task.SetSeaTask(logger, 2, 3, SeaTaskSvr.TaskNotifyType_ENUM_TASK_NOTIFY_TYPE_ADD, &SeaTaskSvr.TaskInfo{
		UserId:  proto.Uint64(123),
		Time:    proto.Uint32(2),
		Type:    proto.Uint32(234),
		Context: []byte("sss"),
	})

	time.Sleep(30 * time.Second)
	task.OnStop(logger)
	task.OnFinish(logger)
}
