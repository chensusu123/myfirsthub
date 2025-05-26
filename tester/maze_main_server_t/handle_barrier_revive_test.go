/*
 * @Author: majian
 * @Date: 2025-03-25 14:19:57
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-25 17:49:10
 */
package maze_main_server_t

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"maze_game_server/io/redis/mazeuserbarrierredis"
	"maze_game_server/pb/common/Common"
	"maze_game_server/pb/common/MazeBarrierCache"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/servers/maze_main_server/process/game"
)

func TestOnMazeBarrierRebornRQ(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	req := &MazeGame.MazeBarrierRebornRQ{}
	req.Header = &Common.PacketHeader{}
	req.Header.Session = proto.String("fdsfdsfd")
	req.BarrierId = proto.Int32(1)
	req.RebornAck = proto.Int32(1)
	res := &MazeGame.MazeBarrierRebornRS{}
	e := game.OnMazeBarrierRebornRQ(ctx, 9003200130206264, req, res)
	l, _ := proto.Marshal(res)
	fmt.Println("l", len(l), "e", e, "res", res)
}

func TestOnMazeBarrierRebornAckRQ(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	req := &MazeGame.MazeBarrierRebornRQ{}
	req.Header = &Common.PacketHeader{}
	req.Header.Session = proto.String("fdsfdsfd")
	req.BarrierId = proto.Int32(1)
	req.RebornAck = proto.Int32(2)
	req.RebornCost = append(req.RebornCost, &MazeCommon.MazeItem{ItemId: proto.Int32(46900001),
		Count: proto.Int64(500)})
	res := &MazeGame.MazeBarrierRebornRS{}
	e := game.OnMazeBarrierRebornRQ(ctx, 9003200130206264, req, res)
	l, _ := proto.Marshal(res)
	fmt.Println("l", len(l), "e", e, "res", res)
}

func TestSetBarrier(t *testing.T) {
	SetBarrireInfo(gTestLogger, 9003200130206264, 1)
}

func SetBarrireInfo(logger fklog.FKLogI, userId uint64, barrierId int32) error {
	info := &MazeBarrierCache.MazeBarrierCache{}
	info.BarrierId = proto.Int32(barrierId)
	info.StartTime = proto.Int64(time.Now().Unix())
	info.EndTime = proto.Int64(0)
	return mazeuserbarrierredis.SetUserBarrierInfo(logger, userId, barrierId, info)
}
