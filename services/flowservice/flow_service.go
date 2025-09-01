package flowservice

import (
	"context"
	flowmodel "maze_game_server/model/flowmodel/flow_model"
	"sync"
)

type flowService interface {
	// 发送流水
	SendFlowData(ctx context.Context, record interface{})
	// 处理数据的地方 初始化时创建 后续把所有数据都转化为json来发送
	ProcessFlowData()
}

type service struct {
	sync.RWMutex
	ch         chan *flowmodel.FlowData
	reportTime map[uint64]uint64
}

func NewFlowService() flowService {
	svr := &service{}
	svr.ch = make(chan *flowmodel.FlowData, flowmodel.MAXFLOWCHANSIZE)
	svr.reportTime = make(map[uint64]uint64)
	go svr.ProcessFlowData()
	return svr
}

var GflowService flowService

func init() {
	GflowService = NewFlowService()
}
