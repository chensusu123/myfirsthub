package mazemapeditorconfigidcfgex

import (
	"maze_game_server/config/GMazeMapEditorConfigIdV8Cfg"
	"sync/atomic"
	"unsafe"
)

type MazeMapEditorConfigIdV8ConfigEx struct {
	BarrierConfigMap map[int32][]*GMazeMapEditorConfigIdV8Cfg.MazeMapEditorConfigIdV8ConfigRow
}

func init() {
	GMazeMapEditorConfigIdV8Cfg.RegisterMazeMapEditorConfigIdV8InitCallBack("mazemapeditorconfigidcfgex", loadConfigByBarrier)
}

var gConfigDataEx *MazeMapEditorConfigIdV8ConfigEx

func loadConfigByBarrier(f *GMazeMapEditorConfigIdV8Cfg.MazeMapEditorConfigIdV8Config) {
	temp := &MazeMapEditorConfigIdV8ConfigEx{
		BarrierConfigMap: make(map[int32][]*GMazeMapEditorConfigIdV8Cfg.MazeMapEditorConfigIdV8ConfigRow),
	}
	for _, i := range f.GetAll() {
		barrierConfigs := temp.BarrierConfigMap[i.Level_id]
		barrierConfigs = append(barrierConfigs, i)
		temp.BarrierConfigMap[i.Level_id] = barrierConfigs
	}
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigDataEx)), unsafe.Pointer(temp))
}

// 根据关卡ID获取所有配置
func GetBarrierConfigs(barrier int32) []*GMazeMapEditorConfigIdV8Cfg.MazeMapEditorConfigIdV8ConfigRow {
	return gConfigDataEx.BarrierConfigMap[barrier]
}
