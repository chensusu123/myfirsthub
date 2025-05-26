// @Author pangchenyang 2025/3/21 17:12:00
// @Desc:
package sweep

import (
	"maze_game_server/lib/net/websocket_service"
	"maze_game_server/pb/common/MazeGame"
)

func RegTcpHandler() {
	// start sweep
	_ = websocket_service.RegProcSimple(10471, &MazeGame.StartMazeSweepRQ{},
		10472, &MazeGame.StartMazeSweepRS{},
		OnStartMazeSweepRQ)
}
