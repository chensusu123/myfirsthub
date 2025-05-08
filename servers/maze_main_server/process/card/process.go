package process

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/custom"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCard"
)

func InitTcp() {
	// 获取迷宫月卡
	tcp_service.RegProcSimple(10430, &MazeCard.GetMazeCardRQ{},
		10431, &MazeCard.GetMazeCardRS{}, GetMazeCardRQ)
}

func InitKafkaConsumer() {
	//_ = kafka_consumer.PlugKafkaConsumer("MazeCardNotifyMsg",
	//	0,
	//	kafka_consumer.WithGroup(fkserver.MonitorName),
	//	kafka_consumer.WithKafkaCustomKeyContent(OnMazeCardChangeProcess),
	//)
}

func InitCustom() {
	custom.AddCustomProc("DealExpirationCard", 1, DealExpirationCardProcess)
}
