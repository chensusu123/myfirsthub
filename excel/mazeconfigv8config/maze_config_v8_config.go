package mazeconfigv8config

import (
	"context"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"maze_game_server/config/GMazeConfigV8Cfg"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/21 19:47
 * @Description:
 */

func GetMazeConfig(configId int32) map[int32]int64 {
	row := GMazeConfigV8Cfg.Get(configId)
	if row == nil {
		return nil
	}

	return row.Value_map
}

func GetBuffSelectCount() int64 {
	//config := GMazeConfigV8Cfg.Get(999)
	//if config == nil {
	//	return 3
	//}
	//
	//return config.Value_int
	return 3
}

func GetBuffSelectTime(ctx context.Context) int32 {
	config := GMazeConfigV8Cfg.GetWithCtx(ctx, 401, config_manager.QueryNullable())
	if config == nil {
		return 60
	}

	return int32(config.Value_int)
}
