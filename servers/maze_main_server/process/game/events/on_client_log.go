package events

import (
	"context"
	"maze_game_server/pb/common/MazeGame"
)

// OnClientLog 客户端日志上报
func OnClientLog(ctx context.Context, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventClientLog) {

}
