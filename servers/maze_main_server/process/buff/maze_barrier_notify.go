package buff

import (
	"go.uber.org/zap"
	"maze_game_server/common/constdef"
	"maze_game_server/common/structsdef"
	"maze_game_server/io/kafka/mazebarrieruserkafka"
	"maze_game_server/io/redis/mazeattrcalcnotifyqueue"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/services/tempbuffservice"

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
	// 删除临时buff武力属性
	err := mazebuffinforedis.DelMazeBuffBySrc(logger, msg.UserId, constdef.MazeBuffSrcSelectBuffForce)
	if err != nil {
		logger.ErrorWF("MazeBarrierNotifyProcess DelMazeBuffBySrc failed", zap.Uint64("userId", msg.UserId), zap.Error(err))
		return
	}

	// 推送属性计算消息
	calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{
		UserId: msg.UserId,
		// FromServer: fmt.Sprintf("%d %s", fkconfig.EnvVal.ServerType, fkconfig.EnvVal.AppName),
		ChgType: constdef.MazeBuffChgForceValue,
		Session: "buff",
		BuffSrc: constdef.MazeBuffSrcSelectBuffForce,
	}
	mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(logger, calcAttrNotify)
	_ = tempbuffservice.GlobalTempBuffService.DelTempBuff(logger, msg.UserId, msg.Barrier)
	return
}
