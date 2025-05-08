package collect

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazebarrieruserkafka"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazeuserlevelkafka"
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
