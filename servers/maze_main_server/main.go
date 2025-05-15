package main

import (
	"github.com/lonng/nano"
	"github.com/lonng/nano/serialize/json"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process"
	"gitlab.ifreetalk.com/maze/maze_game_server/usecase/business"
	"gitlab.ifreetalk.com/maze/maze_game_server/usecase/tasktimer"
)

// 19987	UN_CGK_SVR_TYPE_MAZE_MAIN_SERVER 小程序版迷宫主服务
func main() {
	fkserver.SetMonitorName(fkserver.GroupNameGO, fkserver.ProjectNamePPWD, "maze_main_server")

	// process.RegisterHandler()

	// Nano
	func() {
		nano.Listen(":5997",
			nano.WithDebugMode(),
			nano.WithIsWebsocket(true),
			// nano.WithSerializer(protobuf.NewSerializer()),
			nano.WithSerializer(json.NewSerializer()),
			nano.WithComponents(process.Components()),
		)
	}()

	fkserver.AddBusiness(&business.GCustomBusiness)
	fkserver.AddBusiness(tasktimer.GTaskTimerBusiness)
	fkserver.Run()
}
