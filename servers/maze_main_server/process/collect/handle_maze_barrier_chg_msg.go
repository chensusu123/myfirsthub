package collect

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazebarrieruserkafka"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazecollectredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/funcopencheck"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserinfo"
	"go.uber.org/zap"
)

// 用户迷宫闯关纪录
type MazeBarrierUserGameRecord = mazebarrieruserkafka.MazeBarrierUserGameRecord

func HandleMazeBarrierMsg(logger fklog.FKLogI, msg *MazeBarrierUserGameRecord) {
	// msg := &MazeBarrierUserGameRecord{}
	// err = json.Unmarshal(data, msg)
	// if err != nil {
	// 	logger.ErrorWF("HandleMazeBarrierMsg Unmarshal", zap.Error(err),
	// 		zap.Int("msg's len", len(data)), zap.Uint64("userId", msg.UserId))
	// 	return
	// }

	userId := msg.UserId
	logger.InfoWF("HandleMazeBarrierMsg start", zap.Any("msg", msg),
		zap.Uint64("userId", userId))
	if userId <= 0 {
		return
	}
	if msg.GameRet != 1 {
		return
	}
	// 道具产出信息
	collectInfo, err := mazecollectredis.GetCollectInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("HandleMazeBarrierMsg GetCollectInfo", zap.Error(err))
		return
	}
	if collectInfo != nil {
		return
	}
	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("HandleMazeBarrierMsg GetUserInfoV2 fail", zap.Error(err))
		return
	}
	if userInfo.PassBarrier <= 0 {
		return
	}

	result, err := funcopencheck.IsFuncOpen(1, int32(userInfo.Level))
	if !result.IsOpen {
		return
	}
	err = InitMazeCollectLand(logger, userId, userInfo.PassBarrier)
	if err != nil {
		logger.ErrorWF("HandleMazeBarrierMsg InitMazeCollectLand", zap.Error(err))
		return
	}
	return
}
