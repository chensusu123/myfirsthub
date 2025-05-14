package game

import (
	"github.com/lonng/nano/session"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazeuserlevelkafka"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazecalcattrredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazecommonvalue"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazemoney"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserinfo"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeLevelV8Cfg"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/protodef/MazeGame"
	"go.uber.org/zap"
)

func (g *Game) OnMazeLoginRQ(s *session.Session, req *MazeGame.MazeLoginRQ) (err error) {
	fkprometheus.InfoPMT("OnMazeLoginRQ")()
	logger := fklog.AppLogger().Clone("game")

	res := &MazeGame.MazeLoginRS{}

	logger.InfoWF("OnMazeLoginRQ start", zap.Any("req", req))
	defer func() {
		logger.InfoWF("OnMazeLoginRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.MazeVersion = req.MazeVersion

	userId := req.GetHeader().GetSharding()

	var level, exp, expMax, force, money, extra, extraExp, diamond int64
	var isInit bool

	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, uint64(userId))
	if err != nil {
		logger.ErrorWF("OnMazeLoginRQ GetUserInfo fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	if userInfo.Level == 0 {
		isInit = true
		levelRecord := &mazeuserlevelkafka.MazeUserLevelRecord{
			UserId:   uint64(userId),
			OldLevel: 0,
			NewLevel: 1,
		}
		defer func() {
			if err == nil {
				mazeuserlevelkafka.PushMazeLevelRecord(logger, levelRecord)
			}
		}()
		userInfo.SetLevel(1)
	}
	level = userInfo.Level
	exp = userInfo.Exp
	levelCfg := GMazeLevelV8Cfg.Get(int32(level))
	if levelCfg == nil {
		res.ErrInfo = errors.CONFIG_NOT_FOUND.Wrap("等级表获取失败")
		return s.Response(res)
	}
	expMax = levelCfg.Next_level_need_exp

	money, diamond, err = mazemoney.GetUserMoney(logger, uint64(userId))
	if err != nil {
		logger.ErrorWF("OnMazeLoginRQ GetUserMoney fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return s.Response(res)
	}

	if isInit || (userInfo.UserType != req.GetMazeVersion()) {
		if userInfo.UserType != req.GetMazeVersion() {
			userInfo.SetUserType(req.GetMazeVersion())
		}
		err = mazeuserinfo.SetUserInfoV2(logger, uint64(userId), userInfo)
		if err != nil {
			logger.ErrorWF("OnMazeLoginRQ SetUserInfo fail", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return s.Response(res)
		}
	}

	force, err = mazecalcattrredis.GetMazeForce(logger, uint64(userId))
	if err != nil {
		logger.ErrorWF("OnMazeLoginRQ GetMazeForce fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	extra, err2 := mazecommonvalue.MakeCommonValueExtra(logger, uint64(userId), level, 0)
	if err2 != nil {
		logger.ErrorWF("OnMazeLoginRQ MakeCommonValueExtra fail", zap.Error(err2))
		// res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		// return
	}

	// extraExp, err = MakeCommonValueExtraExp(logger, userId, level, force)
	// if err != nil {
	// 	logger.ErrorWF("OnMazeLoginRQ MakeCommonValueExtraExp fail", zap.Error(err))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }

	commonList := mazecommonvalue.MakeAllCommonValue(logger, uint64(userId), level, exp, expMax, force, money, extra, extraExp, diamond, req.GetHeader().GetSession())

	mazecommonvalue.SendCommonValueIdPack(logger, uint64(userId), commonList)

	// 用户登录成功后，将UserID绑定到会话上，方便后续操作
	// TODO 这个操作应该在验证成功后执行，待调整
	s.Bind(userId)

	return s.Response(res)
}
