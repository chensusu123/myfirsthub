package auth

import (
	"maze_game_server/lib/net/websocket_service"
	"maze_game_server/pb/common/UserLogin"
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
