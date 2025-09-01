package events

import (
	"context"
	"maze_game_server/pb/common/MazeGame"
)

// OnLeaveBarrier 角色离开关卡时触发(服务端直接触发，非客户端上报)
func OnLeaveBarrier(ctx context.Context, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventLeaveBarrier) {
	// 通关
	if event.GetResult() == MazeGame.BarrierResult_PASS {

	}
	// 失败
	if event.GetResult() == MazeGame.BarrierResult_DEATH {

	}
}
