package config_business_t

import (
	"testing"

	"gitlab.ifreetalk.com/maze/maze_game_server/usecase/business"
)

func TestLoad(t *testing.T) {
	logger := gTestLogger.Clone("TestLoad")
	logger.DebugWF("TestLoad")
	business.GCustomBusiness.OnInit(logger, nil)
}
