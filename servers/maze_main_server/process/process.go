package process

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/buff"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/collect"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip_gm"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/game"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/game/energy"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/game/sweep"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/interact"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/item"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/thrift_service"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/web_service"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/attr_calc"
)

func RegisterHandler() {
	// 注册Tcp接口
	tcp_service.PlugTcpService(func() {
		// 主功能接口
		game.RegTcpHandler()
		// 扫荡功能接口
		sweep.RegTcpHandler()
		// 体力相关接口
		energy.RegTcpHandler()
		// 装备功能接口
		equip.RegTcpHandler()
		// 挂机收集接口
		collect.RegTcpHandler()
		item.RegTcpHandler()
		interact.RegTcpHandler()
		// 装备gm接口
		equip_gm.RegTcpHandler()

		buff.RegTcpHandler()

		// 属性计算
		attr_calc.RegTcpHandler()

		// // 挂机收集接口
		// collect.RegTcpHandler()
	})

	// 注册Rpc接口
	thrift_service.PlugThriftRpcService(func() {
		// 装备rpc
		// equip.RegRpcHandler()
		// 属性计算rpc todo 目前看没有地方调用，先注释掉
		// attr_calc.RegRpcHandler()
		// 通用数值
		// common_value.RegisterRpcPackProcessor()
	},
	)

	// 消费队列
	{
		// 主功能队列
		game.RegConsumeHandler()
		// 装备队列
		equip.RegConsumeHandler()
		// 挂机收集队列
		collect.RegConsumeHandler()
		// 属性计算队列 废弃
		// attr_calc.RegConsumeHandler()
		// kafka转发队列 废弃
		// kafka_dispatch.RegConsumeHandler()
	}

	// 注册Web接口
	web_service.PlugWebService(func(logger fklog.FKLogI) {
		// 装备gm
		equip_gm.RegGm(logger)
	})

}
