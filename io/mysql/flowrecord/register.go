package flowrecord

import (
	"maze_game_server/io/kafka/dollequipassmeblekakfa"
	"maze_game_server/io/kafka/dollequipdismantlekafka"
	"maze_game_server/io/kafka/equipposstrengrecordkafka"
	"maze_game_server/io/kafka/mazeattrchgrecord"
	"maze_game_server/io/kafka/mazebarrieruserkafka"
	"maze_game_server/io/kafka/mazecollectrecord"
	"maze_game_server/io/kafka/mazeequipbagrecord"
	"maze_game_server/io/kafka/mazemoneykafka"
	"maze_game_server/io/kafka/mazerebornkafka"
	"maze_game_server/io/kafka/mazetempbuffchgmsg"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
)

// 注册流水处理
func RegConsumeHandler() {
	mazeattrchgrecord.Watch(SaveAttrChgRecord)
	mazebarrieruserkafka.Watch(SaveBarrierUserRecord)
	mazecollectrecord.Watch(SaveCollectChgRecord)
	dollequipassmeblekakfa.Watch(SaveEquipAssembleRecord)
	mazeequipbagrecord.Watch(SaveEquipBagRecord)
	dollequipdismantlekafka.Watch(SaveEquipDismantRecord)
	equipposstrengrecordkafka.Watch(SaveEquipPosStrengRecord)
	mazemoneykafka.Watch(SaveMoneyRecord)
	mazerebornkafka.Watch(SaveRebornRecord)
	mazebarrieruserkafka.Watch(SaveSweepRecord)
	mazetempbuffchgmsg.Watch(SaveTempBuffChgRecord)
	mazeuserlevelkafka.Watch(SaveUserLevelRecord)
}
