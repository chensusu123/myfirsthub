package events

import (
	"maze_game_server/pb/common/MazeGame"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

// OnBeAttacked 角色被攻击时触发
func OnBeAttacked(logger fklog.FKLogI, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventAttack) {

}
