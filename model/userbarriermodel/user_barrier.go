package userbarriermodel

import (
	"maze_game_server/io/redis/mazeuserbarrierredis"
	"maze_game_server/pb/server/MazeBarrierCache"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type UserBarrierModel struct {
	UserID        uint64 // 用户ID
	BarrierID     int32  `json:"barrier_id,omitempty"`     //关卡id
	BarrierStatus int32  `json:"barrier_status,omitempty"` //关卡状态 1-未进入 2-正在闯关 3-已通关
	RebornCount   int32  `json:"reborn_count,omitempty"`   //复活计数
	StartTime     int64  `json:"start_time,omitempty"`     //进入时间
	EndTime       int64  `json:"end_time,omitempty"`       //本次探险有结果的时间
}

func NewUserBarrierModel(logger fklog.FKLogI, userID uint64) (ub *UserBarrierModel, err error) {
	ub = &UserBarrierModel{
		UserID: userID,
	}
	if err = ub.load(logger); err != nil {
		return nil, err
	}
	return ub, nil
}

// load
func (ub *UserBarrierModel) load(logger fklog.FKLogI) (err error) {
	data, err := mazeuserbarrierredis.GetUserBarrierInfo(logger, ub.UserID, 0)
	if err != nil {
		logger.ErrorWF("load GetUserBarrierInfo fail", zap.Error(err), zap.Uint64("UserID", ub.UserID))
		return err
	}
	ub.BarrierID = data.GetBarrierId()
	ub.BarrierStatus = data.GetBarrierStatus()
	ub.RebornCount = data.GetRebornCount()
	ub.StartTime = data.GetStartTime()
	ub.EndTime = data.GetEndTime()
	return
}

// Save
func (ub *UserBarrierModel) Save(logger fklog.FKLogI) (err error) {
	data := &MazeBarrierCache.MazeBarrierCache{
		BarrierId:     proto.Int32(ub.BarrierID),
		BarrierStatus: proto.Int32(ub.BarrierStatus),
		RebornCount:   proto.Int32(ub.RebornCount),
		StartTime:     proto.Int64(ub.StartTime),
		EndTime:       proto.Int64(ub.EndTime),
	}
	err = mazeuserbarrierredis.SetUserBarrierInfo(logger, ub.UserID, 0, data)
	if err != nil {
		logger.ErrorWF("Save SetUserBarrierInfo fail", zap.Error(err), zap.Uint64("UserID", ub.UserID), zap.Any("data", data))
	}
	return
}
