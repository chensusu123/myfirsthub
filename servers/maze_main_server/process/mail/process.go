package mail

import (
	"maze_game_server/lib/nano/component"
)

type Mail struct {
	component.Base
}

func NewMail() *Mail {
	return &Mail{}
}

func RegTcpHandler() {
}

func RegRpcHandler() {
}

func RegConsumeHandler() {
}
