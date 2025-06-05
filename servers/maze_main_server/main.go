package main

import (
	"maze_game_server/io/mysql"
	"maze_game_server/lib/net/polarismessvc"
	"maze_game_server/servers/maze_main_server/process"
	"maze_game_server/usecase/business"
	"maze_game_server/usecase/naming"
	"maze_game_server/usecase/tasktimer"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager/loadconfigapi"
	"gitlab.ifreetalk.com/maze-plate/freetk/mockio"
	namingI "gitlab.ifreetalk.com/maze-plate/freetk/pkg/naming"
)

// 19987	UN_CGK_SVR_TYPE_MAZE_MAIN_SERVER 小程序版迷宫主服务
func main() {
	fkserver.SetMonitorName(fkserver.GroupNameGO, fkserver.ProjectNamePPWD, "maze_main_server")

	process.RegisterHandler()

	fkserver.AppServer.AddBasicService(&process.NanoInitService{})

	var namingSvc namingI.NamingI
	if fkconfig.EnvVal.IsLocalDev {
		fkserver.AddBusiness(&business.GCustomBusiness)
		loadconfigapi.SetLoadConfigFunc(business.GCustomBusiness.LoadCacheConfig)
		loadconfigapi.SetInitConfigCacheFunc(business.GCustomBusiness.Init)
	} else {
		namingSvc = naming.NewClientSuite("./conf.d/polaris.yaml")
		mockio.SetNaming(namingSvc)
		// loadconfigapi.InitConfigRpcClient()
		fkserver.AddBusiness(&business.GCustomBusiness)
		loadconfigapi.SetLoadConfigFunc(business.GCustomBusiness.LoadCacheConfig)
		loadconfigapi.SetInitConfigCacheFunc(business.GCustomBusiness.Init)
		polarismessSvc := polarismessvc.NewPolarismesSvc()
		fkserver.AddBusiness(polarismessSvc)
	}

	// fkserver.AddBusiness(&business.GCustomBusiness)
	fkserver.AddBusiness(tasktimer.GTaskTimerBusiness)
	// 初始化mysql
	if fkconfig.EnvVal.IsLocalDev {
		mysql.InitMysql(nil)
	} else {
		mysql.InitMysql(namingSvc)
	}

	fkserver.Run()
}
