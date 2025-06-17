package equip

import (
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeActInfoV8Cfg"
	"maze_game_server/config/GMazeAttrSkillV8Cfg"
	"maze_game_server/config/GMazeSkillActV8Cfg"
	"maze_game_server/config/GMazeSkillAutoConditionV8Cfg"
	"maze_game_server/config/GMazeSkillAutoReleaseV8Cfg"
	"maze_game_server/config/GMazeSkillInfoV8Cfg"
	"maze_game_server/config/GMazeSkilleffectV8Cfg"
	"maze_game_server/io/redis/mazebarriertempbuffredis"
	"maze_game_server/io/redis/mazecalcattrredis"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeAIBattle"
	"maze_game_server/pb/server/MazeEquipCache"
	"regexp"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func GetUserSkillTotalInfo(logger fklog.FKLogI, userAttrMap map[int32]int64, baseAttrs []*MazeEquipCache.BaseAttrInfo) (skillTotalInfo *MazeAIBattle.MazeAISkillTotalInfo, changed bool, err error) {
	skillTotalInfo = &MazeAIBattle.MazeAISkillTotalInfo{}
	for _, attr := range baseAttrs {
		for _, realAttr := range attr.GetRealAttrList() {
			attrSkillCfg := GMazeAttrSkillV8Cfg.Get(realAttr.GetAttrId())
			if attrSkillCfg == nil {
				continue
			}
			if attrSkillCfg.Skill_id > 0 {
				skillInfo, _, err := GetUserBattleSkillInfo(logger, attrSkillCfg.Skill_id, userAttrMap)
				if err != nil {
					return nil, false, err
				}
				skillTotalInfo.SkillInfoList = append(skillTotalInfo.SkillInfoList, skillInfo)
				changed = true
			}
			if attrSkillCfg.Auto_skill_id > 0 {
				autoSkillInfo, err := GetMazeAIAutoSkillInfo(logger, attrSkillCfg.Auto_skill_id, userAttrMap)
				if err != nil {
					return nil, false, err
				}
				skillTotalInfo.AutoSkillInfoList = append(skillTotalInfo.AutoSkillInfoList, autoSkillInfo)
				changed = true
			}
		}
	}
	return skillTotalInfo, changed, nil
}

func GetUserBattleSkillInfo(logger fklog.FKLogI, skillId int32, attrMap map[int32]int64) (*MazeAIBattle.MazeAISkillInfo, []*MazeAIBattle.MazeAIActAttackValue, error) {
	skillCfg := GMazeSkillInfoV8Cfg.Get(skillId)
	if skillCfg == nil {
		logger.ErrorWF("GetMazeAIMonsterConfig GMazeSkillInfoV8Cfg err", zap.Any("skillId", skillId))
		return nil, nil, errors.New("配置不存在")
	}
	skillActCfg := GMazeSkillActV8Cfg.Get(skillId)
	if skillActCfg == nil {
		logger.ErrorWF("GetMazeAIMonsterConfig GMazeSkillActV8Cfg err", zap.Any("skillId", skillId))
		return nil, nil, errors.New("配置不存在")
	}
	actDamageConfigs := make([]*MazeAIBattle.MazeAIActAttackValue, 0)
	if len(skillActCfg.Act_id) > 0 {
		for _, actId := range skillActCfg.Act_id {
			if actId == 0 {
				continue
			}
			mazeActCfg := GMazeActInfoV8Cfg.Get(actId)
			if mazeActCfg == nil {
				logger.ErrorWF("GetMazeAIMonsterConfig GMazeActInfoV8Cfg err", zap.Any("skillId", skillId), zap.Any("actId", actId))
				return nil, nil, errors.New("配置不存在")
			}
			actDamageConfig := &MazeAIBattle.MazeAIActAttackValue{
				ActId:           proto.Int32(actId),
				ToughBrokeValue: proto.Int32(mazeActCfg.Tough_broke_value),
				ToughTempValue:  proto.Int32(mazeActCfg.Temp_tough),
				InitSkillCd:     proto.Int32(skillCfg.Initial_cool_time),
				SkillCd:         proto.Int32(skillCfg.Skill_cool_time),
			}
			actDamageConfigs = append(actDamageConfigs, actDamageConfig)
			for k, v := range mazeActCfg.Attack_point_damage_ratio {
				actDamageRatio := &MazeAIBattle.ActDamageRatioInfo{
					Index: proto.Int32(k),
				}
				actDamageRatio.DamageRatio = &MazeAIBattle.MazeAIAttrInfo{
					UserValue:     proto.Int32(v),
					UserValueType: proto.Int32(1),
				}
				actDamageConfig.ActDamageRatios = append(actDamageConfig.ActDamageRatios, actDamageRatio)
			}

		}
	}
	// todo 缺少触发cd
	skillInfo := &MazeAIBattle.MazeAISkillInfo{
		SkillId:    proto.Int32(skillCfg.Id),
		SkillGroup: proto.Int32(skillCfg.Group),
		// TriggerType: proto.Int32(), // TODO 待配置表补充
		// CampType:                     proto.Int32(skillCfg.Target_type),
		// TargetType:                   proto.Int32(skillCfg.Scope_type),
		RangeRadius:     proto.Int32(skillCfg.Scope_param1),
		ReleaseDistance: proto.Int32(skillCfg.Distance_max),
		ReleaseCd:       proto.Int32(0),
		TargetMaxCount:  proto.Int32(skillCfg.Target_num),
		// CanReleaseState:              skillCfg.Is_allow,
		// CanReleaseTargetState:        skillCfg.Is_target,
		MainTargetDamageRate:         proto.Int32(skillCfg.Main_target_damage),
		SecondTargetDamageRate:       proto.Int32(skillCfg.Second_target_damage),
		SkillDamageFixed:             proto.Int32(skillCfg.Main_target_damage_fix),
		SkillMappingActionId:         skillActCfg.Act_id,
		Level:                        proto.Int32(skillCfg.Level),
		SkillType:                    proto.Int32(skillCfg.Type),
		SkillMappingEffectId:         proto.Int32(skillActCfg.Effect_id),
		SecondTargetSkillDamageFixed: proto.Int32(skillCfg.Second_target_damage_fix),
		IsNoTarget:                   proto.Int32(skillCfg.Is_no_target),
		DamageElement:                skillCfg.Damage_element,
		DamageType:                   proto.Int32(skillCfg.Damage_type),
		InitialCoolTime:              proto.Int32(skillCfg.Initial_cool_time),
		PublicCoolTime:               proto.Int32(skillCfg.Public_cool_time),
		SkillCoolTime:                proto.Int32(skillCfg.Skill_cool_time),
		DistanceMin:                  proto.Int32(skillCfg.Distance_min),
		IsBreak:                      proto.Int32(skillCfg.Is_break),
		ScopeType:                    proto.Int32(skillCfg.Scope_type),
		TargetType:                   proto.Int32(skillCfg.Target_type),
	}
	// for k, v := range skillCfg.Target_effect_pro {
	// 	if k == 0 {
	// 		continue
	// 	}
	// 	skillInfo.RateSourceList = append(skillInfo.RateSourceList, &MazeAIBattle.MazeAISkillEffectRateSource{
	// 		SourceId: proto.Int32(k),
	// 		TargetId: proto.Int32(v),
	// 	})
	// }
	for _, effectId := range skillCfg.Target_effect {
		if effectId == 0 {
			continue
		}
		effectCfg := GMazeSkilleffectV8Cfg.Get(effectId)
		if effectCfg == nil {
			logger.ErrorWF("GetMazeAIMonsterConfig GMazeSkilleffectV8Cfg err", zap.Any("effectId", effectId))
			return nil, nil, errors.New("配置不存在")
		}
		skillEffectOther := &MazeAIBattle.MazeAISkillEffectConfigInfo{
			EffectId:      proto.Int32(effectId),
			EffectGroup:   proto.Int32(effectCfg.Effect_group),
			InGroupWeight: proto.Int32(effectCfg.In_group_weight),
			LastTime:      proto.Int32(effectCfg.Last_time),
			BaseHitrate:   proto.Int32(effectCfg.Base_hitrate),
			CoolDown:      proto.Int32(effectCfg.Cool_down),
			AttrId:        proto.Int32(effectCfg.Attr),
			Value_4:       effectCfg.Attr_value_4,
		}

		skillEffectOther.ValueList = append(skillEffectOther.ValueList, &MazeAIBattle.MazeAIEffectValueInfo{
			Value:     proto.Int64(GetEffectAttrValue(effectCfg.Attr_value, effectCfg.Attr_value_variable_id, attrMap)),
			ValueType: proto.Int32(effectCfg.Attr_value_type),
			Index:     proto.Int32(1),
		})
		skillEffectOther.ValueList = append(skillEffectOther.ValueList, &MazeAIBattle.MazeAIEffectValueInfo{
			Value:     proto.Int64(GetEffectAttrValue(effectCfg.Attr_value_2, effectCfg.Attr_value_2_variable_id, attrMap)),
			ValueType: proto.Int32(effectCfg.Attr_value_2_type),
			Index:     proto.Int32(2),
		})
		skillEffectOther.ValueList = append(skillEffectOther.ValueList, &MazeAIBattle.MazeAIEffectValueInfo{
			Value:     proto.Int64(GetEffectAttrValue(effectCfg.Attr_value_3, effectCfg.Attr_value_3_variable_id, attrMap)),
			ValueType: proto.Int32(effectCfg.Attr_value_3_type),
			Index:     proto.Int32(3),
		})
		skillEffectOther.IntervalTime = proto.Int32(int32(GetEffectAttrValue(effectCfg.Attr_value_8, effectCfg.Attr_value_8_variable_id, attrMap)))
		skillInfo.SkillEffectOther = append(skillInfo.SkillEffectOther, skillEffectOther)
	}

	for _, effectId := range skillCfg.Self_effect {
		if effectId == 0 {
			continue
		}
		effectCfg := GMazeSkilleffectV8Cfg.Get(effectId)
		if effectCfg == nil {
			logger.ErrorWF("GetMazeAIMonsterConfig GMazeSkilleffectV8Cfg err", zap.Any("effectId", effectId))
			return nil, nil, errors.New("配置不存在")
		}
		SkillEffectSelf := &MazeAIBattle.MazeAISkillEffectConfigInfo{
			EffectId:      proto.Int32(effectId),
			EffectGroup:   proto.Int32(effectCfg.Effect_group),
			InGroupWeight: proto.Int32(effectCfg.In_group_weight),
			LastTime:      proto.Int32(effectCfg.Last_time),
			BaseHitrate:   proto.Int32(effectCfg.Base_hitrate),
			CoolDown:      proto.Int32(effectCfg.Cool_down),
			AttrId:        proto.Int32(effectCfg.Attr),
			Value_4:       effectCfg.Attr_value_4,
		}
		SkillEffectSelf.ValueList = append(SkillEffectSelf.ValueList, &MazeAIBattle.MazeAIEffectValueInfo{
			Value:     proto.Int64(GetEffectAttrValue(effectCfg.Attr_value, effectCfg.Attr_value_variable_id, attrMap)),
			ValueType: proto.Int32(effectCfg.Attr_value_type),
			Index:     proto.Int32(1),
		})
		SkillEffectSelf.ValueList = append(SkillEffectSelf.ValueList, &MazeAIBattle.MazeAIEffectValueInfo{
			Value:     proto.Int64(GetEffectAttrValue(effectCfg.Attr_value_2, effectCfg.Attr_value_2_variable_id, attrMap)),
			ValueType: proto.Int32(effectCfg.Attr_value_2_type),
			Index:     proto.Int32(2),
		})
		SkillEffectSelf.ValueList = append(SkillEffectSelf.ValueList, &MazeAIBattle.MazeAIEffectValueInfo{
			Value:     proto.Int64(GetEffectAttrValue(effectCfg.Attr_value_3, effectCfg.Attr_value_3_variable_id, attrMap)),
			ValueType: proto.Int32(effectCfg.Attr_value_3_type),
			Index:     proto.Int32(3),
		})
		SkillEffectSelf.IntervalTime = proto.Int32(int32(GetEffectAttrValue(effectCfg.Attr_value_8, effectCfg.Attr_value_8_variable_id, attrMap)))
		skillInfo.SkillEffectSelf = append(skillInfo.SkillEffectSelf, SkillEffectSelf)
	}
	return skillInfo, actDamageConfigs, nil
}

func GetMazeAIAutoSkillInfo(logger fklog.FKLogI, skillId int32, attrMap map[int32]int64) (*MazeAIBattle.MazeAIAutoSkillInfo, error) {
	skillAutoCfg := GMazeSkillAutoReleaseV8Cfg.Get(skillId)
	if skillAutoCfg == nil {
		logger.ErrorWF("GetMazeAIAutoSkillInfo GMazeSkillAutoReleaseV8Cfg err", zap.Any("skillId", skillId))
		return nil, errors.New("配置不存在")
	}
	skillConfigInfo := &MazeAIBattle.MazeAIAutoSkillInfo{
		SkillId:      proto.Int32(skillId),
		BaseHitrate:  proto.Int32(int32(GetAttrValue(skillAutoCfg.Release_ratio, skillAutoCfg.Release_ratio_variable, attrMap))),
		TriggerCount: proto.Int32(int32(GetAttrValue(skillAutoCfg.Release_num, skillAutoCfg.Release_num_variable, attrMap))),
		ProtectCd:    proto.Int32(int32(GetAttrValue(skillAutoCfg.Auto_release_protect_cd, skillAutoCfg.Auto_release_protect_cd_variable, attrMap))),
		LimitCount:   proto.Int32(int32(GetAttrValue(skillAutoCfg.Max_release_limit, skillAutoCfg.Max_release_limit_variable, attrMap))),
		TriggerType:  proto.Int32(skillAutoCfg.Release_time),
		Condition:    proto.String(skillAutoCfg.Release_condition),
	}
	// 触发技能列表
	for _, triggerSkillId := range skillAutoCfg.Skill_id {
		skillConfigInfo.TriggerSkillId = append(skillConfigInfo.TriggerSkillId, triggerSkillId)
	}
	// 技能条件
	conditionIDs, err := filterConditionIDs(skillAutoCfg.Release_condition)
	if err != nil {
		logger.ErrorWF("GetMazeAIAutoSkillInfo filterConditionIDs err", zap.Error(err), zap.Any("skillId", skillId))
		return nil, errors.New("配置错误")
	}
	for _, conditionID := range conditionIDs {
		condition, err := GetMazeSkillConditionInfo(logger, conditionID, attrMap)
		if err != nil {
			logger.ErrorWF("GetMazeAIAutoSkillInfo GetMazeSkillConditionInfo err", zap.Error(err), zap.Any("skillId", skillId), zap.Any("conditionID", conditionID))
			return nil, err
		}
		skillConfigInfo.ConditionConfig = append(skillConfigInfo.ConditionConfig, condition)
	}
	return skillConfigInfo, nil
}

func GetMazeSkillConditionInfo(logger fklog.FKLogI, conditionID int32, attrMap map[int32]int64) (*MazeAIBattle.MazeSkillCondition, error) {
	cfg := GMazeSkillAutoConditionV8Cfg.Get(conditionID)
	if cfg == nil {
		logger.ErrorWF("GetMazeSkillConditionInfo GMazeSkillAutoConditionV8Cfg err", zap.Any("conditionID", conditionID))
		return nil, errors.New("配置不存在")
	}
	condition := &MazeAIBattle.MazeSkillCondition{
		ConditionId: proto.Int32(conditionID),
		Type:        proto.Int32(cfg.Type),
		Symbol:      proto.Int32(cfg.Symbol),
		Value:       proto.Int32(int32(GetAttrValue(cfg.Value, cfg.Value_variable, attrMap))),
	}
	return condition, nil
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

var (
	regexpnum = regexp.MustCompile(`[1-9][0-9]+`)
)

// filterConditionIDs 用来提取条件字符串中的所有技能条件ID
func filterConditionIDs(condition string) (conditionIDs []int32, err error) {
	if condition == "" {
		return nil, nil
	}
	for _, v := range regexpnum.FindAllString(condition, -1) {
		conditionIDs = append(conditionIDs, fkutil.ToInt32(v))
	}
	return
}

// GetEquipSkillInfoChange 获取装备变化引起的技能变化
func GetEquipSkillInfoChange(logger fklog.FKLogI, userID uint64, oldEquip, newEquip *MazeEquipCache.MazeEquipPosInfo) (ret *MazeAIBattle.MazeUserSkillInfoChangeID, changed bool, err error) {
	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userID)
	if err != nil {
		logger.ErrorWF("GetEquipSkillInfoChange GetUserInfoV2 fail", zap.Error(err), zap.Uint64("userID", userID))
		return
	}

	userAttrMap, err := GetUserAttrMap(logger, userID)
	if err != nil {
		logger.ErrorWF("GetMazeBattleData GetUserAttrMap err", zap.Error(err))
		return nil, false, err
	}

	tempBuffInfo, err := mazebarriertempbuffredis.GetBarrierTempBuff(logger, userID, userInfo.Barrier)
	if err != nil {
		logger.ErrorWF("GetMazeBattleData GetBarrierTempBuff err", zap.Error(err))
		return nil, false, err
	}
	for _, buffInfo := range tempBuffInfo.TotalBuff {
		userAttrMap[buffInfo.GetBuffId()] += buffInfo.GetBuffValue()
	}

	// 脱下的装备会删除技能
	if oldEquip != nil {
		equip := oldEquip.GetEquipInfo()
		ret.DelSkillInfoList, changed, err = GetUserSkillTotalInfo(logger, userAttrMap, equip.GetBaseAttrs())
	}
	// 穿戴的装备会增加技能
	if newEquip != nil {
		equip := newEquip.GetEquipInfo()
		ret.AddSkillInfoList, changed, err = GetUserSkillTotalInfo(logger, userAttrMap, equip.GetBaseAttrs())
	}
	return ret, changed, nil
}

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
