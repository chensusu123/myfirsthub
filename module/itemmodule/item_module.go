package itemmodule

import (
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeItemSvr"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/tradeno"
	itemProcess "gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/item"
	"go.uber.org/zap"
)

/**
* @Description:
* @Author: wangyongliang
* @Date: 2025/3/24 17:03
**/

const CostRefreshType = 694 // 刷新buff

// DeductItems 扣除物品
func DeductItems(logger fklog.FKLogI, userId uint64, opType int32, items []*MazeCommon.MazeItem) error {
	if len(items) == 0 {
		return nil
	}

	req := &MazeItemSvr.ConsumeItemRQ{
		UserId:      proto.Uint64(userId),
		Items:       items,
		OpType:      proto.Int32(opType),
		TradeNumber: proto.Uint64(tradeno.GetTradeNum()),
	}
	res := &MazeItemSvr.ConsumeItemRS{}
	// err := mazeitemrpc.DeductItemsRQ(logger, req, res)
	err := itemProcess.OnConsumeItemRQ(logger, req, res)
	if err != nil {
		logger.ErrorWF("DeductItems DeductItemsRQ", zap.Any("req", req),
			zap.Any("res", res), zap.Error(err))
		return errors.New("扣除物品失败")
	}

	if res.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
		err = errors.New(string(res.GetErrInfo().ErrMsg))
		logger.ErrorWF("DeductItems ErrorInfo", zap.Any("req", req), zap.Any("res", res))
		return err
	}

	logger.InfoWF("DeductItems end", zap.Any("req", req), zap.Any("res", res))
	return nil
}

func AddItems(logger fklog.FKLogI, userId uint64, opType int32, items []*MazeCommon.MazeItem) error {
	req := &MazeItemSvr.AddItemRQ{
		UserId:      proto.Uint64(userId),
		Items:       items,
		OpType:      proto.Int32(opType),
		TradeNumber: proto.Uint64(tradeno.GetTradeNum()),
	}
	res := &MazeItemSvr.AddItemRS{}
	// err := mazeitemrpc.AddItemsRQ(logger, req, res)
	err := itemProcess.OnAddItemRQ(logger, req, res)
	if err != nil {
		logger.ErrorWF("AddItems AddItemsRQ", zap.Any("req", req),
			zap.Any("res", res), zap.Error(err))
		return errors.New("加物品失败")
	}

	if res.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
		err = errors.New(string(res.GetErrInfo().ErrMsg))
		logger.ErrorWF("AddItems ErrorInfo", zap.Any("req", req), zap.Any("res", res))
		return err
	}

	logger.InfoWF("AddItems end", zap.Any("req", req), zap.Any("res", res))
	return nil
}

func QueryItem(logger fklog.FKLogI, userId uint64, item []*MazeCommon.MazeItem) error {
	req := &MazeItemSvr.QueryItemRQ{
		UserId: proto.Uint64(userId),
		Items:  item,
	}
	res := &MazeItemSvr.QueryItemRS{}
	// err := mazeitemrpc.QueryItemSvrRQ(logger, req, res)
	err := itemProcess.OnQueryItemRQ(logger, req, res)
	if err != nil {
		logger.ErrorWF("QueryItem QueryItemSvrRQ", zap.Any("req", req),
			zap.Any("res", res), zap.Error(err))
		return errors.New("查询物品失败")
	}

	if res.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
		err = errors.New(string(res.GetErrInfo().ErrMsg))
		logger.ErrorWF("QueryItem ErrorInfo", zap.Any("req", req), zap.Any("res", res))
		return err
	}

	logger.InfoWF("QueryItem end", zap.Any("req", req), zap.Any("res", res))
	return nil
}

// func CheckAddItem(logger fklog.FKLogI, userId uint64, opType int32, item []*MazeCommon.MazeItem) error {
// 	req := &MazeItemSvr.CheckAddItemRQ{
// 		UserId: proto.Uint64(userId),
// 		Items:  item,
// 		OpType: proto.Int32(opType),
// 	}
//
// 	res := &MazeItemSvr.CheckAddItemRS{}
// 	err := mazeitemrpc.CheckAddItemRQ(logger, req, res)
// 	if err != nil {
// 		logger.ErrorWF("CheckAddItem QueryItemSvrRQ", zap.Any("req", req),
// 			zap.Any("res", res), zap.Error(err))
// 		return errors.New("检查物品失败")
// 	}
//
// 	if res.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
// 		err = errors.New(string(res.GetErrInfo().ErrMsg))
// 		logger.ErrorWF("CheckAddItem ErrorInfo", zap.Any("req", req), zap.Any("res", res))
// 		return err
// 	}
//
// 	logger.InfoWF("CheckAddItem end", zap.Any("req", req), zap.Any("res", res))
// 	return nil
// }
