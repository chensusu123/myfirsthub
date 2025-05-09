// @Author pangchenyang 2025/3/21 17:12:00
// @Desc: 
package process

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
	"gitlab.ifreetalk.com/plate/protodef/MazeGame"
)

func RegTcpHandler() {
	// start sweep
	_ = websocket_service.RegProcSimple(16255, &MazeGame.StartMazeSweepRQ{},
		16256, &MazeGame.StartMazeSweepRS{},
		OnStartMazeSweepRQ)
}
