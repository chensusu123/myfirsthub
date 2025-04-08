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

// 用户迷宫闯关纪录
type MazeBarrierUserGameRecord struct {
	UserId     uint64 `json:"user_id"`     // 用户id
	Barrier    int32  `json:"barrier"`     // 关卡id
	GameRet    int32  `json:"game_ret"`    // 用户闯关结果 1-通关成功 2-死亡失败
	GroupID    uint32 `json:"group_id"`    // 组id
	CreateTime int64  `json:"create_time"` // 操作时间
}

func HandleMazeBarrierMsg(ctx context.Context, logger fklog.FKLogI, index int, key, data []byte) (err error) {
	msg := &MazeBarrierUserGameRecord{}
	err = json.Unmarshal(data, msg)
	if err != nil {
		logger.ErrorWF("HandleMazeBarrierMsg Unmarshal", zap.Error(err),
			zap.Int("msg's len", len(data)), zap.Uint64("userId", msg.UserId))
		return
	}

	userId := msg.UserId
	logger.InfoWF("HandleMazeBarrierMsg start", zap.Any("msg", msg),
		zap.Uint64("userId", userId))
	if userId <= 0 {
		return
	}
	if msg.GameRet != 1 {
		return nil
	}
	// 道具产出信息
	collectInfo, err := mazecollectredis.GetCollectInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("HandleMazeBarrierMsg GetCollectInfo", zap.Error(err))
		return
	}
	if collectInfo != nil {
		return nil
	}
	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("HandleMazeBarrierMsg GetUserInfoV2 fail", zap.Error(err))
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
		logger.ErrorWF("HandleMazeBarrierMsg InitMazeCollectLand", zap.Error(err))
		return err
	}
	return nil
}
