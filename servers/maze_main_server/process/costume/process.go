package costume

import (
	"maze_game_server/lib/nano/component"
)

type CostumeComponent struct {
	component.Base
}

func NewCostumeComponent() *CostumeComponent {
	return &CostumeComponent{}
}
