/*
 * @Author: majian
 * @Date: 2025-03-11 16:38:22
 * @Last Modified by: majian
 * @Last Modified time: 2025-04-01 15:18:56
 */
package commonlogic

import (
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeAttributeV8Cfg"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeBuffData"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/simpleset"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/vardef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/dollattr"
)

func GetSrcName(src int) string {
	if desc, ok := vardef.MazeBuffSrcDescMap[int32(src)]; ok {
		return desc
	}
	if desc, ok := DollAttrSrcMapCfg.Load(int32(src)); ok {
		return desc.(string)
	}
	return ""
}

// 判断是否是迷宫计算属性
func IsMazeCalcAttr(attrId int32) bool {
	row := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(attrId)
	if row != nil {
		if row.Is_into_buff == 0 {
			return false
		}
		if dollattr.IsDollCalcAttr(row.Type) ||
			dollattr.IsDollShowAttr(row.Type) ||
			dollattr.IsDollForceAttr(row.Type) {
			return true
		}
	}
	return false
}

func AddtionMazeAttr(in map[int32]int64, attr *MazeBuffData.MazeBuffAttr) {
	in[attr.GetAttrId()] += attr.GetAttrVal()
}

func AddtionMazeAttrKv(in map[int32]int64, attrId int32, val int64) {
	in[attrId] += val
}

// 6 人偶生效属性 7 人偶展示属性 5 武力值
func GetAllMazeAttrByType(attrs *simpleset.Set) []int32 {
	var aimAttrs []int32
	allRows := GMazeAttributeV8Cfg.GetAll()
	for _, attrRow := range allRows {
		if attrRow.Is_into_buff == 0 {
			continue
		}
		if attrs.Has(attrRow.Type) {
			aimAttrs = append(aimAttrs, attrRow.Id)
		}
	}
	return aimAttrs
}
