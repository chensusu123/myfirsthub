package collect

import (
	"context"
	"fmt"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/common/tradeno"
	"maze_game_server/config/GMazeBarriesOnHookV8Cfg"
	"maze_game_server/io/kafka/mazecollectrecord"
	"maze_game_server/io/redis/mazecollectredis"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazecollect"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeCollect"
	"maze_game_server/pb/server/MazeCollectCache"
	"maze_game_server/services/itemservice"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 挂机道具领取
func (c *Collect) OnMazeCollectItemReceiveRQ_10467_10468(s *session.Session, req *MazeCollect.MazeCollectItemReceiveRQ) (err error) {

	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeCollect.MazeCollectItemReceiveRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userId := uint64(s.UID())

	logger.CtxInfo(ctx, "OnMazeCollectItemReceiveRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnMazeCollectItemReceiveRQ end", zap.Any("res", res))
	}()

	// 是否已经初始化
	//collectInfo, err := mazecollectredis.GetCollectInfo(logger, userId)
	//if err != nil {
	//	logger.CtxError(ctx,"OnMazeCollectItemReceiveRQ GetCollectInfo error",
	//		zap.Any("userId", userId),
	//		zap.Error(err))
	//	return
	//}
	cInfo := mazecollect.NewCollectInfo(ctx, userId)
	collectInfo := cInfo.GetCollectInfo(ctx)
	if collectInfo == nil {
		logger.CtxError(ctx, "OnMazeCollectItemReceiveRQ GetCollectInfo error", zap.Any("userId", userId))
		return
	}

	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnMazeCollectItemReceiveRQ GetUserInfoV2", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return err
	}

	// 确保领取时道具产出是正确的
	if !CheckReceiveItems(collectInfo) {
		logger.CtxError(ctx, "OnMazeCollectItemReceiveRQ CheckReceiveItems",
			zap.Int64("nowSecond", time.Now().Unix()), zap.Any("collectInfo", collectInfo))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	rqItems := req.GetItems()

	// 已收集道具
	collectItems, remainItems := GetUserItemsAndRemains(ctx, collectInfo.GetItems())
	if err != nil {
		logger.CtxError(ctx, "OnMazeCollectItemReceiveRQ GetUserItemsAndRemains", zap.Error(err))
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
		logger.CtxWarn(ctx, "OnMazeCollectItemReceiveRQ no items",
			zap.Any("rqItems", rqItems), zap.Any("collectInfo", collectInfo))
		res.ErrInfo = errors.MODULE_ERROR.Wrap("没有可领取道具")
		return
	}

	// 校验领取的道具
	if len(rqItems) > 0 {
		for _, rqItem := range rqItems {
			rqItemId, rqItemCount := rqItem.GetItemId(), rqItem.GetCount()
			if svrItemCount, ok := addItems[rqItemId]; !ok || svrItemCount < rqItemCount {
				logger.CtxError(ctx, "OnMazeCollectItemReceiveRQ rq items invalid", zap.Any("itemId", rqItemId),
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
		logger.CtxError(ctx, "OnMazeCollectItemReceiveRQ ResetMazeCollect", zap.Error(err))
		res.ErrInfo = errors.DB_SAVE_ERROR.Wrap("服务器繁忙")
		return
	}
	logger.CtxInfo(ctx, "OnMazeCollectItemReceiveRQ ResetMazeCollect",
		zap.Any("old collect info", collectInfo), zap.Any("reset collect info", resetCollectInfo))

	res.FreshTime = proto.Int64(GetFreshTime(resetCollectInfo))

	tradeNo := tradeno.GetTradeNum()
	items := Map2Common(addItems)
	awardItems := itemutil.Map2ItemInfo(addItems)
	errInfo := itemservice.GlobalItemService.AddItem(context.TODO(), userId, itemservice.ItemOpTypeCollect, tradeNo, awardItems...)
	if errInfo != nil {
		logger.CtxError(ctx, "GetAllEquipDismantleAward AddItemEx fail", zap.Any("items", items))
		//res.ErrInfo = errInfo
		//return nil
	}

	mazeCollectInfoPb, err := MazeCollectToCliPB(ctx, resetCollectInfo, userInfo.PassBarrier)
	if err != nil {
		logger.CtxError(ctx, "ItemCollect MazeCollectToCliPB err", zap.Any("collectInfo", collectInfo), zap.Error(err))
		return
	}
	res.MazeCollectInfo = mazeCollectInfoPb
	err = PushDollMazeCollectInfoLog(ctx, userId, resetCollectInfo, collectInfo.GetLastTime(), 0, mazecollectrecord.MazeCollectReceive, tradeNo, items, 0)
	if err != nil {
		logger.CtxError(ctx, "ItemCollect PushDollMazeCollectInfoLog err", zap.Any("collectInfo", collectInfo), zap.Error(err))
		return
	}
	return
}

// 重置挂机收集
func ResetMazeCollect(ctx context.Context, userId uint64, startTime int64, collectInfo *MazeCollectCache.MazeCollectInfo, remainItems []*MazeCollectCache.ItemInfo) (resetCollectInfo *MazeCollectCache.MazeCollectInfo, err error) {
	logger := fklog.ContextAppLogger(ctx)
	cfg := GMazeBarriesOnHookV8Cfg.GetWithCtx(ctx, collectInfo.GetBarrierId())
	if cfg == nil {
		logger.CtxError(ctx, "ResetMazeCollect error",
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
		logger.CtxInfo(ctx, "ResetMazeCollect onhook end receive items",
			zap.Int64("oldEndTime", oldEndTime), zap.Int64("oldLastTime", oldLastTime),
			zap.Int64("startTime", startTime), zap.Int64("newLastTime", newLastTime),
			zap.Int64("validRemainTime", oldEndTime-oldLastTime))
	} else { // 挂机未结束领取，继续产出，无需额外处理
		resetCollectInfo.LastTime = proto.Int64(oldLastTime)
		logger.CtxInfo(ctx, "ResetMazeCollect onhooking receive items",
			zap.Int64("startTime", startTime), zap.Int64("newLastTime", oldLastTime))
	}

	err = mazecollectredis.SetCollectInfo(ctx, userId, resetCollectInfo)
	if err != nil {
		err = fmt.Errorf("SetCollectInfo err: %v", err)
		return
	}

	// 如果产出结束了设置定时器
	if IsLastCollect(collectInfo) {
		err1 := SetCollectTimer(ctx, userId, resetCollectInfo)
		if err1 != nil { // 设置定时器失败不中断领取流程，用户依然能领取道具
			logger.CtxError(ctx, "ResetMazeCollect SetCollectTimer", zap.Error(err1))
			return
		}
	}

	return
}
