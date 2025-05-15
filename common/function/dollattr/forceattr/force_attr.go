/*
 * @Author: majian
 * @Date: 2024-11-11 16:35:05
 * @Last Modified by: majian
 * @Last Modified time: 2025-04-01 15:15:19
 */

package forceattr

import (
	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeAttributeFormulaV8Cfg"
)

// 根据配表判断是否是展示武力值属性
func GetShowForceAttrs(AttrId int32) map[int32]struct{} {
	resultMap := make(map[int32]struct{})
	row := GMazeAttributeFormulaV8Cfg.GetMazeAttributeFormulaV8Config(801)
	if row != nil {
		resultMap[row.Parameter_1] = struct{}{}
		resultMap[row.Parameter_2] = struct{}{}
		for _, attr := range row.Parameter_8 {
			resultMap[attr] = struct{}{}
		}
		resultMap[row.Parameter_3] = struct{}{}
		resultMap[row.Parameter_4] = struct{}{}
		resultMap[row.Parameter_5] = struct{}{}
		for _, attr := range row.Parameter_9 {
			resultMap[attr] = struct{}{}
		}
		for _, attr := range row.Parameter_10 {
			resultMap[attr] = struct{}{}
		}
		resultMap[row.Parameter_11] = struct{}{}
		resultMap[row.Parameter_12] = struct{}{}
	}
	return resultMap
}

func IsShowForceParam(attrId int32) bool {
	return attrId/1000 == 9
}
