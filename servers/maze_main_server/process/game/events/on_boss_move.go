package events

import (
	"maze_game_server/pb/common/MazeGame"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

// OnBossMove 猪妖移动时触发
func OnBossMove(logger fklog.FKLogI, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventBossMove) {

}
