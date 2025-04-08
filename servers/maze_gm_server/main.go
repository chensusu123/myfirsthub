package main

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_gm_server/process"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/exportlogservice"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/web_service"
	"gitlab.ifreetalk.com/plate/io_interface/kafka_interface/exportlogkafka"
)

// 20013	UN_CGK_SVR_TYPE_MAZE_GM_SERVER	0	是	迷宫gm服务
func main() {
	fkserver.SetMonitorName(fkserver.GroupNameGO, fkserver.ProjectNamePPWD, "maze_gm_server")
	exportlogservice.PlugExportLogService(exportlogkafka.GetProducer())

	web_service.PlugWebService(process.RegGm)

	fkserver.Run()
}
