// @Author pangchenyang 2025/3/21 17:12:00
// @Desc:
package sweep

import (
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/plate/protodef/MazeGame"
)

func RegTcpHandler() {
	// start sweep
	_ = tcp_service.RegProcSimple(16255, &MazeGame.StartMazeSweepRQ{},
		16256, &MazeGame.StartMazeSweepRS{},
		OnStartMazeSweepRQ)
}
