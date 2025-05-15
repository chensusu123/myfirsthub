package auth

import (
	"gitlab.ifreetalk.com/maze-plate/protodef/UserLogin"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
)

func RegisterHandler() {
	websocket_service.RegProcSimple(
		5183, &UserLogin.UserLoginRq{},
		5184, &UserLogin.UserLoginRs{},
		OnLoginRQ)
	websocket_service.RegProcSimple(
		5149, &UserLogin.UserLiveRq{},
		5150, &UserLogin.UserLiveRs{},
		OnLiveRQ)
}
