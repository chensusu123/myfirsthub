package demconstdef

const (
	DollEquipAssembleOpMaskDbSave        int32 = 1   // 保存装配数据失败
	DollEquipAssembleOpMaskDbGet         int32 = 2   // 查询装配数据失败
	DollEquipAssembleOpMaskForce         int32 = 4   // 计算武力值失败
	DollEquipAssembleOpMaskForcePreview  int32 = 8   // 武力值预览失败
	DollEquipAssembleOpMaskNotifyBag     int32 = 16  // 通知装备背包
	DollEquipAssembleOpMaskForceNotify   int32 = 32  // 通知武力值计算
	DollEquipAssembleOpMaskCalcBuff      int32 = 64  // 计算装备buff
	DollEquipAssembleOpMaskNonForceBuff  int32 = 128 // 计算非武力值buff
	DollEquipAssembleOpMaskCalcAttr      int32 = 256 // 计算人偶属性通知
	DollEquipAssembleOpMaskHandSkillCalc int32 = 512 // 手势技能生效计算
)

var MySvr = "maze_equip_main_server(19578)"
