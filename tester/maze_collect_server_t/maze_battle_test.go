package doll_maze_server_t

import (
	"context"
	"fmt"
	"testing"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/common/errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/protodef/Common"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCollect"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/collect"
)

func TestMazeCollectInfoQueryRQ(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	req := &MazeCollect.MazeCollectInfoQueryRQ{}
	req.Header = &Common.PacketHeader{}
	req.Header.Session = proto.String("fdsfdsfd")
	res := &MazeCollect.MazeCollectInfoQueryRS{}
	res.ErrInfo = errors.NO_ERROR
	e := collect.OnMazeCollectInfoQueryRQ(ctx, 9003200130064576, req, res)
	fmt.Println("e", e, "res", res)
}

func TestMazeCollectItemReceiveRQ(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	req1 := &MazeCollect.MazeCollectInfoQueryRQ{}
	req1.Header = &Common.PacketHeader{}
	req1.Header.Session = proto.String("fdsfdsfd")
	res1 := &MazeCollect.MazeCollectInfoQueryRS{}
	res1.ErrInfo = errors.NO_ERROR
	e1 := collect.OnMazeCollectInfoQueryRQ(ctx, 9003200130064576, req1, res1)
	fmt.Println("e", e1, "res", res1)
	req := &MazeCollect.MazeCollectItemReceiveRQ{}
	req.Header = &Common.PacketHeader{}
	req.Header.Session = proto.String("fdsfdsfd")
	req.Items = res1.MazeCollectInfo.Items
	res := &MazeCollect.MazeCollectItemReceiveRS{}
	res.ErrInfo = errors.NO_ERROR
	e := collect.OnMazeCollectItemReceiveRQ(ctx, 9003200130064576, req, res)
	fmt.Println("e", e, "res", res)
}
