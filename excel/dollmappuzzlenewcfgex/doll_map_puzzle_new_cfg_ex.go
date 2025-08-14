package dollmappuzzlenewcfgex

import (
	"maze_game_server/config/GDollMapPuzzleNewV8Cfg"
	"sync/atomic"
	"unsafe"
)

type AreaInfo struct {
	AreaId    int32 `json:"areaId,omitempty"`
	AreaIndex int32 `json:"areaIndex,omitempty"`
}

type DollMapPuzzleNewV8ConfigEx struct {
	BarrierConfigMap map[int32][]*GDollMapPuzzleNewV8Cfg.DollMapPuzzleNewV8ConfigRow
}

func init() {
	GDollMapPuzzleNewV8Cfg.RegisterDollMapPuzzleNewV8InitCallBack("dollmappuzzlenewcfgex", loadConfigByBarrier)
}

var gConfigDataEx *DollMapPuzzleNewV8ConfigEx

func loadConfigByBarrier(f *GDollMapPuzzleNewV8Cfg.DollMapPuzzleNewV8Config) {
	temp := &DollMapPuzzleNewV8ConfigEx{
		BarrierConfigMap: make(map[int32][]*GDollMapPuzzleNewV8Cfg.DollMapPuzzleNewV8ConfigRow),
	}
	for _, i := range f.GetAll() {
		barrierConfigs := temp.BarrierConfigMap[i.Level]
		barrierConfigs = append(barrierConfigs, i)
		temp.BarrierConfigMap[i.Level] = barrierConfigs
	}
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigDataEx)), unsafe.Pointer(temp))
}

// 根据关卡ID获取所有配置
func GetBarrierConfigs(barrier int32) []*GDollMapPuzzleNewV8Cfg.DollMapPuzzleNewV8ConfigRow {
	return gConfigDataEx.BarrierConfigMap[barrier]
}
