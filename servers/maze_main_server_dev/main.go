package main

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"maze_game_server/servers/maze_main_server/process"
	"maze_game_server/usecase/business"
	"maze_game_server/usecase/tasktimer"
)

// 19987	UN_CGK_SVR_TYPE_MAZE_MAIN_SERVER 小程序版迷宫主服务
func main() {
	fkserver.SetMonitorName(fkserver.GroupNameGO, fkserver.ProjectNamePPWD, "maze_main_server_dev")

	process.RegisterHandler()

	fkserver.AddBusiness(&business.GCustomBusiness)
	fkserver.AddBusiness(tasktimer.GTaskTimerBusiness)

	fkserver.Run()
}
