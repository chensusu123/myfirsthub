package mazeenergyresetcostv8config

import "gitlab.ifreetalk.com/plate/excel/auto/GMazeEnergyResetCostV8Cfg"

/**
 * @Author: liushuhang
 * @Date: 2025/3/24 15:07
 * @Description:
 */

func GetEnergyResetCostConfig(refreshCount int32) *GMazeEnergyResetCostV8Cfg.MazeEnergyResetCostV8ConfigRow {
	for _, config := range GMazeEnergyResetCostV8Cfg.GetAll() {
		if config.Mun_min <= refreshCount && refreshCount <= config.Mun_max {
			return config
		}
	}

	return nil
}
