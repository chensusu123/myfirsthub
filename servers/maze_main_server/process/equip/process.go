package equip

import (
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/plate/protodef/KafkaMsgNotify"
	"gitlab.ifreetalk.com/plate/protodef/MazeGameEquip"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/kafka_consumer"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/limiter"
	"fmt"
)

var GtcpLimiter = limiter.NewLimiter("tcpLimiter")

func MakeLimiterKey(uid uint64, packId int32) string {
	return fmt.Sprintf("%d_%d", uid, packId)
}

func RegTcpHandler() {
	// 查询迷宫装配信息
	_ = tcp_service.RegProcSimple(16179, &MazeGameEquip.GetMazeGameAssembleInfoRQ{},
		16180, &MazeGameEquip.GetMazeGameAssembleInfoRS{}, OnGetMazeAssembleRQ)

	// 更换装备预览
	_ = tcp_service.RegProcSimple(16183, &MazeGameEquip.MazeDressEquipPreviewRQ{},
		16184, &MazeGameEquip.MazeDressEquipPreviewRS{}, OnDressEquipPreviewRQ)

	// 更换装备
	_ = tcp_service.RegProcSimple(16181, &MazeGameEquip.MazeDressEquipRQ{},
		16182, &MazeGameEquip.MazeDressEquipRS{}, OnDressMazeEquipRQ)

	// 多选一穿装备
	_ = tcp_service.RegProcSimple(16222, &MazeGameEquip.SelectDressMazeEquipRQ{},
		16223, &MazeGameEquip.SelectDressMazeEquipRS{}, OnSelectDressMazeEquipRQ)

	// 处理kafkatcp消息
	_ = tcp_service.RegProcSimple(20989, &KafkaMsgNotify.KafkaMsgDistributeRQ{},
		20990, &KafkaMsgNotify.KafkaMsgDistributeRS{}, OnKafkaTcpMsgRQ)
}

func RegRpcHandler() {

}

func RegConsumeHandler() {
	_ = kafka_consumer.PlugKafkaConsumer("maze_lv_chg",
		1001084,
		kafka_consumer.WithGroup(fkserver.MonitorName),
		kafka_consumer.WithKafkaCustomKeyContent(HandleMazeLvChg))
}

func WebHandler(logger fklog.FKLogI) {

}
