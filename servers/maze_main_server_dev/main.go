package main

import (
	"fmt"
	"os"

	"maze_game_server/lib/net/polarismessvc"
	"maze_game_server/servers/maze_main_server/process"
	"maze_game_server/usecase/business"
	"maze_game_server/usecase/cmdbconfig"
	"maze_game_server/usecase/localconfig"
	"maze_game_server/usecase/localnaming"
	"maze_game_server/usecase/naming"
	"maze_game_server/usecase/tasktimer"

	"maze_game_server/io/mysql"
	_ "maze_game_server/io/mysql_t"
	cfg2 "maze_game_server/lib/nano/cfg"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager/loadconfigapi"
	"gitlab.ifreetalk.com/maze-plate/freetk/mockio"
	namingI "gitlab.ifreetalk.com/maze-plate/freetk/pkg/naming"
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
	fkserver.SetMonitorName(fkserver.GroupNameGO, fkserver.ProjectNamePPWD, "maze_main_server_dev")

	process.RegisterHandler()

	var err error
	var namingSvc namingI.NamingI

	var cfgSvr cfg2.CfgSvr

	envMode, err := GetEnv("mode")
	switch envMode {
	case "dev":
		fkconfig.EnvVal.IsLocalDev = true
		cfgSvr = localconfig.New("./conf.d/localconfig.yaml")
		namingSvc = localnaming.NewLocalNaming("./conf.d/config.ini", []namingI.InitCfgFunc{
			mysql.InitMysqlEx,
		})
	case "docker":
		fkconfig.EnvVal.IsLocalDev = true
		cfgSvr = localconfig.New("./conf.d/dockerconfig.yaml")
		namingSvc = localnaming.NewLocalNaming("./conf.d/config.ini", []namingI.InitCfgFunc{
			mysql.InitMysqlEx,
		})
	default:
		cfgSvr = cmdbconfig.New("./conf.d/polaris.yaml")
		namingSvc = naming.NewClientSuite("./conf.d/polaris.yaml", []namingI.InitCfgFunc{
			mysql.InitMysqlEx,
		})
		polarismessSvc := polarismessvc.NewPolarismesSvc()
		fkserver.AddBusiness(polarismessSvc)
	}

	mockio.SetNaming(namingSvc)
	fkserver.AddBusiness(&business.GCustomBusiness)
	loadconfigapi.SetLoadConfigFunc(business.GCustomBusiness.LoadCacheConfig)
	loadconfigapi.SetInitConfigCacheFunc(business.GCustomBusiness.Init)

	myBiz := mysql.BizFlow{}
	err = myBiz.Init(cfgSvr)
	if err != nil {
		panic("mysql init err:" + err.Error())
	}
	fkserver.AppServer.AddBasicService(&process.NanoInitService{})
	fkserver.AddBusiness(tasktimer.GTaskTimerBusiness)

	fkserver.Run()
}
