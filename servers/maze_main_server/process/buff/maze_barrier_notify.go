package buff

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazebarrieruserkafka"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazetempbuffredis"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

/**
* @Description:
* @Author: wangyongliang
* @Date: 2025/3/21 21:28
**/

// MazeBarrierUserGameRecord 用户迷宫闯关纪录
type MazeBarrierUserGameRecord = mazebarrieruserkafka.MazeBarrierUserGameRecord

func MazeBarrierNotifyProcess(logger fklog.FKLogI, msg *MazeBarrierUserGameRecord) {
	// defer fkprometheus.DebugPMT("MazeBarrierNotifyProcess")()
	// msg := &MazeBarrierUserGameRecord{}
	// err := json.Unmarshal(data, msg)
	// if err != nil {
	// 	logger.ErrorWF("MazeBarrierNotifyProcess Unmarshal", zap.ByteString("data", data), zap.Error(err))
	// 	return nil
	// }
	if msg.GameRet != 1 && msg.GameRet != 2 {
		return
	}
	_ = mazetempbuffredis.DelMazeTempBuff(logger, msg.UserId, msg.Barrier)
	return
}
