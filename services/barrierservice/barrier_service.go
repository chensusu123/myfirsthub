package barrierservice

import (
	"context"
	"maze_game_server/pb/common/Common"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeEnergy"
	"maze_game_server/pb/common/MazeGame"
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
	SweepBarrier(ctx context.Context, header *Common.PacketHeader, userID uint64, barrierID int32) (
		energyInfo *MazeEnergy.EnergyInfo, remainVal int32, gameID uint64, awardItem, rareItem []*MazeCommon.MazeItem, errinfo *MessageType.ErrorInfo)

	// GetBarrierInfos 获取用户对应关卡列表信息，返回挑战次数与挑战刷新时间
	//
	// 参数：
	//	- userID: 用户ID
	GetBarrierInfos(logger fklog.FKLogI, userID uint64) (barrierInfos []*MazeGame.MazeBarrierInfo, errinfo *MessageType.ErrorInfo)

	// // BarrierEnter 进入指定关卡，返回关卡战斗相关信息与配置
	// //
	// // TODO 该功能接口依赖比较多，延后
	// BarrierEnter(logger fklog.FKLogI, userID uint64, barrierID int32) (errinfo *MessageType.ErrorInfo)

	// BarrierPass 通关指定关卡
	//
	// 参数：
	//	- userID: 用户ID
	BarrierPass(ctx context.Context, header *Common.PacketHeader, userID uint64, barrierID int32, foeExp int32) (
		killMonsterNum int32, totalDamage int64, awards, rareAwards []*MazeCommon.MazeItem, errinfo *MessageType.ErrorInfo)

	// BarrierDeath 在指定关卡中死亡
	//
	// 参数：
	//	- userID: 用户ID
	BarrierDeath(ctx context.Context, header *Common.PacketHeader, userID uint64, barrierID int32, foeExp int32) (
		killMonsterNum int32, totalDamage int64, awards []*MazeCommon.MazeItem, errinfo *MessageType.ErrorInfo)

	// GuardDeath 关卡中击杀(守卫)怪物掉落奖励(注意：GuardDeath接口不负责增加奖励，增加操作由调用方处理)
	//
	// 返回值：
	// 	- kongfu: 宝箱掉落通关值
	// 	- equips: 宝箱掉落装备
	// 	- items: 掉落道具
	GuardDeath(logger fklog.FKLogI, userID uint64, barrierID int32, monsterID int32, monsterGuid int32) (
		kongfu int32, equips map[int32]int32, items map[int32]int64, errinfo *MessageType.ErrorInfo)

	// OpenBox 关卡中打开宝箱(注意：OpenBox接口不负责增加奖励，增加操作由调用方处理)
	//
	// 返回值：
	// 	- kongfu: 宝箱掉落通关值
	// 	- equips: 宝箱掉落装备
	// 	- items: 掉落道具
	OpenBox(logger fklog.FKLogI, userID uint64, barrierID int32, boxID int32) (
		kongfu int32, equips map[int32]int32, items map[int32]int64, errinfo *MessageType.ErrorInfo)

	// // GetUserBarrierInfo 获取用户指定关卡的存储信息
	// GetUserBarrierInfo(logger fklog.FKLogI, userID uint64, barrierID int32) (barrierInfo *MazeBarrierCache.MazeBarrierCache, err error)

	// // SetUserBarrierInfo 设置用户指定关卡的存储信息
	// SetUserBarrierInfo(logger fklog.FKLogI, userID uint64, barrierID int32, barrierInfo *MazeBarrierCache.MazeBarrierCache) (err error)

	// // GetRebornCostByCount 获取指定复活次数的复活消耗
	// //
	// // 参数：
	// //	- count: 复活次数，第几次复活
	// GetRebornCostByCount(logger fklog.FKLogI, count int32) (cost map[int32]int64, maxReborn int32, canReborn bool)
}

var (
	Global BarrierService = &barrier{}
)

type barrier struct {
}
