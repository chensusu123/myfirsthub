package events

import (
	"maze_game_server/pb/common/MazeGame"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

// OnAttackMonster 角色攻击怪物时触发
func OnAttackMonster(logger fklog.FKLogI, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventAttack) {

}
