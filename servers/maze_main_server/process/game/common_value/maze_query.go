package common_value

import (
	"time"

	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazemoney"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkrpc"
	"gitlab.ifreetalk.com/plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/plate/protodef/MazeCommonValueSvr"
	"go.uber.org/zap"
)

func OnMazeCommonValueQueryRQ(ctx fkrpc.RPCContext, shardingID int64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	defer fkprometheus.DebugPMT("OnMazeCommonValueQueryRQ")()

	req := rqMsg.(*MazeCommonValueSvr.MazeCommonValueQueryRQ)
	res := rsMsg.(*MazeCommonValueSvr.MazeCommonValueQueryRS)
	res.ErrInfo = errors.NO_ERROR
	res.UserId = req.UserId

	logger := ctx.FKLogI
	logger.WarnWF("OnMazeCommonValueQueryRQ with", zap.Any("rq", req))

	addStartTime := time.Now()
	defer func() {
		costTime := time.Since(addStartTime).Seconds()
		logger.WarnWF("OnMazeCommonValueQueryRQ end ", zap.Any("req", req), zap.Any("res", res), zap.Float64("costTime", costTime))
		if costTime >= 0.5 {
			logger.ErrorWF("OnMazeCommonValueQueryRQ timeout", zap.Any("req", req), zap.Any("res", res), zap.Float64("costTime", costTime))
		}
	}()

	query := req.GetQueryItems()

	coin, diamond, err := mazemoney.GetUserMoney(logger, req.GetUserId())
	if err != nil {
		logger.ErrorWF("OnMazeCommonValueQueryRQ GetUserMoney fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	for _, item := range query {
		if item.GetItemId() == constdef.MazeCommonItemCoin {
			res.Items = append(res.Items, &MazeCommon.MazeItem{ItemId: proto.Int32(constdef.MazeCommonItemCoin), Count: proto.Int64(coin)})
		} else if item.GetItemId() == constdef.MazeCommonItemDiamond {
			res.Items = append(res.Items, &MazeCommon.MazeItem{ItemId: proto.Int32(constdef.MazeCommonItemDiamond), Count: proto.Int64(diamond)})
		}
	}

	return
}
