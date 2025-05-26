package collect

import (
	"time"

	"google.golang.org/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/pb/server/SeaTaskSvr"
)

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

func ProcessTimeOut(logger fklog.FKLogI, shardingID uint64, req SeaTaskSvr.TaskExpireNotifyRQ) (res SeaTaskSvr.TaskExpireNotifyRS, err error) {
	res.TaskInfo = &SeaTaskSvr.TaskInfo{UserId: req.TaskInfo.UserId}
	res.ErrInfo = errors.NO_ERROR

	defer func() {
		logger.InfoWF("ProcessTimeOut end", zap.Any("res", res))
	}()
	logger.InfoWF("ProcessTimeOut with ", zap.Any("Msg", req))

	taskInfo := req.GetTaskInfo()
	if req.GetTaskInfo() == nil {
		logger.ErrorWF("ProcessTimeOut taskinfo nil", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("task info nil")
		return
	}
	if uint32(time.Now().Unix()) < taskInfo.GetTime() {
		logger.ErrorWF("ProcessTimeOut check time failed. time not touch,call later.",
			zap.Uint64("uid", taskInfo.GetUserId()),
			zap.Stringer("task", taskInfo), zap.Uint64("shardingId", shardingID),
		)
		res.ErrInfo = errors.ARGS_NOT_MATCH.Wrap("时间还没到")
		return
	}
	if taskInfo.GetType() != uint32(234) {
		logger.ErrorWF("ProcessTimeOut task typ not match", zap.Any("req", req), zap.Uint32("myType", 234))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("task typ not match")
		return
	}
	err = ItemCollectCallback(logger, shardingID, taskInfo.GetContext())
	return
}
