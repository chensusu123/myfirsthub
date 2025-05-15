package mazeequiptyperesv8

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeEquipTypeResV8Cfg"
)

type MazeEquipTypeResV8ConfigEx struct {
	EquipTypeResCfgMap map[int32]map[int32][]*GMazeEquipTypeResV8Cfg.MazeEquipTypeResV8ConfigRow
	lock               sync.RWMutex
}

func init() {
	GMazeEquipTypeResV8Cfg.RegisterMazeEquipTypeResV8InitCallBack("gdollequipaffixsppoolv8cfgex", loadMazeEquipTypeResV8ConfigEx)
}

var gConfigData *MazeEquipTypeResV8ConfigEx

// 重新构建配置
func loadMazeEquipTypeResV8ConfigEx(cfg *GMazeEquipTypeResV8Cfg.MazeEquipTypeResV8Config) {
	g := &MazeEquipTypeResV8ConfigEx{}
	g.EquipTypeResCfgMap = make(map[int32]map[int32][]*GMazeEquipTypeResV8Cfg.MazeEquipTypeResV8ConfigRow, 0)
	for _, value := range cfg.ConfigRows {
		if g.EquipTypeResCfgMap[value.Equipment_id] == nil {
			g.EquipTypeResCfgMap[value.Equipment_id] = make(map[int32][]*GMazeEquipTypeResV8Cfg.MazeEquipTypeResV8ConfigRow, 0)
		}
		g.EquipTypeResCfgMap[value.Equipment_id][value.Pos_sub_type] = append(g.EquipTypeResCfgMap[value.Equipment_id][value.Pos_sub_type], value)
	}
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(g))
}

func GetEquipTypeResCfg(equipId, equipType, attrId int32) *GMazeEquipTypeResV8Cfg.MazeEquipTypeResV8ConfigRow {
	posAttrCfgMap, ok := gConfigData.EquipTypeResCfgMap[equipId]
	if !ok {
		return nil
	}
	attrCfgList, ok := posAttrCfgMap[equipType]
	if !ok {
		return nil
	}
	for _, cfg := range attrCfgList {
		if cfg.Attr_id == attrId {
			return cfg
		}
	}
	return nil
}
