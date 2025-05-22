package mazeequipaffixrolltypev8

import (
	"sort"
	"sync"
	"sync/atomic"
	"unsafe"

	"gitlab.ifreetalk.com/maze/maze_game_server/config/GMazeEquipAffixRollTypeV8Cfg"
)

type MazeEquipAffixModPoolV8CfgEx struct {
	lock        sync.RWMutex
	rollTypeMap map[int32][]*GMazeEquipAffixRollTypeV8Cfg.MazeEquipAffixRollTypeV8ConfigRow
}

func init() {
	GMazeEquipAffixRollTypeV8Cfg.RegisterMazeEquipAffixRollTypeV8InitCallBack("gdollequipaffixrolltypev8cfgex",
		loadGMazeEquipAffixRollTypeV8CfgEx)
}

var gConfigData *MazeEquipAffixModPoolV8CfgEx

func loadGMazeEquipAffixRollTypeV8CfgEx(
	cfg *GMazeEquipAffixRollTypeV8Cfg.MazeEquipAffixRollTypeV8Config) {
	tmpCfg := &MazeEquipAffixModPoolV8CfgEx{}
	tmpCfg.rollTypeMap = make(map[int32][]*GMazeEquipAffixRollTypeV8Cfg.MazeEquipAffixRollTypeV8ConfigRow)
	for _, row := range cfg.ConfigRows {
		tmpCfg.rollTypeMap[row.Roll_type] = append(tmpCfg.rollTypeMap[row.Roll_type], row)
	}
	for _, rollList := range tmpCfg.rollTypeMap {
		sort.Slice(rollList, func(i, j int) bool {
			if rollList[i].Get_score_min != rollList[j].Get_score_min {
				return rollList[i].Get_score_min < rollList[j].Get_score_min
			} else {
				return rollList[i].Order < rollList[j].Order
			}
		})
	}
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(tmpCfg))
}

// 根据类型查询，返回列表，用新变量返回结果，防止外面修改排序
func GetMazeEquipRollCfgByRollType(rollType int32) []*GMazeEquipAffixRollTypeV8Cfg.MazeEquipAffixRollTypeV8ConfigRow {
	rtList := gConfigData.rollTypeMap[rollType]
	if len(rtList) == 0 {
		return nil
	}
	rs := make([]*GMazeEquipAffixRollTypeV8Cfg.MazeEquipAffixRollTypeV8ConfigRow, 0, len(rtList))
	rs = append(rs, rtList...)
	return rs
}

func GetMazeEquipRollCfgByRollTypeAndScore(rollType, score int32) map[int32]int32 {
	rs := gConfigData.rollTypeMap[rollType]
	rollWeightMap := make(map[int32]int32)
	for _, row := range rs {
		if row.Get_score_min > score { // 从小到大有序
			break
		}
		if row.Get_score_min <= score && score <= row.Get_score_max {
			rollWeightMap[row.Order] = row.Weight
		}
	}
	return rollWeightMap
}
