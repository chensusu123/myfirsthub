package game

import (
	"context"
	"fmt"

	"maze_game_server/common/constdef"
	"maze_game_server/common/structsdef"
	"maze_game_server/config/GMazeAttrSkillV8Cfg"
	"maze_game_server/config/GMazeAttributeV8Cfg"
	"maze_game_server/io/redis/mazecalcattrredis"
	"maze_game_server/pb/common/MazeAIBattle"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func GetUserAttrMap(ctx context.Context, userId uint64) (map[int32]int64, error) {
	//attrIds := GetAttrIds()
	//skillAttrIds := GetSkillAttrIds()
	//if len(skillAttrIds) > 0 {
	//	attrIds = append(attrIds, skillAttrIds...)
	//}
	logger := fklog.ContextAppLogger(ctx)
	attrDbs, err := mazecalcattrredis.GetAllMazeCalcAttr(ctx, userId)
	if err != nil {
		logger.CtxWarn(ctx, "GetUserBattleAttr BatchGetDollCalcAttr nil", zap.Uint64("userId", userId))
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

func GetUserBattleAttr(ctx context.Context, userId uint64, userAttrMap map[int32]int64) (map[int32]*MazeAIBattle.MazeAIAttrInfo, error) {
	// attrTypeMap := GetAttrType()
	logger := fklog.ContextAppLogger(ctx)
	attrMap := make(map[int32]*MazeAIBattle.MazeAIAttrInfo, 0)
	for attrId, attrVal := range userAttrMap {
		if attrId <= 0 {
			continue
		}
		attrCfg := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(attrId)
		if attrCfg == nil {
			logger.CtxWarn(ctx, "GetUserBattleAttr GetAttributeConfig error", zap.Uint64("userId", userId), zap.Any("attrId", attrId))
			return nil, fmt.Errorf("属性配置不存在: %d", attrId)
		}
		attrMap[attrId] = &MazeAIBattle.MazeAIAttrInfo{
			Type:          proto.Int32(attrId),
			UserValue:     proto.Int32(int32(attrVal)),
			UserValueType: proto.Int32(attrCfg.Figure),
		}
	}
	//if len(skillIds) > 0 {
	//	skillCfg := GMazeSkillInfoV8Cfg.Get(skillIds[0])
	//	if skillCfg != nil {
	//		//if attrMap[constdef.AtkNumber] == nil {
	//		//	attrCfg := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(constdef.AtkNumber)
	//		//	if attrCfg == nil {
	//		//		logger.CtxWarn(ctx,"GetUserBattleAttr GetAttributeConfig error", zap.Uint64("userId", userId), zap.Any("attrId", constdef.AtkNumber))
	//		//		return nil, errors.New("配置不存在")
	//		//	}
	//		//	attrMap[constdef.AtkNumber] = &MazeAIBattle.MazeAIAttrInfo{
	//		//		Type:          proto.Int32(attrTypeMap[constdef.AtkNumber]),
	//		//		UserValue:     proto.Int32(0),
	//		//		UserValueType: proto.Int32(attrCfg.Figure),
	//		//	}
	//		//}
	//		if attrMap[constdef.AtkDis] == nil {
	//			attrCfg := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(constdef.AtkDis)
	//			if attrCfg == nil {
	//				logger.CtxWarn(ctx,"GetUserBattleAttr GetAttributeConfig error", zap.Uint64("userId", userId), zap.Any("attrId", constdef.AtkDis))
	//				return nil, errors.New("配置不存在")
	//			}
	//			attrMap[constdef.AtkDis] = &MazeAIBattle.MazeAIAttrInfo{
	//				Type:          proto.Int32(attrTypeMap[constdef.AtkDis]),
	//				UserValue:     proto.Int32(0),
	//				UserValueType: proto.Int32(attrCfg.Figure),
	//			}
	//		}
	//		//attrMap[constdef.AtkNumber].UserValue = proto.Int32(attrMap[constdef.AtkNumber].GetUserValue())
	//		attrMap[constdef.AtkDis].UserValue = proto.Int32(attrMap[constdef.AtkDis].GetUserValue() + skillCfg.Distance_max)
	//	}
	//}
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
	attrTypeMap[constdef.ContinuousDamageHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_CONTINUOUS_DAMAGE_HURT_VALUE)
	attrTypeMap[constdef.ContinuousDamageBeHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_CONTINUOUS_DAMAGE_BE_HURT_VALUE)
	attrTypeMap[constdef.ContinuousDamageExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_CONTINUOUS_DAMAGE_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.ContinuousDamageExtraBeHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_CONTINUOUS_DAMAGE_EXTRA_BE_HURT_VALUE_ADD)
	attrTypeMap[constdef.LifeStealEffectivenessRatioPer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_LIFE_STEAL_EFFECTIVENESS_RATIO_PER)
	attrTypeMap[constdef.FlatLifeRestore] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FLAT_LIFE_RESTORE)
	attrTypeMap[constdef.RestorePercentageOfMaxHpPer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_RESTORE_PERCENTAGE_OF_MAX_HP_PER)
	attrTypeMap[constdef.KillHealBonusRatioPer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_KILL_HEAL_BONUS_RATIO_PER)
	attrTypeMap[constdef.BasicAttackHitRestoreFlatHp] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_BASIC_ATTACK_HIT_RESTORE_FLAT_HP)
	attrTypeMap[constdef.BasicAttackHitRestorePercentMaxHpPer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_BASIC_ATTACK_HIT_RESTORE_PERCENT_MAX_HP_PER)
	attrTypeMap[constdef.BasicAttackLifeStealBonusRatioPer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_BASIC_ATTACK_LIFE_STEAL_BONUS_RATIO_PER)
	attrTypeMap[constdef.BasicAttackDamage] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_BASIC_ATTACK_DAMAGE)

	attrTypeMap[constdef.IceTagAttrId] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_TAG_ATTR_ID)
	attrTypeMap[constdef.IceAtkAppendElementHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_ATK_APPEND_ELEMENT_HURT_VALUE)
	attrTypeMap[constdef.IceAtkValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_ATK_VALUE_ADD)
	attrTypeMap[constdef.IceAtkValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_ATK_VALUE_PER)
	attrTypeMap[constdef.IceAtkHurtValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_ATK_HURT_VALUE_PER)
	attrTypeMap[constdef.IceBeHurtValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_BE_HURT_VALUE_PER)
	attrTypeMap[constdef.IceExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.IceBeExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_BE_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.IceContinuousDamageHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_CONTINUOUS_DAMAGE_HURT_VALUE)
	attrTypeMap[constdef.IceContinuousDamageBeHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_CONTINUOUS_DAMAGE_BE_HURT_VALUE)
	attrTypeMap[constdef.IceContinuousDamageExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_CONTINUOUS_DAMAGE_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.IceContinuousDamageExtraBeHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_CONTINUOUS_DAMAGE_EXTRA_BE_HURT_VALUE_ADD)
	attrTypeMap[constdef.IceProjectileCount] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_PROJECTILE_COUNT)
	attrTypeMap[constdef.IceProjectileDamage] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_PROJECTILE_DAMAGE)
	attrTypeMap[constdef.IceProjectileAoe] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_PROJECTILE_AOE)
	attrTypeMap[constdef.AdditionalIceProjectilesCount] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ADDITIONAL_ICE_PROJECTILES_COUNT)
	attrTypeMap[constdef.DelayTimeForAdditionalIceProjectilesMs] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_DELAY_TIME_FOR_ADDITIONAL_ICE_PROJECTILES_MS)
	attrTypeMap[constdef.AdditionalIceProjectileDamageModifier] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ADDITIONAL_ICE_PROJECTILE_DAMAGE_MODIFIER)
	attrTypeMap[constdef.IceExplosionDamageModifierOnHit] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_EXPLOSION_DAMAGE_MODIFIER_ON_HIT)
	attrTypeMap[constdef.FreezeChanceOnIceProjectileHitAccuracy] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FREEZE_CHANCE_ON_ICE_PROJECTILE_HIT_ACCURACY)
	attrTypeMap[constdef.FreezeDurationOnIceProjectileHitDuration] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FREEZE_DURATION_ON_ICE_PROJECTILE_HIT_DURATION)
	attrTypeMap[constdef.IceProjectileCooldown] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_PROJECTILE_COOLDOWN)

	attrTypeMap[constdef.FireTagAttrId] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_TAG_ATTR_ID)
	attrTypeMap[constdef.FireAtkAppendElementHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_ATK_APPEND_ELEMENT_HURT_VALUE)
	attrTypeMap[constdef.FireAtkValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_ATK_VALUE_ADD)
	attrTypeMap[constdef.FireAtkValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_ATK_VALUE_PER)
	attrTypeMap[constdef.FireAtkHurtValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_ATK_HURT_VALUE_PER)
	attrTypeMap[constdef.FireBeHurtValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_BE_HURT_VALUE_PER)
	attrTypeMap[constdef.FireExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.FireBeExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_BE_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.FireContinuousDamageHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_CONTINUOUS_DAMAGE_HURT_VALUE)
	attrTypeMap[constdef.FireContinuousDamageBeHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_CONTINUOUS_DAMAGE_BE_HURT_VALUE)
	attrTypeMap[constdef.FireContinuousDamageExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_CONTINUOUS_DAMAGE_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.FireContinuousDamageExtraBeHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_CONTINUOUS_DAMAGE_EXTRA_BE_HURT_VALUE_ADD)
	attrTypeMap[constdef.MeteorExtraAoe] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_METEOR_EXTRA_AOE)
	attrTypeMap[constdef.MeteorDamage] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_METEOR_DAMAGE)
	attrTypeMap[constdef.MeteorCount] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_METEOR_COUNT)
	attrTypeMap[constdef.ShockwaveAoePostMeteorImpact] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_SHOCKWAVE_AOE_POST_METEOR_IMPACT)
	attrTypeMap[constdef.BurnDamageRatioPer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_BURN_DAMAGE_RATIO_PER)
	attrTypeMap[constdef.BurnFlatDamage] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_BURN_FLAT_DAMAGE)
	attrTypeMap[constdef.BurnDurationParam5] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_BURN_DURATION_PARAM5)
	attrTypeMap[constdef.BurnAccuracyParam6] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_BURN_ACCURACY_PARAM6)
	attrTypeMap[constdef.BurnDamageTickInterval] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_BURN_DAMAGE_TICK_INTERVAL)
	attrTypeMap[constdef.MeteorCooldown] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_METEOR_COOLDOWN)

	attrTypeMap[constdef.ElectricityTagAttrId] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_TAG_ATTR_ID)
	attrTypeMap[constdef.ElectricityAtkAppendElementHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_ATK_APPEND_ELEMENT_HURT_VALUE)
	attrTypeMap[constdef.ElectricityAtkValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_ATK_VALUE_ADD)
	attrTypeMap[constdef.ElectricityAtkValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_ATK_VALUE_PER)
	attrTypeMap[constdef.ElectricityAtkHurtValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_ATK_HURT_VALUE_PER)
	attrTypeMap[constdef.ElectricityBeHurtValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_BE_HURT_VALUE_PER)
	attrTypeMap[constdef.ElectricityExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.ElectricityBeExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_BE_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.ElectricityContinuousDamageHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_CONTINUOUS_DAMAGE_HURT_VALUE)
	attrTypeMap[constdef.ElectricityContinuousDamageBeHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_CONTINUOUS_DAMAGE_BE_HURT_VALUE)
	attrTypeMap[constdef.ElectricityContinuousDamageExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_CONTINUOUS_DAMAGE_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.ElectricityContinuousDamageExtraBeHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_CONTINUOUS_DAMAGE_EXTRA_BE_HURT_VALUE_ADD)
	attrTypeMap[constdef.LightningTargetCount] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_LIGHTNING_TARGET_COUNT)
	attrTypeMap[constdef.LightningDamage] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_LIGHTNING_DAMAGE)
	// attrTypeMap[constdef.ParalysisChanceOnLightningHitAccuracy] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_PARALYSIS_CHANCE_ON_LIGHTNING_HIT_ACCURACY)
	// attrTypeMap[constdef.ParalysisDurationOnLightningHitParam5] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_PARALYSIS_DURATION_ON_LIGHTNING_HIT_DURATION_PARAM5)
	attrTypeMap[constdef.LightningCooldown] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_LIGHTNING_COOLDOWN)

	attrTypeMap[constdef.PoisonTagAttrId] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_TAG_ATTR_ID)
	attrTypeMap[constdef.PoisonAtkAppendElementHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_ATK_APPEND_ELEMENT_HURT_VALUE)
	attrTypeMap[constdef.PoisonAtkValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_ATK_VALUE_ADD)
	attrTypeMap[constdef.PoisonAtkValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_ATK_VALUE_PER)
	attrTypeMap[constdef.PoisonAtkHurtValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_ATK_HURT_VALUE_PER)
	attrTypeMap[constdef.PoisonBeHurtValuePer] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_BE_HURT_VALUE_PER)
	attrTypeMap[constdef.PoisonExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.PoisonBeExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_BE_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.PoisonContinuousDamageHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_CONTINUOUS_DAMAGE_HURT_VALUE)
	attrTypeMap[constdef.PoisonContinuousDamageBeHurtValue] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_CONTINUOUS_DAMAGE_BE_HURT_VALUE)
	attrTypeMap[constdef.PoisonContinuousDamageExtraHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_CONTINUOUS_DAMAGE_EXTRA_HURT_VALUE_ADD)
	attrTypeMap[constdef.PoisonContinuousDamageExtraBeHurtValueAdd] = int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_CONTINUOUS_DAMAGE_EXTRA_BE_HURT_VALUE_ADD)
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

func SendMazeBarrierChgPack(ctx context.Context, userId uint64, mazeBattleInfo *MazeAIBattle.MazeBarrierInfo) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	moneyPack := &MazeAIBattle.MazeBarrierInfoChangeID{
		MazeBarrierInfo: mazeBattleInfo,
	}
	logger.CtxInfo(ctx, "SendMazeBarrierChgPack send client with", zap.Uint64("userId", userId), zap.Any("moneyPack", moneyPack))
	return online.ClusterPush(context.TODO(), uint64(userId), 10485, moneyPack)
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

func GetAttrValue(attrValue int32, attrValueVariableId map[int32]int32, userAttrMap map[int32]int64) int64 {
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

func GetElementAttrValue(attrValue map[int32]int32, attrValueVariableId []int32, userAttrMap map[int32]int64) (attrValues []*MazeAIBattle.MazeAIAttrInfo) {
	_, value := GetSkillAttrValue(attrValue, userAttrMap)
	for _, v := range attrValueVariableId {
		if v == 0 { // 属性ID为0则给默认值
			attrValues = append(attrValues, &MazeAIBattle.MazeAIAttrInfo{
				// Type:      proto.Int32(attrID),
				UserValue: proto.Int32(value),
			})
		} else {
			if userAttrMap[v] <= 0 {
				attrValues = append(attrValues, &MazeAIBattle.MazeAIAttrInfo{
					// Type:      proto.Int32(attrID),
					UserValue: proto.Int32(0),
				})
			} else {
				rate := int32(float64(value) * float64(userAttrMap[v]) / 10000.0) // 原值 x (属性值 / 10000)
				attrValues = append(attrValues, &MazeAIBattle.MazeAIAttrInfo{
					// Type:      proto.Int32(attrID),
					UserValue: proto.Int32(rate),
				})
			}
		}
	}
	return attrValues
}

// FillElementAttrValue 用默认值填充各元素属性值
func FillElementAttrValue(attrValue map[int32]int32, count int, userAttrMap map[int32]int64) (attrValues []*MazeAIBattle.MazeAIAttrInfo) {
	_, value := GetSkillAttrValue(attrValue, userAttrMap)
	for i := 0; i < count; i++ {
		attrValues = append(attrValues, &MazeAIBattle.MazeAIAttrInfo{
			// Type:      proto.Int32(attrID),
			UserValue: proto.Int32(value),
		})
	}
	return
}

func FilterSliceZeroValue[T int | int32 | int64](values []T) []T {
	if len(values) == 1 && values[0] == 0 {
		return make([]T, 0)
	}
	return values
}
