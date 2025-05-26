package card

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/custom"
	"maze_game_server/lib/net/websocket_service"
	"maze_game_server/pb/common/MazeCard"
)

func RegTcpHandler() {
	// 获取迷宫月卡
	websocket_service.RegProcSimple(10430, &MazeCard.GetMazeCardRQ{},
		10431, &MazeCard.GetMazeCardRS{}, GetMazeCardRQ)
}

func InitKafkaConsumer() {
	// _ = kafka_consumer.PlugKafkaConsumer("MazeCardNotifyMsg",
	// 	1001098,
	// 	kafka_consumer.WithGroup(fkserver.MonitorName),
	// 	kafka_consumer.WithKafkaCustomKeyContent(OnMazeCardChangeProcess),
	// )
}

func InitCustom() {
	custom.AddCustomProc("DealExpirationCard", 1, DealExpirationCardProcess)
}
