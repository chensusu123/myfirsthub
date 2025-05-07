package collect

import (
	"encoding/json"
	"fmt"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/gentradeno"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazecollectrecord"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazecollectredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/funcopencheck"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserinfo"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_collect_server/process"
	"gitlab.ifreetalk.com/plate/definition/uncgkconst"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeBarriesOnHookV8Cfg"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/kafka_consumer"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/plate/protodef/MazeCollect"
	"gitlab.ifreetalk.com/plate/protodef/MazeCollectCache"
	"gitlab.ifreetalk.com/plate/protodef/SeaTaskSvr"
	"go.uber.org/zap"
	"time"
)

func RegTcpHandler() {
	// 迷宫挂机查询
	_ = tcp_service.RegProcSimple(16257, &MazeCollect.MazeCollectInfoQueryRQ{},
		16258, &MazeCollect.MazeCollectInfoQueryRS{}, OnMazeCollectInfoQueryRQ)

	_ = tcp_service.RegProcSimple(16259, &MazeCollect.MazeCollectItemReceiveRQ{},
		16260, &MazeCollect.MazeCollectItemReceiveRS{}, OnMazeCollectItemReceiveRQ)

	_ = tcp_service.RegProcSimple(uncgkconst.UN_TCP_PACK_SVR_SEA_TASK_EXPIRE_NOTIFY_RQ, &SeaTaskSvr.TaskExpireNotifyRQ{},
		uncgkconst.UN_TCP_PACK_SVR_SEA_TASK_EXPIRE_NOTIFY_RS, &SeaTaskSvr.TaskExpireNotifyRS{}, OnTimeOut)
}

func RegisgterKafa() {
	kafka_consumer.PlugKafkaConsumer("maze_barrier_chg_msg",
		1001105,
		kafka_consumer.WithGroup(fkserver.MonitorName),
		kafka_consumer.WithKafkaCustomKeyContent(process.HandleMazeBarrierMsg))

	kafka_consumer.PlugKafkaConsumer("maze_level_chg_msg",
		1001084,
		kafka_consumer.WithGroup(fkserver.MonitorName),
		kafka_consumer.WithKafkaCustomKeyContent(process.HandleMazeLevelMsg))
}

// 用户迷宫闯关纪录
type MazeBarrierUserGameRecord struct {
	UserId     uint64 `json:"user_id"`     // 用户id
	Barrier    int32  `json:"barrier"`     // 关卡id
	GameRet    int32  `json:"game_ret"`    // 用户闯关结果 1-通关成功 2-死亡失败
	GroupID    uint32 `json:"group_id"`    // 组id
	CreateTime int64  `json:"create_time"` // 操作时间
}

func HandleMazeBarrierMsg(ctx context.Context, logger fklog.FKLogI, index int, key, data []byte) (err error) {
	msg := &MazeBarrierUserGameRecord{}
	err = json.Unmarshal(data, msg)
	if err != nil {
		logger.ErrorWF("HandleMazeBarrierMsg Unmarshal", zap.Error(err),
			zap.Int("msg's len", len(data)), zap.Uint64("userId", msg.UserId))
		return
	}

	userId := msg.UserId
	logger.InfoWF("HandleMazeBarrierMsg start", zap.Any("msg", msg),
		zap.Uint64("userId", userId))
	if userId <= 0 {
		return
	}
	if msg.GameRet != 1 {
		return nil
	}
	// 道具产出信息
	collectInfo, err := mazecollectredis.GetCollectInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("HandleMazeBarrierMsg GetCollectInfo", zap.Error(err))
		return
	}
	if collectInfo != nil {
		return nil
	}
	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("HandleMazeBarrierMsg GetUserInfoV2 fail", zap.Error(err))
		return
	}
	if userInfo.PassBarrier <= 0 {
		return nil
	}

	result, err := funcopencheck.IsFuncOpen(1, int32(userInfo.Level))
	if !result.IsOpen {
		return nil
	}
	err = InitMazeCollectLand(logger, userId, userInfo.PassBarrier)
	if err != nil {
		logger.ErrorWF("HandleMazeBarrierMsg InitMazeCollectLand", zap.Error(err))
		return err
	}
	return nil
}

// 用户等级变化流水
type MazeUserLevelRecord struct {
	UserId     uint64 `json:"user_id"`     // 用户id
	OldLevel   int32  `json:"old_level"`   // 旧等级
	NewLevel   int32  `json:"new_level"`   // 新等级
	GroupID    uint32 `json:"group_id"`    // 组id
	CreateTime int64  `json:"create_time"` // 操作时间 毫秒
}

func HandleMazeLevelMsg(ctx context.Context, logger fklog.FKLogI, index int, key, data []byte) (err error) {
	msg := &MazeUserLevelRecord{}
	err = json.Unmarshal(data, msg)
	if err != nil {
		logger.ErrorWF("HandleMazeLevelMsg Unmarshal", zap.Error(err),
			zap.Int("msg's len", len(data)), zap.Uint64("userId", msg.UserId))
		return
	}

	userId := msg.UserId
	logger.InfoWF("HandleMazeLevelMsg start", zap.Any("msg", msg),
		zap.Uint64("userId", userId))
	if userId <= 0 {
		return
	}
	// 道具产出信息
	collectInfo, err := mazecollectredis.GetCollectInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("HandleMazeLevelMsg GetCollectInfo", zap.Error(err))
		return
	}
	if collectInfo != nil {
		return nil
	}
	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("HandleMazeLevelMsg GetUserInfoV2 fail", zap.Error(err))
		return
	}
	if userInfo.PassBarrier <= 0 {
		return nil
	}

	result, err := funcopencheck.IsFuncOpen(1, int32(userInfo.Level))
	if !result.IsOpen {
		return nil
	}
	err = InitMazeCollectLand(logger, userId, userInfo.PassBarrier)
	if err != nil {
		logger.ErrorWF("HandleMazeLevelMsg InitMazeCollectLand", zap.Error(err))
		return err
	}
	return nil
}

// 道具收集查询
func OnMazeCollectInfoQueryRQ(ctx fknet.TCPContext, userId uint64, rq proto.Message, rs proto.Message) (err error) {
	req := rq.(*MazeCollect.MazeCollectInfoQueryRQ)
	res := rs.(*MazeCollect.MazeCollectInfoQueryRS)
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.QueryType = req.QueryType

	ctx.InfoWF("OnMazeCollectInfoQueryRQ start", zap.Any("req", req))
	defer func() {
		ctx.InfoWF("OnMazeCollectInfoQueryRQ end", zap.Any("res", res))
	}()

	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		ctx.ErrorWF("OnMazeCollectInfoQueryRQ GetUserInfoV2", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return err
	}

	// 道具产出信息
	collectInfo, err := mazecollectredis.GetCollectInfo(ctx, userId)
	if err != nil {
		ctx.ErrorWF("OnMazeCollectInfoQueryRQ GetCollectInfo", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	if collectInfo == nil {
		if userInfo.PassBarrier <= 0 {
			return nil
		}
		err = InitMazeCollectLand(ctx, userId, userInfo.PassBarrier)
		if err != nil {
			ctx.ErrorWF("OnMazeCollectInfoQueryRQ InitMazeCollectLand", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return err
		}
		// 道具产出信息
		collectInfo, err = mazecollectredis.GetCollectInfo(ctx, userId)
		if err != nil {
			ctx.ErrorWF("OnMazeCollectInfoQueryRQ GetCollectInfo", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return err
		}
	}
	if collectInfo == nil {
		return nil
	}
	if IsTimerLoss(collectInfo) {
		// 定时器丢失修复道具产出
		ctx.InfoWF("OnMazeCollectInfoQueryRQ fix ItemCollect start", zap.Any("collectInfo", collectInfo))
		err = ItemCollect(ctx, userId, collectInfo)
		if err != nil {
			ctx.ErrorWF("OnPetCollectInfoQueryRQ fix ItemCollect", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}
		ctx.InfoWF("OnPetCollectInfoQueryRQ fix ItemCollect end", zap.Any("collectInfo", collectInfo))
	}

	mazeCollectInfoPb, err := MazeCollectToCliPB(ctx, collectInfo, userInfo.PassBarrier)
	if err != nil {
		ctx.ErrorWF("ItemCollect MazeCollectToCliPB err", zap.Any("collectInfo", collectInfo), zap.Error(err))
		return err
	}
	res.MazeCollectInfo = mazeCollectInfoPb
	res.FreshTime = proto.Int64(GetFreshTime(collectInfo))
	return
}

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

func OnTimeOut(ctx fknet.TCPContext, shardingID uint64, request proto.Message, response proto.Message) (err error) {
	defer fkprometheus.InfoPMT("OnTimeOut")()
	req := request.(*SeaTaskSvr.TaskExpireNotifyRQ)
	res := response.(*SeaTaskSvr.TaskExpireNotifyRS)
	res.TaskInfo = &SeaTaskSvr.TaskInfo{UserId: req.TaskInfo.UserId}
	res.ErrInfo = errors.NO_ERROR

	defer func() {
		ctx.InfoWF("OnTimeOut end", zap.Any("res", res))
	}()
	ctx.InfoWF("OnTimeOut with ", zap.Any("Msg", req))
	agent := fkserver.NewUserContext(ctx.Context, uint64(shardingID), ctx.FKLogI)

	taskInfo := req.GetTaskInfo()
	if req.GetTaskInfo() == nil {
		agent.ErrorWF("OnTimeOut taskinfo nil", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("task info nil")
		return
	}
	if uint32(time.Now().Unix()) < taskInfo.GetTime() {
		ctx.ErrorWF("OnTimeOut check time failed. time not touch,call later.",
			zap.Uint64("uid", taskInfo.GetUserId()),
			zap.Stringer("task", taskInfo), zap.Uint64("shardingId", shardingID),
		)
		res.ErrInfo = errors.ARGS_NOT_MATCH.Wrap("时间还没到")
		return
	}
	if taskInfo.GetType() != uint32(234) {
		agent.ErrorWF("OnTimeOut task typ not match", zap.Any("req", req), zap.Uint32("myType", 234))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("task typ not match")
		return
	}
	return ItemCollectCallback(agent, agent.UserID, taskInfo.GetContext())
}

// 初始化迷宫挂机
func InitMazeCollectLand(logger fklog.FKLogI, userId uint64, barrierId int32) (err error) {
	cfg := GMazeBarriesOnHookV8Cfg.Get(barrierId)
	if cfg == nil {
		logger.ErrorWF("InitMazeCollectLand error",
			zap.Any("barrierId", barrierId))
		return errors.New("配置不存在")
	}

	// 是否已经初始化
	collectInfo, err := mazecollectredis.GetCollectInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("InitMazeCollectLand error",
			zap.Any("barrierId", barrierId),
			zap.Error(err))
		return
	}
	now := time.Now().Unix()
	collectInfo = &MazeCollectCache.MazeCollectInfo{}
	collectInfo.StartTime = proto.Int64(now)
	collectInfo.LastTime = proto.Int64(now)
	collectInfo.PeriodTime = proto.Int32(cfg.Cycle_time)
	collectInfo.EndTime = proto.Int64(now + int64(cfg.Cycle_time*cfg.Maxlimit_cycle))
	collectInfo.Items = []*MazeCollectCache.ItemInfo{}
	collectInfo.BarrierId = proto.Int32(barrierId)
	collectInfo.AvailableTime = proto.Int64(now + int64(cfg.Can_receive_time))
	logger.InfoWF("InitMazeCollectLand init pet", zap.Any("collectInfo", collectInfo))

	err = mazecollectredis.SetCollectInfo(logger, userId, collectInfo)
	if err != nil {
		logger.ErrorWF("InitMazeCollectLand SetCollectInfo error", zap.Any("collectInfo", collectInfo), zap.Error(err))
		return
	}

	// 设置下一周期定时器
	err = SetCollectTimer(logger, userId, collectInfo)
	if err != nil {
		logger.ErrorWF("InitMazeCollectLand SetCollectTimer err", zap.Any("collectInfo", collectInfo), zap.Error(err))
		return
	}
	PushDollMazeCollectInfoLog(logger, userId, collectInfo, collectInfo.GetLastTime(), 0, mazecollectrecord.MazeCollectInit, 0, nil, 0)

	pack := &MazeCollect.MazeCollectOpenID{
		FreshTime: proto.Int64(GetFreshTime(collectInfo)),
	}
	mazeCollectInfoPb, err := MazeCollectToCliPB(logger, collectInfo, collectInfo.GetBarrierId())
	if err != nil {
		logger.ErrorWF("InitMazeCollectLand MazeCollectToCliPB err", zap.Any("collectInfo", collectInfo), zap.Error(err))
		return
	}
	pack.MazeCollectInfo = mazeCollectInfoPb
	err = MustArriveRedis.SendArrivePacket(userId, 16261, pack)
	if err != nil {
		logger.ErrorWF("InitMazeCollectLand SendArrivePacket", zap.Any("pack", pack), zap.Error(err))
		return
	}
	logger.InfoWF("InitMazeCollectLand SendArrivePacket", zap.Any("pack", pack))
	return
}
