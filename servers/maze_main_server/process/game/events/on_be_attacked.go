package events

import (
	"context"
	"maze_game_server/pb/common/MazeGame"
)

// OnBeAttacked 角色被攻击时触发
func OnBeAttacked(ctx context.Context, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventAttack) {

}
