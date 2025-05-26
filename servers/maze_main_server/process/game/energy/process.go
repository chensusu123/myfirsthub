/*
 * @Author: majian
 * @Date: 2025-03-21 21:18:08
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-27 14:14:43
 */
package energy

import (
	"maze_game_server/lib/net/websocket_service"
	"maze_game_server/pb/common/MazeEnergy"
)

func RegTcpHandler() {
	// 迷宫体力查询
	websocket_service.RegProcSimple(10469, &MazeEnergy.QueryMazeEnergyRQ{},
		10470, &MazeEnergy.QueryMazeEnergyRS{}, OnQueryMazeEnergyRQ)

}

// func RegRpcHandler() {
// 	// 扣迷宫体力
// 	thrift_service.RegisterTwowaySimple(131445, &MazeEnergySvr.SubMazeEnergyRQ{},
// 		131446, &MazeEnergySvr.SubMazeEnergyRS{}, OnSubMazeEnergyRQ)

// 	// 加迷宫体力
// 	thrift_service.RegisterTwowaySimple(131455, &MazeEnergySvr.AddMazeEnergyRQ{},
// 		131456, &MazeEnergySvr.AddMazeEnergyRS{}, OnAddMazeEnergyRQ)
// }
