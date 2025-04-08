package main

import (
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/exportlogservice"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/kafka_consumer"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/web_service"
	"gitlab.ifreetalk.com/plate/io_interface/kafka_interface/exportlogkafka"
	"gitlab.ifreetalk.com/servers/maze_game_server/servers/maze_collect_server/process"
)

// 20008	UN_CGK_SVR_TYPE_MAZE_COLLECT_SERVER	0	是	迷宫挂机服务
func main() {
	fkserver.SetMonitorName(fkserver.GroupNameGO, fkserver.ProjectNamePPWD, "maze_collect_server")
	exportlogservice.PlugExportLogService(exportlogkafka.GetProducer())
	tcp_service.PlugTcpService(process.RegTcpHandler)
	web_service.PlugWebService(emptyService)

	kafka_consumer.PlugKafkaConsumer("maze_barrier_chg_msg",
		1001105,
		kafka_consumer.WithGroup(fkserver.MonitorName),
		kafka_consumer.WithKafkaCustomKeyContent(process.HandleMazeBarrierMsg))

	kafka_consumer.PlugKafkaConsumer("maze_level_chg_msg",
		1001084,
		kafka_consumer.WithGroup(fkserver.MonitorName),
		kafka_consumer.WithKafkaCustomKeyContent(process.HandleMazeLevelMsg))
	fkserver.Run()
}

func emptyService(logger fklog.FKLogI) {
	_ = logger
}
