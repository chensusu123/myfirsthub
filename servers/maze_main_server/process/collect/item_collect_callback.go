package collect

import (
	"context"
	"fmt"
	"strings"
	"time"

	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/settimer"
	"maze_game_server/config/GMazeBarriesOnHookV8Cfg"
	"maze_game_server/io/kafka/mazecollectrecord"
	"maze_game_server/io/redis/mazecollectredis"
	"maze_game_server/module/mazecollect"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeCollect"
	"maze_game_server/pb/server/MazeCollectCache"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type CollectMsg struct {
	UserId    uint64 `json:"user_id"`
	TimerType int32  `json:"timer_type"`
}

func (a *CollectMsg) Marshal() []byte {
	return []byte(fmt.Sprintf("%d,%d", a.UserId, a.TimerType))
}

func (a *CollectMsg) Unmarshal(data []byte) error {
	fields := strings.Split(string(data), ",")
	if len(fields) != 2 {
		return errors.New("collectMsg before start info not match")
	}

	a.UserId = fkutil.ToUint64(fields[0])
	a.TimerType = fkutil.ToInt32(fields[1])
	return nil
}

func ItemCollectCallback(ctx context.Context, userId uint64, bs []byte) error {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "ItemCollectCallback start warship skin callback", zap.Any("bs", bs))

	msg := &CollectMsg{}
	err := msg.Unmarshal(bs)
	if err != nil {
		logger.CtxError(ctx, "ItemCollectCallback json unmarshal err", zap.Error(err))
		return err
	}

	if msg.UserId != userId {
		logger.CtxError(ctx, "ItemCollectCallback petCfgId no match", zap.Any("timerUserId", msg.UserId), zap.Any("userId", userId))
		return errors.New("用户id不匹配")
	}

	//// 判断用户是否移民，移民直接丢弃定时器，移民用户登录时自修复道具产出数据
	//groupId, err := usergroup.GetUserGroupId(ctx, userId)
	//if err != nil {
	//	ctx.ErrorWF("OnItemCollectRQ GetUserGroupId err", zap.Error(err), zap.Any("userId", userId))
	//	return
	//}
	//if fkconfig.EnvVal.GroupID != uint32(groupId) {
	//	ctx.WarnWF("OnItemCollectRQ user groupId no match, maybe migrate",
	//		zap.Uint32("srcGroupId", fkconfig.EnvVal.GroupID), zap.Int64("currentGroupId", groupId))
	//	return
	//}

	// 是否已经初始化
	//collectInfo, err := mazecollectredis.GetCollectInfo(logger, userId)
	//if err != nil {
	//	logger.CtxError(ctx,"ItemCollectCallback GetCollectInfo error",
	//		zap.Any("userId", userId),
	//		zap.Error(err))
	//	return err
	//}
	cInfo := mazecollect.NewCollectInfo(ctx, userId)
	collectInfo := cInfo.GetCollectInfo(ctx)
	if collectInfo == nil {
		logger.CtxError(ctx, "ItemCollectCallback collectInfo not exist",
			zap.Any("userId", userId),
			zap.Error(err))
		return nil
	}
	err = ItemCollect(ctx, userId, collectInfo)
	if err != nil {
		logger.CtxError(ctx, "ItemCollectCallback ItemCollect", zap.Error(err))
		return err
	}
	return nil
}

func ItemCollect(ctx context.Context, userId uint64, collectInfo *MazeCollectCache.MazeCollectInfo) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	startTime := collectInfo.GetStartTime() // 收集开始时间
	lastTime := collectInfo.GetLastTime()   // 上次收集结算时间
	endTime := collectInfo.GetEndTime()     // 收集停止时间
	periodSeconds := int64(collectInfo.GetPeriodTime())

	var newLastTime int64             // 本次收集结算时间
	var collectTimes int64            // 本次产出周期数、结算周期数
	addItems := make(map[int32]int64) // 增加道具

	defer func() {
		logger.CtxInfo(ctx, "ItemCollect end",
			zap.Any("err", err), zap.Any("collect times", collectTimes),
			zap.Int64("startTime", startTime), zap.Int64("endTime", endTime),
			zap.Int64("lastTime", lastTime), zap.Int64("newLastTime", newLastTime),
			zap.Any("total collect periods", (endTime-startTime)/periodSeconds),
			zap.Any("already collect periods", (newLastTime-startTime)/periodSeconds),
			zap.Any("add items", addItems), zap.Any("collectInfo", collectInfo))
	}()

	if collectInfo.GetLastTime() >= collectInfo.GetEndTime() {
		logger.CtxError(ctx, "ItemCollect PetCollectIsEnd")
		return
	}
	cfg := GMazeBarriesOnHookV8Cfg.GetWithCtx(ctx, collectInfo.GetBarrierId())
	if cfg == nil {
		logger.CtxError(ctx, "ItemCollect error",
			zap.Any("barrierId", collectInfo.GetBarrierId()))
		return errors.New("配置不存在")
	}

	collectTime := time.Now().Unix()
	if collectTime > endTime { // 控制道具产出上限
		collectTime = endTime
	}

	timeCap := collectTime - lastTime      // 距离上次结算的时间
	collectTimes = timeCap / periodSeconds // 道具产出周期数
	if collectTimes == 0 {                 // 定时器回调提前了，不结算收集
		logger.CtxWarn(ctx, "ItemCollect timer ahead", zap.Int64("collectTime", collectTime), zap.Int64("lastTime", lastTime),
			zap.Int64("startTime", startTime), zap.Int64("endTime", endTime),
			zap.Int64("timeCap", timeCap), zap.Int64("periodSeconds", periodSeconds))
		err = SetCollectTimer(ctx, userId, collectInfo)
		if err != nil {
			err = fmt.Errorf("SetCollectTimer err: %v", err)
			return
		}
		return
	}

	for id, count := range cfg.Cycle_award {
		addItems[id] = count * collectTimes
	}
	for id, count := range cfg.Cycle_award_2 {
		addItems[id] = count * collectTimes
	}

	// 更新已收集道具
	for id, count := range addItems {
		var ok bool
		for _, item := range collectInfo.Items {
			if item.GetId() == id { // 已有道具，增加数量
				item.Count = proto.Int64(item.GetCount() + count)
				ok = true
				break
			}
		}
		if !ok { // 新增道具
			collectInfo.Items = append(collectInfo.Items, &MazeCollectCache.ItemInfo{
				Id:    proto.Int32(id),
				Count: proto.Int64(count),
			})
		}
	}

	newLastTime = lastTime + collectTimes*periodSeconds // 本次结算时间
	collectInfo.LastTime = proto.Int64(newLastTime)
	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "ItemCollect GetUserInfoV2", zap.Error(err))
		return
	}
	if userInfo.PassBarrier != collectInfo.GetBarrierId() {
		newCfg := GMazeBarriesOnHookV8Cfg.GetWithCtx(ctx, userInfo.PassBarrier)
		if newCfg == nil {
			logger.CtxError(ctx, "ItemCollect error",
				zap.Any("barrierId", userInfo.PassBarrier))
			return errors.New("配置不存在")
		}
		collectInfo.PeriodTime = proto.Int32(newCfg.Cycle_time)
		collectInfo.BarrierId = proto.Int32(userInfo.PassBarrier)
	}

	// 更新收集信息
	err = mazecollectredis.SetCollectInfo(ctx, userId, collectInfo)
	if err != nil {
		logger.CtxError(ctx, "ItemCollect SetCollectInfo", zap.Error(err))
		return
	}

	if newLastTime == endTime {
		logger.CtxInfo(ctx, "ItemCollect collect stop", zap.Any("collectInfo", collectInfo))
		return nil
	}

	// 设置下一周期定时器
	err = SetCollectTimer(ctx, userId, collectInfo)
	if err != nil {
		logger.CtxError(ctx, "ItemCollect SetCollectTimer err", zap.Any("collectInfo", collectInfo), zap.Error(err))
		return
	}

	err = PushDollMazeCollectInfoLog(ctx, userId, collectInfo, lastTime, collectTimes, mazecollectrecord.MazeCollectTimeOut, 0, nil, 0)
	if err != nil {
		logger.CtxError(ctx, "ItemCollect PushDollMazeCollectInfoLog err", zap.Any("collectInfo", collectInfo), zap.Error(err))
		return
	}
	return
}

// 设置道具收集定时器
func SetCollectTimer(ctx context.Context, userId uint64, info *MazeCollectCache.MazeCollectInfo) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	lastTime := info.GetLastTime()
	periodSeconds := info.GetPeriodTime() // 产出周期
	endTime := info.GetEndTime()
	expireAt := lastTime + int64(periodSeconds)
	if expireAt > endTime {
		expireAt = endTime
	}
	msg := &CollectMsg{
		UserId:    userId,
		TimerType: constdef.CollectNotice,
	}
	jsonData := msg.Marshal()
	err = settimer.SetTaskExpire(ctx, userId,
		234, expireAt, jsonData)
	if err != nil {
		err = fmt.Errorf("SetTaskExpire err: %v", err)
		return
	}

	logger.CtxInfo(ctx, "SetCollectTimer",
		zap.Any("msg", msg), zap.Int64("expireAt", expireAt),
		zap.Int64("lastTime", lastTime), zap.Int64("endTime", endTime))
	return
}

func NewCollectAfter(ctx context.Context, userId uint64, collectInfo *MazeCollectCache.MazeCollectInfo) {
	logger := fklog.ContextAppLogger(ctx)
	// 设置下一周期定时器
	err := SetCollectTimer(ctx, userId, collectInfo)
	if err != nil {
		logger.CtxError(ctx, "NewCollectAfter SetCollectTimer err", zap.Any("collectInfo", collectInfo), zap.Error(err))
		return
	}
	err = PushDollMazeCollectInfoLog(ctx, userId, collectInfo, collectInfo.GetLastTime(), 0, mazecollectrecord.MazeCollectInit, 0, nil, 0)
	if err != nil {
		logger.CtxError(ctx, "NewCollectAfter PushDollMazeCollectInfoLog err", zap.Any("collectInfo", collectInfo), zap.Error(err))
		return
	}

	pack := &MazeCollect.MazeCollectOpenID{
		FreshTime: proto.Int64(GetFreshTime(collectInfo)),
	}
	mazeCollectInfoPb, err := MazeCollectToCliPB(ctx, collectInfo, collectInfo.GetBarrierId())
	if err != nil {
		logger.CtxError(ctx, "NewCollectAfter MazeCollectToCliPB err", zap.Any("collectInfo", collectInfo), zap.Error(err))
		return
	}
	pack.MazeCollectInfo = mazeCollectInfoPb
	err = online.ClusterPush(context.TODO(), uint64(userId), 10480, pack)
	if err != nil {
		logger.CtxError(ctx, "NewCollectAfter SendArrivePacket", zap.Any("pack", pack), zap.Error(err))
		return
	}
	logger.CtxInfo(ctx, "NewCollectAfter SendArrivePacket", zap.Any("pack", pack))
}
