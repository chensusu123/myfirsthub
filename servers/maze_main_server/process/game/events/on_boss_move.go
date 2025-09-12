package events

import (
	"context"
	"maze_game_server/pb/common/MazeGame"
)

// OnBossMove 猪妖移动时触发
func OnBossMove(ctx context.Context, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventBossMove) {

}
