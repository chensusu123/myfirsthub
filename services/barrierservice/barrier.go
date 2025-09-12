package barrierservice

import (
	"context"
	"maze_game_server/config/GMazeRebornCostV8Cfg"
	"maze_game_server/io/redis/mazeuserbarrierredis"
	"maze_game_server/pb/server/MazeBarrierCache"
	"sort"
)

// type BarrierService interface {

// 	// GetBarrierInfos 获取用户对应关卡列表信息，返回挑战次数与挑战刷新时间
// 	GetBarrierInfos(ctx context.Context, userID uint64) (
// 		barrierInfos []*MazeGame.MazeBarrierInfo, challengeInfo *MazeCommon.MazeCount, challengeRefresh int64, errinfo *MessageType.ErrorInfo)

// 	// BarrierEnter 进入指定关卡，返回关卡战斗相关信息与配置
// 	//
// 	// TODO 该功能接口依赖比较多，延后
// 	BarrierEnter(ctx context.Context, userID uint64, barrierID int32) (errinfo *MessageType.ErrorInfo)

// 	// BarrierPass 通关指定关卡
// 	//
// 	// TODO 该功能接口依赖比较多，延后
// 	BarrierPass(ctx context.Context, userID uint64, barrierID int32) (errinfo *MessageType.ErrorInfo)

// 	// BarrierDeath 在指定关卡中死亡
// 	//
// 	// TODO 该功能接口依赖比较多，延后
// 	BarrierDeath(ctx context.Context, userID uint64, barrierID int32) (errinfo *MessageType.ErrorInfo)

// 	// GuardDeath 关卡中击杀(守卫)怪物掉落奖励(注意：GuardDeath接口不负责增加奖励，增加操作由调用方处理)
// 	//
// 	// 返回值：
// 	// 	- kongfu: 宝箱掉落通关值
// 	// 	- equips: 宝箱掉落装备
// 	// 	- items: 掉落道具
// 	GuardDeath(ctx context.Context, userID uint64, barrierID int32, monsterID int32) (
// 		kongfu int32, equips map[int32]int32, items map[int32]int64, errinfo *MessageType.ErrorInfo)

// 	// OpenBox 关卡中打开宝箱(注意：OpenBox接口不负责增加奖励，增加操作由调用方处理)
// 	//
// 	// 返回值：
// 	// 	- kongfu: 宝箱掉落通关值
// 	// 	- equips: 宝箱掉落装备
// 	// 	- items: 掉落道具
// 	OpenBox(ctx context.Context, userID uint64, barrierID int32, boxID int32) (
// 		kongfu int32, equips map[int32]int32, items map[int32]int64, errinfo *MessageType.ErrorInfo)

// 	// GetUserBarrierInfo 获取用户指定关卡的存储信息
// 	GetUserBarrierInfo(ctx context.Context, userID uint64, barrierID int32) (barrierInfo *MazeBarrierCache.MazeBarrierCache, err error)

// 	// SetUserBarrierInfo 设置用户指定关卡的存储信息
// 	SetUserBarrierInfo(ctx context.Context, userID uint64, barrierID int32, barrierInfo *MazeBarrierCache.MazeBarrierCache) (err error)

// 	// GetRebornCostByCount 获取指定复活次数的复活消耗
// 	//
// 	// 参数：
// 	//	- count: 复活次数，第几次复活
// 	GetRebornCostByCount(ctx context.Context, count int32) (cost map[int32]int64, maxReborn int32, canReborn bool)
// }

// var (
// 	Global BarrierService = &barrier{}
// )

// type barrier struct {
// }

// GetRebornCostByCount implements BarrierService.
func (b *barrier) GetRebornCostByCount(ctx context.Context, count int32) (cost map[int32]int64, maxReborn int32, canReborn bool) {
	allRows := GMazeRebornCostV8Cfg.GetAllMazeRebornCostV8Config()
	if len(allRows) == 0 {
		return nil, 0, false
	}
	sort.Slice(allRows, func(i, j int) bool {
		return allRows[i].Order > allRows[j].Order
	})

	ret := make(map[int32]int64)

	for _, row := range allRows {
		if maxReborn < row.Reborn_max {
			maxReborn = row.Reborn_max
		}
		if count <= row.Reborn_max && count >= row.Reborn_min {
			for k, v := range row.Cost {
				if k <= 0 || v <= 0 {
					continue
				}
				ret[k] = v
			}
			// break
		}
	}

	if len(ret) == 0 {
		return ret, maxReborn, false
	}

	return ret, maxReborn, true
}

// GetUserBarrierInfo implements BarrierService.
func (b *barrier) GetUserBarrierInfo(ctx context.Context, userID uint64, barrierID int32) (barrierInfo *MazeBarrierCache.MazeBarrierCache, err error) {
	return mazeuserbarrierredis.GetUserBarrierInfo(ctx, userID, barrierID)
}

// SetUserBarrierInfo implements BarrierService.
func (b *barrier) SetUserBarrierInfo(ctx context.Context, userID uint64, barrierID int32, barrierInfo *MazeBarrierCache.MazeBarrierCache) (err error) {
	return mazeuserbarrierredis.SetUserBarrierInfo(ctx, userID, barrierID, barrierInfo)
}
