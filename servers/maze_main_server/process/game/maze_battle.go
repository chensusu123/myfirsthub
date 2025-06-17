package game

import (
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeActInfoV8Cfg"
	"maze_game_server/config/GMazeAttrItemAttrV8Cfg"
	"maze_game_server/config/GMazeAttrSkillV8Cfg"
	"maze_game_server/config/GMazeBrushFoeV8Cfg"
	"maze_game_server/config/GMazeFoeV8Cfg"
	"maze_game_server/config/GMazeSkillActV8Cfg"
	"maze_game_server/config/GMazeSkillAutoConditionV8Cfg"
	"maze_game_server/config/GMazeSkillAutoReleaseV8Cfg"
	"maze_game_server/config/GMazeSkillInfoV8Cfg"
	"maze_game_server/config/GMazeSkilleffectV8Cfg"
	"maze_game_server/io/redis/mazebarriertempbuffredis"
	"maze_game_server/io/redis/mazecalcattrredis"
	"maze_game_server/pb/common/MazeAIBattle"
	"regexp"
	"sort"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 获取迷宫战斗数据
func GetMazeBattleData(logger fklog.FKLogI, userId uint64, barrierId int32) (mazeBattleInfo *MazeAIBattle.MazeBarrierInfo, err error) {
	defer fkprometheus.InfoPMT("GetMazeBattleData")()
	mazeBattleInfo = &MazeAIBattle.MazeBarrierInfo{
		BarrierId: proto.Int32(barrierId),
	}
	userAttrMap, err := GetUserAttrMap(logger, userId)
	if err != nil {
		logger.ErrorWF("GetMazeBattleData GetUserAttrMap err", zap.Error(err))
		return nil, err
	}

	force, err := mazecalcattrredis.GetMazeForce(logger, userId)
	if err != nil {
		logger.ErrorWF("OnMazeLoginRQ GetMazeForce fail", zap.Error(err))
		return nil, err
	}

	tempBuffInfo, err := mazebarriertempbuffredis.GetBarrierTempBuff(logger, userId, barrierId)
	if err != nil {
		logger.ErrorWF("GetMazeBattleData GetBarrierTempBuff err", zap.Error(err))
		return nil, err
	}
	for _, buffInfo := range tempBuffInfo.TotalBuff {
		userAttrMap[buffInfo.GetBuffId()] += buffInfo.GetBuffValue()
	}

	userStiffRatio := userAttrMap[constdef.MazeAttr3000101]
	areaInfos, err := GetFoeAreaInfos(logger, userId, force, barrierId, userStiffRatio)
	if err != nil {
		logger.ErrorWF("GetMazeBattleData GetFoeAreaInfos err", zap.Any("barrierId", barrierId), zap.Error(err))
		return nil, err
	}
	mazeBattleInfo.AreaInfos = areaInfos
	userAttrInfo, err := GetUserAttrInfo(logger, userId, userAttrMap)
	if err != nil {
		logger.ErrorWF("GetMazeBattleData GetUserAttrInfo err", zap.Error(err), zap.Any("userId", userId))
		return nil, err
	}
	mazeBattleInfo.RoleConfigInfo = userAttrInfo
	for _, cfg := range GMazeFoeV8Cfg.GetAll() {
		if cfg.In_barries_id != barrierId {
			continue
		}
		if cfg.Foe_type == 1 {
			continue
		}
		eliteMonsterConfig, err := GetMazeAIMonsterConfig(logger, userId, force, cfg.Order, userStiffRatio)
		if err != nil {
			logger.ErrorWF("GetMazeBattleData GetMazeAIMonsterConfig err", zap.Any("foeId", cfg.Order))
			return nil, err
		}
		mazeBattleInfo.EliteMonsterInfos = append(mazeBattleInfo.EliteMonsterInfos, eliteMonsterConfig)
	}

	foeSkillMap := make(map[int32]struct{})
	for _, areaInfo := range mazeBattleInfo.AreaInfos {
		for _, foeCfg := range areaInfo.MonsterConfigInfos {
			for _, skillInfo := range foeCfg.GetSkillTotalInfo().GetSkillInfoList() {
				foeSkillMap[skillInfo.GetSkillId()] = struct{}{}
			}
		}
	}
	for _, eliteInfo := range mazeBattleInfo.EliteMonsterInfos {
		for _, skillInfo := range eliteInfo.GetSkillTotalInfo().GetSkillInfoList() {
			foeSkillMap[skillInfo.GetSkillId()] = struct{}{}
		}
	}
	for skillId := range foeSkillMap {
		skillConfigInfo, err := GetFoeSkillConfigInfo(logger, skillId)
		if err != nil {
			logger.ErrorWF("GetMazeBattleData GetFoeSkillConfigInfo err", zap.Any("skillId", skillId), zap.Error(err))
			return nil, err
		}
		mazeBattleInfo.SkillConfigInfos = append(mazeBattleInfo.SkillConfigInfos, skillConfigInfo)
	}

	// 道具使用配置
	for _, row := range GMazeAttrItemAttrV8Cfg.GetAll() {
		if row.Add_attr <= 0 {
			continue
		}
		attrSkill := GMazeAttrSkillV8Cfg.Get(row.Add_attr)
		if attrSkill == nil {
			continue
		}
		if attrSkill.Skill_id <= 0 {
			continue
		}
		itemUseInfo := &MazeAIBattle.MazeItemUseInfo{
			ItemId: proto.Int32(row.Order),
		}
		itemUseInfo.SkillIds = append(itemUseInfo.SkillIds, attrSkill.Skill_id)
		mazeBattleInfo.ItemUseInfos = append(mazeBattleInfo.ItemUseInfos, itemUseInfo)
	}

	return mazeBattleInfo, nil
}

func GetFoeAreaInfos(logger fklog.FKLogI, userId uint64, force int64, barrierId int32, userStiffRatio int64) ([]*MazeAIBattle.MazeAIAreaInfo, error) {
	areaInfos := make([]*MazeAIBattle.MazeAIAreaInfo, 0)
	areaFoeMap := make(map[int32]map[int32]struct{})
	for _, cfg := range GMazeBrushFoeV8Cfg.GetAll() {
		if cfg.Barries_id != barrierId {
			continue
		}
		if areaFoeMap[cfg.Brush_area_id] == nil {
			areaFoeMap[cfg.Brush_area_id] = make(map[int32]struct{})
		}
		for _, foeId := range cfg.Monsters_id {
			areaFoeMap[cfg.Brush_area_id][foeId] = struct{}{}
		}
	}
	logger.InfoWF("GetFoeAreaInfos area foes dumps",
		zap.Int32("barrierId", barrierId),
		zap.Any("areaFoeMap", areaFoeMap),
	)
	for areaId, foeMap := range areaFoeMap {
		areaInfo := &MazeAIBattle.MazeAIAreaInfo{
			AreaId: proto.Int32(areaId),
		}
		areaInfo.MonsterConfigInfos = make([]*MazeAIBattle.MazeAIMonsterConfigInfo, 0)
		for foeId := range foeMap {
			monsterConfigInfo, err := GetMazeAIMonsterConfig(logger, userId, force, foeId, userStiffRatio)
			if err != nil {
				logger.ErrorWF("GetMazeBattleData GetMazeAIMonsterConfig err", zap.Any("foeId", foeId))
				return nil, err
			}
			logger.InfoWF("GetFoeAreaInfos GetMazeAIMonsterConfig dumps",
				zap.Int32("foeId", foeId),
				zap.Any("monsterConfigInfo", monsterConfigInfo),
			)
			areaInfo.MonsterConfigInfos = append(areaInfo.MonsterConfigInfos, monsterConfigInfo)
		}
		areaInfos = append(areaInfos, areaInfo)
	}
	return areaInfos, nil
}

// todo 技能公共cd
func GetMazeAIMonsterConfig(logger fklog.FKLogI, userId uint64, force int64, foeId int32, userStiffRatio int64) (*MazeAIBattle.MazeAIMonsterConfigInfo, error) {
	foeCfg := GMazeFoeV8Cfg.Get(foeId)
	if foeCfg == nil {
		logger.ErrorWF("GetMazeAIMonsterConfig GMazeFoeV8Cfg err", zap.Any("foeId", foeId))
		return nil, errors.New("配置不存在")
	}
	monsterConfigInfo := &MazeAIBattle.MazeAIMonsterConfigInfo{
		FoeType: proto.Int32(foeCfg.Foe_type),
	}
	attackValue := &MazeAIBattle.MazeAIAttackValue{
		ConfigId:          proto.Int32(foeId),
		MaxHp:             proto.Int64(int64(foeCfg.Hp_max)),
		MonsterStiffRatio: proto.Int64(int64(foeCfg.Tough_max)),
		// UserStiffRatio:    proto.Int64(userStiffRatio),
		UserStiffRatio: proto.Int64(CalcUserStiffRatio(foeCfg, force)),
	}
	monsterConfigInfo.AttackValue = attackValue
	monsterConfigInfo.AttrInfo = make([]*MazeAIBattle.MazeAIAttrInfo, 0)
	monsterConfigInfo.AttrInfo = append(monsterConfigInfo.AttrInfo, &MazeAIBattle.MazeAIAttrInfo{
		Type:          proto.Int32(int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ROLE_ATK_VALUE)),
		UserValue:     proto.Int32(foeCfg.Attack_max),
		UserValueType: proto.Int32(1),
	})
	monsterConfigInfo.AttrInfo = append(monsterConfigInfo.AttrInfo, &MazeAIBattle.MazeAIAttrInfo{
		Type:          proto.Int32(int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ROLE_DEF_VALUE)),
		UserValue:     proto.Int32(foeCfg.Def_max),
		UserValueType: proto.Int32(1),
	})
	// 怪物攻击速度提高万分比
	monsterConfigInfo.AttrInfo = append(monsterConfigInfo.AttrInfo, &MazeAIBattle.MazeAIAttrInfo{
		Type:          proto.Int32(int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ATK_SPEED)),
		UserValue:     proto.Int32(foeCfg.Attack_speed_pro),
		UserValueType: proto.Int32(2),
	})
	// 怪物受击回复速度提高万分比
	monsterConfigInfo.AttrInfo = append(monsterConfigInfo.AttrInfo, &MazeAIBattle.MazeAIAttrInfo{
		Type:          proto.Int32(int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_RECOVERY_SPEED)),
		UserValue:     proto.Int32(foeCfg.Be_attack_recovery_speed_pro),
		UserValueType: proto.Int32(2),
	})
	// 怪物韧性上限
	monsterConfigInfo.AttrInfo = append(monsterConfigInfo.AttrInfo, &MazeAIBattle.MazeAIAttrInfo{
		Type:          proto.Int32(int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_TOUGH_MAX)),
		UserValue:     proto.Int32(foeCfg.Tough_max),
		UserValueType: proto.Int32(1),
	})
	skillIds := make([]int32, 0)
	if foeCfg.Nor_attack_skill_id > 0 {
		skillIds = append(skillIds, foeCfg.Nor_attack_skill_id)
	}
	for _, skillId := range foeCfg.Passive_skill_id {
		skillIds = append(skillIds, skillId)
	}
	if len(skillIds) > 0 {
		skillCfg := GMazeSkillInfoV8Cfg.Get(skillIds[0])
		if skillCfg == nil {
			logger.ErrorWF("GetMazeAIMonsterConfig GMazeSkillInfoV8Cfg err", zap.Any("skillId", skillIds[0]))
			return nil, errors.New("配置不存在")
		}
		attackValue.SkillCd = proto.Int32(skillCfg.Public_cool_time)
	}
	skillTotalInfo := &MazeAIBattle.MazeAISkillTotalInfo{}
	skillTotalInfo.SkillInfoList = make([]*MazeAIBattle.MazeAISkillInfo, 0)
	attrMap := make(map[int32]int64)
	for _, skillId := range skillIds {
		if skillId == 0 {
			continue
		}
		skillInfo, err := GetFoeBattleSkillInfo(logger, skillId, attrMap)
		if err != nil {
			logger.WarnWF("GetUserBattleAttr BattleSkillTopPb nil", zap.Uint64("userId", userId), zap.Any("skillId", skillId))
			return nil, err
		}
		skillTotalInfo.SkillInfoList = append(skillTotalInfo.SkillInfoList, skillInfo)
		// attackValue.ActDamageConfig = append(attackValue.ActDamageConfig, actDamageConfigs...)
	}
	// 韧性被打空时释放技能
	if foeCfg.Tough_deplete > 0 {
		skillInfo, err := GetFoeBattleSkillInfo(logger, foeCfg.Tough_deplete, attrMap)
		if err != nil {
			logger.WarnWF("GetUserBattleAttr BattleSkillTopPb nil", zap.Uint64("userId", userId), zap.Any("Tough_deplete", foeCfg.Tough_deplete))
			return nil, err
		}
		skillTotalInfo.SkillInfoList = append(skillTotalInfo.SkillInfoList, skillInfo)
	}

	monsterConfigInfo.MonsterSpeed = proto.Int32(foeCfg.Speed)
	monsterConfigInfo.SkillTotalInfo = skillTotalInfo
	return monsterConfigInfo, nil
}

func GetUserAttrInfo(logger fklog.FKLogI, userId uint64, userAttrMap map[int32]int64) (*MazeAIBattle.MazeAIRoleConfigInfo, error) {
	skillIds := make([]int32, 0)
	autoSkillId := make([]int32, 0)
	for _, cfg := range GMazeSkillInfoV8Cfg.GetAll() {
		if cfg.Type != 1 {
			continue
		}
		// if cfg.Level != displayLevel {
		//	continue
		// }
		skillIds = append(skillIds, cfg.Id)
	}

	for _, cfg := range GMazeAttrSkillV8Cfg.GetAll() {
		if userAttrMap[cfg.Attr_id] > 0 {
			skillIds = append(skillIds, cfg.Skill_id)
			if cfg.Auto_skill_id > 0 {
				autoSkillId = append(autoSkillId, cfg.Auto_skill_id)
			}
		}
	}
	roleConfigInfo := &MazeAIBattle.MazeAIRoleConfigInfo{}
	userAttrInfo := &MazeAIBattle.MazeAIUserAttrInfo{}
	attrMap, err := GetUserBattleAttr(logger, userId, userAttrMap)
	if err != nil {
		logger.WarnWF("GetUserBattleAttr BatchGetDollCalcAttr nil", zap.Uint64("userId", userId))
		return nil, err
	}
	if attrMap[constdef.DollFormulaBlood] != nil {
		// userAttrInfo.UserHp = proto.Int64(int64(attrMap[constdef.DollFormulaBlood].GetUserValue()))
		userAttrInfo.UserTotal = proto.Int64(int64(attrMap[constdef.DollFormulaBlood].GetUserValue()))
	}
	for _, attrInfo := range attrMap {
		if attrInfo.GetType() <= 0 {
			continue
		}
		userAttrInfo.AttrInfo = append(userAttrInfo.AttrInfo, attrInfo)
	}
	roleConfigInfo.UserAttrInfo = userAttrInfo
	userSkillInfo := &MazeAIBattle.MazeAISkillTotalInfo{}
	userSkillInfo.SkillInfoList = make([]*MazeAIBattle.MazeAISkillInfo, 0)
	userSkillInfo.AutoSkillInfoList = make([]*MazeAIBattle.MazeAIAutoSkillInfo, 0)
	actDamageConfigList := make([]*MazeAIBattle.MazeAIActAttackValue, 0)
	for _, skillId := range autoSkillId {
		if skillId == 0 {
			continue
		}
		autoSkillInfo, err := GetMazeAIAutoSkillInfo(logger, skillId, userAttrMap)
		if err != nil {
			logger.WarnWF("GetUserBattleAttr GetMazeAIAutoSkillInfo nil", zap.Uint64("userId", userId), zap.Any("skillId", skillId))
			return nil, err
		}
		userSkillInfo.AutoSkillInfoList = append(userSkillInfo.AutoSkillInfoList, autoSkillInfo)
	}
	for _, skillInfo := range userSkillInfo.AutoSkillInfoList {
		if len(skillInfo.TriggerSkillId) <= 0 {
			continue
		}
		skillIds = append(skillIds, skillInfo.TriggerSkillId...)
	}
	for _, skillId := range skillIds {
		if skillId == 0 {
			continue
		}
		skillInfo, actDamageConfigs, err := GetUserBattleSkillInfo(logger, skillId, userAttrMap)
		if err != nil {
			logger.WarnWF("GetUserBattleAttr GetUserBattleSkillInfo nil", zap.Uint64("userId", userId), zap.Any("skillId", skillId))
			return nil, err
		}
		userSkillInfo.SkillInfoList = append(userSkillInfo.SkillInfoList, skillInfo)
		actDamageConfigList = append(actDamageConfigList, actDamageConfigs...)
	}
	// 使用道具后可使用属性技能
	for _, row := range GMazeAttrItemAttrV8Cfg.GetAll() {
		if row.Add_attr <= 0 {
			continue
		}
		attrSkill := GMazeAttrSkillV8Cfg.Get(row.Add_attr)
		if attrSkill == nil {
			continue
		}
		if attrSkill.Skill_id <= 0 {
			continue
		}
		skillInfo, actDamageConfigs, err := GetUserBattleSkillInfo(logger, attrSkill.Skill_id, userAttrMap)
		if err != nil {
			logger.WarnWF("GetUserBattleAttr GetUserBattleSkillInfo nil", zap.Uint64("userId", userId), zap.Any("attrSkillId", attrSkill.Skill_id))
			return nil, err
		}
		userSkillInfo.SkillInfoList = append(userSkillInfo.SkillInfoList, skillInfo)
		actDamageConfigList = append(actDamageConfigList, actDamageConfigs...)
	}

	roleConfigInfo.UserSkillInfo = userSkillInfo
	roleConfigInfo.ActDamageConfig = actDamageConfigList
	return roleConfigInfo, nil
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

func GetFoeBattleSkillInfo(logger fklog.FKLogI, skillId int32, attrMap map[int32]int64) (*MazeAIBattle.MazeAISkillInfo, error) {
	skillCfg := GMazeSkillInfoV8Cfg.Get(skillId)
	if skillCfg == nil {
		logger.ErrorWF("GetMazeAIMonsterConfig GMazeSkillInfoV8Cfg err", zap.Any("skillId", skillId))
		return nil, errors.New("配置不存在")
	}
	skillActCfg := GMazeSkillActV8Cfg.Get(skillId)
	if skillActCfg == nil {
		logger.ErrorWF("GetMazeAIMonsterConfig GMazeSkillActV8Cfg err", zap.Any("skillId", skillId))
		return nil, errors.New("配置不存在")
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
			return nil, errors.New("配置不存在")
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
			return nil, errors.New("配置不存在")
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
	return skillInfo, nil
}

func GetFoeSkillConfigInfo(logger fklog.FKLogI, skillId int32) (*MazeAIBattle.MazeSkillConfigInfo, error) {
	skillCfg := GMazeSkillInfoV8Cfg.Get(skillId)
	if skillCfg == nil {
		logger.ErrorWF("GetMazeAIMonsterConfig GMazeSkillInfoV8Cfg err", zap.Any("skillId", skillId))
		return nil, errors.New("配置不存在")
	}
	skillActCfg := GMazeSkillActV8Cfg.Get(skillId)
	if skillActCfg == nil {
		logger.ErrorWF("GetMazeAIMonsterConfig GMazeSkillActV8Cfg err", zap.Any("skillId", skillId))
		return nil, errors.New("配置不存在")
	}
	skillConfigInfo := &MazeAIBattle.MazeSkillConfigInfo{
		SkillId:     proto.Int32(skillId),
		InitSkillCd: proto.Int32(skillCfg.Initial_cool_time),
		SkillCd:     proto.Int32(skillCfg.Skill_cool_time),
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
				return nil, errors.New("配置不存在")
			}
			actDamageConfig := &MazeAIBattle.MazeAIActAttackValue{
				ActId:           proto.Int32(actId),
				ToughBrokeValue: proto.Int32(mazeActCfg.Tough_broke_value),
				ToughTempValue:  proto.Int32(mazeActCfg.Temp_tough),
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
	skillConfigInfo.ActDamageConfig = actDamageConfigs
	return skillConfigInfo, nil
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

func CalcUserStiffRatio(cfg *GMazeFoeV8Cfg.MazeFoeV8ConfigRow, force int64) (ratio int64) {
	if cfg == nil {
		return 0
	}
	keys := make([]int64, 0, len(cfg.Tough_borke_raitio))
	for k, _ := range cfg.Tough_borke_raitio {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	for i, k := range keys {
		if force < k {
			if i == 0 {
				return 0
			}
			return int64(cfg.Tough_borke_raitio[keys[i-1]])
		}
	}
	if len(keys) > 0 {
		return int64(cfg.Tough_borke_raitio[keys[len(keys)-1]])
	}
	return 0
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
