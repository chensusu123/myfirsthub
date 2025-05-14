package equip

import (
	"fmt"

	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/limiter"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazeuserlevelkafka"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipPos"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGameEquip"
)

var GtcpLimiter = limiter.NewLimiter("tcpLimiter")

func MakeLimiterKey(uid uint64, packId int32) string {
	return fmt.Sprintf("%d_%d", uid, packId)
}

func RegTcpHandler() {
	// 查询迷宫装配信息
	_ = websocket_service.RegProcSimple(16179, &MazeGameEquip.GetMazeGameAssembleInfoRQ{},
		16180, &MazeGameEquip.GetMazeGameAssembleInfoRS{}, OnGetMazeAssembleRQ)

	// 更换装备预览
	_ = websocket_service.RegProcSimple(16183, &MazeGameEquip.MazeDressEquipPreviewRQ{},
		16184, &MazeGameEquip.MazeDressEquipPreviewRS{}, OnDressEquipPreviewRQ)

	// 更换装备
	_ = websocket_service.RegProcSimple(16181, &MazeGameEquip.MazeDressEquipRQ{},
		16182, &MazeGameEquip.MazeDressEquipRS{}, OnDressMazeEquipRQ)

	// 多选一穿装备
	_ = websocket_service.RegProcSimple(16222, &MazeGameEquip.SelectDressMazeEquipRQ{},
		16223, &MazeGameEquip.SelectDressMazeEquipRS{}, OnSelectDressMazeEquipRQ)

	// 处理kafkatcp消息
	// _ = websocket_service.RegProcSimple(20989, &KafkaMsgNotify.KafkaMsgDistributeRQ{},
	// 	20990, &KafkaMsgNotify.KafkaMsgDistributeRS{}, OnKafkaTcpMsgRQ)

	// 拉取背包装备列表
	_ = websocket_service.RegProcSimple(16190, &MazeGameEquip.GetMazeBagEquipListRQ{},
		16191, &MazeGameEquip.GetMazeBagEquipListRS{}, OnGetMazeBagEquipListRQ)

	_ = websocket_service.RegProcSimple(16192, &MazeGameEquip.QueryMazeEquipDetailRQ{},
		16193, &MazeGameEquip.QueryMazeEquipDetailRS{}, OnQueryMazeEquipDetailRQ)

	// 装备分解
	_ = websocket_service.RegProcSimple(16188, &MazeGameEquip.MazeEquipDismantleRQ{},
		16189, &MazeGameEquip.MazeEquipDismantleRS{}, OnDollEquipDismantleRQ)

	// 装备位强化预览
	websocket_service.RegProcSimple(16201, &MazeEquipPos.MazeEquipPosLvUpPreviewRQ{},
		16202, &MazeEquipPos.MazeEquipPosLvUpPreviewRS{}, OnEquipPosLvUpPreviewRQ)

	// 装备位强化
	websocket_service.RegProcSimple(16203, &MazeEquipPos.MazeEquipPosLvUpRQ{},
		16204, &MazeEquipPos.MazeEquipPosLvUpRS{}, OnEquipPosLvUpRQ)
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
