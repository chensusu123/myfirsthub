package frame

import (
	"maze_game_server/lib/nano/component"
)

type Frame struct {
	component.Base
}

func NewFrame() *Frame {
	return &Frame{}
}

func RegTcpHandler() {
}

func RegRpcHandler() {
}

func RegConsumeHandler() {
}
