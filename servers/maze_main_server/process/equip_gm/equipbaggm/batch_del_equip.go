/*
 * @Author: majian
 * @Date: 2024-11-26 14:51:34
 * @Last Modified by: majian
 * @Last Modified time: 2024-11-26 17:22:52
 */
package equipbaggm

import (
	"strings"

	"gitlab.ifreetalk.com/plate/freetk/fkutil"
)

// 解析规则Id
func ParseRules(in string) (ruleMap map[int32]int32) {
	ruleMap = make(map[int32]int32)
	pairs := strings.Split(in, "_")
	for _, kv := range pairs {
		elems := strings.Split(kv, ":")
		if len(elems) == 2 {
			equipId := fkutil.ToInt32(elems[0])
			ruleId := fkutil.ToInt32(elems[1])
			if equipId > 0 && ruleId > 0 {
				ruleMap[equipId] = ruleId
			}
		}
	}
	return ruleMap
}
