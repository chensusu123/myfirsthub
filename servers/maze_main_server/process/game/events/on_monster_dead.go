package events

import (
	"maze_game_server/pb/common/MazeGame"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

// OnMonsterDead 怪物死亡时触发
func OnMonsterDead(logger fklog.FKLogI, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventMonsterDead) {

}
