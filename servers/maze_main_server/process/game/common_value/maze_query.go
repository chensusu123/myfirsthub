package common_value

import (
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/module/mazemoney"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/server/MazeCommonValueSvr"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func MazeCommonValueQueryRQ(logger fklog.FKLogI, userID int64, req *MazeCommonValueSvr.MazeCommonValueQueryRQ, res *MazeCommonValueSvr.MazeCommonValueQueryRS) (err error) {
	res.ErrInfo = errors.NO_ERROR
	res.UserId = req.UserId

	logger.WarnWF("MazeCommonValueQueryRQ with", zap.Any("rq", req))

	addStartTime := time.Now()
	defer func() {
		costTime := time.Since(addStartTime).Seconds()
		logger.WarnWF("MazeCommonValueQueryRQ end ", zap.Any("req", req), zap.Any("res", res), zap.Float64("costTime", costTime))
		if costTime >= 0.5 {
			logger.ErrorWF("MazeCommonValueQueryRQ timeout", zap.Any("req", req), zap.Any("res", res), zap.Float64("costTime", costTime))
		}
	}()

	query := req.GetQueryItems()

	coin, diamond, err := mazemoney.GetUserMoney(logger, req.GetUserId())
	if err != nil {
		logger.ErrorWF("MazeCommonValueQueryRQ GetUserMoney fail", zap.Error(err))
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
