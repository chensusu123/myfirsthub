package process

import (
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGame"
	"gitlab.ifreetalk.com/maze-plate/protodef/SysPackDef"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
)

func RegisterHandler() {
	websocket_service.RegProcSimple(
		10492, &SysPackDef.UserLoginRq{},
		10493, &SysPackDef.UserLoginRs{},
		OnLoginRQ)

	websocket_service.RegProcSimple(
		16212, &MazeGame.BarrierDeathRQ{},
		16213, &MazeGame.BarrierDeathRS{},
		OnTestRQ)
}
