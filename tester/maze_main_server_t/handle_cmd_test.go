package maze_main_server_t

import (
	"context"
	"fmt"
	"testing"

	"maze_game_server/pb/common/Common"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/servers/maze_main_server/process/game"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"google.golang.org/protobuf/proto"
)

func TestOnMazeCmdRQ(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	req := &MazeGame.SendDollMazeCmdRQ{}
	req.Header = &Common.PacketHeader{}
	req.Header.Session = proto.String("fdsfdsfd")
	req.CmdCode = proto.String("1006")
	req.CmdParam = proto.String("add_cnt=100")
	res := &MazeGame.SendDollMazeCmdRS{}
	e := game.OnSendDollMazeCmdRQ(ctx, 9003200130206264, req, res)
	l, _ := proto.Marshal(res)
	fmt.Println("l", len(l), "e", e, "res", res)
}
