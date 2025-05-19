package config_business_t

import (
	"testing"
	"time"

	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeBarriesV8Cfg"
	"gitlab.ifreetalk.com/maze/maze_game_server/usecase/business"
	"go.uber.org/zap"
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
