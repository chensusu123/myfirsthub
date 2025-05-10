package main

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process"
	"gitlab.ifreetalk.com/maze/maze_game_server/usecase/business"
	"gitlab.ifreetalk.com/maze/maze_game_server/usecase/tasktimer"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
)

// 19987	UN_CGK_SVR_TYPE_MAZE_MAIN_SERVER 小程序版迷宫主服务
func main() {
	fkserver.SetMonitorName(fkserver.GroupNameGO, fkserver.ProjectNamePPWD, "maze_main_server")

	process.RegisterHandler()

	fkserver.AddBusiness(&business.GCustomBusiness)
	fkserver.AddBusiness(tasktimer.GTaskTimerBusiness)
	fkserver.Run()
}
