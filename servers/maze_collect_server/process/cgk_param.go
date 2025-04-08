package process

import (
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig/param"
)

var (
	ClientFreshDelaySeconds int64 // 客户端刷新延迟秒数，保证定时器回调并产出道具完成
	TimerDelaySeconds       int64 // 定时器多久没回调算超时，必须小于 ClientFreshDelaySeconds
)

func init() {
	param.Int64P(&ClientFreshDelaySeconds, "ClientFreshDelaySeconds", 5, "客户端刷新延迟秒数")
	param.Int64P(&TimerDelaySeconds, "TimerDelaySeconds", 4, "判定定时器丢失的延迟秒数")
}
