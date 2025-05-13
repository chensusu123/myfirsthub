package mazeenergyaffixlibraryv8config

import (
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeEnergyAffixLibraryV8Cfg"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/25 11:26
 * @Description:
 */

func GetEnergyLibraryAffixList(configId int32) ([]int32, []int32) {
	config := GMazeEnergyAffixLibraryV8Cfg.Get(configId)
	if config == nil {
		return nil, nil
	}

	return config.Affix_id_list, config.Certainly_affix_id_list
}
