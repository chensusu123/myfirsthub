package events

import (
	"maze_game_server/pb/common/MazeGame"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

// OnTriggerTrap 角色触发关卡时，触发
func OnTriggerTrap(logger fklog.FKLogI, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventTriggerTrap) {

}
