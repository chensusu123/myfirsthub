package game

import (
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeLevelV8Cfg"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/io/redis/mazecalcattrredis"
	"maze_game_server/module/mazecommonvalue"
	"maze_game_server/module/mazemoney"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeGame"
)

func OnMazeLoginRQ(logger fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	fkprometheus.InfoPMT("OnMazeLoginRQ")()

	req := rqMsg.(*MazeGame.MazeLoginRQ)
	res := rsMsg.(*MazeGame.MazeLoginRS)

	logger.InfoWF("OnMazeLoginRQ start", zap.Any("req", req))
	defer func() {
		logger.InfoWF("OnMazeLoginRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.MazeVersion = req.MazeVersion

	userId := shardingID

	var level, exp, expMax, force, money, extra, extraExp, diamond int64
	var isInit bool

	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("OnMazeLoginRQ GetUserInfo fail", zap.Error(err))
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
		return
	}
	expMax = levelCfg.Next_level_need_exp

	money, diamond, err = mazemoney.GetUserMoney(logger, userId)
	if err != nil {
		logger.ErrorWF("OnMazeLoginRQ GetUserMoney fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	if isInit || (userInfo.UserType != req.GetMazeVersion()) {
		if userInfo.UserType != req.GetMazeVersion() {
			userInfo.SetUserType(req.GetMazeVersion())
		}
		err = mazeuserinfo.SetUserInfoV2(logger, userId, userInfo)
		if err != nil {
			logger.ErrorWF("OnMazeLoginRQ SetUserInfo fail", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}
	}

	force, err = mazecalcattrredis.GetMazeForce(logger, userId)
	if err != nil {
		logger.ErrorWF("OnMazeLoginRQ GetMazeForce fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	extra, err2 := mazecommonvalue.MakeCommonValueExtra(logger, userId, level, 0)
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

	commonList := mazecommonvalue.MakeAllCommonValue(logger, userId, level, exp, expMax, force, money, extra, extraExp, diamond, req.GetHeader().GetSession())

	mazecommonvalue.SendCommonValueIdPack(logger, userId, commonList)

	return nil
}
