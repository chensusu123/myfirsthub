// @Author pangchenyang 2025/6/12 17:04:00
// @Desc: 
package redisconfig

import (
	cfg2 "maze_game_server/lib/nano/cfg"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 新增全局变量
var globalRedisService *RedisService

func InitGlobalRedis(logger fklog.FKLogI, config cfg2.CfgSvr) error {
	globalRedisService = NewRedisService(config)
	if err := globalRedisService.OnInit(logger, nil); err != nil {
		logger.ErrorWF("InitGlobalRedis OnInit error.", zap.Error(err))
		return err
	}
	if err := globalRedisService.OnStart(logger, nil); err != nil {
		logger.ErrorWF("InitGlobalRedis OnStart error.", zap.Error(err))
		return err
	}
	// todo stop
	return nil
}
