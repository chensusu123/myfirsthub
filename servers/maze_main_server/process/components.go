package process

import (
	"github.com/lonng/nano/component"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/collect"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/game"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/game/energy"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/game/sweep"
)

func Components() *component.Components {
	components := &component.Components{}
	// 注册主功能组件
	components.Register(game.NewGame())
	// 注册体力组件
	components.Register(energy.NewEnergy())
	// 注册挂机组件
	components.Register(collect.NewCollect())
	// 注册扫荡组件
	components.Register(sweep.NewSweep())
	return components
}
