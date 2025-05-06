/*
 * @Author: majian
 * @Date: 2024-08-19 21:12:53
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-17 15:07:37
 * @Desc 传奇属性
 */
package calcassembleattr

import (
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeEquipSuiteAttrV8Cfg"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/maputil"
)

func LegendSuitKey(suitId, cnt int32) int32 {
	return suitId*10000 + cnt
}

// 计算传奇套装属性加成
func CalcLegendSuitBuff(in map[int32]int32) (buffs, showBuffs map[int32]int64) {
	buffs = make(map[int32]int64)
	showBuffs = make(map[int32]int64)
	for k, v := range in {
		row := GetLeastRow(k, v)
		if row == nil {
			continue
		}
		buffs = maputil.Int64MapAppend(buffs, row.Add_attr)
		showBuffs = maputil.Int64MapAppend(showBuffs, row.Add_attr_show)
	}
	return buffs, showBuffs
}

func GetLeastRow(suitId, n int32) *GMazeEquipSuiteAttrV8Cfg.MazeEquipSuiteAttrV8ConfigRow {
	if suitId <= 0 || n <= 0 {
		return nil
	}
	var tmpRow *GMazeEquipSuiteAttrV8Cfg.MazeEquipSuiteAttrV8ConfigRow
	for _, row := range GMazeEquipSuiteAttrV8Cfg.GetAll() {
		if row.Suite_id != suitId {
			continue
		}
		if row.Affix_num > n {
			continue
		}
		if tmpRow == nil {
			tmpRow = row
		} else if tmpRow.Affix_num < row.Affix_num {
			tmpRow = row
		}
	}
	return tmpRow
}
