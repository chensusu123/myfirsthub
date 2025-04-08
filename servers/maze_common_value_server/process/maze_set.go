package process

import (
	"time"

	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkrpc"
	"gitlab.ifreetalk.com/plate/protodef/MazeCommonValueSvr"
	"gitlab.ifreetalk.com/servers/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/servers/maze_game_server/module/mazemoney"
	"go.uber.org/zap"
)

func OnMazeCommonValueSetRQ(ctx fkrpc.RPCContext, shardingID int64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	defer fkprometheus.DebugPMT("OnMazeCommonValueSetRQ")()

	req := rqMsg.(*MazeCommonValueSvr.MazeCommonValueSetRQ)
	res := rsMsg.(*MazeCommonValueSvr.MazeCommonValueSetRS)
	res.ErrInfo = errors.NO_ERROR
	logger := ctx.FKLogI
	logger.WarnWF("OnMazeCommonValueSetRQ with", zap.Any("rq", req))

	addStartTime := time.Now()
	defer func() {
		costTime := time.Since(addStartTime).Seconds()
		logger.WarnWF("OnMazeCommonValueSetRQ end ", zap.Any("req", req), zap.Any("res", res), zap.Float64("costTime", costTime))
		if costTime >= 0.5 {
			logger.ErrorWF("OnMazeCommonValueSetRQ timeout", zap.Any("req", req), zap.Any("res", res), zap.Float64("costTime", costTime))
		}
	}()

	query := req.GetSetItems()

	for _, item := range query {
		if item.GetItemId() == constdef.MazeCommonItemCoin {
			err = mazemoney.SetUserMoney(logger, req.GetUserId(), item.GetCount())
			if err != nil {
				logger.ErrorWF("OnMazeCommonValueSetRQ SetUserMoney fail", zap.Error(err))
				res.ErrInfo = errors.MODULE_ERROR.ToInfo()
				return
			}
		} else if item.GetItemId() == constdef.MazeCommonItemDiamond {
			err = mazemoney.SetUserDiamond(logger, req.GetUserId(), item.GetCount())
			if err != nil {
				logger.ErrorWF("OnMazeCommonValueSetRQ SetUserDiamond fail", zap.Error(err))
				res.ErrInfo = errors.MODULE_ERROR.ToInfo()
				return
			}
		}
	}

	return
}
