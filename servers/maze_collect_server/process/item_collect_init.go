package process

import (
	"time"

	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazecollectredis"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeBarriesOnHookV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/common/MustArriveRedis"
	"gitlab.ifreetalk.com/plate/protodef/MazeCollect"
	"gitlab.ifreetalk.com/plate/protodef/MazeCollectCache"
	"go.uber.org/zap"
)

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
