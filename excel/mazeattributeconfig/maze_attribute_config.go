package mazeattributeconfig

import (
	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeAttributeV8Cfg"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/22 10:51
 * @Description:
 */

func GetMazeAttributeConfig(configId int32) *GMazeAttributeV8Cfg.MazeAttributeV8ConfigRow {
	config := GMazeAttributeV8Cfg.Get(configId)
	if config == nil {
		return nil
	}

	return config
}
