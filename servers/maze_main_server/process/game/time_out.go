package game

import (
	"fmt"
	"strings"

	"gitlab.ifreetalk.com/maze/maze_game_server/common/errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig/param"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
)

var HeartTime int64
var StopProduceTime int64

func init() {
	param.Int64P(&HeartTime, "heart:time", 5, "心跳间隔")
	param.Int64P(&StopProduceTime, "stop:produce:time", 10, "停止生产的最长心跳间隔")
}

type DollMazeProduceMsg struct {
	UserId  uint64
	Barrier int32
	AreaId  int32
}

func (dmp *DollMazeProduceMsg) Marshal() (data []byte) {
	return []byte(fmt.Sprintf("%d-%d-%d", dmp.UserId, dmp.Barrier, dmp.AreaId))
}

func (dmp *DollMazeProduceMsg) Unmarshal(data []byte) (err error) {
	fields := strings.Split(string(data), "-")
	if len(fields) != 3 {
		return errors.New("unmarshal fail")
	}

	dmp.UserId = fkutil.ToUint64(fields[0])
	dmp.Barrier = fkutil.ToInt32(fields[1])
	dmp.AreaId = fkutil.ToInt32(fields[2])
	return
}

// func SetTimer(logger fklog.FKLogI, uid uint64, expireTime int64, barrierId, areaId int32) (err error) {
// 	logger.WarnWF("setTimer with", zap.Uint64("uid", uid), zap.Int32("barrierId", barrierId), zap.Int64("expireTime", expireTime))

// 	msg := &DollMazeProduceMsg{
// 		UserId:  uid,
// 		Barrier: barrierId,
// 		AreaId:  areaId,
// 	}

// 	jsonData := msg.Marshal()
// 	err = settimer.SetTaskExpire(context.TODO(), logger, uid, 233, expireTime, jsonData)
// 	if err != nil {
// 		logger.ErrorWF("setTimer push to delay task queue fail", zap.Error(err), zap.Int64("expireTime", expireTime),
// 			zap.Any("taskInfo", jsonData), zap.Any("jsonData", jsonData))
// 	} else {
// 		logger.WarnWF("setTimer push to delay task queue succ", zap.Int64("expireTime", expireTime),
// 			zap.Any("taskInfo", jsonData), zap.Any("jsonData", jsonData))
// 	}
// 	return
// }

// func RmTimer(logger fklog.FKLogI, uid uint64, barrierId, areaId int32) (err error) {
// 	msg := &DollMazeProduceMsg{
// 		UserId:  uid,
// 		Barrier: barrierId,
// 		AreaId:  areaId,
// 	}
// 	validTime := int64(0)
// 	logger.WarnWF("RmTimer start with", zap.Uint64("uid", uid))
// 	jsonData := msg.Marshal()
// 	err = settimer.RemoveTaskTimer(logger, uid, 233, validTime, jsonData)
// 	if err != nil {
// 		logger.ErrorWF("RmTimer push to delay task queue fail", zap.Error(err), zap.Int64("validTime", validTime), zap.Any("taskInfo", jsonData), zap.Any("jsonData", jsonData))
// 	} else {
// 		logger.WarnWF("RmTimer push to delay task queue succ", zap.Int64("validTime", validTime), zap.Any("taskInfo", jsonData), zap.Any("jsonData", jsonData))
// 	}
// 	return
// }

// func OnTimeOut(ctx fknet.TCPContext, shardingID uint64, request proto.Message, response proto.Message) (err error) {
// 	defer fkprometheus.InfoPMT("OnTimeOut")()
// 	req := request.(*SeaTaskSvr.TaskExpireNotifyRQ)
// 	res := response.(*SeaTaskSvr.TaskExpireNotifyRS)
// 	res.TaskInfo = &SeaTaskSvr.TaskInfo{UserId: req.TaskInfo.UserId}
// 	res.ErrInfo = errors.NO_ERROR

// 	defer func() {
// 		ctx.InfoWF("OnTimeOut end", zap.Any("res", res))
// 	}()
// 	ctx.InfoWF("OnTimeOut with ", zap.Any("Msg", req))
// 	taskInfo := req.GetTaskInfo()
// 	if req.GetTaskInfo() == nil {
// 		ctx.ErrorWF("OnTimeOut taskInfo nil", zap.Any("req", req))
// 		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("task info nil")
// 		return
// 	}
// 	if uint32(time.Now().Unix()) < taskInfo.GetTime() {
// 		ctx.ErrorWF("OnTimeOut check time failed. time not touch,call later.",
// 			zap.Uint64("uid", taskInfo.GetUserId()),
// 			zap.Stringer("task", taskInfo), zap.Uint64("shardingId", shardingID),
// 		)
// 		res.ErrInfo = errors.ARGS_NOT_MATCH.Wrap("时间还没到")
// 		return
// 	}

// 	// switch taskInfo.GetType() {
// 	// case 233: //ENUM_TASK_TYPE_DOLL_MAZE_PRODUCT                      = 233; /// 人偶版本迷宫生产倒计时 --张登元	调用代理17762
// 	// 	err = DollMazeProduceCallBack(ctx, taskInfo.GetContext())
// 	// 	if err != nil {
// 	// 		ctx.ErrorWF("OnTimeOut DollMazeProduceCallBack error", zap.Any("req", req), zap.Uint32("myType", taskInfo.GetType()))
// 	// 		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("deal maze produce func wrong")
// 	// 		return
// 	// 	}
// 	// default:
// 	// 	ctx.ErrorWF("OnTimeOut task typ not match", zap.Any("req", req), zap.Uint32("type", taskInfo.GetType()))
// 	// 	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("task typ not match")
// 	// 	return
// 	// }

// 	return nil
// }

func DollMazeProduceCallBack(logger fklog.FKLogI, bs []byte) (err error) {
	// defer fkprometheus.InfoPMT("DollMazeProduceCallBack")()
	// logger.InfoWF("DollMazeProduceCallBack start", zap.Any("bs", bs))

	// msg := &DollMazeProduceMsg{}
	// err = msg.Unmarshal(bs)
	// if err != nil {
	// 	logger.ErrorWF("DollMazeProduceCallBack unmarshal err", zap.Error(err))
	// 	return
	// }

	// userId := msg.UserId
	// barrier := msg.Barrier
	// areaId := msg.AreaId

	// barrierId, _, highArea, err := dollmazebarrier.GetMazeInfo(logger, userId)
	// if err != nil {
	// 	logger.ErrorWF("DollMazeProduceCallBack GetMazeInfo fail", zap.Error(err))
	// 	return
	// }
	// if barrier != barrierId || areaId != highArea {
	// 	logger.WarnWF("DollMazeProduceCallBack barrier area not mathch", zap.Any("msg", msg),
	// 		zap.Any("barrierId", barrierId), zap.Any("highArea", highArea))
	// 	return
	// }

	// produce, err := dollmazeproduceredis.GetUserProduce(logger, userId)
	// if err != nil {
	// 	logger.ErrorWF("DollMazeProduceCallBack GetUserProduce fail", zap.Error(err))
	// 	return
	// }

	// _, _, tempMax, err := GetCurrProduceCfg(logger, userId)
	// if err != nil {
	// 	logger.ErrorWF("DollMazeProduceCallBack GetCurrProduceCfgByArea fail", zap.Error(err))
	// 	return
	// }
	// if tempMax == 0 {
	// 	logger.ErrorWF("DollMazeProduceCallBack load limit num fail")
	// 	return
	// }

	// if tempMax == produce.GetProduceCount() {
	// 	//如果当前值就是上限 则不开启下个周期
	// 	logger.WarnWF("DollMazeProduceCallBack produce limit", zap.Any("msg", msg),
	// 		zap.Any("tempMax", tempMax), zap.Any("produce", produce))
	// 	return
	// }

	// if produce.GetProduceCount()+produce.GetCycleAddCount() < tempMax {
	// 	produce.ProduceCount = proto.Int64(produce.GetProduceCount() + produce.GetCycleAddCount())
	// } else {
	// 	produce.ProduceCount = proto.Int64(tempMax)
	// }
	// produce.LastSettleTime = proto.Int64(time.Now().Unix())

	// newProduce, isProduce, err := checkProduceContinue(logger, userId, barrier, areaId, produce)
	// if err != nil {
	// 	logger.ErrorWF("DollMazeProduceCallBack checkProduceContinue fail", zap.Error(err))
	// 	return
	// }
	// err = dollmazeproduceredis.SetUserProduce(logger, userId, newProduce)
	// if err != nil {
	// 	logger.ErrorWF("DollMazeProduceCallBack SetUserProduce fail", zap.Error(err), zap.Any("produce", newProduce))
	// 	return
	// }

	// if !isProduce {
	// 	logger.WarnWF("DollMazeProduceCallBack produce no need next timer")
	// 	return
	// }

	// err2 := SetTimer(logger, userId, newProduce.GetLastSettleTime()+newProduce.GetCycleTime(), barrier, newProduce.GetAreaId())
	// if err2 != nil {
	// 	logger.ErrorWF("DollMazeProduceCallBack SetTimer fail", zap.Error(err2), zap.Any("produce", newProduce))
	// 	return
	// }
	return
}
