/*
 * @Author: majian
 * @Date: 2025-03-14 20:29:49
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-21 20:50:30
 */
package constdef

// 迷宫buff来源
const (
	MazeBuffSrcEquip     int32 = 1 // 迷宫装备  马健
	MazeBuffSrcLv        int32 = 2 // 迷宫等级  马健
	MazeBuffSrcMonthCard int32 = 3 // 迷宫月卡  王永亮
	MazeBuffSrcOldBC     int32 = 4 // 旧buff中心 王振虎
	MazeBuffSrcInit      int32 = 5 // 初始化属性 马健
	MazeBuffSrcEquipPos  int32 = 6 // 装备位强化 马健
)

// 迷宫计算组成
const (
	MazePartInitAttr = 100 // 迷宫初始属性
	MazePartBc       = 101 // 迷宫buff中心
	MazePartOldBc    = 102 // 旧buff中心
)

const (
	MazeBuffChgTypeEquipInit  = 101 // 装备初始化
	MazeBuffChgTypeEquipDress = 102 // 穿戴装备
	MazeBuffChgTypeEquipGm    = 103 // gm触发
	MazeBuffEquipFix          = 104 // 穿戴装备修复
	MazeBuffEquipPosUpgrade   = 105 // 装备位强化
	MazeBuffInitAttr          = 201 // 初始化属性
	MazeBuffLvChg             = 301 // 迷宫升级
	MazeBuffChgTypeGm         = 401 // gm计算
	MazeBuffCenter            = 501 // buff中心
	MazeBuffChgMonthCard      = 601 // 月卡
)

const (
	MazeCardChgOpenType       = 1 // 月卡开通
	MazeCardChgRenewType      = 2 // 月卡续费
	MazeCardChgExpirationType = 3 // 月卡过期
)

const (
	MazeBuffChgTypeCardOpen       = 1001 // 月卡开通
	MazeBuffChgTypeCardRenew      = 1002 // 月卡续费
	MazeBuffChgTypeCardExpiration = 1003 // 月卡过期
)
