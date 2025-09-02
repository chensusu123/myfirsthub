package game

import (
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeLevelV8Cfg"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/io/redis/mazecalcattrredis"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazecommonvalue"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/services/allianceservice"
	"maze_game_server/services/barriersavedataservice"
	"maze_game_server/services/moneyservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

func (g *Game) OnMazeLoginRQ_10451_10452(s *session.Session, req *MazeGame.MazeLoginRQ) (err error) {
	defer fkprometheus.InfoPMT("OnMazeLoginRQ")()

	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeGame.MazeLoginRS{}

	logger.CtxInfo(ctx, "OnMazeLoginRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnMazeLoginRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.MazeVersion = req.MazeVersion

	userId := uint64(s.UID())

	var level, exp, expMax, force, money, extra, extraExp, diamond int64
	var isInit bool

	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnMazeLoginRQ GetUserInfo fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	if userInfo.Level == 0 {
		isInit = true
		levelRecord := &mazeuserlevelkafka.MazeUserLevelRecord{
			UserId:   userId,
			OldLevel: 0,
			NewLevel: 1,
		}
		defer func() {
			if err == nil {
				mazeuserlevelkafka.PushMazeLevelRecord(ctx, levelRecord)
			}
		}()
		userInfo.SetLevel(1)
	}
	level = userInfo.Level
	exp = userInfo.Exp
	levelCfg := GMazeLevelV8Cfg.GetWithCtx(ctx, int32(level))
	if levelCfg == nil {
		res.ErrInfo = errors.CONFIG_NOT_FOUND.Wrap("等级表获取失败")
		return
	}
	expMax = levelCfg.Next_level_need_exp

	money, diamond, err = moneyservice.GlobalMoneyService.GetUserMoney(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnMazeLoginRQ GetUserMoney fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	if isInit || (userInfo.UserType != req.GetMazeVersion()) {
		if userInfo.UserType != req.GetMazeVersion() {
			userInfo.SetUserType(req.GetMazeVersion())
		}
		err = mazeuserinfo.SetUserInfoV2(ctx, userId, userInfo)
		if err != nil {
			logger.CtxError(ctx, "OnMazeLoginRQ SetUserInfo fail", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}
	}

	force, err = mazecalcattrredis.GetMazeForce(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnMazeLoginRQ GetMazeForce fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	extra, err2 := mazecommonvalue.MakeCommonValueExtra(ctx, userId, level, 0)
	if err2 != nil {
		logger.CtxError(ctx, "OnMazeLoginRQ MakeCommonValueExtra fail", zap.Error(err2))
		// res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		// return
	}

	// extraExp, err = MakeCommonValueExtraExp(logger, userId, level, force)
	// if err != nil {
	// 	logger.CtxError(ctx,"OnMazeLoginRQ MakeCommonValueExtraExp fail", zap.Error(err))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }

	passValue, err := barriersavedataservice.GlobalBarrierSaveDataService.GetPassValue(ctx, userId, userInfo.Barrier)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return nil
	}
	commonList := mazecommonvalue.MakeAllCommonValue(ctx, userId, level, exp, expMax, force, money, extra, extraExp, diamond, passValue, req.GetHeader().GetSession())

	mazecommonvalue.SendCommonValueIdPack(ctx, userId, commonList)

	allianceInfo, err := allianceservice.GlobalAllianceService.QueryAllianceInfo(ctx, 1)
	if err != nil {
		logger.CtxError(ctx, "OnMazeLoginRQ QueryAllianceInfo Fail",
			zap.Error(err))
	}
	res.AllianceInfo = allianceInfo.DataToAllianceInfoPb()
	return nil
}
