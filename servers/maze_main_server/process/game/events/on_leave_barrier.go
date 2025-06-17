package events

import (
	"maze_game_server/pb/common/MazeGame"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

// OnLeaveBarrier 角色离开关卡时触发(服务端直接触发，非客户端上报)
func OnLeaveBarrier(logger fklog.FKLogI, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventLeaveBarrier) {
	// 通关
	if event.GetResult() == MazeGame.BarrierResult_PASS {

	}
	// 失败
	if event.GetResult() == MazeGame.BarrierResult_DEATH {

	}
}
