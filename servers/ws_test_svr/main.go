package main

import (
	"time"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGame"
	"gitlab.ifreetalk.com/maze-plate/protodef/SysPackDef"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/raw_pkg"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_es/process"
)

var gAddr = "127.0.0.1:9876"

func initLog() {
	logConfig := fklog.LogConfig{
		LogDir:     "/tmp",
		LogLevel:   "debug",
		LogName:    "ws_test_svr",
		LogType:    "zap",
		WithCaller: true,
	}
	fklog.InitAppFkLog(&logConfig)
}

func initService() {
	websocket_service.RegProcSimple(
		5183, &SysPackDef.UserLoginRq{},
		5184, &SysPackDef.UserLoginRs{},
		process.OnLoginRQ)

	websocket_service.RegProcSimple(
		16212, &MazeGame.BarrierDeathRQ{},
		16213, &MazeGame.BarrierDeathRS{},
		process.OnTestRQ)
}

func sendMsgToUser(logger fklog.FKLogI) {
	for i := 0; i < 10; i++ {
		data := makeOtherData()
		websocket_service.SendBytes(logger, 999, 1, data)
		time.Sleep(time.Second * 1)
	}
}

func svr_run(logger fklog.FKLogI) {
	initService()
	websocket_service.MockOnInit(logger, gAddr)
	websocket_service.MockOnStart(logger)

	// go client(logger, gAddr)
	go sendMsgToUser(logger)
	// go client(logger, gAddr)
	time.Sleep(time.Second * 120)
}

func makeLoginData() []byte {
	pbPacket := SysPackDef.UserLoginRq{
		UserID:        proto.Uint64(999),
		UserAuthToken: proto.Uint64(2224),
	}
	pbData, _ := proto.Marshal(&pbPacket)
	pkg := &raw_pkg.StruSvrEsRawBaseHead{}
	pkg.PackType = 5183
	pkg.Data = pbData
	data, _ := pkg.Pack()
	return data
}

func makeOtherData() []byte {
	pbPacket := MazeGame.BarrierDeathRQ{
		BarrierId: proto.Int32(222),
	}
	pbData, _ := proto.Marshal(&pbPacket)
	pkg := &raw_pkg.StruSvrEsRawBaseHead{}
	pkg.PackType = 16212
	pkg.Data = pbData
	data, _ := pkg.Pack()
	return data
}

func main() {
	initLog()
	gTestLogger := fklog.AppLogger().Clone("maze_energy_server_t")

	svr_run(gTestLogger)
}
