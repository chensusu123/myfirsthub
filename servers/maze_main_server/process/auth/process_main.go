package auth

import (
	"github.com/lonng/nano/component"
)

type Auth struct {
	component.Base
}

func NewAuth() *Auth {
	return &Auth{}
}

func RegisterHandler() {
	// websocket_service.RegProcSimple(
	// 	10492, &UserLogin.UserLoginRq{},
	// 	10493, &UserLogin.UserLoginRs{},
	// 	OnLoginRQ)
	// websocket_service.RegProcSimple(
	// 	10494, &UserLogin.UserLiveRq{},
	// 	10495, &UserLogin.UserLiveRs{},
	// 	OnLiveRQ)
	// websocket_service.RegProcSimple(
	// 	10500, &UserLogin.ConfigDataMd5Rq{},
	// 	10501, &UserLogin.ConfigDataMd5Rs{},
	// 	OnConfigDataMd5Rq)
}
