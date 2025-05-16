package auth

import (
	"gitlab.ifreetalk.com/maze-plate/protodef/UserLogin"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
)

func RegisterHandler() {
	websocket_service.RegProcSimple(
		10492, &UserLogin.UserLoginRq{},
		10493, &UserLogin.UserLoginRs{},
		OnLoginRQ)
	websocket_service.RegProcSimple(
		10494, &UserLogin.UserLiveRq{},
		10495, &UserLogin.UserLiveRs{},
		OnLiveRQ)
}
