package mazeenergyaffixfrontv8config

import "maze_game_server/config/GMazeEnergyAffixFrontV8Cfg"

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
