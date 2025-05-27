package gate

import "maze_game_server/lib/nano/component"

var (
	// All services in master server
	Services = &component.Components{}

	bindService = newBindService()
)

func init() {
	Services.Register(bindService)
}
