package events

import (
	"context"
	"maze_game_server/pb/common/MazeGame"
)

// OnTriggerTrap 角色触发关卡时，触发
func OnTriggerTrap(ctx context.Context, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventTriggerTrap) {

}
