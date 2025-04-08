// @Author pangchenyang 2025/3/21 17:11:00
// @Desc:
package main

import (
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/exportlogservice"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/frontcache_service"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/web_service"
	"gitlab.ifreetalk.com/plate/io_interface/kafka_interface/exportlogkafka"
	"gitlab.ifreetalk.com/servers/maze_game_server/servers/maze_sweep_server/process"
)

func main() {
	fkserver.SetMonitorName(fkserver.GroupNameGO, fkserver.ProjectNamePPWD, "MazeSweepServer")
	exportlogservice.PlugExportLogService(exportlogkafka.GetProducer())
	tcp_service.PlugTcpService(process.RegTcpHandler)
	web_service.PlugWebService(regWebHandler)
	frontcache_service.PlugFCService()
	fkserver.Run()
}

func regWebHandler(logger fklog.FKLogI) {
}
