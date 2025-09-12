/*
 * @Author: majian
 * @Date: 2024-07-09 11:45:18
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-25 11:03:43
 */
package mazeconfigv8

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/config/GMazeConfigV8Cfg"
)

//// 体力初始值
//func GetEnergyInitVal() int32 {
//	row := GMazeConfigV8Cfg.GetMazeConfigV8Config(constdef.MazeCfgId302)
//	if row != nil {
//		return int32(row.Value_int)
//	}
//	return 200
//}
//
//// 体力上限
//func GetEnergyMax() int32 {
//	row := GMazeConfigV8Cfg.GetMazeConfigV8Config(constdef.MazeCfgId301)
//	if row != nil {
//		return int32(row.Value_int)
//	}
//	return 200
//}
//
//// 体力恢复效率
//func GetEnergyRate() (costTime, recoverVal int32) {
//	row := GMazeConfigV8Cfg.GetMazeConfigV8Config(constdef.MazeCfgId303)
//	if row != nil {
//		for k, v := range row.Value_map {
//			return k, int32(v)
//		}
//	}
//	return 90, 1
//}
//
//// 迷宫体力回复间隔时间（秒）
//func GetEnergyRecoverCfg() int64 {
//	row := GMazeConfigV8Cfg.GetMazeConfigV8Config(constdef.MazeCfgId303)
//	if row != nil {
//		return row.Value_int
//	}
//	return 30
//}
//
//// 迷宫体力瓶道具id：对应的体力数量
//func GetEnergyItemCfg() (id, count int32) {
//	row := GMazeConfigV8Cfg.GetMazeConfigV8Config(constdef.MazeCfgId304)
//	if row != nil {
//		for k, v := range row.Value_map {
//			return k, int32(v)
//		}
//	}
//	return 49000001, 60
//}

// 获取客户端要接收的属性变化Id
func GetClientCareAttrIds() map[int32]struct{} {
	rs := make(map[int32]struct{})
	row := GMazeConfigV8Cfg.GetMazeConfigV8Config(constdef.MazeCfgId601)
	if row != nil {
		for _, id := range row.Value_list {
			rs[int32(id)] = struct{}{}
		}
	}
	return rs
}

// 获取复活等待时间
func GetUserReviveTime() int64 {
	row := GMazeConfigV8Cfg.GetMazeConfigV8Config(constdef.MazeCfgId701)
	if row != nil {
		return row.Value_int
	}
	return 60
}

// 获取最大buff组数量
func GetMaxBuffGroupCount(ctx context.Context) int64 {
	row := GMazeConfigV8Cfg.GetWithCtx(ctx, constdef.MazeCfgId941)
	if row != nil {
		return row.Value_int
	}
	return 4
}

// 获取不算组数量的组id
func GetSpecialBuffGroupId(ctx context.Context) int64 {
	row := GMazeConfigV8Cfg.GetWithCtx(ctx, constdef.MazeCfgId942)
	if row != nil {
		return row.Value_int
	}
	return 6009
}

// 获取展示道具和实际道具的对应map
func GetReallyItemMap(ctx context.Context, itemType int32) map[int32]int64 {
	row := GMazeConfigV8Cfg.GetWithCtx(ctx, itemType)
	if row != nil {
		return row.Value_map
	}

	return make(map[int32]int64)
}
