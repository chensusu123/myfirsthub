package flowmodel

import (
	"context"
)

var MAXFLOWCHANSIZE = 50000
var WARNFLOWCHANSIZE = 10000
var ERRORFLOWCHANSIZE = 30000

// 系统内流水数据传递使用
type FlowData struct {
	Ctx  context.Context
	Data interface{}
}

func NewFlowData(ctx context.Context, data interface{}) *FlowData {
	return &FlowData{
		Ctx:  ctx,
		Data: data,
	}
}
