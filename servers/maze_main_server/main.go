package main

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/exportlogservice"
	"gitlab.ifreetalk.com/plate/io_interface/kafka_interface/exportlogkafka"
)

// 19987	UN_CGK_SVR_TYPE_MAZE_MAIN_SERVER 小程序版迷宫主服务
func main() {
	fkserver.SetMonitorName(fkserver.GroupNameGO, fkserver.ProjectNamePPWD, "maze_main_server")
	exportlogservice.PlugExportLogService(exportlogkafka.GetProducer())
	process.RegisterHandler()
	fkserver.Run()
}
