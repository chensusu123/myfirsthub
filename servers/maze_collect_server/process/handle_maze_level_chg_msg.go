package process

import (
	"context"
	"encoding/json"

	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazecollectredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/funcopencheck"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserinfo"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 用户等级变化流水
type MazeUserLevelRecord struct {
	UserId     uint64 `json:"user_id"`     // 用户id
	OldLevel   int32  `json:"old_level"`   // 旧等级
	NewLevel   int32  `json:"new_level"`   // 新等级
	GroupID    uint32 `json:"group_id"`    // 组id
	CreateTime int64  `json:"create_time"` // 操作时间 毫秒
}

func HandleMazeLevelMsg(ctx context.Context, logger fklog.FKLogI, index int, key, data []byte) (err error) {
	msg := &MazeUserLevelRecord{}
	err = json.Unmarshal(data, msg)
	if err != nil {
		logger.ErrorWF("HandleMazeLevelMsg Unmarshal", zap.Error(err),
			zap.Int("msg's len", len(data)), zap.Uint64("userId", msg.UserId))
		return
	}

	userId := msg.UserId
	logger.InfoWF("HandleMazeLevelMsg start", zap.Any("msg", msg),
		zap.Uint64("userId", userId))
	if userId <= 0 {
		return
	}
	// 道具产出信息
	collectInfo, err := mazecollectredis.GetCollectInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("HandleMazeLevelMsg GetCollectInfo", zap.Error(err))
		return
	}
	if collectInfo != nil {
		return nil
	}
	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("HandleMazeLevelMsg GetUserInfoV2 fail", zap.Error(err))
		return
	}
	if userInfo.PassBarrier <= 0 {
		return nil
	}

	result, err := funcopencheck.IsFuncOpen(1, int32(userInfo.Level))
	if !result.IsOpen {
		return nil
	}
	err = InitMazeCollectLand(logger, userId, userInfo.PassBarrier)
	if err != nil {
		logger.ErrorWF("HandleMazeLevelMsg InitMazeCollectLand", zap.Error(err))
		return err
	}
	return nil
}
