package friend

import (
	"maze_game_server/lib/nano/component"
)

type FriendComponent struct {
	component.Base
}

func NewFriendComponent() *FriendComponent {
	friendComponent := &FriendComponent{}
	return friendComponent
}
