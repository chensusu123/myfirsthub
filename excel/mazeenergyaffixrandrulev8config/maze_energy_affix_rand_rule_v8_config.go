package mazeenergyaffixrandrulev8config

import "maze_game_server/config/GMazeEnergyAffixRandRuleV8Cfg"

/**
 * @Author: liushuhang
 * @Date: 2025/4/10 16:40
 * @Description:
 */

func GetKey(randId int32, level int32) int32 {
	return randId*10000 + level
}

func GetEnergyAffixRandRuleConfig(configId int32) *GMazeEnergyAffixRandRuleV8Cfg.MazeEnergyAffixRandRuleV8ConfigRow {
	config := GMazeEnergyAffixRandRuleV8Cfg.Get(configId)
	if config == nil {
		return nil
	}

	return config
}
