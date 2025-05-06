/*
 * @Author: majian
 * @Date: 2024-09-18 21:09:56
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-17 14:37:46
 */
package mazeequipinfocfgex

import (
	"sync/atomic"
	"unsafe"

	"gitlab.ifreetalk.com/plate/excel/auto/GMazeEquipInfoV8Cfg"
)

type MazeEquipExInfo struct {
	//lock                  sync.RWMutex
	EquipPosAndQualityMap map[int32]map[int32][]*GMazeEquipInfoV8Cfg.MazeEquipInfoV8ConfigRow
}

var gMazeEquipExMap *MazeEquipExInfo

func init() {
	GMazeEquipInfoV8Cfg.RegisterMazeEquipInfoV8InitCallBack("mazeequipinfocfgex", loadEquipInfoEx)
}

func loadEquipInfoEx(in *GMazeEquipInfoV8Cfg.MazeEquipInfoV8Config) {
	gTmp := &MazeEquipExInfo{}
	gTmp.EquipPosAndQualityMap = make(map[int32]map[int32][]*GMazeEquipInfoV8Cfg.MazeEquipInfoV8ConfigRow)
	for _, row := range in.ConfigRows {
		if gTmp.EquipPosAndQualityMap[row.Pos] == nil {
			gTmp.EquipPosAndQualityMap[row.Pos] = make(map[int32][]*GMazeEquipInfoV8Cfg.MazeEquipInfoV8ConfigRow)
		}
		gTmp.EquipPosAndQualityMap[row.Pos][row.Quality] = append(gTmp.EquipPosAndQualityMap[row.Pos][row.Quality], row)
	}
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gMazeEquipExMap)), unsafe.Pointer(gTmp))
}

func GetEquipCfgRows(pos, qua int32) []*GMazeEquipInfoV8Cfg.MazeEquipInfoV8ConfigRow {
	quRows := gMazeEquipExMap.EquipPosAndQualityMap[pos]
	if quRows != nil {
		return quRows[qua]
	}
	return nil
}
