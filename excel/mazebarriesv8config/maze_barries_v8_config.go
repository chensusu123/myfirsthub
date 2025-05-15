package mazebarriesv8config

import "gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeBarriesV8Cfg"

/**
 * @Author: liushuhang
 * @Date: 2025/3/25 20:42
 * @Description:
 */

func GetStageConfig(configId int32) *GMazeBarriesV8Cfg.MazeBarriesV8ConfigRow {
	config := GMazeBarriesV8Cfg.Get(configId)
	if config == nil {
		return nil
	}

	return config

}
