/*
 * @Author: majian
 * @Date: 2024-07-08 21:21:04
 * @Last Modified by: majian
 * @Last Modified time: 2025-01-09 16:52:41
 */
package mazeattrorderv8

import (
	"sort"
	"sync/atomic"
	"unsafe"

	"gitlab.ifreetalk.com/plate/excel/auto/GMazeAttrListOrderV8Cfg"
)

type MazeAttrListOrderV8ConfigEx struct {
	DollAttrGroupMap map[int32][]*GMazeAttrListOrderV8Cfg.MazeAttrListOrderV8ConfigRow
}

func init() {
	GMazeAttrListOrderV8Cfg.RegisterMazeAttrListOrderV8InitCallBack("mazeattrorderv8", loadDollAttrOrderByGroup)
}

var gConfigData *MazeAttrListOrderV8ConfigEx

func loadDollAttrOrderByGroup(f *GMazeAttrListOrderV8Cfg.MazeAttrListOrderV8Config) {
	gTmp := &MazeAttrListOrderV8ConfigEx{}
	gTmp.DollAttrGroupMap = make(map[int32][]*GMazeAttrListOrderV8Cfg.MazeAttrListOrderV8ConfigRow)
	for _, row := range f.ConfigRows {
		gTmp.DollAttrGroupMap[row.Type] = append(gTmp.DollAttrGroupMap[row.Type], row)
	}
	for _, group := range gTmp.DollAttrGroupMap {
		sort.Slice(group, func(i, j int) bool {
			return group[i].In_type_order <= group[j].In_type_order
		})
	}

	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(gTmp))
}

// 根据分组ID获取属性列表
func GetPanelAttrByGroup(gid int32) []*GMazeAttrListOrderV8Cfg.MazeAttrListOrderV8ConfigRow {
	return gConfigData.DollAttrGroupMap[gid]
}
