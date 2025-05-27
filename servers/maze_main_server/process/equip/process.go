package equip

import (
	"fmt"
	"maze_game_server/common/function/limiter"
	"maze_game_server/io/kafka/mazeuserlevelkafka"

	"github.com/lonng/nano/component"
)

var GtcpLimiter = limiter.NewLimiter("tcpLimiter")

func MakeLimiterKey(uid uint64, packId int32) string {
	return fmt.Sprintf("%d_%d", uid, packId)
}

type Equip struct {
	component.Base
}

func NewEquip() *Equip {
	return &Equip{}
}

func RegTcpHandler() {
	// // 查询迷宫装配信息
	// _ = websocket_service.RegProcSimple(10414, &MazeGameEquip.GetMazeGameAssembleInfoRQ{},
	// 	10415, &MazeGameEquip.GetMazeGameAssembleInfoRS{}, OnGetMazeAssembleRQ)

	// // 更换装备预览
	// _ = websocket_service.RegProcSimple(10416, &MazeGameEquip.MazeDressEquipPreviewRQ{},
	// 	10417, &MazeGameEquip.MazeDressEquipPreviewRS{}, OnDressEquipPreviewRQ)

	// // 更换装备
	// _ = websocket_service.RegProcSimple(10418, &MazeGameEquip.MazeDressEquipRQ{},
	// 	10419, &MazeGameEquip.MazeDressEquipRS{}, OnDressMazeEquipRQ)

	// // 多选一穿装备
	// _ = websocket_service.RegProcSimple(10420, &MazeGameEquip.SelectDressMazeEquipRQ{},
	// 	10421, &MazeGameEquip.SelectDressMazeEquipRS{}, OnSelectDressMazeEquipRQ)

	// // 处理kafkatcp消息
	// // _ = websocket_service.RegProcSimple(20989, &KafkaMsgNotify.KafkaMsgDistributeRQ{},
	// // 	20990, &KafkaMsgNotify.KafkaMsgDistributeRS{}, OnKafkaTcpMsgRQ)

	// // 拉取背包装备列表
	// _ = websocket_service.RegProcSimple(10405, &MazeGameEquip.GetMazeBagEquipListRQ{},
	// 	10406, &MazeGameEquip.GetMazeBagEquipListRS{}, OnGetMazeBagEquipListRQ)

	// _ = websocket_service.RegProcSimple(10407, &MazeGameEquip.QueryMazeEquipDetailRQ{},
	// 	10408, &MazeGameEquip.QueryMazeEquipDetailRS{}, OnQueryMazeEquipDetailRQ)

	// // 装备分解
	// _ = websocket_service.RegProcSimple(10410, &MazeGameEquip.MazeEquipDismantleRQ{},
	// 	10411, &MazeGameEquip.MazeEquipDismantleRS{}, OnDollEquipDismantleRQ)

	// // 装备位强化预览
	// websocket_service.RegProcSimple(10423, &MazeEquipPos.MazeEquipPosLvUpPreviewRQ{},
	// 	10424, &MazeEquipPos.MazeEquipPosLvUpPreviewRS{}, OnEquipPosLvUpPreviewRQ)

	// // 装备位强化
	// websocket_service.RegProcSimple(10425, &MazeEquipPos.MazeEquipPosLvUpRQ{},
	// 	10426, &MazeEquipPos.MazeEquipPosLvUpRS{}, OnEquipPosLvUpRQ)
}

func RegRpcHandler() {
	// thrift_service.RegisterTwowaySimple(131421, &MazeEquipSvr.SvrAddMazeEquipRQ{},
	// 	131422, &MazeEquipSvr.SvrAddMazeEquipRS{}, OnSvrAddMazeEquipRQ)
	// thrift_service.RegisterTwowaySimple(131423, &MazeEquipSvr.SvrMazeEquipAssembleRQ{},
	// 	131424, &MazeEquipSvr.SvrMazeEquipAssembleRS{}, OnSvrMazeEquipAssembleRQ)
	// thrift_service.RegisterTwowaySimple(131425, &MazeEquipSvr.SvrMazeEquipSaleRQ{},
	// 	131426, &MazeEquipSvr.SvrMazeEquipSaleRS{}, OnSvrDollEquipSaleRQ)
}

func RegConsumeHandler() {
	// _ = kafka_consumer.PlugKafkaConsumer("maze_lv_chg",
	// 	1001084,
	// 	kafka_consumer.WithGroup(fkserver.GroupNameGO+"."+fkserver.ProjectNamePPWD+".maze_equip_main_server"),
	// 	kafka_consumer.WithKafkaCustomKeyContent(HandleMazeLvChg))
	mazeuserlevelkafka.Watch(HandleMazeLvChg)

	// 性别变化流水
	// _ = kafka_consumer.PlugKafkaConsumer(constdef.KafkaMDTSexDesc,
	// 	1000159,
	// 	kafka_consumer.WithGroup(fkserver.GroupNameGO+"."+fkserver.ProjectNamePPWD+"."+constdef.KafkaMDTSexDesc),
	// 	kafka_consumer.WithKafkaCustomKeyContent(HandleDollSexChg))
}
