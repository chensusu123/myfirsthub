// @Author pangchenyang 2025/5/6 19:15:00
// @Desc: 
package constdef

// 武力值预览类型
const (
	ForcePreviewEquip         int32 = 32 // 装备预览类型
	ForcePreviewMagicWeapon   int32 = 33 // 法宝预览类型 已废弃
	HaveMountForce            int32 = 37 // 永久坐骑加武力值
	InUseMountForce           int32 = 38 // 出战坐骑加武力值
	EquipPosForce             int32 = 40 // 装备位加武力值
	EquipMagicStrengthenForce int32 = 41 // 法宝强化武力值
	EquipMagicStageForce      int32 = 42 // 法宝阶等级武力值
	KnifeForce                int32 = 43 // 暗器加武力
	EquipMagicEnchantForce    int32 = 44 // 法宝附魔加武力值
)

// pk段位变化来源
const (
	PkLevelChgSrcEquipPosStreng int32 = 1 // 装备位强化
	PkLevelChgSrcGm             int32 = 2 // gm
)

// const EquipPosNum = 8      // 装备位数量
// const CurSuitDef int32 = 1 // 当前套装默认值，默认第一套

const (
	// 属性field
	AttrFieldEquip       = ForcePreviewEquip
	AttrFieldMagicWeapon = ForcePreviewMagicWeapon

	KnifeInitStage = 1 // 暗器初始阶等级
	// MagicWeaponItemOpType = 634 // 法宝在物品服务中的业务类型 已废弃
	KnifeItemOpType = 669 // 暗器在物品服务中的业务类型

	DollEquipAutoSale = 644 // 装备自动出售在物品服务中的业务类型

	KnifeLoadRecord    = 1 // 暗器安装流水
	KnifeStageUpRecord = 2 // 暗器升阶流水

	EquipPosStrengItemOpType = 662 // 装备位强化在物品服务中的业务类型
)

// TODO 扣物品类型统一定义这里
const (
	ItemOpTypeEquipMagicStrengthen = 670 // 人偶法宝强化 UN_CGK_COMMON_BILL_TYPE_670
	ItemOpTypeEquipMagicStage      = 671 // 人偶法宝升阶 UN_CGK_COMMON_BILL_TYPE_671
)

// 功能开启枚举定义
const (
	FuncOpenEquip              int32 = 5800 // 装备开启
	FuncOpenEquipMake          int32 = 5801 // 装备锻造开启
	FuncOpenKnife              int32 = 5900 // 暗器开启
	FuncOpenFaBao              int32 = 7000 // 法宝开启
	FuncOpenFaBaoStrenth       int32 = 7001 // 法宝强化开启
	FuncOpenFaBaoStage         int32 = 7002 // 法宝升阶开启
	FuncOpenFaBaoFuMo          int32 = 7003 // 法宝附魔开启
	FuncOpenCloak              int32 = 5900 // 披风开启
	FuncOpenEnchant            int32 = 5802 // 装备附魔
	FuncOpenIntensify          int32 = 5803 // 附魔属性强化
	FuncOpenDevourTree         int32 = 5804 // 天赋树
	FuncOpenMount              int32 = 6700 // 坐骑
	FuncOpenEquipRefresh       int32 = 5806 // 装备洗练开启
	FuncOpenEquipPosStrengthen int32 = 5807 // 装备pvp等级强化功能的开启等级
	FuncEquipSkill             int32 = 7201 // 技能开启
	FuncEquipCharge            int32 = 7301 // 冲锋技能开启
)

// 装备id包chg_type枚举
// const (
// 	EquipChgTypeAdd    int32 = 1 // 装备新增
// 	EquipChgTypeDel    int32 = 2 // 装备删除
// 	EquipChgTypeMod    int32 = 4 // 装备修改
// 	EquipChgTypeBagCap int32 = 8 // 背包上限变化
// )

// 装备初始化状态枚举
const (
	DollEquipInitStateBag   = 1 // 加入背包
	DollEquipInitStateDress = 2 // 已穿戴
	DollEquipInitDoing      = 4 // 初始化进行中
)

const (
	// 装配信息前缀 多个field的字段 以m_开头 比如 装备位有8个  套装的以s_ 开头
	AssemblePrefixEquip                string = "equip_"
	AssemblePrefixInitEquip            string = "init_state"
	AssemblePrefixCurAssembleSuitIndex string = "cur_suit_index"      // 当前套装索引
	AssemblePrefixSwitchSuitTime       string = "switch_suit_time"    // 最近一次切换套装时间               //坐骑
	AssemblePrefixEquipPos             string = "m_equip_pos"         // 装备位前缀
	AssemblePrefixEquipPosEnSuit       string = "s_equip_pos_en_suit" // 装备位强化套装Id
)

// 装备品质定义
const (
	EquipQualityWhite  = 1 // 白(普通)
	EquipQualityBlue   = 3 // 蓝(精良)
	EquipQualityPurple = 5 // 紫(史诗)
	EquipQualityMyth   = 6 // 橙(神话)
	EquipQualityLegend = 7 // 橙(传奇)
	EquipQualityOrange = 8 // 橙(套装)
)

const (
	AttrId10531 int32 = 10531 // 地下城复活时间
)

// 告警类型定义
const (
// HttpAlertTypeInitAttrExp int32 = 100 // 初始属性异常告警
)

// // 装配激活状态
// const (
// 	FiveElemActivate = 1 // 五行激活标记
// 	HurtSuitActivate = 2 // 伤害套装激活标记
// )
//
// const (
// 	DollMazeMoneyChgTypeAdd   = 1
// 	DollMazeMoneyChgTypeBuy   = 2
// 	DollMazeMoneyChgTypeIdChg = 3
// 	DollMazeMoneyChgTypeDeath = 4
// )
//
// const (
// 	CollectNotice    = 1 // 道具结算
// 	CollectEndNotice = 2 // 挂机结束
// )
//
// const (
// 	MazeCommonItemDiamond = 46900001
// 	MazeCommonItemCoin    = 46200001
// 	MazeCommonItemExp     = 46500001
// )
//
// const (
// 	ItemProducePercent = 10000 // 道具产出比例
// )
