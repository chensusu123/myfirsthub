package main

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/game"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/exportlogservice"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/kafka_consumer"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/web_service"
	"gitlab.ifreetalk.com/plate/io_interface/kafka_interface/exportlogkafka"
)

// 19987	UN_CGK_SVR_TYPE_MAZE_MAIN_SERVER 小程序版迷宫主服务
func main() {
	fkserver.SetMonitorName(fkserver.GroupNameGO, fkserver.ProjectNamePPWD, "maze_main_server")
	exportlogservice.PlugExportLogService(exportlogkafka.GetProducer())

	tcp_service.PlugTcpService(process.RegTcpHandler)

	web_service.PlugWebService(regWebHandler)
	// // 1001071 topic-doll-attr-chg-notify-msg 人偶属性变化通知消息
	// kafka_consumer.PlugKafkaConsumer("user_attr_chg_msg",
	// 	1001071,
	// 	kafka_consumer.WithGroup(fkserver.MonitorName),
	// 	kafka_consumer.WithKafkaCustomKeyContent(process.HandleUserAttrMsg))
	// 1001083 topic-maze-attr-chg-notify-msg 迷宫属性变化通知消息
	kafka_consumer.PlugKafkaConsumer("maze_user_attr_chg_msg",
		1001083,
		kafka_consumer.WithGroup(fkserver.MonitorName),
		kafka_consumer.WithKafkaCustomKeyContent(game.HandleUserAttrMsg))

	kafka_consumer.PlugKafkaConsumer("maze_temp_buff_chg_msg",
		1001104,
		kafka_consumer.WithGroup(fkserver.MonitorName),
		kafka_consumer.WithKafkaCustomKeyContent(game.HandleTempBuffMsg))

	fkserver.Run()
}

func regWebHandler(logger fklog.FKLogI) {

}
