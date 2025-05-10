package main

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_es/process"
	"gitlab.ifreetalk.com/maze/maze_game_server/usecase/business"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/config_manager/loadconfigapi"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/web_service"
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
	loadconfigapi.SetLoadConfigFunc(business.GCustomBusiness.LoadCacheConfig)
	loadconfigapi.SetInitConfigCacheFunc(business.GCustomBusiness.Init)
	fkserver.Run()
}

func InitWeb(logger fklog.FKLogI) {
}
