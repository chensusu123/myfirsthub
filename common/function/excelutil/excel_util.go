package excelutil

func GetMagicWeaponStageCfgOrder(magicWeaponId, stage int32) int32 {
	return magicWeaponId%100000*10000 + stage // 序号是法宝id的后5位*10000+等级
}

func GetCloakStageCfgOrder(cloakId, stage int32) int32 {
	return cloakId%100000*10000 + stage // 序号是披风id的后5位*10000+等级
}

// 装备位强化等级配表 主键
func GetEquipPosEnLevelKey(pos, level int32) int32 {
	return pos*10000 + level
}

func GetFaBaoStageCfgOrder(stage int32) int32 {
	return 10000 + stage
}

// 获取法宝强化key
func GetEquipMagicStrengthenKey(level int32) int32 {
	return 10000 + level
}

// 获取地下城副本主键
func GetDXCKey(chapter, level int32) int32 {
	return chapter*10000 + level
}

// 获取手势技能key
func GetHandSkillKey(skillId, stage, lv int32) int32 {
	return skillId*10000 + stage*100 + lv
}
