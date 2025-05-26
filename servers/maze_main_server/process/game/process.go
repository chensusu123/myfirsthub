package game

import (
	"github.com/lonng/nano/component"
	"maze_game_server/io/kafka/mazeattrmsg"
	"maze_game_server/io/kafka/mazetempbuffchgmsg"
	"maze_game_server/lib/net/websocket_service"
	"maze_game_server/pb/common/MazeGame"
)

type Game struct {
	component.Base
}

func NewGame() *Game {
	return &Game{}
}
func RegTcpHandler() {
	// 人偶版本新手引导关卡信息查询RQ
	// websocket_service.RegProcSimple(
	// 	16131, &DollMazeBarrier.QueryBarrierRQ{},
	// 	16132, &DollMazeBarrier.QueryBarrierRS{},
	// 	OnQueryBarrierRQ)

	// 人偶版本新手引导查询割草怪列表RQ
	// websocket_service.RegProcSimple(
	// 	16133, &DollMazeBarrier.QueryGeCaoFoeRQ{},
	// 	16134, &DollMazeBarrier.QueryGeCaoFoeRS{},
	// 	OnQueryGeCaoFoeRQ)

	// 人偶版本新手引导通关进入下一关RQ
	// websocket_service.RegProcSimple(
	// 	16135, &DollMazeBarrier.PassBarrierRQ{},
	// 	16136, &DollMazeBarrier.PassBarrierRS{},
	// 	OnPassBarrierRQ)

	// 人偶版本新手引导击败敌人获得奖励RQ
	// websocket_service.RegProcSimple(
	// 	16137, &DollMazeBarrier.GetFoeAwardRQ{},
	// 	16138, &DollMazeBarrier.GetFoeAwardRS{},
	// 	OnGetFeoAwardRQ)

	// 人偶版本新手引导开宝箱RQ
	websocket_service.RegProcSimple(
		10445, &MazeGame.BarrierOpenBoxRQ{},
		10446, &MazeGame.BarrierOpenBoxRS{},
		OnBarrierOpenBoxRQ)

	// 人偶版本迷宫进出关卡RQ
	websocket_service.RegProcSimple(
		10447, &MazeGame.MazeBarrierEnterRQ{},
		10448, &MazeGame.MazeBarrierEnterRS{},
		OnMazeBarrierEnterRQ)

	// 人偶版本迷宫进出关卡区域RQ
	// websocket_service.RegProcSimple(
	// 	16148, &DollMazeBarrier.MoveBarrierAreaRQ{},
	// 	16149, &DollMazeBarrier.MoveBarrierAreaRS{},
	// 	OnMoveBarrierAreaRQ)

	// 人偶版本迷宫关卡区域心跳RQ
	// websocket_service.RegProcSimple(
	// 	16150, &DollMazeBarrier.BarrierAreaHeartRQ{},
	// 	16151, &DollMazeBarrier.BarrierAreaHeartRS{},
	// 	OnBarrierAreaHeartRQ)

	// 人偶版本迷宫关卡死亡RQ
	websocket_service.RegProcSimple(
		10449, &MazeGame.BarrierDeathRQ{},
		10450, &MazeGame.BarrierDeathRS{},
		OnMazeBarrierDeathRQ)

	// 人偶版本迷宫查询关卡区域RQ
	// websocket_service.RegProcSimple(
	// 	16154, &DollMazeBarrier.QueryBarrierAreaRQ{},
	// 	16155, &DollMazeBarrier.QueryBarrierAreaRS{},
	// 	OnQueryBarrierAreaRQ)

	// timer
	// _ = websocket_service.RegProcSimple(20147, &SeaTaskSvr.TaskExpireNotifyRQ{},
	// 	20148, &SeaTaskSvr.TaskExpireNotifyRS{}, OnTimeOut)

	// 处理装备命令
	_ = websocket_service.RegProcSimple(10463, &MazeGame.SendDollMazeCmdRQ{},
		10464, &MazeGame.SendDollMazeCmdRS{}, OnSendDollMazeCmdRQ)

	// 人偶版本迷宫上报闲置装备数量RQ
	// websocket_service.RegProcSimple(
	// 	16165, &DollMazeBarrier.MazeAreaEquipNumRQ{},
	// 	16166, &DollMazeBarrier.MazeAreaEquipNumRS{},
	// 	OnMazeAreaEquipNumRQ)

	// 迷宫登陆RQ
	websocket_service.RegProcSimple(
		10451, &MazeGame.MazeLoginRQ{},
		10452, &MazeGame.MazeLoginRS{},
		OnMazeLoginRQ)

	// 上报人物等级和关卡
	websocket_service.RegProcSimple(
		10453, &MazeGame.ReportDataRQ{},
		10454, &MazeGame.ReportDataRS{},
		OnReportDataRQ)

	// 关卡事件上报
	websocket_service.RegProcSimple(
		10496, &MazeGame.ReportBattleEventRQ{},
		10497, &MazeGame.ReportBattleEventRS{},
		OnMazeReportBattleEventRQ)

	// 打怪上报申请加装备
	websocket_service.RegProcSimple(
		10455, &MazeGame.ReportAwardFoeEquipRQ{},
		10456, &MazeGame.ReportAwardFoeEquipRS{},
		OnReportAwardFoeEquipRQ)

	// 关卡列表
	websocket_service.RegProcSimple(
		10457, &MazeGame.MazeBarrierListRQ{},
		10458, &MazeGame.MazeBarrierListRS{},
		OnMazeBarrierListRQ)

	// 通关
	websocket_service.RegProcSimple(
		10459, &MazeGame.MazeBarrierPassRQ{},
		10460, &MazeGame.MazeBarrierPassRS{},
		OnMazeBarrierPassRQ)

	// 挑战复活
	websocket_service.RegProcSimple(
		10461, &MazeGame.MazeBarrierRebornRQ{},
		10462, &MazeGame.MazeBarrierRebornRS{},
		OnMazeBarrierRebornRQ)
}

func RegisterRpcPackProcessor() {
	// thrift_service.RegisterTwowaySimple(131421, &MazeEquipSvr.SvrAddMazeEquipRQ{},
	// 	131422, &MazeEquipSvr.SvrAddMazeEquipRS{}, OnSvrAddMazeEquipRQ)
}

func RegConsumeHandler() {
	// // 1001071 topic-doll-attr-chg-notify-msg 人偶属性变化通知消息
	// kafka_consumer.PlugKafkaConsumer("user_attr_chg_msg",
	// 	1001071,
	// 	kafka_consumer.WithGroup(fkserver.GroupNameGO+"."+fkserver.ProjectNamePPWD+".maze_main_server"),
	// 	kafka_consumer.WithKafkaCustomKeyContent(process.HandleUserAttrMsg))
	// 1001083 topic-maze-attr-chg-notify-msg 迷宫属性变化通知消息
	// kafka_consumer.PlugKafkaConsumer("maze_user_attr_chg_msg",
	// 	1001083,
	// 	kafka_consumer.WithGroup(fkserver.GroupNameGO+"."+fkserver.ProjectNamePPWD+".maze_main_server"),
	// 	kafka_consumer.WithKafkaCustomKeyContent(HandleUserAttrMsg))
	mazeattrmsg.Watch(HandleUserAttrMsg)

	// kafka_consumer.PlugKafkaConsumer("maze_temp_buff_chg_msg",
	// 	1001104,
	// 	kafka_consumer.WithGroup(fkserver.GroupNameGO+"."+fkserver.ProjectNamePPWD+".maze_main_server"),
	// 	kafka_consumer.WithKafkaCustomKeyContent(HandleTempBuffMsg))
	mazetempbuffchgmsg.Watch(HandleTempBuffMsg)
}
