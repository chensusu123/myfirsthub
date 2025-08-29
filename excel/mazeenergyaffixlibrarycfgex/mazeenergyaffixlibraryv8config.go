package mazeenergyaffixlibrarycfgex

import (
	"maze_game_server/config/GMazeEnergyAffixLibraryV8Cfg"
	"sync/atomic"
	"unsafe"
)

type GMazeEnergyAffixLibraryV8CfgEx struct {
	BuffIdLibraryMap map[int32][]int32 // buff归属于哪几个库
}

func init() {
	GMazeEnergyAffixLibraryV8Cfg.RegisterMazeEnergyAffixLibraryV8InitCallBack("mazeenergyaffixlibrarycfgex", loadConfig)
}

var gConfigDataEx *GMazeEnergyAffixLibraryV8CfgEx

func loadConfig(f *GMazeEnergyAffixLibraryV8Cfg.MazeEnergyAffixLibraryV8Config) {
	temp := &GMazeEnergyAffixLibraryV8CfgEx{
		BuffIdLibraryMap: make(map[int32][]int32),
	}
	for _, i := range f.GetAll() {
		for _, buffId := range i.Affix_id_list {
			temp.BuffIdLibraryMap[buffId] = append(temp.BuffIdLibraryMap[buffId], i.Order)
		}
	}

	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigDataEx)), unsafe.Pointer(temp))
}

// 根据buffId获取库id
func GetLibrary(buffId int32) []int32 {
	return gConfigDataEx.BuffIdLibraryMap[buffId]
}
