package main

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_shop_server/process"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/exportlogservice"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/web_service"
	"gitlab.ifreetalk.com/plate/io_interface/kafka_interface/exportlogkafka"
)

// 19988	UN_CGK_SVR_TYPE_MAZE_SHOP_SERVER 小程序版迷宫商店服务
func main() {
	fkserver.SetMonitorName(fkserver.GroupNameGO, fkserver.ProjectNamePPWD, "maze_shop_server")
	exportlogservice.PlugExportLogService(exportlogkafka.GetProducer())

	tcp_service.PlugTcpService(process.RegTcpHandler)
	web_service.PlugWebService(regWebHandler)

	fkserver.Run()
}

func regWebHandler(logger fklog.FKLogI) {

}
