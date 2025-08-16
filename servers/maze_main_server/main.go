package main

import (
	"fmt"
	"os"

	"maze_game_server/io/mysql"
	"maze_game_server/servers/maze_main_server/process"
	"maze_game_server/usecase/business"
	"maze_game_server/usecase/tasktimer"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/serverdepend"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager/loadconfigapi"
)

func GetEnv(name string) (string, error) {
	val, exists := os.LookupEnv(name)
	if !exists {
		return "", fmt.Errorf("environment %s not set", name)
	}
	return val, nil
}

// 19987	UN_CGK_SVR_TYPE_MAZE_MAIN_SERVER 小程序版迷宫主服务
func main() {
	fkserver.SetMonitorName(fkserver.GroupNameGO, fkserver.ProjectNamePPWD, "maze_main_server")

	process.RegisterHandler()

	fkserver.AddBusiness(&business.GCustomBusiness)
	loadconfigapi.SetLoadConfigFunc(business.GCustomBusiness.LoadCacheConfig)
	loadconfigapi.SetInitConfigCacheFunc(business.GCustomBusiness.Init)

	myBiz := mysql.NewBizGorm("BizCfg", "BizGorm")
	// 注册到服务依赖里面.初始化由框架进行调用
	serverdepend.RegisterDepend(myBiz)

	fkserver.AppServer.AddBasicService(&process.NanoInitService{})
	// fkserver.AddBusiness(&business.GCustomBusiness)
	fkserver.AddBusiness(tasktimer.GTaskTimerBusiness)

	fkserver.Run()
}
