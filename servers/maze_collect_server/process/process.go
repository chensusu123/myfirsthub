package process

import (
	"gitlab.ifreetalk.com/plate/definition/uncgkconst"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
	"gitlab.ifreetalk.com/plate/protodef/MazeCollect"
	"gitlab.ifreetalk.com/plate/protodef/SeaTaskSvr"
)

func RegTcpHandler() {
	// 迷宫挂机查询
	_ = websocket_service.RegProcSimple(16257, &MazeCollect.MazeCollectInfoQueryRQ{},
		16258, &MazeCollect.MazeCollectInfoQueryRS{}, OnMazeCollectInfoQueryRQ)

	_ = websocket_service.RegProcSimple(16259, &MazeCollect.MazeCollectItemReceiveRQ{},
		16260, &MazeCollect.MazeCollectItemReceiveRS{}, OnMazeCollectItemReceiveRQ)

	_ = websocket_service.RegProcSimple(uncgkconst.UN_TCP_PACK_SVR_SEA_TASK_EXPIRE_NOTIFY_RQ, &SeaTaskSvr.TaskExpireNotifyRQ{},
		uncgkconst.UN_TCP_PACK_SVR_SEA_TASK_EXPIRE_NOTIFY_RS, &SeaTaskSvr.TaskExpireNotifyRS{}, OnTimeOut)
}
