// @Author pangchenyang 2025/3/21 17:12:00
// @Desc:
package sweep

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGame"
)

func RegTcpHandler() {
	// start sweep
	_ = websocket_service.RegProcSimple(10471, &MazeGame.StartMazeSweepRQ{},
		10472, &MazeGame.StartMazeSweepRS{},
		OnStartMazeSweepRQ)
}
