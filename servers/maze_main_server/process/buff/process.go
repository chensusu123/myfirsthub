package buff

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazebarrieruserkafka"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeTempBuff"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/21 14:16
 * @Description:
 */

func RegTcpHandler() {
	// 查询迷宫buff列表
	websocket_service.RegProcSimple(16242, &MazeTempBuff.GetMazeTempBuffListRQ{},
		16243, &MazeTempBuff.GetMazeTempBuffListRS{}, GetMazeTempBuffListRQ)

	// 查询迷宫可选buff列表
	websocket_service.RegProcSimple(16244, &MazeTempBuff.GetOptionalMazeTempBuffListRQ{},
		16245, &MazeTempBuff.GetOptionalMazeTempBuffListRS{}, GetOptionalMazeTempBuffListRQ)

	// 刷新迷宫可选buff列表
	websocket_service.RegProcSimple(16263, &MazeTempBuff.RefreshOptionalMazeTempBuffListRQ{},
		16264, &MazeTempBuff.RefreshOptionalMazeTempBuffListRS{}, RefreshOptionalMazeTempBuffListRQ)

	// 选择迷宫buff
	websocket_service.RegProcSimple(16246, &MazeTempBuff.SelectMazeTempBuffRQ{},
		16247, &MazeTempBuff.SelectMazeTempBuffRS{}, SelectMazeTempBuffRQ)
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
