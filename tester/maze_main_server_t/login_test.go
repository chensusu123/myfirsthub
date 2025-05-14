package maze_main_server_t

import (
	"context"
	"testing"
	"time"

	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/game"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGame"
)

var gTestUser uint64 = 9003200130206333

func TestLogin(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	req := &MazeGame.MazeLoginRQ{
		MazeVersion: proto.Int32(1),
	}
	res := &MazeGame.MazeLoginRS{}
	game.OnMazeLoginRQ(ctx, gTestUser, req, res)
	time.Sleep(time.Second)
}

func TestBarrierList(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	req := &MazeGame.MazeBarrierListRQ{}
	res := &MazeGame.MazeBarrierListRS{}
	game.OnMazeBarrierListRQ(ctx, gTestUser, req, res)
	time.Sleep(time.Second)
}

func TestEnter(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	req := &MazeGame.MazeBarrierEnterRQ{
		BarrierId: proto.Int32(1),
	}
	res := &MazeGame.MazeBarrierEnterRS{}
	game.OnMazeBarrierEnterRQ(ctx, gTestUser, req, res)
	time.Sleep(time.Second)
}

func TestDeath(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	req := &MazeGame.BarrierDeathRQ{
		BarrierId: proto.Int32(1),
		FoeExp:    proto.Int32(5),
	}
	res := &MazeGame.BarrierDeathRS{}
	game.OnMazeBarrierDeathRQ(ctx, gTestUser, req, res)
	time.Sleep(time.Second)
}

func TestReborn(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	req := &MazeGame.MazeBarrierRebornRQ{
		BarrierId: proto.Int32(1),
		RebornAck: proto.Int32(2),
		RebornCost: []*MazeCommon.MazeItem{
			&MazeCommon.MazeItem{ItemId: proto.Int32(46900001), Count: proto.Int64(500)},
		},
	}
	res := &MazeGame.MazeBarrierRebornRS{}
	game.OnMazeBarrierRebornRQ(ctx, gTestUser, req, res)
	time.Sleep(time.Second)
}
func TestPass(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	req := &MazeGame.MazeBarrierPassRQ{
		BarrierId: proto.Int32(1),
		FoeExp:    proto.Int32(100),
	}
	res := &MazeGame.MazeBarrierPassRS{}
	game.OnMazeBarrierPassRQ(ctx, gTestUser, req, res)
	time.Sleep(time.Second)
}

func TestReport(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	req := &MazeGame.ReportDataRQ{
		UserInfo: &MazeGame.ReportUserInfo{
			ReportMask: proto.Int64(0),
			ExpTotal:   proto.Int64(10),
			MoneyCount: proto.Int64(20),
			EquipPoint: proto.Int64(0),
		},
	}
	res := &MazeGame.ReportDataRS{}
	game.OnMazeLoginRQ(ctx, gTestUser, req, res)
	time.Sleep(time.Second)
}

func TestReportAddEquip(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	req := &MazeGame.ReportAwardFoeEquipRQ{
		UserLevel: proto.Int32(1),
		BarrierId: proto.Int32(1),
		EquipNum:  proto.Int32(1),
	}
	res := &MazeGame.ReportAwardFoeEquipRS{}
	game.OnReportAwardFoeEquipRQ(ctx, gTestUser, req, res)
	time.Sleep(time.Second)
}
