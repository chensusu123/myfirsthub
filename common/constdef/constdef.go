package constdef

const EquipPosNum = 8      // 装备位数量
const CurSuitDef int32 = 1 // 当前套装默认值，默认第一套

// 装备id包chg_type枚举
const (
	EquipChgTypeAdd    int32 = 1 // 装备新增
	EquipChgTypeDel    int32 = 2 // 装备删除
	EquipChgTypeMod    int32 = 4 // 装备修改
	EquipChgTypeBagCap int32 = 8 // 背包上限变化
)

// 告警类型定义
const (
	HttpAlertTypeInitAttrExp int32 = 100 // 初始属性异常告警
)

// 装配激活状态
const (
	FiveElemActivate = 1 // 五行激活标记
	HurtSuitActivate = 2 // 伤害套装激活标记
)

const (
	DollMazeMoneyChgTypeAdd   = 1
	DollMazeMoneyChgTypeBuy   = 2
	DollMazeMoneyChgTypeIdChg = 3
	DollMazeMoneyChgTypeDeath = 4
)

const (
	CollectNotice    = 1 // 道具结算
	CollectEndNotice = 2 // 挂机结束
)

const (
	MazeCommonItemDiamond = 46900001
	MazeCommonItemCoin    = 46200001
	MazeCommonItemExp     = 46500001
)

const (
	ItemProducePercent = 10000 // 道具产出比例
)
