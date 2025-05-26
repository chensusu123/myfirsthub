/*
 * @Author: majian
 * @Date: 2024-07-13 15:20:16
 * @Last Modified by: majian
 * @Last Modified time: 2024-07-13 15:25:51
 */
package dollattr

import "maze_game_server/common/constdef"

// 是否人偶计算属性
func IsDollCalcAttr(typ int32) bool {
	return typ == constdef.DollAttrTypeCalc
}

// 是否人偶展示属性
func IsDollShowAttr(typ int32) bool {
	return typ == constdef.DollAttrTypeShow
}

// 是否人偶武力属性
func IsDollForceAttr(typ int32) bool {
	return typ == constdef.DollAttrTypeForce
}
