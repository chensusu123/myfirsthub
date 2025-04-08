package process

import (
	"fmt"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/settimer"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazecollectredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserinfo"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeBarriesOnHookV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkutil"
	"gitlab.ifreetalk.com/plate/protodef/MazeCollectCache"
	"go.uber.org/zap"
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

func ItemCollectCallback(logger fklog.FKLogI, userId uint64, bs []byte) error {
	logger.InfoWF("ItemCollectCallback start warship skin callback", zap.Any("bs", bs))

	msg := &CollectMsg{}
	err := msg.Unmarshal(bs)
	if err != nil {
		logger.ErrorWF("ItemCollectCallback json unmarshal err", zap.Error(err))
		return err
	}

	if msg.UserId != userId {
		logger.ErrorWF("ItemCollectCallback petCfgId no match", zap.Any("timerUserId", msg.UserId), zap.Any("userId", userId))
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
	collectInfo, err := mazecollectredis.GetCollectInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("ItemCollectCallback GetCollectInfo error",
			zap.Any("userId", userId),
			zap.Error(err))
		return err
	}
	if collectInfo == nil {
		logger.ErrorWF("ItemCollectCallback collectInfo not exist",
			zap.Any("userId", userId),
			zap.Error(err))
		return nil
	}
	err = ItemCollect(logger, userId, collectInfo)
	if err != nil {
		logger.ErrorWF("ItemCollectCallback ItemCollect", zap.Error(err))
		return err
	}
	return nil
}

func ItemCollect(logger fklog.FKLogI, userId uint64, collectInfo *MazeCollectCache.MazeCollectInfo) (err error) {
	startTime := collectInfo.GetStartTime() // 收集开始时间
	lastTime := collectInfo.GetLastTime()   // 上次收集结算时间
	endTime := collectInfo.GetEndTime()     // 收集停止时间
	periodSeconds := int64(collectInfo.GetPeriodTime())

	var newLastTime int64             // 本次收集结算时间
	var collectTimes int64            // 本次产出周期数、结算周期数
	addItems := make(map[int32]int64) // 增加道具

	defer func() {
		logger.InfoWF("ItemCollect end",
			zap.Any("err", err), zap.Any("collect times", collectTimes),
			zap.Int64("startTime", startTime), zap.Int64("endTime", endTime),
			zap.Int64("lastTime", lastTime), zap.Int64("newLastTime", newLastTime),
			zap.Any("total collect periods", (endTime-startTime)/periodSeconds),
			zap.Any("already collect periods", (newLastTime-startTime)/periodSeconds),
			zap.Any("add items", addItems), zap.Any("collectInfo", collectInfo))
	}()

	if collectInfo.GetLastTime() >= collectInfo.GetEndTime() {
		logger.ErrorWF("ItemCollect PetCollectIsEnd")
		return
	}
	cfg := GMazeBarriesOnHookV8Cfg.Get(collectInfo.GetBarrierId())
	if cfg == nil {
		logger.ErrorWF("ItemCollect error",
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
		logger.WarnWF("ItemCollect timer ahead", zap.Int64("collectTime", collectTime), zap.Int64("lastTime", lastTime),
			zap.Int64("startTime", startTime), zap.Int64("endTime", endTime),
			zap.Int64("timeCap", timeCap), zap.Int64("periodSeconds", periodSeconds))
		err = SetCollectTimer(logger, userId, collectInfo)
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
	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("ItemCollect GetUserInfoV2", zap.Error(err))
		return
	}
	if userInfo.PassBarrier != collectInfo.GetBarrierId() {
		newCfg := GMazeBarriesOnHookV8Cfg.Get(userInfo.PassBarrier)
		if newCfg == nil {
			logger.ErrorWF("ItemCollect error",
				zap.Any("barrierId", userInfo.PassBarrier))
			return errors.New("配置不存在")
		}
		collectInfo.PeriodTime = proto.Int32(newCfg.Cycle_time)
		collectInfo.BarrierId = proto.Int32(userInfo.PassBarrier)
	}

	// 更新收集信息
	err = mazecollectredis.SetCollectInfo(logger, userId, collectInfo)
	if err != nil {
		logger.ErrorWF("ItemCollect SetCollectInfo", zap.Error(err))
		return
	}

	if newLastTime == endTime {
		logger.InfoWF("ItemCollect collect stop", zap.Any("collectInfo", collectInfo))
		return nil
	}

	// 设置下一周期定时器
	err = SetCollectTimer(logger, userId, collectInfo)
	if err != nil {
		logger.ErrorWF("ItemCollect SetCollectTimer err", zap.Any("collectInfo", collectInfo), zap.Error(err))
		return
	}

	return
}

// 设置道具收集定时器
func SetCollectTimer(logger fklog.FKLogI, userId uint64, info *MazeCollectCache.MazeCollectInfo) (err error) {
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
	err = settimer.SetTaskExpire(logger, userId,
		234, expireAt, jsonData)
	if err != nil {
		err = fmt.Errorf("SetTaskExpire err: %v", err)
		return
	}

	logger.InfoWF("SetCollectTimer",
		zap.Any("msg", msg), zap.Int64("expireAt", expireAt),
		zap.Int64("lastTime", lastTime), zap.Int64("endTime", endTime))
	return
}
