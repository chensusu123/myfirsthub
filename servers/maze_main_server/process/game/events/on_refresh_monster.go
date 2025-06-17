package events

import (
	"maze_game_server/pb/common/MazeGame"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

// OnRefreshMonster 怪物刷新时触发
func OnRefreshMonster(logger fklog.FKLogI, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventRefreshMonster) {

}
