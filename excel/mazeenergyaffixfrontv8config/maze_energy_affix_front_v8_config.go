package mazeenergyaffixfrontv8config

import "gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeEnergyAffixFrontV8Cfg"

/**
 * @Author: liushuhang
 * @Date: 2025/3/25 11:47
 * @Description:
 */

func GetMazeEnergyAffixFrontConfig(configId int32) *GMazeEnergyAffixFrontV8Cfg.MazeEnergyAffixFrontV8ConfigRow {
	config := GMazeEnergyAffixFrontV8Cfg.Get(configId)
	if config == nil {
		return nil
	}

	return config
}
