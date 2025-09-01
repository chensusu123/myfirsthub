package events

import (
	"context"
	"maze_game_server/pb/common/MazeGame"
)

// OnMonsterDead 怪物死亡时触发
func OnMonsterDead(ctx context.Context, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventMonsterDead) {

}
