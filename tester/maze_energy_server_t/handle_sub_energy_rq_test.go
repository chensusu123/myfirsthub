/*
 * @Author: majian
 * @Date: 2025-03-22 11:42:07
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-22 11:54:24
 */
package maze_energy_server_t

import (
	"testing"

	"google.golang.org/protobuf/proto"
	"maze_game_server/pb/common/MazeEnergySvr"
	"maze_game_server/servers/maze_main_server/process/game/energy"
)

func TestSubEnergyRq(t *testing.T) {
	rq := &MazeEnergySvr.SubMazeEnergyRQ{}
	rs := &MazeEnergySvr.SubMazeEnergyRS{}
	rq.UserId = proto.Uint64(TestUid)
	rq.SubVal = proto.Int32(10)
	rq.OpType = proto.Int32(1)
	rq.TradeNumber = proto.Uint64(131223)
	energy.SubMazeEnergyRQ(gTestLogger, TestUid, rq, rs)
	// mazeenergyrpc.SubMazeEnergyRQ(gTestLogger, rq, rs)
}
