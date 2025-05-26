/*
 * @Author: majian
 * @Date: 2025-03-22 11:37:52
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-22 11:46:32
 */
package maze_energy_server_t

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"maze_game_server/common/errors"
	"maze_game_server/pb/common/Common"
	"maze_game_server/pb/common/MazeEnergy"
	"maze_game_server/servers/maze_main_server/process/game/energy"
)

var TestUid uint64 = 9003200130206264

func TestQueryMazeEnergyRQ(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	req := &MazeEnergy.QueryMazeEnergyRQ{}
	req.Header = &Common.PacketHeader{}
	req.Header.Session = proto.String("fdsfdsfd")
	res := &MazeEnergy.QueryMazeEnergyRS{}
	res.ErrInfo = errors.NO_ERROR
	e := energy.OnQueryMazeEnergyRQ(ctx, TestUid, req, res)
	l, _ := proto.Marshal(res)
	fmt.Println("l", len(l), "e", e, "res", res)
}
