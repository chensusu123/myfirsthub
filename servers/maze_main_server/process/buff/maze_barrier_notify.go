package buff

import (
	"context"
	"encoding/json"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazetempbuffredis"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

/**
* @Description:
* @Author: wangyongliang
* @Date: 2025/3/21 21:28
**/

// MazeBarrierUserGameRecord 用户迷宫闯关纪录
type MazeBarrierUserGameRecord struct {
	UserId     uint64 `json:"user_id"`     // 用户id
	Barrier    int32  `json:"barrier"`     // 关卡id
	GameRet    int32  `json:"game_ret"`    // 用户闯关结果 1-通关成功 2-死亡失败
	GroupID    uint32 `json:"group_id"`    // 组id
	CreateTime int64  `json:"create_time"` // 操作时间
}

func MazeBarrierNotifyProcess(c context.Context, logger fklog.FKLogI, index int, key, data []byte) error {
	defer fkprometheus.DebugPMT("MazeBarrierNotifyProcess")()
	msg := &MazeBarrierUserGameRecord{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		logger.ErrorWF("MazeBarrierNotifyProcess Unmarshal", zap.ByteString("data", data), zap.Error(err))
		return nil
	}
	if msg.GameRet != 1 && msg.GameRet != 2 {
		return nil
	}
	_ = mazetempbuffredis.DelMazeTempBuff(logger, msg.UserId, msg.Barrier)
	return nil
}
