package mazeequipaffixrandpoolv8

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"gitlab.ifreetalk.com/maze/maze_game_server/config/GMazeEquipAffixRandPoolV8Cfg"
)

type MazeEquipAffixRandPoolV8CfgEx struct {
	EquipPoolWeightMap map[int32]map[int32]*GMazeEquipAffixRandPoolV8Cfg.MazeEquipAffixRandPoolV8ConfigRow
	EquipPoolGroupMap  map[int32]map[int32][]int32
	EquipAttrGroupMap  map[int32]map[int32]struct{}
	lock               sync.RWMutex
}

func init() {
	GMazeEquipAffixRandPoolV8Cfg.RegisterMazeEquipAffixRandPoolV8InitCallBack("mazeequipaffixrandpoolv8", loadMazeEquipAffixRandPoolV8CfgEx)
}

var gConfigData *MazeEquipAffixRandPoolV8CfgEx

// 重新构建配置
func loadMazeEquipAffixRandPoolV8CfgEx(cfg *GMazeEquipAffixRandPoolV8Cfg.MazeEquipAffixRandPoolV8Config) {
	g := &MazeEquipAffixRandPoolV8CfgEx{}
	g.EquipPoolWeightMap = make(map[int32]map[int32]*GMazeEquipAffixRandPoolV8Cfg.MazeEquipAffixRandPoolV8ConfigRow, 0)
	g.EquipPoolGroupMap = make(map[int32]map[int32][]int32)
	g.EquipAttrGroupMap = make(map[int32]map[int32]struct{})
	for _, value := range cfg.ConfigRows {
		if g.EquipPoolWeightMap[value.Pool_id] == nil {
			g.EquipPoolWeightMap[value.Pool_id] = make(map[int32]*GMazeEquipAffixRandPoolV8Cfg.MazeEquipAffixRandPoolV8ConfigRow, 0)
		}
		g.EquipPoolWeightMap[value.Pool_id][value.Affix_id] = value
		if g.EquipPoolGroupMap[value.Pool_id] == nil {
			g.EquipPoolGroupMap[value.Pool_id] = make(map[int32][]int32, 0)
		}
		g.EquipPoolGroupMap[value.Pool_id][value.Group] = append(g.EquipPoolGroupMap[value.Pool_id][value.Group], value.Affix_id)

		for attrId := range value.Show_attr_min {
			if attrId == 0 {
				continue
			}
			if g.EquipAttrGroupMap[attrId] == nil {
				g.EquipAttrGroupMap[attrId] = make(map[int32]struct{}, 0)
			}
			g.EquipAttrGroupMap[attrId][value.Affix_id] = struct{}{}
		}
	}
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(g))
}

func GetEquipPoolWeightCfg(poolId int32) map[int32]*GMazeEquipAffixRandPoolV8Cfg.MazeEquipAffixRandPoolV8ConfigRow {
	poolWeight, ok := gConfigData.EquipPoolWeightMap[poolId]
	if !ok {
		return nil
	}
	return poolWeight
}

func GetPoolGroupCfg(poolId int32) map[int32][]int32 {
	groupMap, ok := gConfigData.EquipPoolGroupMap[poolId]
	if !ok {
		return nil
	}
	return groupMap
}

func GetPoolLimitCfg(attrLimits []int32) map[int32]struct{} {
	poolLimitMap := make(map[int32]struct{})
	for _, attrId := range attrLimits {
		groupMap, ok := gConfigData.EquipAttrGroupMap[attrId]
		if !ok {
			continue
		}
		for poolId := range groupMap {
			poolLimitMap[poolId] = struct{}{}
		}
	}
	return poolLimitMap
}
