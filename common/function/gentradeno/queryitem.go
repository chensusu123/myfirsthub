package gentradeno

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/io/rpc/mazeitemrpc"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/plate/protodef/MazeItemSvr"
	"gitlab.ifreetalk.com/plate/protodef/MessageType"
	"go.uber.org/zap"
)

func QueryItems(logger fklog.FKLogI, userId uint64, items ...*MazeCommon.MazeItem) (queryItems []*MazeCommon.MazeItem, errInfo *MessageType.ErrorInfo) {
	rq := &MazeItemSvr.QueryItemRQ{
		UserId: proto.Uint64(userId),
		Items:  items,
	}
	rq.Items = append(rq.Items, items...)

	rs := &MazeItemSvr.QueryItemRS{ErrInfo: errors.NO_ERROR}
	err := mazeitemrpc.QueryItemsRQ(logger, rq, rs)
	if err != nil {
		logger.ErrorWF("QueryItems QueryItemsRQ net error", zap.Error(err), zap.Any("rq", rq), zap.Any("rs", rs))
		errInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return
	}
	if rs.ErrInfo == nil {
		rs.ErrInfo = errors.MODULE_ERROR.ToInfo()
	}
	if rs.GetErrInfo().GetErrCode() == errors.NO_ERROR_CODE {
		queryItems = rs.GetItems()
		return
	}
	errInfo = rs.GetErrInfo()
	logger.ErrorWF("QueryItems query items with err", zap.Any("errInfo", errInfo), zap.Any("rq", rq), zap.Any("rs", rs))
	return
}
