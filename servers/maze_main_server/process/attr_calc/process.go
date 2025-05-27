package attr_calc

import (
	"github.com/lonng/nano/component"
)

func RegRpcHandler() {
	// // 人偶属性预览
	// thrift_service.RegisterTwowaySimple(100000, &MazeAttrCalcSvr.MazeAttrPreviewRQ{},
	// 	100001, &MazeAttrCalcSvr.MazeAttrPreviewRS{}, OnMazeAttrPreviewRQ)

	// // 人偶属性成对预览
	// thrift_service.RegisterTwowaySimple(100002, &MazeAttrCalcSvr.MazeAttrPairPreviewRQ{},
	// 	100003, &MazeAttrCalcSvr.MazeAttrPairPreviewRS{}, OnMazeAttrPairPreviewRQ)
}

type Property struct {
	component.Base
}

func NewProperty() *Property {
	return &Property{}
}

func RegTcpHandler() {
	// // 查询属性面板
	// _ = websocket_service.RegProcSimple(10427, &MazePropertyPanel.QueryMazePropertyPanelRQ{},
	// 	10428, &MazePropertyPanel.QueryMazePropertyPanelRS{}, OnQueryPropertyPanelRQ)
}

func RegConsumeHandler() {
	// _ = redis_consumer.PlugRedisConsumer("doll_calc_attr_consumer",
	// 	21645,
	// 	redis_consumer.WithThreadCount(4),
	// 	redis_consumer.WithErrorWaitTime(500*time.Millisecond),
	// 	redis_consumer.WithListJSONContent(&structsdef.MazeCalcAttrNotifyMsg{}, OnMazeAttrCalcMsg))
}
