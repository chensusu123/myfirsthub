/*
 * @Author: majian
 * @Date: 2024-07-13 15:23:56
 * @Last Modified by: majian
 * @Last Modified time: 2025-04-01 19:48:38
 */
package constdef

const (
	// DollAttrAttack = 19 // 攻击
	// DollAttrDefend = 21 // 防御
	// DollAttrBlood  = 23 // 生命

	DollFormulaAttack = 3010000 // 攻击公式Id
	DollFormulaDefend = 3020000 // 防御公式Id
	DollFormulaBlood  = 3030000 // 生命公式Id
	MazeForce         = 9001    // 迷宫武力值

	MazeMoney          int32 = 10257 // 钱币公式ID
	MazeMoneyBuff10258 int32 = 10258 // 怪物掉落钱币
	MazeMoneyBuff10259 int32 = 10259 // 掉落钱币加成

	MazeExp          int32 = 10260 // 经验公式ID
	MazeExpBuff10261 int32 = 10261 // 怪物掉落经验
	MazeExpBuff10262 int32 = 10262 // 掉落经验加成

	Critical         = 3041301 // 暴击伤害
	CriticalRatio    = 3041201 // 暴击几率
	MoveSpeed        = 3041401 // 移动速度（要和客户端配置的模型数据合并计算）
	AtkSpeed         = 3040901 //攻击速度
	AtkNumber        = 3041501 // 攻击目标数量（要和技能中攻击目标数量合并计算）
	AtkDis           = 3041001 // 攻击距离（要和技能中技能释放最大距离合并计算）
	MazeAttr3000101  = 3000101 // 人物初始韧性上限
	HitReturnBlood   = 3041101 //击中回血
	HurtTotalValue   = 3040701 // 伤害总加成
	BeHurtTotalValue = 3040801 // 受到伤害总加成
	DefValueAdd      = 3040201 // 防御固定值
	DefValuePer      = 3040202 // 防御万分比
	AtkValueAdd      = 3040101 // 攻击值
	AtkValuePer      = 3040102 // 攻击值
	AtkHurtValue     = 3040301 // 攻击伤害
	BeHurtValue      = 3040401 // 受到伤害
	ExtraHurtValue   = 3040501 // 额外伤害
	ExtraBeHurtValue = 3040601 // 受到额外伤害
	ContinuousDamageHurtValue = 3041601 //无元素持续伤害提高，万分比
	ContinuousDamageBeHurtValue = 3041701 //受到无元素持续伤害提高，万分比
	ContinuousDamageExtraHurtValueAdd = 3041801 //额外无元素持续伤害，固定值
	ContinuousDamageExtraBeHurtValueAdd = 3041901 //受到额外无元素持续伤害，固定值

	//冰属性id列表
	IceTagAttrId                 = 3050000 // 冰属性标记Id
	IceAtkAppendElementHurtValue = 3050001 // 攻击附加电元素
	IceAtkValueAdd               = 3050101 // 冰元素攻击固定值
	IceAtkValuePer               = 3050102 // 冰元素攻击万分比
	IceAtkHurtValuePer           = 3050301 // 冰元素伤害万分比
	IceBeHurtValuePer            = 3050401 // 受到冰元素伤害
	IceExtraHurtValueAdd         = 3050501 // 额外冰伤害
	IceBeExtraHurtValueAdd       = 3050601 // 受到额外冰伤害
	IceContinuousDamageHurtValue = 3050701 //冰元素持续伤害提高，万分比
	IceContinuousDamageBeHurtValue = 3050801 //受到冰元素持续伤害提高，万分比
	IceContinuousDamageExtraHurtValueAdd = 3050901 //额外冰元素持续伤害，固定值
	IceContinuousDamageExtraBeHurtValueAdd = 3051001 //受到额外冰元素持续伤害，固定值

	//火属性id列表
	FireTagAttrId                 = 3060000 // 火属性标记Id
	FireAtkAppendElementHurtValue = 3060001 // 攻击附加火元素
	FireAtkValueAdd               = 3060101 // 火元素攻击
	FireAtkValuePer               = 3060102 // 火元素攻击
	FireAtkHurtValuePer           = 3060301 // 火元素伤害
	FireBeHurtValuePer            = 3060401 // 受到火元素伤害
	FireExtraHurtValueAdd         = 3060501 // 额外火伤害
	FireBeExtraHurtValueAdd       = 3060601 // 受到额外火伤害
	FireContinuousDamageHurtValue = 3060701 //火元素持续伤害提高，万分比
	FireContinuousDamageBeHurtValue = 3060801 //受到火元素持续伤害提高，万分比
	FireContinuousDamageExtraHurtValueAdd = 3060901 //额外火元素持续伤害，固定值
	FireContinuousDamageExtraBeHurtValueAdd = 3061001 //受到额外火元素持续伤害，固定值

	//电属性id列表
	ElectricityTagAttrId                 = 3070000 // 电属性标记Id
	ElectricityAtkAppendElementHurtValue = 3070001 // 攻击附加电元素
	ElectricityAtkValueAdd               = 3070101 // 电元素攻击
	ElectricityAtkValuePer               = 3070102 // 电元素攻击
	ElectricityAtkHurtValuePer           = 3070301 // 电元素伤害
	ElectricityBeHurtValuePer            = 3070401 // 受到电元素伤害
	ElectricityExtraHurtValueAdd         = 3070501 // 额外电伤害
	ElectricityBeExtraHurtValueAdd       = 3070601 // 受到额外电伤害
	ElectricityContinuousDamageHurtValue = 3070701 //电元素持续伤害提高，万分比
	ElectricityContinuousDamageBeHurtValue = 3070801 //受到电元素持续伤害提高，万分比
	ElectricityContinuousDamageExtraHurtValueAdd = 3070901 //额外电元素持续伤害，固定值
	ElectricityContinuousDamageExtraBeHurtValueAdd = 3071001 //受到额外电元素持续伤害，固定值

	//毒属性id列表
	PoisonTagAttrId                 = 3080000 // 毒属性标记Id
	PoisonAtkAppendElementHurtValue = 3080001 // 攻击附加毒元素
	PoisonAtkValueAdd               = 3080101 // 毒元素攻击
	PoisonAtkValuePer               = 3080102 // 毒元素攻击
	PoisonAtkHurtValuePer           = 3080301 // 毒元素伤害
	PoisonBeHurtValuePer            = 3080401 // 受到毒元素伤害
	PoisonExtraHurtValueAdd         = 3080501 // 额外毒伤害
	PoisonBeExtraHurtValueAdd       = 3080601 // 受到额外毒伤害
	PoisonContinuousDamageHurtValue = 3080701 //毒元素持续伤害提高，万分比
	PoisonContinuousDamageBeHurtValue = 3080801 //受到毒元素持续伤害提高，万分比
	PoisonContinuousDamageExtraHurtValueAdd = 3080901 //额外毒元素持续伤害，固定值
	PoisonContinuousDamageExtraBeHurtValueAdd = 3081001 //受到额外毒元素持续伤害，固定值
)
