// @Author pangchenyang 2025/5/7 19:50:00
// @Desc: 
package kafka_dispatch

import (
	"gitlab.ifreetalk.com/plate/freetk/fkserver/kafka_consumer"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
)

func RegConsumeHandler() {
	// 性别变化流水
	_ = kafka_consumer.PlugKafkaConsumer(constdef.KafkaMDTSexDesc,
		1000159,
		kafka_consumer.WithGroup(fkserver.GroupNameGO+"."+fkserver.ProjectNamePPWD+"."+constdef.KafkaMDTSexDesc),
		kafka_consumer.WithKafkaCustomKeyContent(HandleSexMsg))
}
