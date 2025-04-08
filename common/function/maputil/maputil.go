/*
 * @Author: majian
 * @Date: 2024-08-19 21:21:32
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-17 20:52:25
 */
package maputil

import (
	"bytes"
	"fmt"
	"sort"

	"gitlab.ifreetalk.com/maze/maze_game_server/common/structsdef"
)

func Int64MapAppend32(in map[int32]int64, adds map[int32]int32) map[int32]int64 {
	for k, v := range adds {
		if k <= 0 { //忽略key=0 v 不能忽略
			continue
		}
		if src, ok := in[k]; ok {
			in[k] = int64(v) + src
		} else {
			in[k] = int64(v)
		}
	}
	return in
}

func Int64MapAppend(in, adds map[int32]int64) map[int32]int64 {
	for k, v := range adds {
		if k <= 0 { //忽略key=0 v 不能忽略
			continue
		}
		if src, ok := in[k]; ok {
			in[k] = v + src
		} else {
			in[k] = v
		}
	}
	return in
}

// map数据转字符串，一般落流水用
func MapToString(items map[int32]int64) string {
	var kvp structsdef.KvPairs
	for k, v := range items {
		kvp = append(kvp, &structsdef.KvPair{K: k, V: v})
	}
	sort.Slice(kvp, func(i, j int) bool {
		return kvp[i].K <= int32(kvp[j].K)
	})

	var bs bytes.Buffer
	for i, kv := range kvp {
		if i != len(kvp)-1 {
			bs.WriteString(fmt.Sprintf("%d:%d,", kv.K, kv.V))
		} else {
			bs.WriteString(fmt.Sprintf("%d:%d", kv.K, kv.V))
		}
	}
	return bs.String()
}

func MapToString32(items map[int32]int32) string {
	item64 := make(map[int32]int64)
	for k, v := range items {
		item64[k] = int64(v)
	}
	return MapToString(item64)
}
