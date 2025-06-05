package game

import (
	"maze_game_server/common/errors"
	"maze_game_server/io/redis/mazebarriereventredis"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeGame"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// OnMazeReportBattleEventRQ 关卡事件上报
func (g *Game) OnMazeReportBattleEventRQ_10496_10497(s *session.Session, req *MazeGame.ReportBattleEventRQ) (err error) {
	defer fkprometheus.InfoPMT("OnMazeReportBattleEventRQ")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeGame.ReportBattleEventRS{}

	logger.InfoWF("OnMazeReportBattleEventRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnMazeReportBattleEventRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userID := uint64(s.UID())

	for _, event := range req.Events {
		eventType := event.GetType()
		eventData := (proto.Message)(nil)
		switch eventType {
		// 攻击怪物
		case MazeGame.BattleEventType_ATTACK_MONSTER:
			eventData = event.GetAttack()
		// 被怪物攻击
		case MazeGame.BattleEventType_BE_ATTACKED:
			eventData = event.GetAttack()
		// 怪物死亡
		case MazeGame.BattleEventType_MONSTER_DEAD:
			eventData = event.GetMonsterDead()
		// 刷新怪物
		case MazeGame.BattleEventType_REFRESH_MONSTER:
			eventData = event.GetRefreshMonster()
		// 触发机关
		case MazeGame.BattleEventType_TRIGGER_TRAP:
			eventData = event.GetTriggerTrap()
		// 人物移动
		case MazeGame.BattleEventType_ROLE_MOVE:
			eventData = event.GetRoleMove()
		}
		if eventData == nil {
			logger.ErrorWF("OnMazeReportBattleEventRQ event type not supported", zap.Int32("eventType", int32(eventType)))
		} else {
			err = mazebarriereventredis.TriggerBarrierEvent(logger, userID, event.GetEventFrame(), event.GetEventTimeMs(), eventType, eventData)
			if err != nil {
				logger.ErrorWF("OnMazeReportBattleEventRQ TriggerBarrierEvent fail", zap.Error(err), zap.Int32("eventType", int32(eventType)))
			}
		}
	}
	return
}
