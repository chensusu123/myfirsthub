package collect

import (
	"context"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/module/funcopencheck"
	"maze_game_server/module/mazecollect"
	"maze_game_server/module/mazeuserinfo"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 用户等级变化流水
type MazeUserLevelRecord = mazeuserlevelkafka.MazeUserLevelRecord

func HandleMazeLevelMsg(ctx context.Context, msg *MazeUserLevelRecord) {
	logger := fklog.ContextAppLogger(ctx)
	// msg := &MazeUserLevelRecord{}
	// err = json.Unmarshal(data, msg)
	// if err != nil {
	// 	logger.CtxError(ctx,"HandleMazeLevelMsg Unmarshal", zap.Error(err),
	// 		zap.Int("msg's len", len(data)), zap.Uint64("userId", msg.UserId))
	// 	return
	// }

	userId := msg.UserId
	logger.CtxInfo(ctx, "HandleMazeLevelMsg start", zap.Any("msg", msg),
		zap.Uint64("userId", userId))
	if userId <= 0 {
		return
	}
	// 道具产出信息
	//collectInfo, err := mazecollectredis.GetCollectInfo(logger, userId)
	//if err != nil {
	//	logger.CtxError(ctx,"HandleMazeLevelMsg GetCollectInfo", zap.Error(err))
	//	return
	//}
	//collectInfo := mazecollect.GetCollectInfo(logger, userId)
	//if collectInfo == nil {
	//	logger.CtxError(ctx,"HandleMazeLevelMsg GetCollectInfo is nil")
	//	return
	//}

	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "HandleMazeLevelMsg GetUserInfoV2 fail", zap.Error(err))
		return
	}
	if userInfo.PassBarrier <= 0 {
		return
	}

	result, err := funcopencheck.IsFuncOpen(ctx, 1, int32(userInfo.Level))
	if !result.IsOpen {
		return
	}
	//err = InitMazeCollectLand(logger, userId, userInfo.PassBarrier)
	cInfo := mazecollect.NewCollectInfo(ctx, userId)
	err, collectInfo := cInfo.NewMazeCollectInfo(ctx, userInfo.PassBarrier)
	if err != nil {
		logger.CtxError(ctx, "HandleMazeLevelMsg InitMazeCollectLand", zap.Error(err))
		return
	}
	NewCollectAfter(ctx, userId, collectInfo)
	return
}
