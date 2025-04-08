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

	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_energy_server/process"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/protodef/Common"
	"gitlab.ifreetalk.com/plate/protodef/MazeEnergy"
)

var TestUid uint64 = 9003200130206264

func TestQueryMazeEnergyRQ(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	req := &MazeEnergy.QueryMazeEnergyRQ{}
	req.Header = &Common.PacketHeader{}
	req.Header.Session = proto.String("fdsfdsfd")
	res := &MazeEnergy.QueryMazeEnergyRS{}
	res.ErrInfo = errors.NO_ERROR
	e := process.OnQueryMazeEnergyRQ(ctx, TestUid, req, res)
	l, _ := proto.Marshal(res)
	fmt.Println("l", len(l), "e", e, "res", res)
}
