package vardef

import "maze_game_server/common/constdef"

var (
	MazeBuffChgTypeDesc = map[int32]string{
		constdef.MazeBuffChgTypeEquipInit:  "装备初始化",
		constdef.MazeBuffChgTypeEquipDress: "装备替换",
		constdef.MazeBuffChgTypeEquipGm:    "装备gm",
		constdef.MazeBuffEquipFix:          "装备修复",
		constdef.MazeBuffInitAttr:          "属性初始化",
		constdef.MazeBuffLvChg:             "等级变化",
		constdef.MazeBuffChgTypeGm:         "GM修复",
		constdef.MazeBuffCenter:            "buff中心变化",
		constdef.MazeBuffEquipPosUpgrade:   "装备位强化",
	}
)

var (
	DollEquipQualityMap = map[int32]string{
		constdef.EquipQualityWhite:  "白(普通)",
		constdef.EquipQualityBlue:   "蓝(精良)",
		constdef.EquipQualityPurple: "紫(史诗)",
		constdef.EquipQualityMyth:   "橙(神话)",
		constdef.EquipQualityLegend: "橙(传奇)",
		constdef.EquipQualityOrange: "橙(套装)",
	}
)

var (
	MazeBuffSrcDescMap = map[int32]string{
		constdef.MazeBuffSrcEquip:           "迷宫装备",
		constdef.MazeBuffSrcLv:              "迷宫等级",
		constdef.MazeBuffSrcMonthCard:       "迷宫月卡",
		constdef.MazeBuffSrcOldBC:           "旧buff中心",
		constdef.MazeBuffSrcInit:            "初始化属性",
		constdef.MazeBuffSrcEquipPos:        "装备位强化",
		constdef.MazeBuffSrcSelectBuffForce: "三选一buff加武力",
	}

	EquipSuitAttrMap = map[int32]string{
		2: "火属性",
		1: "冰属性",
		4: "电属性",
		3: "毒属性",
	}
)

var SrcNameMap = map[int]string{
	constdef.MazePartOldBc:    "旧buff中心",
	constdef.MazePartBc:       "迷宫buff中心",
	constdef.MazePartInitAttr: "迷宫初始属性",
}

var ForceSrcMap = map[int32]string{
	constdef.ForcePreviewEquip:         "穿戴装备", // 装备预览类型
	constdef.KnifeForce:                "暗器",     // 暗器预览类型
	constdef.HaveMountForce:            "坐骑永久", // 永久坐骑加武力值
	constdef.InUseMountForce:           "坐骑出战", // 出战坐骑加武力值
	constdef.EquipPosForce:             "装备位",   // 装备位加武力值
	constdef.EquipMagicStrengthenForce: "法宝强化", // 法宝强化加武力值
	constdef.EquipMagicStageForce:      "法宝进阶", // 法宝升阶加武力值
	constdef.EquipMagicEnchantForce:    "法宝附魔", // 法宝附魔加武力值
}

var BuffCenterSrcMap = map[int32]string{
	801: "时装buff 王君凯",
	802: "随从buff 辛宇",
}

var FiveElemMap = map[int32]string{
	1: "金",
	2: "木",
	3: "水",
	4: "火",
	5: "土",
}
