package main

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/web_service"
	"maze_game_server/lib/net/websocket_service"
	"maze_game_server/servers/maze_es/process"
	"maze_game_server/usecase/business"
)

// 19987	UN_CGK_SVR_TYPE_MAZE_MAIN_SERVER 小程序版迷宫主服务
func main() {
	fkserver.SetMonitorName(fkserver.GroupNameGO, fkserver.ProjectNamePPWD, "maze_es")
	// exportlogservice.PlugExportLogService(exportlogkafka.GetProducer())

	websocket_service.PlugTcpRawService(process.RegisterHandler)
	web_service.PlugWebService(InitWeb)
	tcp_service.PlugTcpService(func() {
	})
	fkserver.AddBusiness(&business.GCustomBusiness)

	fkserver.Run()
}

func InitWeb(logger fklog.FKLogI) {
}
