package gentradeno

import (
	"strings"

	"maze_game_server/common/errors"
	"maze_game_server/pb/common/Common"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MessageType"
	"maze_game_server/pb/server/MazeItemSvr"
	"maze_game_server/servers/maze_main_server/process/item"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var AddGoodOpDescMap = map[int32]string{
	631: "人偶版装备出售加黄金",
	// 632: "人偶版装备铸造消耗",
}

func AddItemEx(logger fklog.FKLogI, userId uint64, opType int32, tradeNo uint64, header *Common.PacketHeader, items ...*MazeCommon.MazeItem) (errInfo *MessageType.ErrorInfo) {
	rq := &MazeItemSvr.AddItemRQ{
		UserId: proto.Uint64(userId),
		OpType: proto.Int32(int32(opType)),
		Header: header,
	}
	rq.Items = append(rq.Items, items...)
	rq.TradeNumber = proto.Uint64(tradeNo)

	rs := &MazeItemSvr.AddItemRS{ErrInfo: errors.NO_ERROR}
	// monitorName := fmt.Sprintf("ItemsNewRpc.AddItemsRQ_%d", fkconfig.GetServerConfig().ServerTypeID)
	// m1 := monitor.GetMonitor(monitorName)
	// m1.StartV1()
	// defer m1.EndV1()

	// err := mazeitemrpc.AddItemsRQ(logger, rq, rs)
	err := item.OnAddItemRQ(logger, rq, rs)
	if err != nil {
		logger.ErrorWF("AddItemEx AddItemsRQ net error", zap.Error(err), zap.Any("rq", rq), zap.Any("rs", rs))
		// 记录异常
		// if IsTimeOut(err.Error()) {
		// 	RecordErrorItem(logger, userId, int32(opType), addType, tradeNo, GoodOpDescMap[opType], items, nil)
		// } else {
		// 	RecordErrorItem(logger, userId, int32(opType), addType, tradeNo, GoodOpDescMap[opType], nil, items)
		// }
		errInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return
	}
	if rs.ErrInfo == nil {
		rs.ErrInfo = errors.MODULE_ERROR.ToInfo()
	}
	if rs.GetErrInfo().GetErrCode() == errors.NO_ERROR_CODE {
		return
	}
	errInfo = rs.GetErrInfo()
	logger.ErrorWF("AddItemEx deduct items with err", zap.Any("errInfo", errInfo), zap.Any("rq", rq), zap.Any("rs", rs))
	// RecordErrorItem(logger, userId, int32(opType), addType, tradeNo, GoodOpDescMap[opType], rs.GetTimeoutItems(), rs.GetLeftItem())
	return
}

func IsTimeOut(errMsg string) bool {
	return strings.Contains(errMsg, "i/o timeout")
}

func RecordErrorItem(logger fklog.FKLogI, uid uint64, opType int32, addType int32, tradeNo uint64, tradeDesc string, timeoutItems, leftItems []*MazeCommon.MazeItem) {
	// logMsg := error_record_kafka.NewGoodFailRecord(uid, tradeNo, tradeDesc, fkconfig.EnvVal.ServerID, opType, addType, error_record_kafka.TradeTypeForAdd)
	// isChange := false
	// if len(timeoutItems) > 0 {
	// 	for _, item := range timeoutItems {
	// 		timeoutItem := &error_record_kafka.TimeOutInfo{
	// 			ItemId: item.GetItemId(),
	// 			Count:  item.GetCount(),
	// 		}
	// 		logMsg.TimeOutJson = append(logMsg.TimeOutJson, timeoutItem)
	// 	}
	// 	isChange = true
	// }
	// if len(leftItems) > 0 {
	// 	for _, item := range leftItems {
	// 		timeoutItem := &error_record_kafka.FailInfo{
	// 			ItemId: item.GetItemId(),
	// 			Count:  item.GetCount(),
	// 		}
	// 		logMsg.FailJson = append(logMsg.FailJson, timeoutItem)
	// 	}
	// 	isChange = true
	// }
	// if isChange {
	// 	innerErr := error_record_kafka.PushGoodFailRecord(logger, logMsg)
	// 	if innerErr == nil {
	// 		logger.WarnWF("RecordErrorItem PushGoodFailRecord succ",
	// 			zap.Any("logMsg", logMsg))
	// 	} else {
	// 		logger.ErrorWF("RecordErrorItem PushGoodFailRecord fail",
	// 			zap.Error(innerErr),
	// 			zap.Any("logMsg", logMsg))
	// 	}
	// }
}
