package gentradeno

import (
	"maze_game_server/common/errors"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MessageType"
	"maze_game_server/pb/server/MazeItemSvr"
	"maze_game_server/servers/maze_main_server/process/item"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var (
// GoodOpDescMap = map[ItemSvr.ENUM_OP_TYPE]string{
// MazeItemSvr.ENUM_OP_TYPE(661): "人偶版装备洗炼",
// MazeItemSvr.ENUM_OP_TYPE(675): "法宝附魔消耗",
// }
)

func DeductItemsEx(logger fklog.FKLogI, userId uint64, opType int32, tradeNo uint64, items ...*MazeCommon.MazeItem) (errInfo *MessageType.ErrorInfo) {
	rq := &MazeItemSvr.ConsumeItemRQ{}
	rq.UserId = proto.Uint64(userId)
	rq.OpType = proto.Int32(int32(opType))
	rq.Items = append(rq.Items, items...)
	rq.TradeNumber = proto.Uint64(tradeNo)

	rs := &MazeItemSvr.ConsumeItemRS{ErrInfo: errors.NO_ERROR}

	// err := mazeitemrpc.DeductItemsRQ(logger, rq, rs)
	err := item.OnConsumeItemRQ(logger, rq, rs)
	if err != nil {
		logger.ErrorWF("DeductItemsEx1 net error", zap.Error(err), zap.Any("rq", rq), zap.Any("rs", rs))

		// netErrRecord := GenSubGoodFailRecord(ctx, tradeNo, int32(opType), 1, GoodOpDescMap[opType])
		// if strings.Contains(err.Error(), "i/o timeout") {
		// 	netErrRecord.TimeOutJson = ItemsToTimeoutRecord(rq.GetItems())
		// } else {
		// 	netErrRecord.FailJson = ItemsToFailRecord(rq.GetItems())
		// }
		// innerErr := error_record_kafka.PushGoodFailRecord(ctx, netErrRecord)
		// if innerErr == nil {
		// 	ctx.WarnWF("DeductItemsEx PushGoodFailRecord succ",
		// 		zap.Any("netErrRecord", netErrRecord))
		// } else {
		// 	ctx.ErrorWF("DeductItemsEx PushGoodFailRecord fail",
		// 		zap.Error(innerErr),
		// 		zap.Any("netErrRecord", netErrRecord))
		// }
		errInfo = errors.COMMON_ERROR_TIPS.Wrap("操作失败")
		return
	}
	if rs.ErrInfo == nil {
		rs.ErrInfo = errors.MODULE_ERROR.ToInfo()
	}
	if rs.GetErrInfo().GetErrCode() == errors.NO_ERROR_CODE {
		return
	}
	errInfo = rs.GetErrInfo()

	errCode := rs.GetErrInfo().GetErrCode()
	if errCode == errors.ITEM_CHECK_TIME_OUT.Code ||
		errCode == errors.ITEM_CHECK_ERROR.Code ||
		errCode == errors.ITEM_CHECK_ITEM_NOT_ENOUGH.Code {
		return
	}
	// errRecord := GenSubGoodFailRecord(ctx, tradeNo, int32(opType), 1, GoodOpDescMap[opType])

	// //if len(rs.FailedItems) > 0 { // 有扣除成功的物品需要处理  test
	// //	errRecord.FailJson = ItemsToFailRecord(rs.FailedItems)
	// //}

	// if len(rs.Items) > 0 { // 有扣除成功的物品需要处理
	// 	errRecord.FailJson = ItemsToFailRecord(rs.Items)
	// }
	// if len(rs.TimeoutItems) > 0 { // 有扣除成功的物品需要处理
	// 	errRecord.TimeOutJson = ItemsToTimeoutRecord(rs.TimeoutItems)
	// }
	// if len(errRecord.FailJson) > 0 || len(errRecord.TimeOutJson) > 0 {
	// 	innerErr := error_record_kafka.PushGoodFailRecord(ctx, errRecord)
	// 	if innerErr == nil {
	// 		ctx.WarnWF("DeductItemsEx PushGoodFailRecord succ",
	// 			zap.Any("errRecord", errRecord))
	// 	} else {
	// 		ctx.ErrorWF("DeductItemsEx PushGoodFailRecord fail",
	// 			zap.Error(innerErr),
	// 			zap.Any("errRecord", errRecord))
	// 	}
	// }
	return
}
