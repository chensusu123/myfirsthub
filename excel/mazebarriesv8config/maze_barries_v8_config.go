package mazebarriesv8config

import (
	"context"
	"maze_game_server/config/GMazeBarriesV8Cfg"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/25 20:42
 * @Description:
 */

func GetStageConfig(ctx context.Context, configId int32) *GMazeBarriesV8Cfg.MazeBarriesV8ConfigRow {
	config := GMazeBarriesV8Cfg.GetWithCtx(ctx, configId)
	if config == nil {
		return nil
	}

	return config

}
