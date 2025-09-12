package flowutil

import (
	"fmt"
	"maze_game_server/services/itemservice"
	"reflect"
	"strings"
)

func Data2Flow(data map[int64]int64) (res string) {
	first := true
	for k, v := range data {
		if first {
			res = fmt.Sprintf("%d:%d", k, v)
			first = false
		} else {
			res = fmt.Sprintf("%s_%d:%d", res, k, v)
		}
	}
	return
}

func UserType2Flow(opts ...any) (res string) {
	var parts []string
	for _, opt := range opts {
		// 将每个参数转换为字符串
		parts = append(parts, fmt.Sprintf("%v", opt))
	}
	// 用下划线拼接所有部分
	res = strings.Join(parts, "_")
	return res
}

func AnyArray2String(opts any) (res string) {
	val := reflect.ValueOf(opts)

	if val.Kind() != reflect.Slice {
		return ""
	}

	var parts []string
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i).Interface()
		parts = append(parts, fmt.Sprintf("%v", elem))
	}

	res = strings.Join(parts, "_")
	return res
}

func ItemInfo2String(optss ...[]*itemservice.ItemInfo) (res string) {
	var parts []string
	for _, opts := range optss {
		for _, opt := range opts {
			parts = append(parts, fmt.Sprintf("%d:%d", opt.ItemId, opt.Count))
		}
	}
	return strings.Join(parts, "_")
}
