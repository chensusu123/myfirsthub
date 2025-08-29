package events

import (
	"maze_game_server/pb/common/MazeGame"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

// OnPauseGame 恢复/暂停游戏时触发
func OnPauseGame(logger fklog.FKLogI, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventPauseGame) {

}
