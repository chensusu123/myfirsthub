package events

import (
	"context"
	"maze_game_server/pb/common/MazeGame"
)

// OnRefreshMonster 怪物刷新时触发
func OnRefreshMonster(ctx context.Context, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventRefreshMonster) {

}
