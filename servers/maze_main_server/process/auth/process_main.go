package auth

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
	"gitlab.ifreetalk.com/maze-plate/protodef/SysPackDef"
)

func RegisterHandler() {
	websocket_service.RegProcSimple(
		5183, &SysPackDef.UserLoginRq{},
		5184, &SysPackDef.UserLoginRs{},
		OnLoginRQ)
}
