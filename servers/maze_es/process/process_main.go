package process

import (
	"maze_game_server/lib/net/websocket_service"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/pb/common/SysPackDef"
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
