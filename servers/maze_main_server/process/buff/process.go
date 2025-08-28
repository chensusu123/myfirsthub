package buff

import (
	"maze_game_server/lib/nano/component"
)

type Buff struct {
	component.Base
}

func NewBuff() *Buff {
	return &Buff{}
}
