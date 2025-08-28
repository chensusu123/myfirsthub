package card

import (
	"context"
	"time"

	"maze_game_server/io/redis/mazecardlistgroupredis"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig/param"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"go.uber.org/zap"
)

var tickerTime int64 = 1

func init() {
	param.Int64P(&tickerTime, "ticker:time", tickerTime, "定时器执行间隔")
}

func DealExpirationCardProcess(ctx context.Context, index int, logger fklog.FKLogI) error {
	shardingID := appconfig.GlobalConfig().Global.ShardingID
	go func(ctx context.Context, logger fklog.FKLogI) {
		if shardingID != 1 {
			return
		}
		lg := logger.Clone("")
		newTicker := time.NewTicker(time.Second * time.Duration(tickerTime))
		defer newTicker.Stop()
		for {
			select {
			case <-newTicker.C:
				lg.SetLogId(time.Now().UnixMilli())
				dealExpirationCardProcess(ctx)
			}
		}
	}(ctx, logger)
	return nil
}

func dealExpirationCardProcess(ctx context.Context) {
	logger := fklog.ContextAppLogger(ctx)
	// 查询过期的月卡用户
	userList, err := mazecardlistgroupredis.GetMazeCardExpirationList(logger)
	if err != nil {
		logger.ErrorWF("dealExpirationCardProcess GetMazeCardExpirationList failed", zap.Error(err))
		return
	}

	var okCount int32
	for _, userId := range userList {
		err = DeleteMazeCard(ctx, uint64(userId))
		if err != nil {
			logger.ErrorWF("dealExpirationCardProcess DeleteMazeCard failed",
				zap.Int64("userId", userId), zap.Error(err))
			continue
		}

		okCount++
	}

	logger.InfoWF("dealExpirationCardProcess end", zap.Int32("okCount", okCount),
		zap.Int("TotalCount", len(userList)))
}
