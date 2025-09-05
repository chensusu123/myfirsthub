package im

import (
	"maze_game_server/lib/nano/component"
)

type IM struct {
	component.Base
}

func NewIM() *IM {
	return &IM{}
}
