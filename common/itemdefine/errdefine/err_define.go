/*
@Author: xiaobo
@Date: 2023/11/30 14:14
@Description: 错误码定义
*/

package errdefine

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/common/errors"
	"gitlab.ifreetalk.com/maze-plate/protodef/MessageType"
	"strings"
)

// IsTimeOut 是否超时
func IsTimeOut(errMsg string) bool {
	return strings.Contains(errMsg, "i/o timeout")
}

// IsBreakConn 断连接
func IsBreakConn(errMsg string) bool {
	return strings.Contains(errMsg, "EOF")
}

// IsCheckErrorCode 检查阶段错误,检查加和检查扣
func IsCheckErrorCode(defaultErr *MessageType.ErrorInfo) bool {
	errCode := defaultErr.GetErrCode()
	if errCode == errors.ITEM_CHECK_TIME_OUT.Code || // 道具超时
		errCode == errors.ItemCheckFailure.Code || // 检查添加上限
		errCode == errors.ITEM_CHECK_ERROR.Code || // 检查逻辑错误,比如redis异常，添加金币、钻石等没有添加映射
		errCode == errors.ITEM_CHECK_ITEM_NOT_ENOUGH.Code { // 检查扣道具不足
		return true
	}
	return false
}
