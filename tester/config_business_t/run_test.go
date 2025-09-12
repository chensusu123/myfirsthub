package config_business_t

import (
	"testing"
	"time"

	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/usecase/business"

	"go.uber.org/zap"
)

func TestLoad(t *testing.T) {
	logger := gTestLogger.Clone("TestLoad")
	logger.CtxDebug(ctx, "TestLoad")
	business.GCustomBusiness.Init(logger)
	time.Sleep(5 * time.Second)
	var index int32 = 3
	barrierCfg := GMazeBarriesV8Cfg.GetWithCtx(ctx, index)
	if barrierCfg == nil {
		logger.CtxError(ctx, "OnMazeBarrierListRQ get barrier cfg fail", zap.Any("barrier", index))

		return
	}
}
