package mazebarriereventredis

import (
	"context"
	"fmt"
	"maze_game_server/pb/common/MazeGame"
	"time"

	jsoniter "github.com/json-iterator/go"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	FrameMax int64 = 999999
)

var (
	cli  = &fkredis.FkRedis{}
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type BarrierEvent struct {
	Frame int64       `json:"frame,omitempty"`
	Time  int64       `json:"time,omitempty"`
	Type  int32       `json:"type,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

func init() {
	fkconfig.RegisterNameNode("mazebarriereventredis", 21644, cli)
}

// getKey 获取缓存操作key。
func getKey(args ...interface{}) string {
	return fmt.Sprintf("u:%d:barrier:event", args...)
}

func getBackupKey(args ...interface{}) string {
	return fmt.Sprintf("u:%d:barrier:event:%d", args...)
}

// TriggerBarrierEvent 触发关卡事件
func TriggerBarrierEvent(logger fklog.FKLogI, userID uint64, frame int64, eventTime int64, eventType MazeGame.BattleEventType, eventData proto.Message) (err error) {
	var (
		key = getKey(userID)
		ctx = context.Background()
	)
	event := BarrierEvent{
		Frame: frame,
		Time:  eventTime,
		Type:  int32(eventType),
		Data:  eventData,
	}
	data, err := json.Marshal(event)
	if err != nil {
		logger.ErrorWF("TriggerBarrierEvent Marshal fail",
			zap.Error(err),
			zap.Uint64("userID", userID),
			zap.Int64("frame_model", frame),
			zap.Int64("eventTime", eventTime),
			zap.Int32("eventType", int32(eventType)),
			zap.Any("eventData", eventData),
		)
		return err
	}
	// 有frame则用frame，否则用时间
	score := eventTime
	// 进入关卡的帧序号是0
	if frame > 0 || eventType == MazeGame.BattleEventType_ENTER_BARRIER {
		score = frame
	}
	_, err = cli.Do(ctx, "ZADD", key, score, data)
	if err != nil {
		logger.ErrorWF("TriggerBarrierEvent ZADD fail",
			zap.Error(err),
			zap.Any("key", key),
		)
		return
	}
	logger.DebugWF("TriggerBarrierEvent success", zap.Uint64("userID", userID), zap.Int64("eventTime", eventTime), zap.Any("eventData", eventData))
	return
}

// EnterBarrier
func EnterBarrier(logger fklog.FKLogI, userID uint64, barrierID int32) (err error) {
	// 备份事件流
	enterTime, err := GetBarrierEnterTime(logger, userID)
	if err != nil {
		// TODO 错误处理
		logger.ErrorWF("EnterBarrier GetBarrierEnterTime fail", zap.Error(err), zap.Uint64("userID", userID), zap.Int32("barrierID", barrierID))
	} else if enterTime > 0 {
		BackupBarrierEvents(logger, userID, enterTime)
	}
	// 进入事件
	event := &MazeGame.BattleEventEnterBarrier{
		BarrierId: proto.Int32(barrierID),
	}
	return TriggerBarrierEvent(logger, userID, 0, time.Now().UnixMilli(), MazeGame.BattleEventType_ENTER_BARRIER, event)
}

// LeaveBarrier
func LeaveBarrier(logger fklog.FKLogI, userID uint64, barrierID int32, passed bool) (err error) {
	// 进入事件
	event := &MazeGame.BattleEventLeaveBarrier{
		BarrierId: proto.Int32(barrierID),
	}
	if passed {
		event.Result = MazeGame.BarrierResult_PASS.Enum()
	} else {
		event.Result = MazeGame.BarrierResult_DEATH.Enum()
	}
	err = TriggerBarrierEvent(logger, userID, FrameMax, time.Now().UnixMilli(), MazeGame.BattleEventType_LEAVE_BARRIER, event)
	if err != nil {
		logger.ErrorWF("EnterBarrier TriggerBarrierEvent fail", zap.Error(err), zap.Uint64("userID", userID), zap.Int32("barrierID", barrierID), zap.Any("passed", passed))
		return
	}
	// 备份事件流
	enterTime, err := GetBarrierEnterTime(logger, userID)
	if err != nil {
		// TODO 错误处理
		logger.ErrorWF("EnterBarrier GetBarrierEnterTime fail", zap.Error(err), zap.Uint64("userID", userID))
	} else if enterTime > 0 {
		BackupBarrierEvents(logger, userID, enterTime)
	}
	return
}

// GetBarrierEnterTime 获取关卡上报记录的第一条事件记录时间
func GetBarrierEnterTime(logger fklog.FKLogI, userID uint64) (enterTime int64, err error) {
	var (
		key = getKey(userID)
		ctx = context.Background()
	)
	memberAndScore, err := redis.Strings(cli.Do(ctx, "ZRANGE", key, 0, 0, "WITHSCORES"))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			logger.ErrorWF("GetBarrierEnterTime ZRANGE fail",
				zap.Error(err),
				zap.Any("key", key),
			)
			return 0, err
		}
	}
	if len(memberAndScore) == 2 {
		enterTime = fkutil.ToInt64(memberAndScore[1])
	}
	return
}

// BackupBarrierEvents
func BackupBarrierEvents(logger fklog.FKLogI, userID uint64, enterTime int64) (err error) {
	var (
		key = getKey(userID)
		new = getBackupKey(userID, enterTime)
		ctx = context.Background()
	)
	_, err = cli.Do(ctx, "RENAME", key, new)
	if err != nil {
		logger.ErrorWF("BackupBarrierEvents RENAME fail",
			zap.Error(err),
			zap.Any("key", key),
		)
	}
	// 备份关卡事件数据
	logger.DebugWF("BackupBarrierEvents success", zap.String("key", new), zap.Int64("enterTime", enterTime))
	return
}

// ClearBarrierEvents
func ClearBarrierEvents(logger fklog.FKLogI, userID uint64) (err error) {
	var (
		key = getKey(userID)
		ctx = context.Background()
	)
	_, err = cli.Do(ctx, "DEL", key)
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			logger.ErrorWF("ClearBarrierEvents command fail",
				zap.Error(err),
				zap.Any("key", key),
			)
			return err
		}
	}
	logger.DebugWF("ClearBarrierEvents success", zap.Uint64("userID", userID))
	return
}
