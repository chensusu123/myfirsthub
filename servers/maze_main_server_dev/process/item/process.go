package item

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/common/additemdefine"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/iteminterface"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/itemutil"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/structdefine"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebagdb"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
	"gitlab.ifreetalk.com/plate/protodef/MazeBag"
	"gitlab.ifreetalk.com/plate/protodef/MazeItemSvr"
	"gitlab.ifreetalk.com/plate/protodef/MessageType"
	"go.uber.org/zap"
	"time"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"context"
)

func RegTcpHandler() {
	// 获取迷宫背包列表
	_ = websocket_service.RegProcSimple(16248, &MazeBag.MazeBagListRQ{},
		16249, &MazeBag.MazeBagListRS{}, OnMazeBagListRQ)

	// 重置迷宫背包列表
	_ = websocket_service.RegProcSimple(16273, &MazeBag.ResetMazeBagRQ{},
		16274, &MazeBag.ResetMazeBagRS{}, OnResetMazeBagRQ)
}

func OnMazeBagListRQ(ctx fknet.TCPContext, uid uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	req, ok := rqMsg.(*MazeBag.MazeBagListRQ)
	if !ok {
		ctx.ErrorWF("OnMazeBagListRQ pb is wrong", zap.Any("rqMsg", rqMsg))
		return
	}
	res, ok := rsMsg.(*MazeBag.MazeBagListRS)
	if !ok {
		ctx.ErrorWF("OnMazeBagListRQ pb is wrong", zap.Any("rsMsg", rsMsg))
		return
	}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	defer fkprometheus.DebugPMT("OnMazeBagListRQ")()
	defer ctx.InfoWF("OnMazeBagListRQ end", zap.Any("req", req), zap.Any("res", res))

	bagItemMap, err := mazebagdb.GetAllBagItem(ctx, uid, 300)
	if err != nil {
		ctx.ErrorWF("OnMazeBagListRQ GetAllBagItem err", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	res.Items = make([]*MazeBag.MazeBagItem, 0, len(bagItemMap))
	for id, count := range bagItemMap {
		if id <= 0 || count <= 0 {
			continue
		}
		res.Items = append(res.Items, itemutil.BuildMazeBagItem(ctx, id, count))
	}
	return
}

func OnResetMazeBagRQ(ctx fknet.TCPContext, uid uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	req, ok := rqMsg.(*MazeBag.ResetMazeBagRQ)
	if !ok {
		ctx.ErrorWF("OnResetMazeBagRQ pb is wrong", zap.Any("rqMsg", rqMsg))
		return
	}
	res, ok := rsMsg.(*MazeBag.ResetMazeBagRS)
	if !ok {
		ctx.ErrorWF("OnResetMazeBagRQ pb is wrong", zap.Any("rsMsg", rsMsg))
		return
	}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	defer fkprometheus.DebugPMT("OnResetMazeBagRQ")()
	defer ctx.InfoWF("OnResetMazeBagRQ end", zap.Any("req", req), zap.Any("res", res))

	err = mazebagdb.DelKey(ctx, uid)
	if err != nil {
		ctx.ErrorWF("OnResetMazeBagRQ DelKey err", zap.Error(err))
		res.ErrInfo = errors.DB_SAVE_ERROR.ToInfo()
		return
	}
	return
}

var GloRegIns = additemdefine.NewRegister()

// func RegisterThriftRPCPackage() {
//	thrift_service.RegisterTwowaySimple(131437, &MazeItemSvr.AddItemRQ{},
//		131438, &MazeItemSvr.AddItemRS{}, OnAddItemRQ)
//
//	thrift_service.RegisterTwowaySimple(131439, &MazeItemSvr.ConsumeItemRQ{},
//		131440, &MazeItemSvr.ConsumeItemRS{}, OnConsumeItemRQ)
//
//	thrift_service.RegisterTwowaySimple(131441, &MazeItemSvr.QueryItemRQ{},
//		131442, &MazeItemSvr.QueryItemRS{}, OnQueryItemRQ)
//
//	thrift_service.RegisterTwowaySimple(131443, &MazeItemSvr.CheckAddItemRQ{},
//		131444, &MazeItemSvr.CheckAddItemRS{}, OnCheckAddItemRQ)
// }

func OnAddItemRQ(ctx fklog.FKLogI, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	req, ok := rqMsg.(*MazeItemSvr.AddItemRQ)
	if !ok {
		ctx.ErrorWF("OnAddItemRQ pb is wrong", zap.Any("rqMsg", rqMsg))
		return
	}
	res, ok := rsMsg.(*MazeItemSvr.AddItemRS)
	if !ok {
		ctx.ErrorWF("OnAddItemRQ pb is wrong", zap.Any("rsMsg", rsMsg))
		return
	}
	res.ErrInfo = errors.NO_ERROR
	res.UserId = req.UserId

	startTime := time.Now()
	defer ctx.InfoWF("OnAddItemRQ end", zap.Any("req", req), zap.Any("res", res), zap.Duration("costTime", time.Since(startTime)))

	defer fkprometheus.DebugPMT("OnAddItemRQ")()

	uid, opType, tradeNo := req.GetUserId(), req.GetOpType(), req.GetTradeNumber()
	if uid <= 0 || opType <= 0 || tradeNo <= 0 {
		ctx.WarnWF("OnAddItemRQ invalid args", zap.Any("req", req))
		res.ErrInfo = errors.ARGS_NOT_MATCH.ToInfo()
		return
	}

	userCtx := itemutil.WrapUserContext(context.TODO(), uid, ctx, req.GetHeader())

	// 检查并合并物品
	realAddItems := itemutil.CheckAndMergeItem(req.GetItems())
	if len(realAddItems) <= 0 {
		userCtx.WarnWF("OnAddItemRQ add items nil", zap.Any("addItems", req.GetItems()))
		res.ErrInfo = errors.ITEM_CHECK_ERROR.Wrap("添加道具为空")
		return
	}

	option := &additemdefine.AddItemOption{
		OpType:      opType,
		TradeNumber: tradeNo,
		RegIns:      GloRegIns,
	}

	var addRes *structdefine.AddItemRes
	addRes, err = iteminterface.GatherItems(userCtx, option, realAddItems...)
	if err != nil {
		userCtx.WarnWF("OnAddItemRQ GatherItems err", zap.Any("gatherItem", realAddItems),
			zap.Any("addRes", addRes), zap.Error(err),
		)

		var rer *errors.CodeError
		rer, ok = err.(*errors.CodeError)
		if ok {
			res.ErrInfo = rer.ToInfo()
		} else {
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		}
	}

	res.FailItems = addRes.FailItem
	res.TimeoutItems = addRes.TimeoutItem
	res.SucItems = addRes.SucItem
	return
}

func OnConsumeItemRQ(ctx fklog.FKLogI, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	req, ok := rqMsg.(*MazeItemSvr.ConsumeItemRQ)
	if !ok {
		ctx.ErrorWF("OnConsumeItemRQ pb is wrong", zap.Any("rqMsg", rqMsg))
		return
	}
	res, ok := rsMsg.(*MazeItemSvr.ConsumeItemRS)
	if !ok {
		ctx.ErrorWF("OnConsumeItemRQ pb is wrong", zap.Any("rsMsg", rsMsg))
		return
	}
	res.ErrInfo = errors.NO_ERROR
	res.UserId = req.UserId

	startTime := time.Now()
	defer ctx.InfoWF("OnConsumeItemRQ end", zap.Any("req", req), zap.Any("res", res), zap.Duration("costTime", time.Since(startTime)))
	defer fkprometheus.DebugPMT("OnConsumeItemRQ")()

	uid, opType, tradeNo := req.GetUserId(), req.GetOpType(), req.GetTradeNumber()
	if uid <= 0 || opType <= 0 || tradeNo <= 0 {
		ctx.WarnWF("OnConsumeItemRQ invalid args", zap.Any("req", req))
		res.ErrInfo = errors.ARGS_NOT_MATCH.ToInfo()
		return
	}

	userCtx := itemutil.WrapUserContext(context.TODO(), uid, ctx, nil)

	// 检查并合并物品
	realSubItems := itemutil.CheckAndMergeItem(req.GetItems())
	if len(realSubItems) == 0 {
		userCtx.WarnWF("OnConsumeItemRQ sub items nil", zap.Any("req", req))
		res.ErrInfo = errors.ITEM_CHECK_ERROR.Wrap("扣道具数量为空")
		return
	}

	option := &additemdefine.AddItemOption{
		OpType:      opType,
		TradeNumber: req.GetTradeNumber(),
		RegIns:      GloRegIns,
	}

	subRes, errInfo := iteminterface.DeductItems(userCtx, option, realSubItems...)
	if errInfo != nil && errInfo.GetErrCode() != errors.NO_ERROR_CODE {
		userCtx.WarnWF("OnConsumeItemRQ DeductItems fail", zap.Any("errInfo", errInfo), zap.Any("realSubItems", realSubItems), zap.Any("subRes", subRes))
		res.ErrInfo = errInfo
	}

	res.SucItems = subRes.SucItem
	res.TimeoutItems = subRes.TimeoutItem
	res.LessItems = subRes.LessItem
	return
}

func OnQueryItemRQ(ctx fklog.FKLogI, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	req, ok := rqMsg.(*MazeItemSvr.QueryItemRQ)
	if !ok {
		ctx.ErrorWF("OnQueryItemRQ pb is wrong", zap.Any("rqMsg", rqMsg))
		return
	}
	res, ok := rsMsg.(*MazeItemSvr.QueryItemRS)
	if !ok {
		ctx.ErrorWF("OnQueryItemRQ pb is wrong", zap.Any("rsMsg", rsMsg))
		return
	}
	res.ErrInfo = errors.NO_ERROR
	res.UserId = req.UserId

	startTime := time.Now()
	defer ctx.InfoWF("OnQueryItemRQ end", zap.Any("req", req), zap.Any("res", res), zap.Duration("costTime", time.Since(startTime)))

	defer fkprometheus.DebugPMT("OnQueryItemRQ")()

	uid := req.GetUserId()
	if uid <= 0 {
		ctx.WarnWF("OnQueryItemRQ invalid args", zap.Any("req", req))
		res.ErrInfo = errors.ARGS_NOT_MATCH.ToInfo()
		return
	}

	userCtx := itemutil.WrapUserContext(context.TODO(), uid, ctx, nil)

	queryItems := itemutil.MergeQueryItems(req.GetItems())
	if len(queryItems) == 0 {
		ctx.WarnWF("OnQueryItemSvrRQ MergeQueryItems items nil", zap.Any("items", req.GetItems()))
		return
	}

	option := &additemdefine.AddItemOption{
		RegIns: GloRegIns,
	}

	var errInfo *MessageType.ErrorInfo
	res.Items, errInfo = iteminterface.QueryItems(userCtx, option, queryItems)
	if errInfo != nil {
		userCtx.WarnWF("OnQueryItemSvrRQ QueryItems err", zap.Error(err))
		res.ErrInfo = errInfo
		return
	}

	return
}

func OnCheckAddItemRQ(ctx fklog.FKLogI, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	req, ok := rqMsg.(*MazeItemSvr.CheckAddItemRQ)
	if !ok {
		ctx.ErrorWF("OnCheckAddItemRQ pb is wrong", zap.Any("rqMsg", rqMsg))
		return
	}
	res, ok := rsMsg.(*MazeItemSvr.CheckAddItemRS)
	if !ok {
		ctx.ErrorWF("OnCheckAddItemRQ pb is wrong", zap.Any("rsMsg", rsMsg))
		return
	}
	res.ErrInfo = errors.NO_ERROR
	res.UserId = req.UserId

	startTime := time.Now()
	defer ctx.InfoWF("OnCheckAddItemRQ end", zap.Any("req", req), zap.Any("res", res), zap.Duration("costTime", time.Since(startTime)))
	defer fkprometheus.DebugPMT("OnCheckAddItemRQ")()

	uid := req.GetUserId()
	if uid <= 0 {
		ctx.WarnWF("OnCheckAddItemRQ invalid args", zap.Any("req", req))
		res.ErrInfo = errors.ARGS_NOT_MATCH.ToInfo()
		return
	}

	userCtx := itemutil.WrapUserContext(context.TODO(), uid, ctx, nil)

	mergeItems := itemutil.CheckAndMergeItem(req.GetItems())
	if len(mergeItems) == 0 {
		userCtx.WarnWF("OnCheckAddItemRQ check add item nil", zap.Any("req", req))
		res.ErrInfo = errors.NewCommonCodeError("检查添加道具数量为空")
		return nil
	}

	option := &additemdefine.AddItemOption{
		RegIns: GloRegIns,
	}
	checkRes, err := iteminterface.CheckAddItems(userCtx, option, mergeItems)
	if err != nil {
		userCtx.WarnWF("OnCheckAddItemRQ error", zap.Error(err))

		var errInfo *errors.CodeError
		errInfo, ok = err.(*errors.CodeError)
		if ok {
			res.ErrInfo = errInfo.ToInfo()
		} else {
			res.ErrInfo = errors.ITEM_CHECK_ERROR.ToInfo()
		}
	}

	res.TimeoutItems = checkRes.TimeoutItem
	res.FailItems = checkRes.FailItem
	res.LimitItems = checkRes.LimitItem
	return
}
