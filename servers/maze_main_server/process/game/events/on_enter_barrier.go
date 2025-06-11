package events

import (
	"maze_game_server/pb/common/MazeGame"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

// OnEnterBarrier 角色进入关卡时触发(服务端直接触发，非客户端上报)
func OnEnterBarrier(logger fklog.FKLogI, userID uint64, frame int64, timeMs int64, event *MazeGame.BattleEventEnterBarrier) {

}
