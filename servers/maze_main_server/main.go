package main

import (
	"maze_game_server/servers/maze_main_server/process"
	"maze_game_server/usecase/business"
	"maze_game_server/usecase/tasktimer"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager/loadconfigapi"
)

// 19987	UN_CGK_SVR_TYPE_MAZE_MAIN_SERVER 小程序版迷宫主服务
func main() {
	fkserver.SetMonitorName(fkserver.GroupNameGO, fkserver.ProjectNamePPWD, "maze_main_server")
	// process.RegisterHandler()

	process.RegisterHandler()
	// if fkconfig.EnvVal.IsLocalDev {
	// 	fkserver.AddBusiness(&business.GCustomBusiness)
	// 	loadconfigapi.SetLoadConfigFunc(business.GCustomBusiness.LoadCacheConfig)
	// 	loadconfigapi.SetInitConfigCacheFunc(business.GCustomBusiness.Init)
	// } else {
	// 	loadconfigapi.InitConfigRpcClient()
	// }
	fkserver.AddBusiness(&business.GCustomBusiness)
	loadconfigapi.SetLoadConfigFunc(business.GCustomBusiness.LoadCacheConfig)
	loadconfigapi.SetInitConfigCacheFunc(business.GCustomBusiness.Init)
	// fkserver.AddBusiness(&business.GCustomBusiness)
	fkserver.AddBusiness(tasktimer.GTaskTimerBusiness)
	fkserver.Run()
}
