package tasktimer

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"maze_game_server/pb/server/SeaTaskSvr"
	"maze_game_server/usecase/tasktimer/delay/redisdelay"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/discoveryutil"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"go.uber.org/zap"
)

type TaskTimerBusiness struct {
	delay     *redisdelay.BucketTicker
	cancel    context.CancelFunc
	ctx       context.Context
	namespace string
}

// OnShutdown implements fkcore.FKServiceI.
func (tb *TaskTimerBusiness) OnShutdown(fklog.FKLogI) error {
	return nil
}

// FKServiceI 服务接口
func (tb *TaskTimerBusiness) Name() string {
	return "TaskTimerService"
}

func (tb *TaskTimerBusiness) OnInit(logger fklog.FKLogI, cfg fkconfig.FkConfigerI) (err error) {
	logger.InfoWF("TaskTimerBusiness OnInit")
	loopLogger := logger.Clone("loop")
	tb.namespace = appconfig.GlobalConfig().Global.Namespace
	redisCfg, err := discoveryutil.GetRedisCfg(tb.namespace, "maze_main_server.redis")
	if err != nil {
		logger.ErrorWF("TaskTimerBusiness GetRedisCfg failed", zap.Error(err))
		return err
	}
	redisAddr := redisCfg.Addr // 从环境变量中获取 Redis 地址
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
		logger.WarnWF("Redis address not found in environment variables, using default address", zap.String("default_address", redisAddr))
	}
	delay, err := redisdelay.New(logger, 1*time.Second, "TaskTimerBusiness", redisAddr, func(data interface{}) bool {
		taskLogger := loopLogger.Clone("task")
		taskLogger.SetLogId(time.Now().UnixNano())
		pack := SeaTaskSvr.TaskInfo{}
		taskErr := json.Unmarshal([]byte(data.(string)), &pack)
		if taskErr != nil {
			taskLogger.ErrorWF("unmarshal task failed", zap.Error(taskErr), zap.Any("data", data))
			return false
		}
		taskLogger.SetUid(pack.GetUserId())
		if gOnTimeout == nil {
			taskLogger.ErrorWF("gOnTimeout is nil", zap.Any("pack", pack))
			return true
		}
		ctx := context.Background()
		gOnTimeout(ctx, pack.GetUserId(), SeaTaskSvr.TaskExpireNotifyRQ{TaskInfo: &pack})
		taskLogger.InfoWF("do task ", zap.Any("pack", pack))
		return true
	})
	if err != nil {
		logger.ErrorWF("New redis delay failed", zap.Error(err))
		return err
	}
	tb.delay = delay
	tb.ctx, tb.cancel = context.WithCancel(context.Background())
	return nil
}

func (tb *TaskTimerBusiness) OnStart(logger fklog.FKLogI, cfg fkconfig.FkConfigerI) (err error) {
	logger.InfoWF("TaskTimerBusiness OnStart")
	if tb.delay != nil {
		loggerDelay := logger.Clone("delay")
		go tb.delay.Start(loggerDelay, tb.ctx)
	}
	return
}

func (tb *TaskTimerBusiness) OnStop(logger fklog.FKLogI) (err error) {
	logger.InfoWF("TaskTimerBusiness OnStop")

	if tb.cancel != nil {
		tb.cancel()
	}
	return
}

func (tb *TaskTimerBusiness) OnFinish(logger fklog.FKLogI) (err error) {
	// 关闭buffer channel
	logger.InfoWF("TaskTimerBusiness OnFinish")
	return
}

func (tb *TaskTimerBusiness) SetSeaTask(logger fklog.FKLogI, uid, session uint64, op SeaTaskSvr.TaskNotifyType, info *SeaTaskSvr.TaskInfo) (err error) {
	if gOnTimeout == nil {
		return errors.New("gOnTimeout is nil")
	}

	if tb.delay == nil {
		return errors.New("delay is nil")
	}
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	logger.DebugWF("SetSeaTask", zap.String("data", string(data)))
	if op == SeaTaskSvr.TaskNotifyType_ENUM_TASK_NOTIFY_TYPE_DEL {
		return tb.delay.DelTask(&redisdelay.Task{
			Id: string(info.GetContext()),
		})
	} else {
		return tb.delay.AddTask(&redisdelay.Task{
			Id:        string(info.GetContext()),
			Data:      string(data),
			Delay:     time.Duration(info.GetTime()) * time.Second,
			Timestamp: int(time.Now().Unix()),
		})
	}
}

var GTaskTimerBusiness = &TaskTimerBusiness{}
