package game

import (
	"maze_game_server/common/errors"
	"maze_game_server/io/kafka/mazebarrieruserkafka"
	"maze_game_server/io/redis/mazebarriereventredis"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/servers/maze_main_server/process/game/events"
	"maze_game_server/services/barrierservice"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (g *Game) OnMazeBarrierDeathRQ_10449_10450(s *session.Session, req *MazeGame.BarrierDeathRQ) (err error) {
	defer fkprometheus.InfoPMT("OnMazeBarrierDeathRQ")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeGame.BarrierDeathRS{}

	logger.InfoWF("OnMazeBarrierDeathRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnMazeBarrierDeathRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.BarrierId = req.BarrierId

	userId := uint64(s.UID())

	if req.GetBarrierId() <= 0 {
		logger.ErrorWF("OnMazeBarrierDeathRQ req barrier invalid", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡id未设置")
		return
	}
	//if req.GetFoeExp() < 0 {
	//	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("经验设置错误")
	//	return
	//}

	killMonsterNum, totalDamage, awards, errinfo := barrierservice.Global.BarrierDeath(logger, req.GetHeader(), userId, req.GetBarrierId(), req.GetFoeExp())
	if errinfo != nil {
		res.ErrInfo = errinfo
		return
	} else {
		res.KillMonsterNum = proto.Int32(killMonsterNum)
		res.TotalDamage = proto.Int64(totalDamage)
		res.BarrierAward = awards
	}

	err = mazebarriereventredis.LeaveBarrier(logger, userId, req.GetBarrierId(), false)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierPassRQ LeaveBarrier fail", zap.Error(err), zap.Any("barrier", req.GetBarrierId()))
	}

	// 触发离开关卡事件
	events.OnLeaveBarrier(logger, userId, 0, time.Now().UnixMilli(), &MazeGame.BattleEventLeaveBarrier{BarrierId: proto.Int32(req.GetBarrierId()), Result: MazeGame.BarrierResult_DEATH.Enum()})

	passRecord := &mazebarrieruserkafka.MazeBarrierUserGameRecord{
		UserId:  userId,
		Barrier: req.GetBarrierId(),
		GameRet: mazebarrieruserkafka.GameRetDeath,
		Awards:  getAwards(awards),
	}

	mazebarrieruserkafka.PushMazeBarrierUserRecord(logger, passRecord)

	return nil
}

// func GetDeathPunish(logger fklog.FKLogI, userId uint64) (subDeathPer int64, lostMin, lostMax int64, err error) {
// 	// forceVal, err := mazecalcattrredis.GetMazeForce(logger, userId)
// 	// if err != nil {
// 	// 	logger.ErrorWF("GetDeathPunish GetMazeForce fail", zap.Error(err))
// 	// 	return
// 	// }

// 	// allDeathCfg := GMazeKongfuDeathV8Cfg.GetAll()
// 	// if len(allDeathCfg) > 0 {
// 	// 	for _, cfg := range allDeathCfg {
// 	// 		if forceVal >= cfg.Kongfu_min && forceVal <= cfg.Kongfu_max {
// 	// 			return int64(cfg.Death_less_pro), cfg.Death_less_min, cfg.Death_less_max, nil
// 	// 		}
// 	// 	}
// 	// }

// 	// err = errors.New("not found death cfg")
// 	return
// }

// func getmapAwards(awardMap map[int32]int64, awardEquip map[int32]int32) string {
// 	awardStr := make([]string, 0)

// 	for k, v := range awardMap {
// 		awardStr = append(awardStr, fmt.Sprintf("%d:%d", k, v))
// 	}

// 	for k, v := range awardEquip {
// 		awardStr = append(awardStr, fmt.Sprintf("%d:%d", k, v))
// 	}
// 	return strings.Join(awardStr, "_")
// }
