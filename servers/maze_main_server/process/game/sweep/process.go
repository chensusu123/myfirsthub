// @Author pangchenyang 2025/3/21 17:12:00
// @Desc:
package sweep

import (
	"maze_game_server/lib/nano/component"
)

type Sweep struct {
	component.Base
}

func NewSweep() *Sweep {
	return &Sweep{}
}

func RegTcpHandler() {
	// // start sweep
	// _ = websocket_service.RegProcSimple(10471, &MazeGame.StartMazeSweepRQ{},
	// 	10472, &MazeGame.StartMazeSweepRS{},
	// 	OnStartMazeSweepRQ)
}
