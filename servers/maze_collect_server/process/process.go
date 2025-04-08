package process

import (
	"gitlab.ifreetalk.com/plate/definition/uncgkconst"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/plate/protodef/MazeCollect"
	"gitlab.ifreetalk.com/plate/protodef/SeaTaskSvr"
)

func RegTcpHandler() {
	// 迷宫挂机查询
	_ = tcp_service.RegProcSimple(16257, &MazeCollect.MazeCollectInfoQueryRQ{},
		16258, &MazeCollect.MazeCollectInfoQueryRS{}, OnMazeCollectInfoQueryRQ)

	_ = tcp_service.RegProcSimple(16259, &MazeCollect.MazeCollectItemReceiveRQ{},
		16260, &MazeCollect.MazeCollectItemReceiveRS{}, OnMazeCollectItemReceiveRQ)

	_ = tcp_service.RegProcSimple(uncgkconst.UN_TCP_PACK_SVR_SEA_TASK_EXPIRE_NOTIFY_RQ, &SeaTaskSvr.TaskExpireNotifyRQ{},
		uncgkconst.UN_TCP_PACK_SVR_SEA_TASK_EXPIRE_NOTIFY_RS, &SeaTaskSvr.TaskExpireNotifyRS{}, OnTimeOut)
}
