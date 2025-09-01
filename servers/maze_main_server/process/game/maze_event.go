package game

import (
	"encoding/json"
	"maze_game_server/common/errors"
	"maze_game_server/io/redis/mazebarriereventredis"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/flowmodel/reportdatamodel"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/servers/maze_main_server/process/game/events"
	"maze_game_server/services/flowservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// OnMazeReportBattleEventRQ 关卡事件上报
func (g *Game) OnMazeReportBattleEventRQ_10496_10497(s *session.Session, req *MazeGame.ReportBattleEventRQ) (err error) {
	defer fkprometheus.InfoPMT("OnMazeReportBattleEventRQ")()

	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)

	res := &MazeGame.ReportBattleEventRS{}
	logger.CtxInfo(ctx, "OnMazeReportBattleEventRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnMazeReportBattleEventRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userID := uint64(s.UID())

	for _, event := range req.Events {
		eventType := event.GetType()
		eventData := (proto.Message)(nil)
		triggerFn := (func())(nil)
		switch eventType {
		// 攻击怪物
		case MazeGame.BattleEventType_ATTACK_MONSTER:
			eventData = event.GetAttack()
			triggerFn = func() {
				events.OnAttackMonster(s.Context(), userID, event.GetEventFrame(), event.GetEventTimeMs(), event.GetAttack())
			}
		// 被怪物攻击
		case MazeGame.BattleEventType_BE_ATTACKED:
			eventData = event.GetAttack()
			triggerFn = func() {
				events.OnBeAttacked(s.Context(), userID, event.GetEventFrame(), event.GetEventTimeMs(), event.GetAttack())
			}
		// 怪物死亡
		case MazeGame.BattleEventType_MONSTER_DEAD:
			eventData = event.GetMonsterDead()
			triggerFn = func() {
				events.OnMonsterDead(s.Context(), userID, event.GetEventFrame(), event.GetEventTimeMs(), event.GetMonsterDead())
			}
		// 刷新怪物
		case MazeGame.BattleEventType_REFRESH_MONSTER:
			eventData = event.GetRefreshMonster()
			triggerFn = func() {
				events.OnRefreshMonster(s.Context(), userID, event.GetEventFrame(), event.GetEventTimeMs(), event.GetRefreshMonster())
			}
		// 触发机关
		case MazeGame.BattleEventType_TRIGGER_TRAP:
			eventData = event.GetTriggerTrap()
			triggerFn = func() {
				events.OnTriggerTrap(s.Context(), userID, event.GetEventFrame(), event.GetEventTimeMs(), event.GetTriggerTrap())
			}
		// 人物移动
		case MazeGame.BattleEventType_ROLE_MOVE:
			eventData = event.GetRoleMove()
			triggerFn = func() {
				events.OnRoleMove(s.Context(), userID, event.GetEventFrame(), event.GetEventTimeMs(), event.GetRoleMove())
			}
		// 猪妖移动
		case MazeGame.BattleEventType_BOSS_MOVE:
			eventData = event.GetBossMove()
			triggerFn = func() {
				events.OnBossMove(s.Context(), userID, event.GetEventFrame(), event.GetEventTimeMs(), event.GetBossMove())
			}
		// 恢复/暂停游戏
		case MazeGame.BattleEventType_PAUSE_GAME:
			eventData = event.GetPauseGame()
			triggerFn = func() {
				events.OnPauseGame(s.Context(), userID, event.GetEventFrame(), event.GetEventTimeMs(), event.GetPauseGame())
			}
		}
		if eventData == nil {
<<<<<<< HEAD
			logger.CtxError(s.Context(), "OnMazeReportBattleEventRQ event type not supported", zap.Int32("eventType", int32(eventType)))
		} else {
			err = mazebarriereventredis.TriggerBarrierEvent(s.Context(), userID, event.GetEventFrame(), event.GetEventTimeMs(), eventType, eventData)
			if err != nil {
				logger.CtxError(s.Context(), "OnMazeReportBattleEventRQ TriggerBarrierEvent fail", zap.Error(err), zap.Int32("eventType", int32(eventType)))
=======
			logger.CtxError(ctx, "OnMazeReportBattleEventRQ event type not supported", zap.Int32("eventType", int32(eventType)))
		} else {
			err = mazebarriereventredis.TriggerBarrierEvent(ctx, userID, event.GetEventFrame(), event.GetEventTimeMs(), eventType, eventData)
			if err != nil {
				logger.CtxError(ctx, "OnMazeReportBattleEventRQ TriggerBarrierEvent fail", zap.Error(err), zap.Int32("eventType", int32(eventType)))
>>>>>>> remotes/origin/dev_human_robot_0315_env
			}
			// 触发事件
			triggerFn()
		}
		dataJson, err := json.Marshal(eventData)
		if err != nil {
			logger.CtxError(ctx, "OnMazeReportBattleEventRQ Marshal fail",
				zap.Any("eventDta", eventData),
				zap.Error(err))
			continue
		}
		data := reportdatamodel.NewReportData(ctx, event.GetType(), userID, uint64(event.GetEventFrame()), uint64(event.GetEventTimeMs()), string(dataJson))
		flowservice.GflowService.SendFlowData(ctx, data)
	}
	return
}
