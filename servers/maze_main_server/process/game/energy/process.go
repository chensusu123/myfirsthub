/*
 * @Author: majian
 * @Date: 2025-03-21 21:18:08
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-27 14:14:43
 */
package energy

import (
	"github.com/lonng/nano/component"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEnergy"
)

type Energy struct {
	component.Base
}

func NewEnergy() *Energy {
	return &Energy{}
}

func RegTcpHandler() {
	// 迷宫体力查询
	websocket_service.RegProcSimple(16240, &MazeEnergy.QueryMazeEnergyRQ{},
		16241, &MazeEnergy.QueryMazeEnergyRS{}, nil /* OnQueryMazeEnergyRQ */)

}

// func RegRpcHandler() {
// 	// 扣迷宫体力
// 	thrift_service.RegisterTwowaySimple(131445, &MazeEnergySvr.SubMazeEnergyRQ{},
// 		131446, &MazeEnergySvr.SubMazeEnergyRS{}, OnSubMazeEnergyRQ)

// 	// 加迷宫体力
// 	thrift_service.RegisterTwowaySimple(131455, &MazeEnergySvr.AddMazeEnergyRQ{},
// 		131456, &MazeEnergySvr.AddMazeEnergyRS{}, OnAddMazeEnergyRQ)
// }
