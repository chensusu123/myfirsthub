package process

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/structsdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazecalcattrredis"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeAttrSkillV8Cfg"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeAttributeV8Cfg"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeSkillInfoV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/common/commonmustarriveredis"
	"gitlab.ifreetalk.com/plate/protodef/MazeAIBattle"
	"go.uber.org/zap"
)

func GetUserAttrMap(logger fklog.FKLogI, userId uint64) (map[int32]int64, error) {
	//attrIds := GetAttrIds()
	//skillAttrIds := GetSkillAttrIds()
	//if len(skillAttrIds) > 0 {
	//	attrIds = append(attrIds, skillAttrIds...)
	//}
	attrDbs, err := mazecalcattrredis.GetAllMazeCalcAttr(logger, userId)
	if err != nil {
		logger.WarnWF("GetUserBattleAttr BatchGetDollCalcAttr nil", zap.Uint64("userId", userId))
		return nil, err
	}
	attrMap := make(map[int32]int64, 0)
	for attrId, attrVal := range attrDbs {
		//attrCfg := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(attrId)
		//if attrCfg.Type != 1{
		//	continue
		//}
		attrMap[attrId] = attrVal
	}
	return attrMap, nil
}
func GetUserBattleAttr(logger fklog.FKLogI, userId uint64, skillIds []int32, userAttrMap map[int32]int64) (map[int32]*MazeAIBattle.MazeAIAttrInfo, error) {
	attrTypeMap := GetAttrType()
	attrMap := make(map[int32]*MazeAIBattle.MazeAIAttrInfo, 0)
	for attrId, attrVal := range userAttrMap {
		if attrId <= 0 {
			continue
		}
		attrCfg := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(attrId)
		if attrCfg == nil {
			logger.WarnWF("GetUserBattleAttr GetAttributeConfig error", zap.Uint64("userId", userId), zap.Any("attrId", attrId))
			return nil, errors.New("配置不存在")
		}
		attrMap[attrId] = &MazeAIBattle.MazeAIAttrInfo{
			Type:          proto.Int32(attrTypeMap[attrId]),
			UserValue:     proto.Int32(int32(attrVal)),
			UserValueType: proto.Int32(attrCfg.Figure),
		}
	}
	if len(skillIds) > 0 {
		skillCfg := GMazeSkillInfoV8Cfg.Get(skillIds[0])
		if skillCfg != nil {
			if attrMap[constdef.AtkNumber] == nil {
				attrCfg := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(constdef.AtkNumber)
				if attrCfg == nil {
					logger.WarnWF("GetUserBattleAttr GetAttributeConfig error", zap.Uint64("userId", userId), zap.Any("attrId", constdef.AtkNumber))
					return nil, errors.New("配置不存在")
				}
				attrMap[constdef.AtkNumber] = &MazeAIBattle.MazeAIAttrInfo{
					Type:          proto.Int32(attrTypeMap[constdef.AtkNumber]),
					UserValue:     proto.Int32(0),
					UserValueType: proto.Int32(attrCfg.Figure),
				}
			}
			if attrMap[constdef.AtkDis] == nil {
				attrCfg := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(constdef.AtkDis)
				if attrCfg == nil {
					logger.WarnWF("GetUserBattleAttr GetAttributeConfig error", zap.Uint64("userId", userId), zap.Any("attrId", constdef.AtkDis))
					return nil, errors.New("配置不存在")
				}
				attrMap[constdef.AtkDis] = &MazeAIBattle.MazeAIAttrInfo{
					Type:          proto.Int32(attrTypeMap[constdef.AtkDis]),
					UserValue:     proto.Int32(0),
					UserValueType: proto.Int32(attrCfg.Figure),
				}
			}
			attrMap[constdef.AtkNumber].UserValue = proto.Int32(100000)
			attrMap[constdef.AtkDis].UserValue = proto.Int32(attrMap[constdef.AtkDis].GetUserValue() + skillCfg.Distance_max)
		}
	}
	return attrMap, nil
}

func GetAttrType() map[int32]int32 {
	attrTypeMap := make(map[int32]int32, 0)
	attrTypeMap[constdef.Critical] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_CRITICAL_PER)
	attrTypeMap[constdef.CriticalRatio] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_CRITICAL_RATIO_PER)
	attrTypeMap[constdef.MoveSpeed] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_MOVE_SPEED_PER)
	attrTypeMap[constdef.AtkSpeed] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ATK_SPEED)
	attrTypeMap[constdef.AtkNumber] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ATK_NUMBER)
	attrTypeMap[constdef.AtkDis] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ATK_DIS)
	attrTypeMap[constdef.DollFormulaAttack] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ROLE_ATK_VALUE)
	attrTypeMap[constdef.DollFormulaDefend] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ROLE_DEF_VALUE)
	attrTypeMap[constdef.HitReturnBlood] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_HIT_RETURN_BLOOD)
	attrTypeMap[constdef.HurtTotalValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_HURT_TOTAL_VALUE)
	attrTypeMap[constdef.BeHurtTotalValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_BE_HURT_TOTAL_VALUE)
	attrTypeMap[constdef.DefValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_DEF_VALUE_ADD)
	attrTypeMap[constdef.DefValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_DEF_VALUE_PER)

	attrTypeMap[constdef.AtkValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ATK_VALUE_ADD)
	attrTypeMap[constdef.AtkValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ATK_VALUE_PER)
	attrTypeMap[constdef.AtkHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ATK_HURT_VALUE_PER)
	attrTypeMap[constdef.BeHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_BE_HURT_VALUE_PER)
	attrTypeMap[constdef.ExtraHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.ExtraBeHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_EXTRA_BE_HURT_VALUE_ADD)

	attrTypeMap[constdef.IceTagAttrId] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_TAG_ATTR_ID)
	attrTypeMap[constdef.IceAtkAppendElementHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_ATK_APPEND_ELEMENT_HURT_VALUE)
	attrTypeMap[constdef.IceAtkValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_ATK_VALUE_ADD)
	attrTypeMap[constdef.IceAtkValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_ATK_VALUE_PER)
	attrTypeMap[constdef.IceAtkHurtValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_ATK_HURT_VALUE_PER)
	attrTypeMap[constdef.IceBeHurtValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_BE_HURT_VALUE_PER)
	attrTypeMap[constdef.IceExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.IceBeExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_BE_EXTRA_HURT_VALUE_ADD)

	attrTypeMap[constdef.FireTagAttrId] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_TAG_ATTR_ID)
	attrTypeMap[constdef.FireAtkAppendElementHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_ATK_APPEND_ELEMENT_HURT_VALUE)
	attrTypeMap[constdef.FireAtkValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_ATK_VALUE_ADD)
	attrTypeMap[constdef.FireAtkValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_ATK_VALUE_PER)
	attrTypeMap[constdef.FireAtkHurtValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_ATK_HURT_VALUE_PER)
	attrTypeMap[constdef.FireBeHurtValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_BE_HURT_VALUE_PER)
	attrTypeMap[constdef.FireExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.FireBeExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_BE_EXTRA_HURT_VALUE_ADD)

	attrTypeMap[constdef.ElectricityTagAttrId] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_TAG_ATTR_ID)
	attrTypeMap[constdef.ElectricityAtkAppendElementHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_ATK_APPEND_ELEMENT_HURT_VALUE)
	attrTypeMap[constdef.ElectricityAtkValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_ATK_VALUE_ADD)
	attrTypeMap[constdef.ElectricityAtkValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_ATK_VALUE_PER)
	attrTypeMap[constdef.ElectricityAtkHurtValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_ATK_HURT_VALUE_PER)
	attrTypeMap[constdef.ElectricityBeHurtValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_BE_HURT_VALUE_PER)
	attrTypeMap[constdef.ElectricityExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.ElectricityBeExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_BE_EXTRA_HURT_VALUE_ADD)

	attrTypeMap[constdef.PoisonTagAttrId] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_TAG_ATTR_ID)
	attrTypeMap[constdef.PoisonAtkAppendElementHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_ATK_APPEND_ELEMENT_HURT_VALUE)
	attrTypeMap[constdef.PoisonAtkValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_ATK_VALUE_ADD)
	attrTypeMap[constdef.PoisonAtkValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_ATK_VALUE_PER)
	attrTypeMap[constdef.PoisonAtkHurtValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_ATK_HURT_VALUE_PER)
	attrTypeMap[constdef.PoisonBeHurtValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_BE_HURT_VALUE_PER)
	attrTypeMap[constdef.PoisonExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.PoisonBeExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_BE_EXTRA_HURT_VALUE_ADD)
	return attrTypeMap
}

//func GetAttrIds() []int32 {
//	attrIds := []int32{
//		constdef.DollFormulaAttack,
//		constdef.DollFormulaDefend,
//		constdef.DollFormulaBlood,
//		constdef.MazeAttr10137,
//		constdef.MazeAttr10136,
//		constdef.MazeAttr10140,
//		constdef.MazeAttr10141,
//		constdef.MazeAttr10143,
//		constdef.MazeAttr10152,
//	}
//	return attrIds
//}

func HasBattleAttr(chgAttrs []*structsdef.AttrChgInfo) bool {
	//for _, attr := range chgAttrs {
	//	attrCfg := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(attr.AttrId)
	//	if attrCfg.Type == 1{
	//		return true
	//	}
	//}
	return true
}

func GetSkillAttrIds() []int32 {
	attrIds := make([]int32, 0)
	for _, cfg := range GMazeAttrSkillV8Cfg.GetAll() {
		attrIds = append(attrIds, cfg.Attr_id)
	}
	return attrIds
}

func SendMazeBarrierChgPack(logger fklog.FKLogI, userId uint64, mazeBattleInfo *MazeAIBattle.MazeBarrierInfo) (err error) {
	moneyPack := &MazeAIBattle.MazeBarrierInfoChangeID{
		MazeBarrierInfo: mazeBattleInfo,
	}
	logger.InfoWF("SendMazeBarrierChgPack send client with", zap.Any("moneyPack", moneyPack))
	return commonmustarriveredis.SendArrivePacketFix(userId, 16172, moneyPack)
}

func GetEffectAttrValue(attrValue int32, attrValueVariableId map[int32]int32, userAttrMap map[int32]int64) int64 {
	effectAttrValue := attrValue
	for k, v := range attrValueVariableId {
		if v == 1 {
			effectAttrValue += int32(userAttrMap[k])
		} else if v == 2 {
			effectAttrValue -= int32(userAttrMap[k])
		}
	}
	return int64(effectAttrValue)
}
