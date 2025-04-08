package main

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_common_value_server/process"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/exportlogservice"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/thrift_service"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/web_service"
	"gitlab.ifreetalk.com/plate/io_interface/kafka_interface/exportlogkafka"
)

// 20009	UN_CGK_SVR_TYPE_MAZE_COMMON_VALUE_SERVER	0	是	迷宫通用数值服务
func main() {
	fkserver.SetMonitorName(fkserver.GroupNameGO, fkserver.ProjectNamePPWD, "maze_common_value_server")
	exportlogservice.PlugExportLogService(exportlogkafka.GetProducer())

	thrift_service.PlugThriftRpcService(process.RegisterRpcPackProcessor)

	web_service.PlugWebService(regWebHandler)

	fkserver.Run()
}

func regWebHandler(logger fklog.FKLogI) {

}
