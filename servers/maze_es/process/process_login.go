package process

import (
	_ "gitlab.ifreetalk.com/plate/excel/auto/GMazeEquipPosLvV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/protodef/MazeGame"
	"gitlab.ifreetalk.com/plate/protodef/SysPackDef"

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
