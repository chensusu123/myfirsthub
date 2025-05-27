package main

import (
	"maze_game_server/servers/maze_main_server/process"
	"maze_game_server/usecase/tasktimer"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
)

// 19987	UN_CGK_SVR_TYPE_MAZE_MAIN_SERVER 小程序版迷宫主服务
func main() {
	fkserver.SetMonitorName(fkserver.GroupNameGO, fkserver.ProjectNamePPWD, "maze_main_server")
	// process.RegisterHandler()

	process.RegisterHandler()
	// fkserver.AddBusiness(&business.GCustomBusiness)
	fkserver.AddBusiness(tasktimer.GTaskTimerBusiness)
	fkserver.Run()
}
