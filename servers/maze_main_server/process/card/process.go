package card

import (
	"maze_game_server/lib/nano/component"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/custom"
)

type Card struct {
	component.Base
}

func NewCard() *Card {
	return &Card{}
}

func RegTcpHandler() {
	// // 获取迷宫月卡
	// websocket_service.RegProcSimple(10430, &MazeCard.GetMazeCardRQ{},
	// 	10431, &MazeCard.GetMazeCardRS{}, GetMazeCardRQ)
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
