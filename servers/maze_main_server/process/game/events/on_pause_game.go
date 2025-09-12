package events

import (
	"context"
	"maze_game_server/pb/common/MazeGame"
)

// OnPauseGame 恢复/暂停游戏时触发
func OnPauseGame(ctx context.Context, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventPauseGame) {

}
