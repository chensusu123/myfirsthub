package events

import (
	"context"
	"maze_game_server/pb/common/MazeGame"
)

// OnAttackMonster 角色攻击怪物时触发
func OnAttackMonster(ctx context.Context, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventAttack) {

}
