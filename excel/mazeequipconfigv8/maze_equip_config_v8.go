/*
 * @Author: majian
 * @Date: 2025-03-13 16:40:36
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-13 16:48:49
 */
package mazeequipconfigv8

import (
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeEquipConfigV8Cfg"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
)

// 获取最大装备套数量
func GetMaxEquipSuitNum() int32 {
	rowCfg := GMazeEquipConfigV8Cfg.GetMazeEquipConfigV8Config(constdef.DollEquipCfgSuitNum)
	if rowCfg != nil {
		return int32(rowCfg.Value_int)
	}
	return 1
}

// 获取切换装备cd
func GetSwitchSuitCd() int32 {
	rowCfg := GMazeEquipConfigV8Cfg.GetMazeEquipConfigV8Config(constdef.DollEquipCfgSwitchCd)
	if rowCfg != nil {
		return int32(rowCfg.Value_int)
	}
	return 15
}

// 获取装备套装配置
func GetEquipSuitCfg() map[int32]int32 {
	rowCfg := GMazeEquipConfigV8Cfg.GetMazeEquipConfigV8Config(constdef.DollEquipCfg801)
	if rowCfg != nil {
		cfgMap := make(map[int32]int32)
		for k, v := range rowCfg.Value_map {
			cfgMap[k] = int32(v)
		}
		return cfgMap
	}
	return nil
}

// 获取套装效果的类型
func GetSuitEffectType(attrId int32) int32 {
	rowCfg := GMazeEquipConfigV8Cfg.GetMazeEquipConfigV8Config(constdef.DollEquipCfg901)
	if rowCfg != nil {
		return int32(rowCfg.Value_map[attrId])

	}
	return 0
}

// 获取背包上限
func GetEquipBagLimit() int {
	equipConfig := GMazeEquipConfigV8Cfg.GetMazeEquipConfigV8Config(constdef.DollEquipCfg201)
	if equipConfig != nil {
		return int(equipConfig.Value_int)
	}
	return 50
}

// 是否是特效品质
func IsTxQuality(quality int32) bool {
	equipConfig := GMazeEquipConfigV8Cfg.GetMazeEquipConfigV8Config(constdef.DollEquipCfg1001)
	if equipConfig != nil {
		if _, ok := equipConfig.Value_map[quality]; ok {
			return true
		}
	}
	return false
}

// 需要展示武器效果的套装数量
func IsShowSuitEffect(suitCnt int32) bool {
	equipConfig := GMazeEquipConfigV8Cfg.GetMazeEquipConfigV8Config(constdef.DollEquipCfg1001)
	if equipConfig != nil {
		for k := range equipConfig.Value_map {
			if suitCnt >= k {
				return true
			}
		}
	}
	return false
}

// 抗性词条属性列表
func GetEquipGoodAttrIds() []int64 {
	equipConfig := GMazeEquipConfigV8Cfg.GetMazeEquipConfigV8Config(constdef.DollEquipCfg1101)
	if equipConfig != nil {
		return equipConfig.Value_list
	}
	return nil
}

// 垃圾词条属性列表
func GetEquipBadAttrIds() []int64 {
	equipConfig := GMazeEquipConfigV8Cfg.GetMazeEquipConfigV8Config(constdef.DollEquipCfg1102)
	if equipConfig != nil {
		return equipConfig.Value_list
	}
	return nil
}

// 装备封印描述
func GetEquipSealDesc() string {
	equipConfig := GMazeEquipConfigV8Cfg.GetMazeEquipConfigV8Config(constdef.DollEquipCfg1301)
	if equipConfig != nil {
		return equipConfig.Value_str
	}
	return ""
}

// 获取装备打孔消耗配置
func GetSkillSlotCost() map[int32]int64 {
	rowCfg := GMazeEquipConfigV8Cfg.GetMazeEquipConfigV8Config(constdef.DollEquipCfg1401)
	if rowCfg != nil {
		cfgMap := make(map[int32]int64)
		for k, v := range rowCfg.Value_map {
			cfgMap[k] = v
		}
		return cfgMap
	}
	return nil
}

// 指定位置获取五行升华激活部位依赖关系
func GetFiveElemActiveLinkRs(pos int32) int32 {
	rowCfg := GMazeEquipConfigV8Cfg.GetMazeEquipConfigV8Config(constdef.DollEquipCfg1)
	if rowCfg != nil {
		return int32(rowCfg.Value_map[pos])
	}
	return 0
}
