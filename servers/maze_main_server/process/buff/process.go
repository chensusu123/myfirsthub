package buff

import (
	"maze_game_server/io/kafka/mazebarrieruserkafka"
	"maze_game_server/lib/nano/component"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/21 14:16
 * @Description:
 */

type Buff struct {
	component.Base
}

func NewBuff() *Buff {
	return &Buff{}
}

func RegTcpHandler() {
	// // 查询迷宫buff列表
	// websocket_service.RegProcSimple(10433, &MazeTempBuff.GetMazeTempBuffListRQ{},
	// 	10434, &MazeTempBuff.GetMazeTempBuffListRS{}, GetMazeTempBuffListRQ)

	// // 查询迷宫可选buff列表
	// websocket_service.RegProcSimple(10435, &MazeTempBuff.GetOptionalMazeTempBuffListRQ{},
	// 	10436, &MazeTempBuff.GetOptionalMazeTempBuffListRS{}, GetOptionalMazeTempBuffListRQ)

	// // 刷新迷宫可选buff列表
	// websocket_service.RegProcSimple(10439, &MazeTempBuff.RefreshOptionalMazeTempBuffListRQ{},
	// 	10440, &MazeTempBuff.RefreshOptionalMazeTempBuffListRS{}, RefreshOptionalMazeTempBuffListRQ)

	// // 选择迷宫buff
	// websocket_service.RegProcSimple(10437, &MazeTempBuff.SelectMazeTempBuffRQ{},
	// 	10438, &MazeTempBuff.SelectMazeTempBuffRS{}, SelectMazeTempBuffRQ)
}

func InitKafkaConsumer() {
	// 1001105 topic-maze-barrier-user-record 用户迷宫闯关纪录
	// _ = kafka_consumer.PlugKafkaConsumer("MazeBarrierNotify",
	// 	1001105,
	// 	kafka_consumer.WithGroup(fkserver.MonitorName),
	// 	kafka_consumer.WithKafkaCustomKeyContent(MazeBarrierNotifyProcess),
	// )

	mazebarrieruserkafka.Watch(MazeBarrierNotifyProcess)

}
