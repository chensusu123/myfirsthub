package config_business_t

import (
	"testing"
	"time"

	"go.uber.org/zap"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/usecase/business"
)

func TestLoad(t *testing.T) {
	logger := gTestLogger.Clone("TestLoad")
	logger.DebugWF("TestLoad")
	business.GCustomBusiness.Init(logger)
	time.Sleep(5 * time.Second)
	var index int32 = 3
	barrierCfg := GMazeBarriesV8Cfg.Get(index)
	if barrierCfg == nil {
		logger.ErrorWF("OnMazeBarrierListRQ get barrier cfg fail", zap.Any("barrier", index))

		return
	}
}
