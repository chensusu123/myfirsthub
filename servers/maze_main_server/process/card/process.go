package card

import (
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/custom"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/kafka_consumer"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
	"gitlab.ifreetalk.com/plate/protodef/MazeCard"
)

func RegTcpHandler() {
	// 获取迷宫月卡
	websocket_service.RegProcSimple(16198, &MazeCard.GetMazeCardRQ{},
		16199, &MazeCard.GetMazeCardRS{}, GetMazeCardRQ)
}

func InitKafkaConsumer() {
	_ = kafka_consumer.PlugKafkaConsumer("MazeCardNotifyMsg",
		1001098,
		kafka_consumer.WithGroup(fkserver.MonitorName),
		kafka_consumer.WithKafkaCustomKeyContent(OnMazeCardChangeProcess),
	)
}

func InitCustom() {
	custom.AddCustomProc("DealExpirationCard", 1, DealExpirationCardProcess)
}
