package process

import (
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/plate/protodef/MazeGame"
)

func RegTcpHandler() {
	// 人偶版本新手引导关卡信息查询RQ
	// tcp_service.RegProcSimple(
	// 	16131, &DollMazeBarrier.QueryBarrierRQ{},
	// 	16132, &DollMazeBarrier.QueryBarrierRS{},
	// 	OnQueryBarrierRQ)

	// 人偶版本新手引导查询割草怪列表RQ
	// tcp_service.RegProcSimple(
	// 	16133, &DollMazeBarrier.QueryGeCaoFoeRQ{},
	// 	16134, &DollMazeBarrier.QueryGeCaoFoeRS{},
	// 	OnQueryGeCaoFoeRQ)

	// 人偶版本新手引导通关进入下一关RQ
	// tcp_service.RegProcSimple(
	// 	16135, &DollMazeBarrier.PassBarrierRQ{},
	// 	16136, &DollMazeBarrier.PassBarrierRS{},
	// 	OnPassBarrierRQ)

	// 人偶版本新手引导击败敌人获得奖励RQ
	// tcp_service.RegProcSimple(
	// 	16137, &DollMazeBarrier.GetFoeAwardRQ{},
	// 	16138, &DollMazeBarrier.GetFoeAwardRS{},
	// 	OnGetFeoAwardRQ)

	// 人偶版本新手引导开宝箱RQ
	tcp_service.RegProcSimple(
		16208, &MazeGame.BarrierOpenBoxRQ{},
		16209, &MazeGame.BarrierOpenBoxRS{},
		OnBarrierOpenBoxRQ)

	// 人偶版本迷宫进出关卡RQ
	tcp_service.RegProcSimple(
		16210, &MazeGame.MazeBarrierEnterRQ{},
		16211, &MazeGame.MazeBarrierEnterRS{},
		OnMazeBarrierEnterRQ)

	// 人偶版本迷宫进出关卡区域RQ
	// tcp_service.RegProcSimple(
	// 	16148, &DollMazeBarrier.MoveBarrierAreaRQ{},
	// 	16149, &DollMazeBarrier.MoveBarrierAreaRS{},
	// 	OnMoveBarrierAreaRQ)

	// 人偶版本迷宫关卡区域心跳RQ
	// tcp_service.RegProcSimple(
	// 	16150, &DollMazeBarrier.BarrierAreaHeartRQ{},
	// 	16151, &DollMazeBarrier.BarrierAreaHeartRS{},
	// 	OnBarrierAreaHeartRQ)

	// 人偶版本迷宫关卡死亡RQ
	tcp_service.RegProcSimple(
		16212, &MazeGame.BarrierDeathRQ{},
		16213, &MazeGame.BarrierDeathRS{},
		OnMazeBarrierDeathRQ)

	// 人偶版本迷宫查询关卡区域RQ
	// tcp_service.RegProcSimple(
	// 	16154, &DollMazeBarrier.QueryBarrierAreaRQ{},
	// 	16155, &DollMazeBarrier.QueryBarrierAreaRS{},
	// 	OnQueryBarrierAreaRQ)

	// timer
	// _ = tcp_service.RegProcSimple(20147, &SeaTaskSvr.TaskExpireNotifyRQ{},
	// 	20148, &SeaTaskSvr.TaskExpireNotifyRS{}, OnTimeOut)

	// 处理装备命令
	_ = tcp_service.RegProcSimple(16218, &MazeGame.SendDollMazeCmdRQ{},
		16219, &MazeGame.SendDollMazeCmdRS{}, OnSendDollMazeCmdRQ)

	// 人偶版本迷宫上报闲置装备数量RQ
	// tcp_service.RegProcSimple(
	// 	16165, &DollMazeBarrier.MazeAreaEquipNumRQ{},
	// 	16166, &DollMazeBarrier.MazeAreaEquipNumRS{},
	// 	OnMazeAreaEquipNumRQ)

	// 迷宫登陆RQ
	tcp_service.RegProcSimple(
		16206, &MazeGame.MazeLoginRQ{},
		16207, &MazeGame.MazeLoginRS{},
		OnMazeLoginRQ)

	// 上报人物等级和关卡
	tcp_service.RegProcSimple(
		16214, &MazeGame.ReportDataRQ{},
		16215, &MazeGame.ReportDataRS{},
		OnReportDataRQ)

	// 打怪上报申请加装备
	tcp_service.RegProcSimple(
		16216, &MazeGame.ReportAwardFoeEquipRQ{},
		16217, &MazeGame.ReportAwardFoeEquipRS{},
		OnReportAwardFoeEquipRQ)

	// 关卡列表
	tcp_service.RegProcSimple(
		16251, &MazeGame.MazeBarrierListRQ{},
		16252, &MazeGame.MazeBarrierListRS{},
		OnMazeBarrierListRQ)

	// 通关
	tcp_service.RegProcSimple(
		16253, &MazeGame.MazeBarrierPassRQ{},
		16254, &MazeGame.MazeBarrierPassRS{},
		OnMazeBarrierPassRQ)

	// 挑战复活
	tcp_service.RegProcSimple(
		16269, &MazeGame.MazeBarrierRebornRQ{},
		16270, &MazeGame.MazeBarrierRebornRS{},
		OnMazeBarrierRebornRQ)
}

func RegisterRpcPackProcessor() {
	// thrift_service.RegisterTwowaySimple(131421, &MazeEquipSvr.SvrAddMazeEquipRQ{},
	// 	131422, &MazeEquipSvr.SvrAddMazeEquipRS{}, OnSvrAddMazeEquipRQ)
}
