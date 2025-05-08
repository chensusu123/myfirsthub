package mazeenergylevelv8config

import "gitlab.ifreetalk.com/plate/excel/auto/GMazeEnergyLevelV8Cfg"

/**
 * @Author: liushuhang
 * @Date: 2025/3/25 14:26
 * @Description:
 */

func GetKey(energyId int32, level int32) int32 {
	return energyId*10000 + level
}

func GetEnergyLevelConfig(configId int32) *GMazeEnergyLevelV8Cfg.MazeEnergyLevelV8ConfigRow {
	config := GMazeEnergyLevelV8Cfg.Get(configId)
	if config == nil {
		return nil
	}

	return config
}
