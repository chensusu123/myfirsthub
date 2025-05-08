package mazeenergyaffixlvv8config

import (
	"sync/atomic"
	"unsafe"

	"gitlab.ifreetalk.com/plate/excel/auto/GMazeEnergyAffixV8Cfg"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/22 15:21
 * @Description:
 */

var gConfigData *LordMazeEnergyAffixLvV8Config

type LordMazeEnergyAffixLvV8Config struct {
	ConfigRowMap map[int32][]int32
}

func init() {
	GMazeEnergyAffixV8Cfg.RegisterMazeEnergyAffixV8InitCallBack("GLordMazeEnergyAffixLvConfig", LordMazeEnergyAffixLvConfig)
}

func LordMazeEnergyAffixLvConfig(cfg *GMazeEnergyAffixV8Cfg.MazeEnergyAffixV8Config) {
	g := new(LordMazeEnergyAffixLvV8Config)
	affixMap := make(map[int32][]int32, 20)
	for _, row := range cfg.ConfigRows {
		affixMap[row.Affix_group_id] = append(affixMap[row.Affix_group_id], row.Order)
	}

	g.ConfigRowMap = affixMap
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(g))
}

func GetAffixConfig(configId int32) *GMazeEnergyAffixV8Cfg.MazeEnergyAffixV8ConfigRow {
	config := GMazeEnergyAffixV8Cfg.Get(configId)
	if config == nil {
		return nil
	}

	return config
}
