package process

import (
	"maze_game_server/servers/maze_main_server/process/collect"
	"maze_game_server/servers/maze_main_server/process/equip"
	"maze_game_server/servers/maze_main_server/process/game"
	"maze_game_server/servers/maze_main_server/process/pay"
	"maze_game_server/services/gmservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/web_service"
)

func RegisterHandler() {
	// 注册Tcp接口
	// tcp_service.PlugTcpService(func() {
	// 	// 目前只有一个Timer触发接口

	// })
	collect.RegTcpHandler()

	// // 注册WebSocket接口
	// websocket_service.PlugTcpRawService(func() {
	// 	auth.RegisterHandler()
	// 	// 主功能接口
	// 	game.RegTcpHandler()
	// 	// 扫荡功能接口
	// 	sweep.RegTcpHandler()
	// 	// 体力相关接口
	// 	energy.RegTcpHandler()
	// 	// 装备功能接口
	// 	equip.RegTcpHandler()
	// 	// 挂机收集接口
	// 	collect.RegWsHandler()
	// 	item.RegTcpHandler()
	// 	// interact.RegTcpHandler()
	// 	// 装备gm接口
	// 	equip_gm.RegTcpHandler()

	// 	buff.RegTcpHandler()

	// 	// 属性计算
	// 	attr_calc.RegTcpHandler()
	// 	card.RegTcpHandler()

	// 	rob.RegTcpHandler()
	// })

	// 注册Rpc接口
	// thrift_service.PlugThriftRpcService(func() {
	// 	// 装备rpc
	// 	// equip.RegRpcHandler()
	// 	// 属性计算rpc todo 目前看没有地方调用，先注释掉
	// 	// attr_calc.RegRpcHandler()
	// 	// 通用数值
	// 	// common_value.RegisterRpcPackProcessor()
	// },
	// )

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
		gmservice.GmService.RegHttp(logger)
		// 充值发货
		pay.RegPayDelivery(logger)
	})
}
