package process

import (
	"github.com/lonng/nano/component"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/game"
)

func Components() (components *component.Components) {
	components = &component.Components{}
	// 注册主服务
	components.Register(game.NewGame())
	return
}
