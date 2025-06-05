package game

import (
	"fmt"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/io/kafka/mazebarrieruserkafka"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/io/redis/mazebarriereventredis"
	"maze_game_server/io/redis/mazeuserbarrierredis"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazecommonvalue"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeGame"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
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
	if req.GetFoeExp() < 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("经验设置错误")
		return
	}

	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierDeathRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	if req.GetBarrierId() != userInfo.Barrier {
		logger.ErrorWF("OnMazeBarrierDeathRQ req barrier lt pass barrier", zap.Any("req", req), zap.Int32("save", userInfo.Barrier))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("请求的关卡id和存储的不一致")
		return
	}

	userBarrier, err := mazeuserbarrierredis.GetUserBarrierInfo(logger, userId, req.GetBarrierId())
	if err != nil {
		logger.ErrorWF("OnMazeBarrierDeathRQ GetUserBarrierInfo fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	//更新等级经验
	oldLevel := userInfo.Level
	oldExp := userInfo.TotalExp
	err = userInfo.AddExp(int64(req.GetFoeExp()))
	if err != nil {
		logger.ErrorWF("OnMazeBarrierDeathRQ addExp fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	newLevel := userInfo.Level
	err = mazeuserinfo.SetUserInfoV2(logger, userId, userInfo)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierDeathRQ SetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	mazecommonvalue.HandleUserLevelExpChg(logger, userId, userInfo.Level, userInfo.Exp, req.GetHeader().GetSession())

	defer func() {
		if oldLevel != newLevel {
			levelRecord := &mazeuserlevelkafka.MazeUserLevelRecord{
				UserId:      userId,
				OldLevel:    int32(oldLevel),
				OldTotalExp: oldExp,
				NewLevel:    int32(newLevel),
				NewTotalExp: int32(userInfo.TotalExp),
			}
			mazeuserlevelkafka.PushMazeLevelRecord(logger, levelRecord)
		}
	}()

	// 死亡之后是否需要清空复活次数
	userBarrier.BarrierId = proto.Int32(req.GetBarrierId())
	userBarrier.RebornCount = proto.Int32(0)
	userBarrier.BarrierStatus = proto.Int32(1)
	userBarrier.EndTime = proto.Int64(time.Now().Unix())
	err = mazeuserbarrierredis.SetUserBarrierInfo(logger, userId, 0, userBarrier)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierDeathRQ SetUserBarrierInfo fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	err = mazebarriereventredis.LeaveBarrier(logger, userId, req.GetBarrierId(), false)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierPassRQ LeaveBarrier fail", zap.Error(err))
	}

	if req.GetFoeExp() > 0 {
		res.BarrierAward = []*MazeCommon.MazeItem{&MazeCommon.MazeItem{ItemId: proto.Int32(constdef.MazeCommonItemExp), Count: proto.Int64(int64(req.GetFoeExp()))}}
	}

	passRecord := &mazebarrieruserkafka.MazeBarrierUserGameRecord{
		UserId:  userId,
		Barrier: req.GetBarrierId(),
		GameRet: mazebarrieruserkafka.GameRetDeath,
		Awards:  fmt.Sprintf("%d:%d", constdef.MazeCommonItemExp, req.GetFoeExp()),
	}

	mazebarrieruserkafka.PushMazeBarrierUserRecord(logger, passRecord)

	return nil
}

func GetDeathPunish(logger fklog.FKLogI, userId uint64) (subDeathPer int64, lostMin, lostMax int64, err error) {
	// forceVal, err := mazecalcattrredis.GetMazeForce(logger, userId)
	// if err != nil {
	// 	logger.ErrorWF("GetDeathPunish GetMazeForce fail", zap.Error(err))
	// 	return
	// }

	// allDeathCfg := GMazeKongfuDeathV8Cfg.GetAll()
	// if len(allDeathCfg) > 0 {
	// 	for _, cfg := range allDeathCfg {
	// 		if forceVal >= cfg.Kongfu_min && forceVal <= cfg.Kongfu_max {
	// 			return int64(cfg.Death_less_pro), cfg.Death_less_min, cfg.Death_less_max, nil
	// 		}
	// 	}
	// }

	// err = errors.New("not found death cfg")
	return
}
