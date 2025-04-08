package maze_main_server_t

import (
	"context"
	"fmt"
	"testing"

	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/protodef/Common"
	"gitlab.ifreetalk.com/plate/protodef/MazeGame"
	"gitlab.ifreetalk.com/servers/maze_game_server/servers/maze_main_server/process"
)

func TestOnMazeCmdRQ(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	req := &MazeGame.SendDollMazeCmdRQ{}
	req.Header = &Common.PacketHeader{}
	req.Header.Session = proto.String("fdsfdsfd")
	req.CmdCode = proto.String("1006")
	req.CmdParam = proto.String("add_cnt=100")
	res := &MazeGame.SendDollMazeCmdRS{}
	e := process.OnSendDollMazeCmdRQ(ctx, 9003200130206264, req, res)
	l, _ := proto.Marshal(res)
	fmt.Println("l", len(l), "e", e, "res", res)
}
