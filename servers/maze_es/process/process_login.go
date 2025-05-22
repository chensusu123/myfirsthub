package process

import (
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGame"
	"gitlab.ifreetalk.com/maze-plate/protodef/SysPackDef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/errors"
	"gitlab.ifreetalk.com/maze/maze_game_server/config/GMazeBarriesV8Cfg"
	_ "gitlab.ifreetalk.com/maze/maze_game_server/config/GMazeEquipPosLvV8Cfg"

	_ "gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeequipgetnumredis"

	"go.uber.org/zap"
)

func OnLoginRQ(ctx fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	fkprometheus.InfoPMT("OnLoginRQ")()

	req := rqMsg.(*SysPackDef.UserLoginRq)
	res := rsMsg.(*SysPackDef.UserLoginRs)

	logger := ctx

	defer func() {
		logger.InfoWF("OnLoginRQ end", zap.Any("req", req), zap.Any("res", res))
	}()
	// var index int32 = 3
	barrierCfg := GMazeBarriesV8Cfg.GetAll()
	// if barrierCfg == nil {
	// 	logger.ErrorWF("OnMazeBarrierListRQ get barrier cfg fail", zap.Any("barrier", index))
	// 	return
	// }
	if barrierCfg != nil {
		for _, v := range barrierCfg {
			logger.InfoWF("OnMazeBarrierListRQ get barrier cfg fail", zap.Any("barrierCfg", v.Order))
		}
	}
	// logger.InfoWF("OnMazeBarrierListRQ get barrier cfg fail", zap.Any("barrierCfg", barrierCfg))
	return
	// 认证
	// 失败直接返回

	res.Error = errors.NO_ERROR

	// 认证成功设置用户ID, 底层会处理
	userId := req.GetUserID()
	ctx.SetTag("userID", userId)
	return nil
}

func OnTestRQ(logger fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	fkprometheus.InfoPMT("OnTestRQ")()

	req := rqMsg.(*MazeGame.BarrierDeathRQ)
	res := rsMsg.(*MazeGame.BarrierDeathRS)

	defer func() {
		logger.InfoWF("OnTestRQ end", zap.Uint64("userID", shardingID),
			zap.Any("req", req.String()), zap.Any("res", res.String()))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	return nil
}
