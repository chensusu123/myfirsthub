package mazecollect

import (
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeBarriesOnHookV8Cfg"
	"maze_game_server/io/redis/mazecollectredis"
	"maze_game_server/pb/server/MazeCollectCache"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type CollectInfo struct {
	UserId uint64
	FKLogI fklog.FKLogI
}

func NewCollectInfo(logger fklog.FKLogI, userId uint64) *CollectInfo {
	ret := &CollectInfo{
		UserId: userId,
		FKLogI: logger,
	}
	return ret
}

// 获取挂机信息
func (c *CollectInfo) GetCollectInfo() *MazeCollectCache.MazeCollectInfo {
	// 是否已经初始化
	collectInfo, err := mazecollectredis.GetCollectInfo(c.FKLogI, c.UserId)
	if err != nil {
		c.FKLogI.ErrorWF("GetCollectInfo error",
			zap.Any("userId", c.UserId),
			zap.Error(err))
		return nil
	}
	return collectInfo
}

// 初始化迷宫挂机
func (c *CollectInfo) NewMazeCollectInfo(barrierId int32) (err error, info *MazeCollectCache.MazeCollectInfo) {
	cfg := GMazeBarriesOnHookV8Cfg.Get(barrierId)
	if cfg == nil {
		c.FKLogI.ErrorWF("InitMazeCollectLand error",
			zap.Any("barrierId", barrierId))
		return errors.New("配置不存在"), nil
	}

	// 是否已经初始化
	collectInfo := c.GetCollectInfo()
	if collectInfo != nil {
		return nil, collectInfo
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
	c.FKLogI.InfoWF("InitMazeCollectLand init pet", zap.Any("collectInfo", collectInfo))

	err = mazecollectredis.SetCollectInfo(c.FKLogI, c.UserId, collectInfo)
	if err != nil {
		c.FKLogI.ErrorWF("InitMazeCollectLand SetCollectInfo error", zap.Any("collectInfo", collectInfo), zap.Error(err))
		return
	}

	return nil, collectInfo
}
