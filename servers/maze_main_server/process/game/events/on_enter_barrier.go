package events

import (
	"context"
	"maze_game_server/pb/common/MazeGame"
)

// OnEnterBarrier 角色进入关卡时触发(服务端直接触发，非客户端上报)
func OnEnterBarrier(ctx context.Context, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventEnterBarrier) {

}
