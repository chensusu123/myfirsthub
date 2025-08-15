package barrierservice

import (
	"maze_game_server/pb/common/Common"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeEnergy"
	"maze_game_server/pb/common/MessageType"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

type BarrierService interface {

	// SweepBarrier 扫荡关卡
	//
	// 参数：
	//	- header: 请求包头
	//	- userID: 用户ID
	//	- barrierID: 关卡ID
	SweepBarrier(logger fklog.FKLogI, header *Common.PacketHeader, userID uint64, barrierID int32) (
		energyInfo *MazeEnergy.EnergyInfo, remainVal int32, gameID uint64, awardItem, rareItem []*MazeCommon.MazeItem, errinfo *MessageType.ErrorInfo)
}

var (
	Global BarrierService = &barrier{}
)

type barrier struct {
}
