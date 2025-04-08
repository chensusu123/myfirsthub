/*
 * @Author: majian
 * @Date: 2024-11-30 13:50:01
 * @Last Modified by: majian
 * @Last Modified time: 2024-11-30 13:53:32
 */
package toastmsgtipexcel

import "gitlab.ifreetalk.com/plate/excel/auto/GMazeToastMsgInfoCfg"

// 从配表获取提示内容
func GetToastMsgTip(toastId int32, def string) string {
	row := GMazeToastMsgInfoCfg.GetMazeToastMsgInfoConfig(toastId)
	if row != nil {
		return row.Msg_text
	}
	return def
}
