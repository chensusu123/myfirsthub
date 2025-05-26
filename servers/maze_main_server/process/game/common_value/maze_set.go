package common_value

import (
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/module/mazemoney"
	"maze_game_server/pb/server/MazeCommonValueSvr"
)

func MazeCommonValueSetRQ(logger fklog.FKLogI, userID int64, req *MazeCommonValueSvr.MazeCommonValueSetRQ, res *MazeCommonValueSvr.MazeCommonValueSetRS) (err error) {
	res.ErrInfo = errors.NO_ERROR

	logger.WarnWF("MazeCommonValueSetRQ with", zap.Any("rq", req))

	addStartTime := time.Now()
	defer func() {
		costTime := time.Since(addStartTime).Seconds()
		logger.WarnWF("MazeCommonValueSetRQ end ", zap.Any("req", req), zap.Any("res", res), zap.Float64("costTime", costTime))
		if costTime >= 0.5 {
			logger.ErrorWF("MazeCommonValueSetRQ timeout", zap.Any("req", req), zap.Any("res", res), zap.Float64("costTime", costTime))
		}
	}()

	query := req.GetSetItems()

	for _, item := range query {
		if item.GetItemId() == constdef.MazeCommonItemCoin {
			err = mazemoney.SetUserMoney(logger, req.GetUserId(), item.GetCount())
			if err != nil {
				logger.ErrorWF("MazeCommonValueSetRQ SetUserMoney fail", zap.Error(err))
				res.ErrInfo = errors.MODULE_ERROR.ToInfo()
				return
			}
		} else if item.GetItemId() == constdef.MazeCommonItemDiamond {
			err = mazemoney.SetUserDiamond(logger, req.GetUserId(), item.GetCount())
			if err != nil {
				logger.ErrorWF("MazeCommonValueSetRQ SetUserDiamond fail", zap.Error(err))
				res.ErrInfo = errors.MODULE_ERROR.ToInfo()
				return
			}
		}
	}

	return
}
