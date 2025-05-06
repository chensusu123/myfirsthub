package process

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/game"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/web_service"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip"
)

func RegisterHandler() {
	// 注册Tcp接口
	tcp_service.PlugTcpService(func() {
		// 游戏主功能接口
		game.RegTcpHandler()
		equip.RegTcpHandler()
	})

	// 主功能消费队列
	game.RegConsumeHandler()
	equip.RegConsumeHandler()

	// 注册Web接口
	web_service.PlugWebService(func(logger fklog.FKLogI) {

	})
}
