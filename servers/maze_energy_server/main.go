package main

import (
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/exportlogservice"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/frontcache_service"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/thrift_service"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/web_service"
	"gitlab.ifreetalk.com/plate/io_interface/kafka_interface/exportlogkafka"
	"gitlab.ifreetalk.com/servers/maze_game_server/servers/maze_energy_server/process"
)

// 20007	UN_CGK_SVR_TYPE_MAZE_ENERGY_SERVER	0	是	迷宫体力服务
func main() {
	fkserver.SetMonitorName(fkserver.GroupNameGO, fkserver.ProjectNamePPWD, "maze_energy_server")
	exportlogservice.PlugExportLogService(exportlogkafka.GetProducer()) // 需要紧跟在SetMonitorName之后

	web_service.PlugWebService(emptyWeb)
	tcp_service.PlugTcpService(process.RegTcpHandler)
	thrift_service.PlugThriftRpcService(process.RegRpcHandler)
	frontcache_service.PlugFCService()
	fkserver.Run()

}

func emptyWeb(_ fklog.FKLogI) {}
