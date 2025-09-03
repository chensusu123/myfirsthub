package flowutil

import (
	"fmt"
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
	fmt.Println(res)
	return res
}
