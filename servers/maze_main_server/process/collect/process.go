package collect

import (
	"maze_game_server/io/kafka/mazebarrieruserkafka"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/usecase/tasktimer"

	"github.com/lonng/nano/component"
)

type Collect struct {
	component.Base
}

func NewCollect() *Collect {
	return &Collect{}
}

func RegWsHandler() {
	// // 迷宫挂机查询
	// _ = websocket_service.RegProcSimple(10465, &MazeCollect.MazeCollectInfoQueryRQ{},
	// 	10466, &MazeCollect.MazeCollectInfoQueryRS{}, OnMazeCollectInfoQueryRQ)

	// _ = websocket_service.RegProcSimple(10467, &MazeCollect.MazeCollectItemReceiveRQ{},
	// 	10468, &MazeCollect.MazeCollectItemReceiveRS{}, OnMazeCollectItemReceiveRQ)
}

func RegTcpHandler() {
	// Timer服务触发接口
	// _ = tcp_service.RegProcSimple(uncgkconst.UN_TCP_PACK_SVR_SEA_TASK_EXPIRE_NOTIFY_RQ, &SeaTaskSvr.TaskExpireNotifyRQ{},
	// 	uncgkconst.UN_TCP_PACK_SVR_SEA_TASK_EXPIRE_NOTIFY_RS, &SeaTaskSvr.TaskExpireNotifyRS{}, OnTimeOut)

	// 本地定时器
	tasktimer.RegOnTimeoutFunc(ProcessTimeOut)
}

func RegConsumeHandler() {
	// kafka_consumer.PlugKafkaConsumer("maze_barrier_chg_msg",
	// 	1001105,
	// 	kafka_consumer.WithGroup(fkserver.GroupNameGO+"."+fkserver.ProjectNamePPWD+".maze_collect_server"),
	// 	kafka_consumer.WithKafkaCustomKeyContent(HandleMazeBarrierMsg))
	mazebarrieruserkafka.Watch(HandleMazeBarrierMsg)

	// kafka_consumer.PlugKafkaConsumer("maze_level_chg_msg",
	// 	1001084,
	// 	kafka_consumer.WithGroup(fkserver.GroupNameGO+"."+fkserver.ProjectNamePPWD+".maze_collect_server"),
	// 	kafka_consumer.WithKafkaCustomKeyContent(HandleMazeLevelMsg))
	mazeuserlevelkafka.Watch(HandleMazeLevelMsg)
}
