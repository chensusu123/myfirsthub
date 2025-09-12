package game

import (
	"context"
	"fmt"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeActInfoV8Cfg"
	"maze_game_server/config/GMazeAttrItemAttrV8Cfg"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/config/GMazeBoxV8Cfg"
	"maze_game_server/config/GMazeBrushFoeV8Cfg"
	"maze_game_server/config/GMazeFoeV8Cfg"
	"maze_game_server/config/GMazeSkillActV8Cfg"
	"maze_game_server/config/GMazeSkillAutoConditionV8Cfg"
	"maze_game_server/config/GMazeSkillInfoV8Cfg"
	"maze_game_server/config/GMazeSkilleffectV8Cfg"
	"maze_game_server/config/GMazeSummonV8Cfg"
	"maze_game_server/io/redis/mazeboxredis"
	"maze_game_server/io/redis/mazecalcattrredis"
	"maze_game_server/pb/common/MazeAIBattle"
	"maze_game_server/services/tempbuffservice"
	"regexp"
	"sort"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 获取迷宫战斗数据
func GetMazeBattleData(ctx context.Context, userId uint64, barrierId int32) (mazeBattleInfo *MazeAIBattle.MazeBarrierInfo, err error) {
	logger := fklog.ContextAppLogger(ctx)
	defer fkprometheus.InfoPMT("GetMazeBattleData")()
	mazeBattleInfo = &MazeAIBattle.MazeBarrierInfo{
		BarrierId: proto.Int32(barrierId),
	}
	userAttrMap, err := GetUserAttrMap(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "GetMazeBattleData GetUserAttrMap err", zap.Error(err))
		return nil, err
	}

	force, err := mazecalcattrredis.GetMazeForce(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnMazeLoginRQ GetMazeForce fail", zap.Error(err))
		return nil, err
	}

	tempBuffInfo, err := tempbuffservice.GlobalTempBuffService.GetTempBuffInfo(ctx, userId, barrierId)
	if err != nil {
		logger.CtxError(ctx, "GetMazeBattleData GetBarrierTempBuff err", zap.Error(err))
		return nil, err
	}
	for _, buffInfo := range tempBuffInfo.TotalBuff {
		userAttrMap[buffInfo.BuffId] += buffInfo.BuffValue
	}

	userStiffRatio := userAttrMap[constdef.MazeAttr3000101]
	areaInfos, err := GetFoeAreaInfos(ctx, userId, force, barrierId, userStiffRatio)
	if err != nil {
		logger.CtxError(ctx, "GetMazeBattleData GetFoeAreaInfos err", zap.Any("barrierId", barrierId), zap.Error(err))
		return nil, err
	}
	mazeBattleInfo.AreaInfos = areaInfos
	userAttrInfo, summonIds, err := GetUserAttrInfo(ctx, userId, userAttrMap)
	if err != nil {
		logger.CtxError(ctx, "GetMazeBattleData GetUserAttrInfo err", zap.Error(err), zap.Any("userId", userId))
		return nil, err
	}
	mazeBattleInfo.RoleConfigInfo = userAttrInfo

	// 精简包结构：客户端已经不使用服务器回包的精英怪数据
	// for _, cfg := range GMazeFoeV8Cfg.GetAll() {
	// 	if cfg.In_barries_id != barrierId {
	// 		continue
	// 	}
	// 	if cfg.Foe_type == 1 {
	// 		continue
	// 	}
	// 	eliteMonsterConfig, err := GetMazeAIMonsterConfig(ctx, userId, force, cfg.Order, userStiffRatio)
	// 	if err != nil {
	// 		logger.CtxError(ctx,"GetMazeBattleData GetMazeAIMonsterConfig err", zap.Any("foeId", cfg.Order))
	// 		return nil, err
	// 	}
	// 	mazeBattleInfo.EliteMonsterInfos = append(mazeBattleInfo.EliteMonsterInfos, eliteMonsterConfig)
	// }

	foeSkillMap := make(map[int32]struct{})
	for _, areaInfo := range mazeBattleInfo.AreaInfos {
		for _, foeCfg := range areaInfo.MonsterConfigInfos {
			for _, skillInfo := range foeCfg.GetSkillTotalInfo().GetSkillInfoList() {
				foeSkillMap[skillInfo.GetSkillId()] = struct{}{}
			}
		}
	}
	// for _, eliteInfo := range mazeBattleInfo.EliteMonsterInfos {
	// 	for _, skillInfo := range eliteInfo.GetSkillTotalInfo().GetSkillInfoList() {
	// 		foeSkillMap[skillInfo.GetSkillId()] = struct{}{}
	// 	}
	// }
	for skillId := range foeSkillMap {
		skillConfigInfo, err := GetFoeSkillConfigInfo(ctx, skillId)
		if err != nil {
			logger.CtxError(ctx, "GetMazeBattleData GetFoeSkillConfigInfo err", zap.Any("skillId", skillId), zap.Error(err))
			return nil, err
		}
		mazeBattleInfo.SkillConfigInfos = append(mazeBattleInfo.SkillConfigInfos, skillConfigInfo)
	}

	summonAttrMap := make(map[int32]int64)

	// 召唤物配置
	for _, summonId := range summonIds {
		summonCfg := GMazeSummonV8Cfg.GetWithCtx(ctx, summonId)
		if summonCfg == nil {
			logger.CtxError(ctx, "GetMazeBattleData GMazeSummonV8Cfg.Get fail", zap.Int32("summonId", summonId))
			return nil, fmt.Errorf("召唤物配置不存在")
		}

		var (
			skillIds = make([]int32, 0)
		)
		summonConfigInfo := &MazeAIBattle.MazeAISummonConfigInfo{
			SummonId: proto.Int32(summonId),
		}

		// 召唤物普通技能
		if summonCfg.Nor_attack_skill_id > 0 {
			skillIds = append(skillIds, summonCfg.Nor_attack_skill_id)
		}
		// // 召唤物技能
		for _, skillId := range summonCfg.Skill_id {
			if skillId <= 0 {
				continue
			}
			skillIds = append(skillIds, skillId)
		}
		for _, skillId := range skillIds {
			// 技能配置
			skillConfigInfo, err := GetFoeSkillConfigInfo(ctx, skillId)
			if err != nil {
				logger.CtxError(ctx, "GetMazeBattleData GetFoeSkillConfigInfo err", zap.Any("skillId", skillId), zap.Error(err))
				return nil, err
			}
			summonConfigInfo.SkillConfigInfos = append(summonConfigInfo.SkillConfigInfos, skillConfigInfo)
			// 技能信息
			skillInfo, err := GetFoeBattleSkillInfo(ctx, skillId, summonAttrMap)
			if err != nil {
				logger.CtxError(ctx, "GetMazeBattleData GetFoeBattleSkillInfo err", zap.Any("skillId", skillId), zap.Error(err))
				return nil, err
			}
			summonConfigInfo.SkillInfos = append(summonConfigInfo.SkillInfos, skillInfo)
		}
		mazeBattleInfo.SummonConfigInfo = append(mazeBattleInfo.SummonConfigInfo, summonConfigInfo)
	}

	// 道具使用配置
	attrItems := make(map[int32][]int32)
	for _, itemAttrCfg := range GMazeAttrItemAttrV8Cfg.GetAll() {
		attrItems[itemAttrCfg.Add_attr] = append(attrItems[itemAttrCfg.Add_attr], itemAttrCfg.Order)
	}
	for _, skillCfg := range GMazeSkillInfoV8Cfg.GetAll() {
		if skillCfg.Skill_attr_id <= 0 {
			continue
		}
		// 检查技能属性
		itemIDs, found := attrItems[skillCfg.Skill_attr_id]
		if found {
			for _, itemID := range itemIDs {
				itemUseInfo := &MazeAIBattle.MazeItemUseInfo{
					ItemId: proto.Int32(itemID),
				}
				itemUseInfo.SkillIds = append(itemUseInfo.SkillIds, skillCfg.Id)
				mazeBattleInfo.ItemUseInfos = append(mazeBattleInfo.ItemUseInfos, itemUseInfo)
			}
		}
	}

	barrierCfg := GMazeBarriesV8Cfg.GetWithCtx(ctx, barrierId)
	if barrierCfg == nil {
		logger.CtxError(ctx, "GetMazeBattleData GMazeBarriesV8Cfg fail", zap.Any("barrierId", barrierId))
		return nil, errors.CONFIG_NOT_FOUND
	}

	// 关卡掉落物品(宝箱结构物品)
	for _, boxID := range barrierCfg.Box_ids {
		if boxID <= 0 {
			continue
		}
		boxCfg := GMazeBoxV8Cfg.GetWithCtx(ctx, boxID)
		if boxCfg == nil {
			logger.CtxError(ctx, "GetMazeBattleData GMazeBoxV8Cfg fail", zap.Any("boxId", boxID), zap.Any("barrierId", barrierId))
			return nil, errors.CONFIG_NOT_FOUND
		}
		opened, err := mazeboxredis.IsOpenedBox(ctx, userId, barrierId, boxID)
		if err != nil {
			logger.CtxError(ctx, "GetMazeBattleData IsOpenedBox fail", zap.Error(err), zap.Any("boxId", boxID), zap.Any("barrierId", barrierId))
			return nil, err
		}

		var resID, resType int32

		if opened > 0 {
			resID = boxCfg.Res_id
			resType = boxCfg.Res_type
		} else {
			resID = boxCfg.Res_id_first
			resType = boxCfg.Res_type_first
		}
		mazeBattleInfo.DropItemBoxInfos = append(mazeBattleInfo.DropItemBoxInfos, &MazeAIBattle.MazeDropItemBoxInfo{
			BoxId:   proto.Int32(boxID),
			ResId:   proto.Int32(resID),
			ResType: proto.Int32(resType),
		})
	}

	return mazeBattleInfo, nil
}

func GetFoeAreaInfos(ctx context.Context, userId uint64, force int64, barrierId int32, userStiffRatio int64) ([]*MazeAIBattle.MazeAIAreaInfo, error) {
	logger := fklog.ContextAppLogger(ctx)
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
			if foeId > 0 {
				areaFoeMap[cfg.Brush_area_id][foeId] = struct{}{}
			}
		}
		for foeId, _ := range cfg.Monsterslist_ids_and_nums {
			if foeId > 0 {
				areaFoeMap[cfg.Brush_area_id][foeId] = struct{}{}
			}
		}
	}
	logger.CtxInfo(ctx, "GetFoeAreaInfos area foes dumps",
		zap.Int32("barrierId", barrierId),
		zap.Any("areaFoeMap", areaFoeMap),
	)
	for areaId, foeMap := range areaFoeMap {
		areaInfo := &MazeAIBattle.MazeAIAreaInfo{
			AreaId: proto.Int32(areaId),
		}
		areaInfo.MonsterConfigInfos = make([]*MazeAIBattle.MazeAIMonsterConfigInfo, 0)
		for foeId := range foeMap {
			monsterConfigInfo, err := GetMazeAIMonsterConfig(ctx, userId, force, foeId, userStiffRatio)
			if err != nil {
				logger.CtxError(ctx, "GetMazeBattleData GetMazeAIMonsterConfig err", zap.Any("foeId", foeId))
				return nil, err
			}
			logger.CtxInfo(ctx, "GetFoeAreaInfos GetMazeAIMonsterConfig dumps",
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
func GetMazeAIMonsterConfig(ctx context.Context, userId uint64, force int64, foeId int32, userStiffRatio int64) (*MazeAIBattle.MazeAIMonsterConfigInfo, error) {
	logger := fklog.ContextAppLogger(ctx)
	foeCfg := GMazeFoeV8Cfg.GetWithCtx(ctx, foeId)
	if foeCfg == nil {
		logger.CtxError(ctx, "GetMazeAIMonsterConfig GMazeFoeV8Cfg err", zap.Any("foeId", foeId))
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
	monsterConfigInfo.AttackedBackRange = proto.Int32(foeCfg.Attacked_back_range)
	monsterConfigInfo.AttackedBackRangeAfter = proto.Int32(foeCfg.Attacked_back_range_after)
	monsterConfigInfo.ThreatValue = proto.Int32(foeCfg.Threat_value)
	monsterConfigInfo.ArrowThreatValue = proto.Int32(foeCfg.Arrow_threat_value)
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
	// 寒冰元素抗性万分比
	monsterConfigInfo.AttrInfo = append(monsterConfigInfo.AttrInfo, &MazeAIBattle.MazeAIAttrInfo{
		Type:          proto.Int32(int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ICE_DEF_VALUE_PER)),
		UserValue:     proto.Int32(foeCfg.Ice_res),
		UserValueType: proto.Int32(2),
	})
	// 火焰元素抗性万分比
	monsterConfigInfo.AttrInfo = append(monsterConfigInfo.AttrInfo, &MazeAIBattle.MazeAIAttrInfo{
		Type:          proto.Int32(int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_FIRE_DEF_VALUE_PER)),
		UserValue:     proto.Int32(foeCfg.Fire_res),
		UserValueType: proto.Int32(2),
	})
	// 剧毒元素抗性万分比
	monsterConfigInfo.AttrInfo = append(monsterConfigInfo.AttrInfo, &MazeAIBattle.MazeAIAttrInfo{
		Type:          proto.Int32(int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_POISON_DEF_VALUE_PER)),
		UserValue:     proto.Int32(foeCfg.Poi_res),
		UserValueType: proto.Int32(2),
	})
	// 闪电元素抗性万分比
	monsterConfigInfo.AttrInfo = append(monsterConfigInfo.AttrInfo, &MazeAIBattle.MazeAIAttrInfo{
		Type:          proto.Int32(int32(MazeAIBattle.MAZE_AI_ATTR_TYPE_ELECTRICITY_DEF_VALUE_PER)),
		UserValue:     proto.Int32(foeCfg.Ele_res),
		UserValueType: proto.Int32(2),
	})
	skillIds := make([]int32, 0)
	if foeCfg.Nor_attack_skill_id > 0 {
		skillIds = append(skillIds, foeCfg.Nor_attack_skill_id)
	}
	if len(foeCfg.Passive_skill_id) > 0 {
		skillIds = append(skillIds, foeCfg.Passive_skill_id...)
	}
	skillTotalInfo := &MazeAIBattle.MazeAISkillTotalInfo{}
	skillTotalInfo.SkillInfoList = make([]*MazeAIBattle.MazeAISkillInfo, 0)
	attrMap := make(map[int32]int64)
	if len(skillIds) > 0 {
		skillCfg := GMazeSkillInfoV8Cfg.GetWithCtx(ctx, skillIds[0])
		if skillCfg == nil {
			logger.CtxError(ctx, "GetMazeAIMonsterConfig GMazeSkillInfoV8Cfg err", zap.Any("skillId", skillIds[0]))
			return nil, errors.New("配置不存在")
		}
		attackValue.SkillCd = GetSkillAttr(skillCfg.Skill_cool_time, attrMap)
	}
	for _, skillId := range skillIds {
		if skillId == 0 {
			continue
		}
		skillInfo, err := GetFoeBattleSkillInfo(ctx, skillId, attrMap)
		if err != nil {
			logger.CtxWarn(ctx, "GetUserBattleAttr BattleSkillTopPb nil", zap.Uint64("userId", userId), zap.Any("skillId", skillId))
			return nil, err
		}
		skillTotalInfo.SkillInfoList = append(skillTotalInfo.SkillInfoList, skillInfo)
		// attackValue.ActDamageConfig = append(attackValue.ActDamageConfig, actDamageConfigs...)
	}
	// 韧性被打空时释放技能
	if foeCfg.Tough_deplete > 0 {
		skillInfo, err := GetFoeBattleSkillInfo(ctx, foeCfg.Tough_deplete, attrMap)
		if err != nil {
			logger.CtxWarn(ctx, "GetUserBattleAttr BattleSkillTopPb nil", zap.Uint64("userId", userId), zap.Any("Tough_deplete", foeCfg.Tough_deplete))
			return nil, err
		}
		skillTotalInfo.SkillInfoList = append(skillTotalInfo.SkillInfoList, skillInfo)
	}

	monsterConfigInfo.MonsterSpeed = proto.Int32(foeCfg.Speed)
	monsterConfigInfo.SkillTotalInfo = skillTotalInfo
	return monsterConfigInfo, nil
}

func GetUserAttrInfo(ctx context.Context, userId uint64, userAttrMap map[int32]int64) (*MazeAIBattle.MazeAIRoleConfigInfo, []int32, error) {
	logger := fklog.ContextAppLogger(ctx)
	skillIds := make([]int32, 0)
	summonIds := make([]int32, 0)
	// 激活人物技能
	for _, cfg := range GMazeSkillInfoV8Cfg.GetAll() {
		if cfg.Skill_attr_id > 0 {
			value, ok := userAttrMap[cfg.Skill_attr_id]
			// 判断是否激活技能
			if ok && value > 0 {
				skillIds = append(skillIds, cfg.Id)
				// 追加召唤物技能
				if cfg.Summon_id > 0 {
					summonIds = append(summonIds, cfg.Summon_id)
				}
			}
		}
	}
	roleConfigInfo := &MazeAIBattle.MazeAIRoleConfigInfo{}
	userAttrInfo := &MazeAIBattle.MazeAIUserAttrInfo{}
	attrMap, err := GetUserBattleAttr(ctx, userId, userAttrMap)
	if err != nil {
		logger.CtxWarn(ctx, "GetUserBattleAttr BatchGetDollCalcAttr nil", zap.Uint64("userId", userId))
		return nil, nil, err
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
	for _, skillId := range skillIds {
		if skillId == 0 {
			continue
		}
		skillInfo, actDamageConfigs, err := GetUserBattleSkillInfo(ctx, skillId, userAttrMap)
		if err != nil {
			logger.CtxWarn(ctx, "GetUserBattleAttr GetUserBattleSkillInfo nil", zap.Uint64("userId", userId), zap.Any("skillId", skillId))
			return nil, nil, err
		}
		userSkillInfo.SkillInfoList = append(userSkillInfo.SkillInfoList, skillInfo)
		actDamageConfigList = append(actDamageConfigList, actDamageConfigs...)
	}
	// // 使用道具后可使用属性技能
	// for _, row := range GMazeAttrItemAttrV8Cfg.GetAll() {
	// 	if row.Add_attr <= 0 {
	// 		continue
	// 	}
	// 	attrSkill := GMazeAttrSkillV8Cfg.GetWithCtx(ctx,row.Add_attr)
	// 	if attrSkill == nil {
	// 		continue
	// 	}
	// 	if attrSkill.Skill_id <= 0 {
	// 		continue
	// 	}
	// 	skillInfo, actDamageConfigs, err := GetUserBattleSkillInfo(ctx, attrSkill.Skill_id, userAttrMap)
	// 	if err != nil {
	// 		logger.CtxWarn(ctx,"GetUserBattleAttr GetUserBattleSkillInfo nil", zap.Uint64("userId", userId), zap.Any("attrSkillId", attrSkill.Skill_id))
	// 		return nil, err
	// 	}
	// 	userSkillInfo.SkillInfoList = append(userSkillInfo.SkillInfoList, skillInfo)
	// 	actDamageConfigList = append(actDamageConfigList, actDamageConfigs...)
	// }

	roleConfigInfo.UserSkillInfo = userSkillInfo
	roleConfigInfo.ActDamageConfig = actDamageConfigList
	return roleConfigInfo, summonIds, nil
}

func GetUserBattleSkillInfo(ctx context.Context, skillId int32, attrMap map[int32]int64) (*MazeAIBattle.MazeAISkillInfo, []*MazeAIBattle.MazeAIActAttackValue, error) {
	logger := fklog.ContextAppLogger(ctx)
	skillCfg := GMazeSkillInfoV8Cfg.GetWithCtx(ctx, skillId)
	if skillCfg == nil {
		logger.CtxError(ctx, "GetUserBattleSkillInfo GMazeSkillInfoV8Cfg err", zap.Any("skillId", skillId))
		return nil, nil, errors.New("配置不存在")
	}
	var (
		actIDs           []int32
		effectID         int32
		actDamageConfigs = make([]*MazeAIBattle.MazeAIActAttackValue, 0)
	)
	skillActCfg := GMazeSkillActV8Cfg.GetWithCtx(ctx, skillId, config_manager.QueryNullable())
	if skillActCfg != nil {
		if len(skillActCfg.Act_id) > 0 {
			for _, actId := range skillActCfg.Act_id {
				if actId == 0 {
					continue
				}
				mazeActCfg := GMazeActInfoV8Cfg.GetWithCtx(ctx, actId)
				if mazeActCfg == nil {
					logger.CtxError(ctx, "GetUserBattleSkillInfo GMazeActInfoV8Cfg err", zap.Any("skillId", skillId), zap.Any("actId", actId))
					return nil, nil, errors.New("配置不存在")
				}
				actDamageConfig := &MazeAIBattle.MazeAIActAttackValue{
					ActId:           proto.Int32(actId),
					ToughBrokeValue: proto.Int32(mazeActCfg.Tough_broke_value),
					ToughTempValue:  proto.Int32(mazeActCfg.Temp_tough),
					// InitSkillCd:     proto.Int32(skillCfg.Initial_cool_time),
					// SkillCd:         proto.Int32(skillCfg.Skill_cool_time),
				}
				// 打击点索引:击退距离
				for index, value := range mazeActCfg.Attack_back_range {
					attackBackRange := &MazeAIBattle.MazeAttackBackRange{
						Index:           proto.Int32(index),
						AttackBackRange: proto.Int32(value),
					}
					actDamageConfig.AttackBackRange = append(actDamageConfig.AttackBackRange, attackBackRange)
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
				// 计算受击档位武力比系数
				for index, value := range mazeActCfg.Kongfu_hit_time_ratio {
					actDamageConfig.KongfuHitTimeRatio = append(actDamageConfig.KongfuHitTimeRatio, &MazeAIBattle.MazeAttackHitTimeRatio{
						Index:        proto.Int32(index),
						HitTimeRatio: proto.Int32(value),
					})
				}
			}
		}
		actIDs = skillActCfg.Act_id
		effectID = skillActCfg.Effect_id
	}
	// 技能基本信息
	skillInfo := &MazeAIBattle.MazeAISkillInfo{
		SkillId:                      proto.Int32(skillCfg.Id),
		RangeRadius:                  GetSkillAttr(skillCfg.Scope_param1, attrMap),
		ReleaseDistance:              GetSkillAttr(skillCfg.Distance_max, attrMap),
		TargetMaxCount:               GetSkillAttr(skillCfg.Target_num, attrMap),
		SkillDamageFixed:             GetSkillAttr(skillCfg.Main_target_damage_fix, attrMap),
		SkillMappingActionId:         actIDs,
		Level:                        GetSkillAttr(skillCfg.Level, attrMap),
		SkillType:                    proto.Int32(skillCfg.Type),
		SkillMappingEffectId:         proto.Int32(effectID),
		SecondTargetSkillDamageFixed: GetSkillAttr(skillCfg.Second_target_damage_fix, attrMap),
		IsNoTarget:                   proto.Int32(skillCfg.Is_no_target),
		DamageElement:                skillCfg.Damage_element,
		DamageType:                   proto.Int32(skillCfg.Damage_type),
		SkillCoolTime:                GetSkillAttr(skillCfg.Skill_cool_time, attrMap),
		DistanceMin:                  GetSkillAttr(skillCfg.Distance_min, attrMap),
		IsBreak:                      proto.Int32(skillCfg.Is_break),
		ScopeType:                    proto.Int32(skillCfg.Scope_type),
		TargetType:                   proto.Int32(skillCfg.Target_type),
		MainTargetDamageRates:        GetElementAttrValue(skillCfg.Main_target_damage, skillCfg.Damage_element_adjust, attrMap),
		SecondTargetDamageRates:      GetElementAttrValue(skillCfg.Second_target_damage, skillCfg.Damage_element_adjust, attrMap),
		Priority:                     proto.Int32(skillCfg.Priority),
		SummonId:                     proto.Int32(skillCfg.Summon_id),
		SummonNum:                    GetSkillAttr(skillCfg.Summon_num, attrMap),
		Duration:                     GetSkillAttr(skillCfg.Duration, attrMap),
		Interval:                     GetSkillAttr(skillCfg.Interval, attrMap),
		DamageAdjustment:             GetSkillAttr(skillCfg.Damage_adjustment, attrMap),
		TrajectoryNum:                GetSkillAttr(skillCfg.Trajectory_num, attrMap),
	}
	// 技能触发时机
	skillInfo.ReleaseTime = proto.Int32(skillCfg.Auto_release_time)

	// 技能触发条件
	conditionIDs, err := filterConditionIDs(skillCfg.Auto_release_condition)
	if err != nil {
		logger.CtxError(ctx, "GetUserBattleSkillInfo filterConditionIDs err", zap.Error(err), zap.Any("skillId", skillId))
		return nil, nil, errors.New("配置错误")
	}
	for _, conditionID := range conditionIDs {
		condition, err := GetMazeSkillConditionInfo(ctx, conditionID, attrMap)
		if err != nil {
			logger.CtxError(ctx, "GetUserBattleSkillInfo GetMazeSkillConditionInfo err", zap.Error(err), zap.Any("skillId", skillId), zap.Any("conditionID", conditionID))
			return nil, nil, err
		}
		skillInfo.ConditionConfig = append(skillInfo.ConditionConfig, condition)
	}
	skillInfo.ReleaseCondition = proto.String(skillCfg.Auto_release_condition)

	for _, effectId := range skillCfg.Target_effect {
		if effectId == 0 {
			continue
		}
		effectCfg := GMazeSkilleffectV8Cfg.GetWithCtx(ctx, effectId)
		if effectCfg == nil {
			logger.CtxError(ctx, "GetUserBattleSkillInfo GMazeSkilleffectV8Cfg err", zap.Any("effectId", effectId))
			return nil, nil, errors.New("配置不存在")
		}
		skillEffectOther := &MazeAIBattle.MazeAISkillEffectConfigInfo{
			EffectId:      proto.Int32(effectId),
			EffectGroup:   proto.Int32(effectCfg.Effect_group),
			InGroupWeight: proto.Int32(effectCfg.In_group_weight),
			LastTime:      proto.Int32(int32(GetEffectAttrValue(effectCfg.Last_time, effectCfg.Last_time_variable_id, attrMap))),
			BaseHitrate:   proto.Int32(int32(GetEffectAttrValue(effectCfg.Base_hitrate, effectCfg.Base_hitrate_variable_id, attrMap))),
			CoolDown:      proto.Int32(effectCfg.Cool_down),
			AttrId:        proto.Int32(effectCfg.Attr),
			Value_4:       effectCfg.Attr_value_4,
			MaxLayer:      proto.Int32(int32(GetEffectAttrValue(effectCfg.Attr_value_7, effectCfg.Attr_value_7_variable_id, attrMap))),
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
		for attrID, value := range effectCfg.Modify_attr_value {
			// 属性是加还是减
			op, ok := effectCfg.Modify_attr_value_variable_id[attrID]
			if !ok {
				logger.CtxError(ctx, "GetUserBattleSkillInfo Modify_attr_value_variable_id invalid",
					zap.Any("effectCfg", effectCfg), zap.Int32("attrID", attrID))
				continue
			}
			// 值类型
			valueType, ok := effectCfg.Modify_attr_value_type[attrID]
			if !ok {
				logger.CtxError(ctx, "GetUserBattleSkillInfo Modify_attr_value_type invalid",
					zap.Any("effectCfg", effectCfg), zap.Int32("attrID", attrID))
				continue
			}
			// 要加成的属性
			targetAttrID, ok := effectCfg.Modify_attr_value_attr_id[attrID]
			if !ok {
				logger.CtxError(ctx, "GetUserBattleSkillInfo Modify_attr_value_attr_id invalid",
					zap.Any("effectCfg", effectCfg), zap.Int32("attrID", attrID))
				continue
			}
			if op == 1 {
				value += int32(attrMap[attrID])
			} else if op == 2 {
				value -= int32(attrMap[attrID])
			}
			skillEffectOther.AttrModifier = append(skillEffectOther.AttrModifier, &MazeAIBattle.MazeAIAttrInfo{
				Type:          proto.Int32(targetAttrID),
				UserValue:     proto.Int32(value),
				UserValueType: proto.Int32(valueType),
			})
		}
		skillEffectOther.IntervalTime = proto.Int32(int32(GetEffectAttrValue(effectCfg.Attr_value_8, effectCfg.Attr_value_8_variable_id, attrMap)))
		skillInfo.SkillEffectOther = append(skillInfo.SkillEffectOther, skillEffectOther)
	}

	for _, effectId := range skillCfg.Self_effect {
		if effectId == 0 {
			continue
		}
		effectCfg := GMazeSkilleffectV8Cfg.GetWithCtx(ctx, effectId)
		if effectCfg == nil {
			logger.CtxError(ctx, "GetUserBattleSkillInfo GMazeSkilleffectV8Cfg err", zap.Any("effectId", effectId))
			return nil, nil, errors.New("配置不存在")
		}
		SkillEffectSelf := &MazeAIBattle.MazeAISkillEffectConfigInfo{
			EffectId:      proto.Int32(effectId),
			EffectGroup:   proto.Int32(effectCfg.Effect_group),
			InGroupWeight: proto.Int32(effectCfg.In_group_weight),
			LastTime:      proto.Int32(int32(GetEffectAttrValue(effectCfg.Last_time, effectCfg.Last_time_variable_id, attrMap))),
			BaseHitrate:   proto.Int32(int32(GetEffectAttrValue(effectCfg.Base_hitrate, effectCfg.Base_hitrate_variable_id, attrMap))),
			CoolDown:      proto.Int32(effectCfg.Cool_down),
			AttrId:        proto.Int32(effectCfg.Attr),
			Value_4:       effectCfg.Attr_value_4,
			MaxLayer:      proto.Int32(int32(GetEffectAttrValue(effectCfg.Attr_value_7, effectCfg.Attr_value_7_variable_id, attrMap))),
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
		for attrID, value := range effectCfg.Modify_attr_value {
			// 属性是加还是减
			op, ok := effectCfg.Modify_attr_value_variable_id[attrID]
			if !ok {
				logger.CtxError(ctx, "GetUserBattleSkillInfo Modify_attr_value_variable_id invalid",
					zap.Any("effectCfg", effectCfg), zap.Int32("attrID", attrID))
				continue
			}
			// 值类型
			valueType, ok := effectCfg.Modify_attr_value_type[attrID]
			if !ok {
				logger.CtxError(ctx, "GetUserBattleSkillInfo Modify_attr_value_type invalid",
					zap.Any("effectCfg", effectCfg), zap.Int32("attrID", attrID))
				continue
			}
			// 要加成的属性
			targetAttrID, ok := effectCfg.Modify_attr_value_attr_id[attrID]
			if !ok {
				logger.CtxError(ctx, "GetUserBattleSkillInfo Modify_attr_value_attr_id invalid",
					zap.Any("effectCfg", effectCfg), zap.Int32("attrID", attrID))
				continue
			}
			if op == 1 {
				value += int32(attrMap[attrID])
			} else if op == 2 {
				value -= int32(attrMap[attrID])
			}
			SkillEffectSelf.AttrModifier = append(SkillEffectSelf.AttrModifier, &MazeAIBattle.MazeAIAttrInfo{
				Type:          proto.Int32(targetAttrID),
				UserValue:     proto.Int32(value),
				UserValueType: proto.Int32(valueType),
			})
		}
		SkillEffectSelf.IntervalTime = proto.Int32(int32(GetEffectAttrValue(effectCfg.Attr_value_8, effectCfg.Attr_value_8_variable_id, attrMap)))
		skillInfo.SkillEffectSelf = append(skillInfo.SkillEffectSelf, SkillEffectSelf)
	}

	return skillInfo, actDamageConfigs, nil
}

func GetFoeBattleSkillInfo(ctx context.Context, skillId int32, attrMap map[int32]int64) (*MazeAIBattle.MazeAISkillInfo, error) {
	logger := fklog.ContextAppLogger(ctx)
	skillCfg := GMazeSkillInfoV8Cfg.GetWithCtx(ctx, skillId)
	if skillCfg == nil {
		logger.CtxError(ctx, "GetFoeBattleSkillInfo GMazeSkillInfoV8Cfg err", zap.Any("skillId", skillId))
		return nil, errors.New("配置不存在")
	}
	skillActCfg := GMazeSkillActV8Cfg.GetWithCtx(ctx, skillId)
	if skillActCfg == nil {
		logger.CtxError(ctx, "GetFoeBattleSkillInfo GMazeSkillActV8Cfg err", zap.Any("skillId", skillId))
		return nil, errors.New("配置不存在")
	}
	// todo 缺少触发cd
	skillInfo := &MazeAIBattle.MazeAISkillInfo{
		SkillId:                      proto.Int32(skillCfg.Id),
		RangeRadius:                  GetSkillAttr(skillCfg.Scope_param1, attrMap),
		ReleaseDistance:              GetSkillAttr(skillCfg.Distance_max, attrMap),
		TargetMaxCount:               GetSkillAttr(skillCfg.Target_num, attrMap),
		SkillDamageFixed:             GetSkillAttr(skillCfg.Main_target_damage_fix, attrMap),
		SkillMappingActionId:         skillActCfg.Act_id,
		Level:                        GetSkillAttr(skillCfg.Level, attrMap),
		SkillType:                    proto.Int32(skillCfg.Type),
		SkillMappingEffectId:         proto.Int32(skillActCfg.Effect_id),
		SecondTargetSkillDamageFixed: GetSkillAttr(skillCfg.Second_target_damage_fix, attrMap),
		IsNoTarget:                   proto.Int32(skillCfg.Is_no_target),
		DamageElement:                skillCfg.Damage_element,
		DamageType:                   proto.Int32(skillCfg.Damage_type),
		SkillCoolTime:                GetSkillAttr(skillCfg.Skill_cool_time, attrMap),
		DistanceMin:                  GetSkillAttr(skillCfg.Distance_min, attrMap),
		IsBreak:                      proto.Int32(skillCfg.Is_break),
		ScopeType:                    proto.Int32(skillCfg.Scope_type),
		TargetType:                   proto.Int32(skillCfg.Target_type),
		MainTargetDamageRates:        FillElementAttrValue(skillCfg.Main_target_damage, 5, attrMap),
		SecondTargetDamageRates:      FillElementAttrValue(skillCfg.Second_target_damage, 5, attrMap),
		Priority:                     proto.Int32(skillCfg.Priority),
		SummonId:                     proto.Int32(skillCfg.Summon_id),
		SummonNum:                    GetSkillAttr(skillCfg.Summon_num, attrMap),
		Duration:                     GetSkillAttr(skillCfg.Duration, attrMap),
		Interval:                     GetSkillAttr(skillCfg.Interval, attrMap),
		DamageAdjustment:             GetSkillAttr(skillCfg.Damage_adjustment, attrMap),
		TrajectoryNum:                GetSkillAttr(skillCfg.Trajectory_num, attrMap),
	}
	// 技能触发时机
	skillInfo.ReleaseTime = proto.Int32(skillCfg.Auto_release_time)

	// 技能触发条件
	conditionIDs, err := filterConditionIDs(skillCfg.Auto_release_condition)
	if err != nil {
		logger.CtxError(ctx, "GetFoeBattleSkillInfo filterConditionIDs err", zap.Error(err), zap.Any("skillId", skillId))
		return nil, errors.New("配置错误")
	}
	for _, conditionID := range conditionIDs {
		condition, err := GetMazeSkillConditionInfo(ctx, conditionID, attrMap)
		if err != nil {
			logger.CtxError(ctx, "GetFoeBattleSkillInfo GetMazeSkillConditionInfo err", zap.Error(err), zap.Any("skillId", skillId), zap.Any("conditionID", conditionID))
			return nil, err
		}
		skillInfo.ConditionConfig = append(skillInfo.ConditionConfig, condition)
	}
	skillInfo.ReleaseCondition = proto.String(skillCfg.Auto_release_condition)

	for _, effectId := range skillCfg.Target_effect {
		if effectId == 0 {
			continue
		}
		effectCfg := GMazeSkilleffectV8Cfg.GetWithCtx(ctx, effectId)
		if effectCfg == nil {
			logger.CtxError(ctx, "GetFoeBattleSkillInfo GMazeSkilleffectV8Cfg err", zap.Any("effectId", effectId))
			return nil, errors.New("配置不存在")
		}
		skillEffectOther := &MazeAIBattle.MazeAISkillEffectConfigInfo{
			EffectId:      proto.Int32(effectId),
			EffectGroup:   proto.Int32(effectCfg.Effect_group),
			InGroupWeight: proto.Int32(effectCfg.In_group_weight),
			LastTime:      proto.Int32(int32(GetEffectAttrValue(effectCfg.Last_time, effectCfg.Last_time_variable_id, attrMap))),
			BaseHitrate:   proto.Int32(int32(GetEffectAttrValue(effectCfg.Base_hitrate, effectCfg.Base_hitrate_variable_id, attrMap))),
			CoolDown:      proto.Int32(effectCfg.Cool_down),
			AttrId:        proto.Int32(effectCfg.Attr),
			Value_4:       effectCfg.Attr_value_4,
			MaxLayer:      proto.Int32(int32(GetEffectAttrValue(effectCfg.Attr_value_7, effectCfg.Attr_value_7_variable_id, attrMap))),
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
		for attrID, value := range effectCfg.Modify_attr_value {
			// 属性是加还是减
			op, ok := effectCfg.Modify_attr_value_variable_id[attrID]
			if !ok {
				logger.CtxError(ctx, "GetFoeBattleSkillInfo Modify_attr_value_variable_id invalid",
					zap.Any("effectCfg", effectCfg), zap.Int32("attrID", attrID))
				continue
			}
			// 值类型
			valueType, ok := effectCfg.Modify_attr_value_type[attrID]
			if !ok {
				logger.CtxError(ctx, "GetFoeBattleSkillInfo Modify_attr_value_type invalid",
					zap.Any("effectCfg", effectCfg), zap.Int32("attrID", attrID))
				continue
			}
			// 要加成的属性
			targetAttrID, ok := effectCfg.Modify_attr_value_attr_id[attrID]
			if !ok {
				logger.CtxError(ctx, "GetFoeBattleSkillInfo Modify_attr_value_attr_id invalid",
					zap.Any("effectCfg", effectCfg), zap.Int32("attrID", attrID))
				continue
			}
			if op == 1 {
				value += int32(attrMap[attrID])
			} else if op == 2 {
				value -= int32(attrMap[attrID])
			}
			skillEffectOther.AttrModifier = append(skillEffectOther.AttrModifier, &MazeAIBattle.MazeAIAttrInfo{
				Type:          proto.Int32(targetAttrID),
				UserValue:     proto.Int32(value),
				UserValueType: proto.Int32(valueType),
			})
		}
		skillEffectOther.IntervalTime = proto.Int32(int32(GetEffectAttrValue(effectCfg.Attr_value_8, effectCfg.Attr_value_8_variable_id, attrMap)))
		skillInfo.SkillEffectOther = append(skillInfo.SkillEffectOther, skillEffectOther)
	}

	for _, effectId := range skillCfg.Self_effect {
		if effectId == 0 {
			continue
		}
		effectCfg := GMazeSkilleffectV8Cfg.GetWithCtx(ctx, effectId)
		if effectCfg == nil {
			logger.CtxError(ctx, "GetFoeBattleSkillInfo GMazeSkilleffectV8Cfg err", zap.Any("effectId", effectId))
			return nil, errors.New("配置不存在")
		}
		SkillEffectSelf := &MazeAIBattle.MazeAISkillEffectConfigInfo{
			EffectId:      proto.Int32(effectId),
			EffectGroup:   proto.Int32(effectCfg.Effect_group),
			InGroupWeight: proto.Int32(effectCfg.In_group_weight),
			LastTime:      proto.Int32(int32(GetEffectAttrValue(effectCfg.Last_time, effectCfg.Last_time_variable_id, attrMap))),
			BaseHitrate:   proto.Int32(int32(GetEffectAttrValue(effectCfg.Base_hitrate, effectCfg.Base_hitrate_variable_id, attrMap))),
			CoolDown:      proto.Int32(effectCfg.Cool_down),
			AttrId:        proto.Int32(effectCfg.Attr),
			Value_4:       effectCfg.Attr_value_4,
			MaxLayer:      proto.Int32(int32(GetEffectAttrValue(effectCfg.Attr_value_7, effectCfg.Attr_value_7_variable_id, attrMap))),
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
		for attrID, value := range effectCfg.Modify_attr_value {
			// 属性是加还是减
			op, ok := effectCfg.Modify_attr_value_variable_id[attrID]
			if !ok {
				logger.CtxError(ctx, "GetFoeBattleSkillInfo Modify_attr_value_variable_id invalid",
					zap.Any("effectCfg", effectCfg), zap.Int32("attrID", attrID))
				continue
			}
			// 值类型
			valueType, ok := effectCfg.Modify_attr_value_type[attrID]
			if !ok {
				logger.CtxError(ctx, "GetFoeBattleSkillInfo Modify_attr_value_type invalid",
					zap.Any("effectCfg", effectCfg), zap.Int32("attrID", attrID))
				continue
			}
			// 要加成的属性
			targetAttrID, ok := effectCfg.Modify_attr_value_attr_id[attrID]
			if !ok {
				logger.CtxError(ctx, "GetFoeBattleSkillInfo Modify_attr_value_attr_id invalid",
					zap.Any("effectCfg", effectCfg), zap.Int32("attrID", attrID))
				continue
			}
			if op == 1 {
				value += int32(attrMap[attrID])
			} else if op == 2 {
				value -= int32(attrMap[attrID])
			}
			SkillEffectSelf.AttrModifier = append(SkillEffectSelf.AttrModifier, &MazeAIBattle.MazeAIAttrInfo{
				Type:          proto.Int32(targetAttrID),
				UserValue:     proto.Int32(value),
				UserValueType: proto.Int32(valueType),
			})
		}
		SkillEffectSelf.IntervalTime = proto.Int32(int32(GetEffectAttrValue(effectCfg.Attr_value_8, effectCfg.Attr_value_8_variable_id, attrMap)))
		skillInfo.SkillEffectSelf = append(skillInfo.SkillEffectSelf, SkillEffectSelf)
	}

	return skillInfo, nil
}

func GetFoeSkillConfigInfo(ctx context.Context, skillId int32) (*MazeAIBattle.MazeSkillConfigInfo, error) {
	logger := fklog.ContextAppLogger(ctx)
	skillCfg := GMazeSkillInfoV8Cfg.GetWithCtx(ctx, skillId)
	if skillCfg == nil {
		logger.CtxError(ctx, "GetMazeAIMonsterConfig GMazeSkillInfoV8Cfg err", zap.Any("skillId", skillId))
		return nil, errors.New("配置不存在")
	}
	skillActCfg := GMazeSkillActV8Cfg.GetWithCtx(ctx, skillId)
	if skillActCfg == nil {
		logger.CtxError(ctx, "GetMazeAIMonsterConfig GMazeSkillActV8Cfg err", zap.Any("skillId", skillId))
		return nil, errors.New("配置不存在")
	}
	skillConfigInfo := &MazeAIBattle.MazeSkillConfigInfo{
		SkillId: proto.Int32(skillId),
		SkillCd: GetSkillAttr(skillCfg.Skill_cool_time, map[int32]int64{}),
	}
	actDamageConfigs := make([]*MazeAIBattle.MazeAIActAttackValue, 0)
	if len(skillActCfg.Act_id) > 0 {
		for _, actId := range skillActCfg.Act_id {
			if actId == 0 {
				continue
			}
			mazeActCfg := GMazeActInfoV8Cfg.GetWithCtx(ctx, actId)
			if mazeActCfg == nil {
				logger.CtxError(ctx, "GetMazeAIMonsterConfig GMazeActInfoV8Cfg err", zap.Any("skillId", skillId), zap.Any("actId", actId))
				return nil, errors.New("配置不存在")
			}
			actDamageConfig := &MazeAIBattle.MazeAIActAttackValue{
				ActId:           proto.Int32(actId),
				ToughBrokeValue: proto.Int32(mazeActCfg.Tough_broke_value),
				ToughTempValue:  proto.Int32(mazeActCfg.Temp_tough),
			}
			// 打击点索引:击退距离
			for index, value := range mazeActCfg.Attack_back_range {
				attackBackRange := &MazeAIBattle.MazeAttackBackRange{
					Index:           proto.Int32(index),
					AttackBackRange: proto.Int32(value),
				}
				actDamageConfig.AttackBackRange = append(actDamageConfig.AttackBackRange, attackBackRange)
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
			// 计算受击档位武力比系数
			for index, value := range mazeActCfg.Kongfu_hit_time_ratio {
				actDamageConfig.KongfuHitTimeRatio = append(actDamageConfig.KongfuHitTimeRatio, &MazeAIBattle.MazeAttackHitTimeRatio{
					Index:        proto.Int32(index),
					HitTimeRatio: proto.Int32(value),
				})
			}
		}
	}
	skillConfigInfo.ActDamageConfig = actDamageConfigs
	return skillConfigInfo, nil
}

// func GetMazeAIAutoSkillInfo(ctx context.Context, skillId int32, attrMap map[int32]int64) (*MazeAIBattle.MazeAIAutoSkillInfo, error) {
// 	skillAutoCfg := GMazeSkillAutoReleaseV8Cfg.GetWithCtx(ctx,skillId)
// 	if skillAutoCfg == nil {
// 		logger.CtxError(ctx,"GetMazeAIAutoSkillInfo GMazeSkillAutoReleaseV8Cfg err", zap.Any("skillId", skillId))
// 		return nil, errors.New("配置不存在")
// 	}
// 	skillConfigInfo := &MazeAIBattle.MazeAIAutoSkillInfo{
// 		SkillId:      proto.Int32(skillId),
// 		BaseHitrate:  proto.Int32(int32(GetAttrValue(skillAutoCfg.Release_ratio, skillAutoCfg.Release_ratio_variable, attrMap))),
// 		TriggerCount: proto.Int32(int32(GetAttrValue(skillAutoCfg.Release_num, skillAutoCfg.Release_num_variable, attrMap))),
// 		ProtectCd:    proto.Int32(int32(GetAttrValue(skillAutoCfg.Auto_release_protect_cd, skillAutoCfg.Auto_release_protect_cd_variable, attrMap))),
// 		LimitCount:   proto.Int32(int32(GetAttrValue(skillAutoCfg.Max_release_limit, skillAutoCfg.Max_release_limit_variable, attrMap))),
// 		TriggerType:  proto.Int32(skillAutoCfg.Release_time),
// 		Condition:    proto.String(skillAutoCfg.Release_condition),
// 	}
// 	// 触发技能列表
// 	for _, triggerSkillId := range skillAutoCfg.Skill_id {
// 		skillConfigInfo.TriggerSkillId = append(skillConfigInfo.TriggerSkillId, triggerSkillId)
// 	}
// 	// 技能条件
// 	conditionIDs, err := filterConditionIDs(skillAutoCfg.Release_condition)
// 	if err != nil {
// 		logger.CtxError(ctx,"GetMazeAIAutoSkillInfo filterConditionIDs err", zap.Error(err), zap.Any("skillId", skillId))
// 		return nil, errors.New("配置错误")
// 	}
// 	for _, conditionID := range conditionIDs {
// 		condition, err := GetMazeSkillConditionInfo(ctx, conditionID, attrMap)
// 		if err != nil {
// 			logger.CtxError(ctx,"GetMazeAIAutoSkillInfo GetMazeSkillConditionInfo err", zap.Error(err), zap.Any("skillId", skillId), zap.Any("conditionID", conditionID))
// 			return nil, err
// 		}
// 		skillConfigInfo.ConditionConfig = append(skillConfigInfo.ConditionConfig, condition)
// 	}
// 	return skillConfigInfo, nil
// }

func GetMazeSkillConditionInfo(ctx context.Context, conditionID int32, attrMap map[int32]int64) (*MazeAIBattle.MazeSkillCondition, error) {
	logger := fklog.ContextAppLogger(ctx)
	cfg := GMazeSkillAutoConditionV8Cfg.GetWithCtx(ctx, conditionID)
	if cfg == nil {
		logger.CtxError(ctx, "GetMazeSkillConditionInfo GMazeSkillAutoConditionV8Cfg err", zap.Any("conditionID", conditionID))
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

// GetSkillAttrValue 根据技能属性配置与用户属性列表计算技能最终属性
func GetSkillAttrValue(skillAttrMap map[int32]int32, userAttrMap map[int32]int64) (attrID int32, value int32) {
	// 取人物属性值
	attrID, ok := skillAttrMap[1]
	if ok {
		value, ok := userAttrMap[attrID]
		if ok {
			return attrID, int32(value)
		}
		return attrID, 0
	}
	// 取默认值
	defaultValue, ok := skillAttrMap[0]
	if ok {
		value = defaultValue
	}
	return
}

// GetSkillAttr 根据技能属性配置与用户属性列表计算技能最终属性
func GetSkillAttr(skillAttrMap map[int32]int32, userAttrMap map[int32]int64) (attr *MazeAIBattle.MazeAIAttrInfo) {
	attrID, value := GetSkillAttrValue(skillAttrMap, userAttrMap)
	return &MazeAIBattle.MazeAIAttrInfo{Type: proto.Int32(attrID), UserValue: proto.Int32(value)}
}
