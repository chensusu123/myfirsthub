package process

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
	"gitlab.ifreetalk.com/plate/protodef/MazeGame"
	"gitlab.ifreetalk.com/plate/protodef/SysPackDef"
)

func RegisterHandler() {
	websocket_service.RegProcSimple(
		5183, &SysPackDef.UserLoginRq{},
		5184, &SysPackDef.UserLoginRs{},
		OnLoginRQ)

	websocket_service.RegProcSimple(
		16212, &MazeGame.BarrierDeathRQ{},
		16213, &MazeGame.BarrierDeathRS{},
		OnTestRQ)
}
