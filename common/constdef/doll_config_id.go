/*
 * @Author: majian
 * @Date: 2024-04-26 11:12:26
 * @Last Modified by: majian
 * @Last Modified time: 2025-04-01 19:25:40
 */
package constdef

// 所有人偶常量配置ID枚举在这里定义
const (
	DollCfgId3301 int32 = 3301 // 3301	马健	角色创建后默认穿戴的装备id
	DollCfgId4501 int32 = 4501 // 4501	马健	1万武力转化为攻防血时，各自属性的比例19=攻击 21=防御 23=血量， 四舍五入	19:4800_21:400_23:24000
	DollCfgId4502 int32 = 4502 // 4502	马健	属性列表中，特定初始值的属性id				10902:10000_11021:10000
	DollCfgId4503 int32 = 4503 // 面板除了属性列表以外，额外需要关注的属性id

	DollCfgId5001 int32 = 5001 // 5001	王振虎 临时背包过期时间
	DollCfgId6701 int32 = 6701 // 切换实力分析面板版本的副本层级
	DollCfgId6901 int32 = 6901 // 冲锋技能的技能id：技能默认可以累计的次数
)

// 所有人偶装备常量配置ID在这里定义
const (
	DollEquipCfgSuitNum   int32 = 401  // 401	马健	一键换装可切换的套装数量
	DollEquipCfgSwitchCd  int32 = 402  // 402	马健	一键换装的cd时间（秒）
	DollResistanceConvert int32 = 602  // 602	马健	展示属性：全抗【value_int】对应的4种抗性属性id【value_map】
	DollEquipCfg501       int32 = 501  // 501    王振虎 攻击距离属性
	DollEquipCfg502       int32 = 502  // 502    王振虎 攻击人数属性
	DollEquipCfg801       int32 = 801  // 801	马健	装备套装id对应的最大套装部件数量	1:5_2:5_3:5_4:5_5:5
	DollEquipCfg701       int32 = 701  // 701    王振虎 词条roll值随机的万分比对应颜色档位
	DollEquipCfg901       int32 = 901  // 901	马健、贾自杰	属性伤害的属性id对应的客户端用属性伤害的枚举id	10101:1_10102:2_10103:3_10104:4
	DollEquipCfg201       int32 = 201  // 201	王振虎 背包容量上限
	DollEquipCfg1001      int32 = 1001 // 战斗有光效的武器品质
	DollEquipCfg1101      int32 = 1101 // 装备随机属性时，被视为抗性词条的属性id
	DollEquipCfg1102      int32 = 1102 // 装备随机属性时，被视为垃圾词条的属性id
	DollEquipCfg1301      int32 = 1301 // 装备封印中显示的文本
	DollEquipCfg1401      int32 = 1401 // 1401 装备开启技能孔消耗的材料：数量
	DollEquipCfg1         int32 = 1    // 装备部位之间的激活关系A：B表示，B激活A
)

// 人偶初始属性定义
const (
	DollInitAttrCfgId  int32 = 1 // 人偶初始属性列表
	DollInitAttrCfgId2 int32 = 2 // 马健	属性列表特定属性id初始展示属性
	DollInitAttrCfgId3 int32 = 3 // 马健	统计属性时需要排除掉初始值的属性id
)

// 迷宫初始属性定义
const (
	MazeInitAttrCfgId int32 = 1 // 迷宫初始属性列表
)

// 展示道具id
const (
	GoldPileItemCfgId            = 46200002 // 金币堆道具id
	StrengthenStonePileItemCfgId = 46700002 // 强化石堆道具id
)

// 展示道具与实际道具对应配置id
const (
	GoldPileShow2RealCfgId            int32 = 901
	StrengthenStonePileShow2RealCfgId int32 = 902
)
