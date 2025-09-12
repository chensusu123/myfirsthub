package events

import (
	"context"
	"maze_game_server/pb/common/MazeGame"
)

// OnRoleMove 角色移动时触发
func OnRoleMove(ctx context.Context, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventRoleMove) {

}
