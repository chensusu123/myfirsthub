/*
 * @Author: majian
 * @Date: 2024-09-13 22:17:44
 * @Last Modified by: majian
 * @Last Modified time: 2024-09-13 22:18:11
 */
package timeutil

import "time"

// 时间格式化输出
func FtTime(t int64) string {
	return time.Unix(int64(t), 0).Format("2006-01-02 15:04:05")
}

func GetTodayZeroTimestamp() int64 {
	t := time.Now()
	zeroTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return zeroTime.Unix()
}
