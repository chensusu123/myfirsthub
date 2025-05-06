package process

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/game"
)

func RegTcpHandler() {
	// 游戏主功能接口
	game.RegTcpHandler()
}
