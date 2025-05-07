package process

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/game"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/interact"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/item"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/web_service"
)

func RegisterHandler() {
	// 注册Tcp接口
	tcp_service.PlugTcpService(func() {
		// 主功能接口
		game.RegTcpHandler()
		// // 体力相关接口
		// energy.RegTcpHandler()
		// 装备功能接口
		equip.RegTcpHandler()
		item.RegTcpHandler()
		interact.RegTcpHandler()
		
		// // 挂机收集接口
		// collect.RegTcpHandler()
	})

	// 消费队列
	{
		// 主功能队列
		game.RegConsumeHandler()
		// 装备队列
		equip.RegConsumeHandler()
		// // 挂机收集队列
		// collect.RegConsumeHandler()
	}

	// 注册Web接口
	web_service.PlugWebService(func(logger fklog.FKLogI) {

	})

}
