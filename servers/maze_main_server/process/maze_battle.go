package process

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebarriertempbuffredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazecalcattrredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazehurtcalc"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeActInfoV8Cfg"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeAttrSkillV8Cfg"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeBrushFoeV8Cfg"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeFoeV8Cfg"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeSkillActV8Cfg"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeSkillInfoV8Cfg"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeSkilleffectV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/protodef/MazeAIBattle"
	"go.uber.org/zap"
)

// 获取迷宫战斗数据
func GetMazeBattleData(logger fklog.FKLogI, userId uint64, barrierId int32) (mazeBattleInfo *MazeAIBattle.MazeBarrierInfo, err error) {
	defer fkprometheus.InfoPMT("GetMazeBattleData")()
	mazeBattleInfo = &MazeAIBattle.MazeBarrierInfo{
		BarrierId: proto.Int32(barrierId),
	}
	forceVal, err := mazecalcattrredis.GetMazeForce(logger, userId)
	if err != nil {
		logger.ErrorWF("GetMazeBattleData GetMazeForce fail", zap.Error(err))
		return nil, err
	}
	userAttrMap, err := GetUserAttrMap(logger, userId)
	if err != nil {
		logger.ErrorWF("GetMazeBattleData GetUserAttrMap err", zap.Error(err))
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
	userTotalBlood := userAttrMap[constdef.DollFormulaBlood]
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
	for areaId, foeMap := range areaFoeMap {
		areaInfo := &MazeAIBattle.MazeAIAreaInfo{
			AreaId: proto.Int32(areaId),
		}
		areaInfo.MonsterConfigInfos = make([]*MazeAIBattle.MazeAIMonsterConfigInfo, 0)
		//areaInfo.BrushConfigInfo = &MazeAIBattle.MazeAIBrushAreaConfigInfo{
		//	FirstFoe: proto.Int32(cfg.First_foe),
		//	MinFoe:   proto.Int32(cfg.Min_foe),
		//	MaxFoe:   proto.Int32(cfg.Max_foe),
		//	BrushCd:  proto.Int32(cfg.Brush_cd),
		//}
		for foeId := range foeMap {
			monsterConfigInfo, err := GetMazeAIMonsterConfig(logger, userId, foeId, forceVal, userStiffRatio, userTotalBlood)
			if err != nil {
				logger.ErrorWF("GetMazeBattleData GetMazeAIMonsterConfig err", zap.Any("foeId", foeId))
				return nil, err
			}
			areaInfo.MonsterConfigInfos = append(areaInfo.MonsterConfigInfos, monsterConfigInfo)
			//areaInfo.BrushConfigInfo.FoePool = append(areaInfo.BrushConfigInfo.FoePool, &MazeAIBattle.BrushFoePoolInfo{
			//	MonsterId:    proto.Int32(foeId),
			//	MonsterCount: proto.Int32(count),
			//})
		}
		mazeBattleInfo.AreaInfos = append(mazeBattleInfo.AreaInfos, areaInfo)
	}
	userAttrInfo, err := GetUserAttrInfo(logger, userId, forceVal, userAttrMap)
	if err != nil {
		logger.ErrorWF("GetMazeBattleData GetUserAttrInfo err", zap.Error(err), zap.Any("userId", userId), zap.Any("forceVal", forceVal))
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
		eliteMonsterConfig, err := GetMazeAIMonsterConfig(logger, userId, cfg.Order, forceVal, userStiffRatio, userTotalBlood)
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
	return mazeBattleInfo, nil
}

// todo 技能公共cd
func GetMazeAIMonsterConfig(logger fklog.FKLogI, userId uint64, foeId int32, forceVal int64, userStiffRatio int64, userTotalBlood int64) (*MazeAIBattle.MazeAIMonsterConfigInfo, error) {
	foeCfg := GMazeFoeV8Cfg.Get(foeId)
	if foeCfg == nil {
		logger.ErrorWF("GetMazeAIMonsterConfig GMazeFoeV8Cfg err", zap.Any("foeId", foeId))
		return nil, errors.New("配置不存在")
	}
	p1 := mazehurtcalc.MazeAttrInfo{}
	p1.Force = forceVal
	p1.Hp = userTotalBlood
	p2 := mazehurtcalc.MazeAttrInfo{}
	p2.Force = int64(foeCfg.Kongfu)
	p2.Hp = int64(foeCfg.Hp_max)
	r, e := mazehurtcalc.CalcLoseHurt(logger, foeCfg.Hp_lose_type, &p1, &p2)
	if e != nil {
		logger.ErrorWF("GetMazeAIMonsterConfig CalcLoseHurt err", zap.Any("foeId", foeId), zap.Error(e))
		return nil, errors.New("配置不存在")
	}
	monsterConfigInfo := &MazeAIBattle.MazeAIMonsterConfigInfo{
		FoeType: proto.Int32(foeCfg.Foe_type),
	}
	attackValue := &MazeAIBattle.MazeAIAttackValue{
		ConfigId:          proto.Int32(foeId),
		MonsterValue:      proto.Int64(r.PlayerLoseHurt),
		UserBaseValue:     proto.Int64(r.FoeLoseHurt),
		MaxHp:             proto.Int64(int64(foeCfg.Hp_max)),
		MonsterStiffRatio: proto.Int64(int64(foeCfg.Tough_max)),
		UserStiffRatio:    proto.Int64(userStiffRatio),
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
		//attackValue.ActDamageConfig = append(attackValue.ActDamageConfig, actDamageConfigs...)
	}
	monsterConfigInfo.MonsterSpeed = proto.Int32(foeCfg.Speed)
	monsterConfigInfo.SkillTotalInfo = skillTotalInfo
	return monsterConfigInfo, nil
}

func GetUserAttrInfo(logger fklog.FKLogI, userId uint64, forceVal int64, userAttrMap map[int32]int64) (*MazeAIBattle.MazeAIRoleConfigInfo, error) {
	//displayLevel := int32(1)
	//for _, cfg := range GMazeKongfuDisplayV8Cfg.GetAll() {
	//	if forceVal >= cfg.Kongfu_min && forceVal <= cfg.Kongfu_max {
	//		displayLevel = cfg.Display_level
	//		break
	//	}
	//}
	skillIds := make([]int32, 0)
	for _, cfg := range GMazeSkillInfoV8Cfg.GetAll() {
		if cfg.Type != 1 {
			continue
		}
		//if cfg.Level != displayLevel {
		//	continue
		//}
		skillIds = append(skillIds, cfg.Id)
	}

	for _, cfg := range GMazeAttrSkillV8Cfg.GetAll() {
		if userAttrMap[cfg.Attr_id] > 0 {
			skillIds = append(skillIds, cfg.Skill_id)
		}
	}
	roleConfigInfo := &MazeAIBattle.MazeAIRoleConfigInfo{}
	userAttrInfo := &MazeAIBattle.MazeAIUserAttrInfo{}
	attrMap, err := GetUserBattleAttr(logger, userId, skillIds, userAttrMap)
	if err != nil {
		logger.WarnWF("GetUserBattleAttr BatchGetDollCalcAttr nil", zap.Uint64("userId", userId), zap.Any("forceVal", forceVal))
		return nil, err
	}
	if attrMap[constdef.DollFormulaBlood] != nil {
		//userAttrInfo.UserHp = proto.Int64(int64(attrMap[constdef.DollFormulaBlood].GetUserValue()))
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
	actDamageConfigList := make([]*MazeAIBattle.MazeAIActAttackValue, 0)
	for _, skillId := range skillIds {
		if skillId == 0 {
			continue
		}
		skillInfo, actDamageConfigs, err := GetUserBattleSkillInfo(logger, skillId, userAttrMap)
		if err != nil {
			logger.WarnWF("GetUserBattleAttr BattleSkillTopPb nil", zap.Uint64("userId", userId), zap.Any("skillId", skillId))
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
	//todo 缺少触发cd
	skillInfo := &MazeAIBattle.MazeAISkillInfo{
		SkillId:                proto.Int32(skillCfg.Id),
		SkillGroup:             proto.Int32(skillCfg.Group),
		CampType:               proto.Int32(skillCfg.Target_type),
		TargetType:             proto.Int32(skillCfg.Scope_type),
		RangeRadius:            proto.Int32(skillCfg.Scope_param1),
		ReleaseDistance:        proto.Int32(skillCfg.Distance_min),
		ReleaseCd:              proto.Int32(0),
		TargetMaxCount:         proto.Int32(skillCfg.Target_num),
		CanReleaseState:        skillCfg.Is_allow,
		CanReleaseTargetState:  skillCfg.Is_target,
		MainTargetDamageRate:   proto.Int32(skillCfg.Main_target_damage),
		SecondTargetDamageRate: proto.Int32(skillCfg.Second_target_damage),
		SkillDamageFixed:       proto.Int32(skillCfg.Main_target_damage_fix),
		SkillMappingActionId:   skillActCfg.Act_id,
		Level:                  proto.Int32(skillCfg.Level),
		SkillType:              proto.Int32(skillCfg.Type),
	}
	for k, v := range skillCfg.Target_effect_pro {
		if k == 0 {
			continue
		}
		skillInfo.RateSourceList = append(skillInfo.RateSourceList, &MazeAIBattle.MazeAISkillEffectRateSource{
			SourceId: proto.Int32(k),
			TargetId: proto.Int32(v),
		})
	}
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
	//todo 缺少触发cd
	skillInfo := &MazeAIBattle.MazeAISkillInfo{
		SkillId:                proto.Int32(skillCfg.Id),
		SkillGroup:             proto.Int32(skillCfg.Group),
		CampType:               proto.Int32(skillCfg.Target_type),
		TargetType:             proto.Int32(skillCfg.Scope_type),
		RangeRadius:            proto.Int32(skillCfg.Scope_param1),
		ReleaseDistance:        proto.Int32(skillCfg.Distance_min),
		ReleaseCd:              proto.Int32(0),
		TargetMaxCount:         proto.Int32(skillCfg.Target_num),
		CanReleaseState:        skillCfg.Is_allow,
		CanReleaseTargetState:  skillCfg.Is_target,
		MainTargetDamageRate:   proto.Int32(skillCfg.Main_target_damage),
		SecondTargetDamageRate: proto.Int32(skillCfg.Second_target_damage),
		SkillDamageFixed:       proto.Int32(skillCfg.Main_target_damage_fix),
		SkillMappingActionId:   skillActCfg.Act_id,
		Level:                  proto.Int32(skillCfg.Level),
		SkillType:              proto.Int32(skillCfg.Type),
	}
	for k, v := range skillCfg.Target_effect_pro {
		if k == 0 {
			continue
		}
		skillInfo.RateSourceList = append(skillInfo.RateSourceList, &MazeAIBattle.MazeAISkillEffectRateSource{
			SourceId: proto.Int32(k),
			TargetId: proto.Int32(v),
		})
	}
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
