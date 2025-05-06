/*
 * @Author: majian
 * @Date: 2024-04-25 11:01:05
 * @Last Modified by: majian
 * @Last Modified time: 2024-04-25 11:05:57
 */
package assembleidpack

import "time"

func GetAssembleToken() int64 {
	return time.Now().UnixNano() / 1000000
}
