// @Author pangchenyang 2025/3/21 17:12:00
// @Desc:
package sweep

import (
	"github.com/lonng/nano/component"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
	"gitlab.ifreetalk.com/plate/protodef/MazeGame"
)

type Sweep struct {
	component.Base
}

func NewSweep() *Sweep {
	return &Sweep{}
}

func RegTcpHandler() {
	// start sweep
	_ = websocket_service.RegProcSimple(16255, &MazeGame.StartMazeSweepRQ{},
		16256, &MazeGame.StartMazeSweepRS{},
		nil /* OnStartMazeSweepRQ */)
}
