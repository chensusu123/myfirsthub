package collect

import (
	"fmt"
	"time"

	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazecollectrecord"

	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/gentradeno"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazecollectredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserinfo"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeBarriesOnHookV8Cfg"
	"gitlab.ifreetalk.com/plate/protodef/MazeCollect"
	"gitlab.ifreetalk.com/plate/protodef/MazeCollectCache"

	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"go.uber.org/zap"
)

// 道具领取
func OnMazeCollectItemReceiveRQ(ctx fknet.TCPContext, userId uint64, rq proto.Message, rs proto.Message) (err error) {
	req := rq.(*MazeCollect.MazeCollectItemReceiveRQ)
	res := rs.(*MazeCollect.MazeCollectItemReceiveRS)
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	userCtx := fkserver.UserContext{UserID: userId, FKLogI: ctx.FKLogI, Header: req.Header, Context: ctx.Context}

	ctx.InfoWF("OnMazeCollectItemReceiveRQ start", zap.Any("req", req))
	defer func() {
		ctx.InfoWF("OnMazeCollectItemReceiveRQ end", zap.Any("res", res))
	}()

	// 是否已经初始化
	collectInfo, err := mazecollectredis.GetCollectInfo(ctx, userId)
	if err != nil {
		ctx.ErrorWF("OnMazeCollectItemReceiveRQ GetCollectInfo error",
			zap.Any("userId", userId),
			zap.Error(err))
		return
	}

	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		ctx.ErrorWF("OnMazeCollectItemReceiveRQ GetUserInfoV2", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return err
	}

	// 确保领取时道具产出是正确的
	if !CheckReceiveItems(collectInfo) {
		ctx.ErrorWF("OnMazeCollectItemReceiveRQ CheckReceiveItems",
			zap.Int64("nowSecond", time.Now().Unix()), zap.Any("collectInfo", collectInfo))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	rqItems := req.GetItems()

	// 已收集道具
	collectItems, remainItems := GetUserItemsAndRemains(ctx, collectInfo.GetItems())
	if err != nil {
		ctx.ErrorWF("OnMazeCollectItemReceiveRQ GetUserItemsAndRemains", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	//// 已收集道具
	//collectItems := collectInfo.GetItems()
	addItems := make(map[int32]int64)
	for _, item := range collectItems {
		addItems[item.GetItemId()] = item.GetCount()
	}
	if len(addItems) == 0 {
		ctx.WarnWF("OnMazeCollectItemReceiveRQ no items",
			zap.Any("rqItems", rqItems), zap.Any("collectInfo", collectInfo))
		res.ErrInfo = errors.MODULE_ERROR.Wrap("没有可领取道具")
		return
	}

	// 校验领取的道具
	if len(rqItems) > 0 {
		for _, rqItem := range rqItems {
			rqItemId, rqItemCount := rqItem.GetItemId(), rqItem.GetCount()
			if svrItemCount, ok := addItems[rqItemId]; !ok || svrItemCount < rqItemCount {
				ctx.ErrorWF("OnMazeCollectItemReceiveRQ rq items invalid", zap.Any("itemId", rqItemId),
					zap.Any("reItemCount", rqItemCount), zap.Any("svrItemCount", svrItemCount), zap.Any("exist", ok))
				res.ErrInfo = errors.MODULE_ERROR.Wrap("道具参数校验失败")
				return
			}
		}
	}

	// 清理已收集道具
	now := time.Now().Unix()
	resetCollectInfo, err := ResetMazeCollect(ctx, userId, now, collectInfo, remainItems)
	if err != nil {
		ctx.ErrorWF("OnMazeCollectItemReceiveRQ ResetMazeCollect", zap.Error(err))
		res.ErrInfo = errors.DB_SAVE_ERROR.Wrap("服务器繁忙")
		return
	}
	ctx.InfoWF("OnMazeCollectItemReceiveRQ ResetMazeCollect",
		zap.Any("old collect info", collectInfo), zap.Any("reset collect info", resetCollectInfo))

	res.FreshTime = proto.Int64(GetFreshTime(resetCollectInfo))

	tradeNo := gentradeno.GetTradeNum()
	items := Map2Common(addItems)
	errInfo := gentradeno.AddItemEx(userCtx, userCtx.UserID, 692, tradeNo, req.Header, items...)
	if errInfo != nil {
		ctx.ErrorWF("GetAllEquipDismantleAward AddItemEx fail", zap.Any("items", items))
		//res.ErrInfo = errInfo
		//return nil
	}

	mazeCollectInfoPb, err := MazeCollectToCliPB(ctx, resetCollectInfo, userInfo.PassBarrier)
	if err != nil {
		ctx.ErrorWF("ItemCollect MazeCollectToCliPB err", zap.Any("collectInfo", collectInfo), zap.Error(err))
		return
	}
	res.MazeCollectInfo = mazeCollectInfoPb
	PushDollMazeCollectInfoLog(userCtx, userId, resetCollectInfo, collectInfo.GetLastTime(), 0, mazecollectrecord.MazeCollectReceive, tradeNo, items, 0)
	return
}

// 重置宠物收集
func ResetMazeCollect(logger fklog.FKLogI, userId uint64, startTime int64, collectInfo *MazeCollectCache.MazeCollectInfo, remainItems []*MazeCollectCache.ItemInfo) (resetCollectInfo *MazeCollectCache.MazeCollectInfo, err error) {
	cfg := GMazeBarriesOnHookV8Cfg.Get(collectInfo.GetBarrierId())
	if cfg == nil {
		logger.ErrorWF("ResetMazeCollect error",
			zap.Any("barrierId", collectInfo.GetBarrierId()))
		return nil, errors.New("配置不存在")
	}
	resetCollectInfo = &MazeCollectCache.MazeCollectInfo{
		StartTime:     proto.Int64(startTime),
		PeriodTime:    proto.Int32(cfg.Cycle_time),
		EndTime:       proto.Int64(startTime + int64(cfg.Cycle_time*cfg.Maxlimit_cycle)),
		BarrierId:     proto.Int32(collectInfo.GetBarrierId()),
		AvailableTime: proto.Int64(startTime + int64(cfg.Can_receive_time)),
		Items:         remainItems,
	}
	oldEndTime := collectInfo.GetEndTime()
	oldLastTime := collectInfo.GetLastTime()
	if startTime > oldEndTime { // 挂机结束领取，保留有效挂机时间（移除挂机挂机结束后道具产出暂停的时间）
		newLastTime := startTime - (oldEndTime - oldLastTime)
		resetCollectInfo.LastTime = proto.Int64(newLastTime)
		logger.InfoWF("ResetMazeCollect onhook end receive items",
			zap.Int64("oldEndTime", oldEndTime), zap.Int64("oldLastTime", oldLastTime),
			zap.Int64("startTime", startTime), zap.Int64("newLastTime", newLastTime),
			zap.Int64("validRemainTime", oldEndTime-oldLastTime))
	} else { // 挂机未结束领取，继续产出，无需额外处理
		resetCollectInfo.LastTime = proto.Int64(oldLastTime)
		logger.InfoWF("ResetMazeCollect onhooking receive items",
			zap.Int64("startTime", startTime), zap.Int64("newLastTime", oldLastTime))
	}

	err = mazecollectredis.SetCollectInfo(logger, userId, resetCollectInfo)
	if err != nil {
		err = fmt.Errorf("SetCollectInfo err: %v", err)
		return
	}

	// 如果产出结束了设置定时器
	if IsLastCollect(collectInfo) {
		err1 := SetCollectTimer(logger, userId, resetCollectInfo)
		if err1 != nil { // 设置定时器失败不中断领取流程，用户依然能领取道具
			logger.ErrorWF("ResetMazeCollect SetCollectTimer", zap.Error(err1))
			return
		}
	}

	return
}
