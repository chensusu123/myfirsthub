package attr_calc

import (
	"gitlab.ifreetalk.com/plate/freetk/fkserver/thrift_service"
	"gitlab.ifreetalk.com/plate/protodef/MazeAttrCalcSvr"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/plate/protodef/MazePropertyPanel"
)

func RegRpcHandler() {
	// 人偶属性预览
	thrift_service.RegisterTwowaySimple(131427, &MazeAttrCalcSvr.MazeAttrPreviewRQ{},
		131428, &MazeAttrCalcSvr.MazeAttrPreviewRS{}, OnMazeAttrPreviewRQ)

	// 人偶属性成对预览
	thrift_service.RegisterTwowaySimple(131429, &MazeAttrCalcSvr.MazeAttrPairPreviewRQ{},
		131430, &MazeAttrCalcSvr.MazeAttrPairPreviewRS{}, OnMazeAttrPairPreviewRQ)
}

func RegTcpHandler() {
	// 查询属性面板
	_ = tcp_service.RegProcSimple(16220, &MazePropertyPanel.QueryMazePropertyPanelRQ{},
		16221, &MazePropertyPanel.QueryMazePropertyPanelRS{}, OnQueryPropertyPanelRQ)
}

func RegConsumeHandler() {
	// _ = redis_consumer.PlugRedisConsumer("doll_calc_attr_consumer",
	// 	21645,
	// 	redis_consumer.WithThreadCount(4),
	// 	redis_consumer.WithErrorWaitTime(500*time.Millisecond),
	// 	redis_consumer.WithListJSONContent(&structsdef.MazeCalcAttrNotifyMsg{}, OnMazeAttrCalcMsg))
}
