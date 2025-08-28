package equipbaggm

import (
	"context"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

func Reg(logger fklog.FKLogI) {
	// gm.SafeHttpRegister(logger, "/AddEquip")

	// gm.SafeHttpRegister(logger, "/ClearBag")

	// gm.SafeHttpRegister(logger, "/BatchClearBag")

	// gm.SafeHttpRegister(logger, "/ClearBagByMap", func(writer http.ResponseWriter, request *http.Request) {
	// 	// 外网线上环境不允许使用GM
	// 	request.ParseForm()
	//
	// 	mapId := fkutil.ToUint64(request.Form.Get("mapId"))
	//
	// 	leagueIDMap, err := WorldLeagueRedis.GetWorldLeagueInfo(logger, mapId)
	// 	if err != nil {
	// 		logger.ErrorWF("ClearBagByMap load world leagueInfo fail",
	// 			zap.Uint64("mapID", mapId),
	// 			zap.Error(err))
	// 		return
	// 	}
	// 	logger.InfoWF("ClearBagByMap map info",
	// 		zap.Uint64("map", mapId),
	// 		zap.Int("leagueLen", len(leagueIDMap)),
	// 	)
	//
	// 	var familyID uint64
	// 	for leagueID := range leagueIDMap {
	// 		// 取联盟下的散人家族
	// 		familyIDs, err := LeagueFamilyRedis.GetAllLeagueFamilyIDs(logger, leagueID)
	// 		if err != nil {
	// 			logger.ErrorWF("ClearBagByMap get league familyIDs fail", zap.Any("leagueID", leagueID), zap.Error(err))
	// 			continue
	// 		}
	//
	// 		for _, family := range familyIDs {
	// 			familyID = fkutil.ToUint64(family)
	// 			// 取家族下所有人
	// 			users, err := FamilyAllocUserRedis.GetAllFamilyUIDSliceFix(logger, familyID)
	// 			if err != nil {
	// 				logger.ErrorWF("ClearBagByMap get family users fail", zap.Error(err))
	// 				continue
	// 			}
	//
	// 			if len(users) == 0 {
	// 				continue
	// 			}
	//
	// 			for _, uid := range users {
	// 				err := ClearUserBag(logger, uid)
	// 				if err != nil {
	// 					logger.ErrorWF("ClearBagByMap ClearUserBag fail", zap.Error(err), zap.Uint64("uid", uid))
	// 					continue
	// 				}
	// 			}
	// 		}
	// 	}
	//
	// 	writer.Write([]byte("ok"))
	//
	// 	return
	// })

	// gm.SafeHttpRegister(logger, "/CheckDollEquipCfg", func(writer http.ResponseWriter, request *http.Request) {
	// 	// 外网线上环境不允许使用GM
	// 	request.ParseForm()

	// for _, cfg := range GDollEquipInfoV8Cfg.GetAll() {
	// 	//校验基础属性
	// 	var maxBaseNum int32
	// 	for baseNum := range cfg.Affix_base_num {
	// 		if baseNum > maxBaseNum {
	// 			maxBaseNum = baseNum
	// 		}
	// 	}
	// 	groupBasePoolMap := make(map[int32]map[int32]int32, 0)
	// 	for poolId, index := range cfg.Affix_base_pool {
	// 		if index == 0 {
	// 			continue
	// 		}
	// 		if groupBasePoolMap[index] == nil {
	// 			groupBasePoolMap[index] = make(map[int32]int32, 0)
	// 		}
	// 		groupBasePoolMap[index][poolId] = 1000
	// 	}
	// 	for index := int32(1); index <= maxBaseNum; index++ {
	// 		basePoolMap, ok := groupBasePoolMap[index]
	// 		if !ok {
	// 			logger.ErrorWF("doll_equip_info_v8【人偶-装备-信息】.xlsx 基础属性条数缺失 ", zap.Any("装备id：", cfg.Equipment_id), zap.Any("缺失条数位置：", index))
	// 			continue
	// 		}
	// 		for poolId := range basePoolMap {
	// 			poolMap := mazeequipaffixrandpoolv8.GetEquipPoolWeightCfg(poolId)
	// 			if len(poolMap) == 0 {
	// 				logger.ErrorWF("doll_equip_affix_pool_v8【人偶-装备-词条随机池】.xlsx 基础词条库缺失 ", zap.Any("装备id:", cfg.Equipment_id), zap.Any("池子id:", poolId))
	// 				continue
	// 			}
	// 		}
	// 	}

	// 	//校验随机属性
	// 	var maxRandNum int32
	// 	for randNum := range cfg.Affix_rand_num {
	// 		if randNum > maxRandNum {
	// 			maxRandNum = randNum
	// 		}
	// 	}
	// 	for poolId := range cfg.Affix_rand_pool {
	// 		if poolId == 0 {
	// 			continue
	// 		}
	// 		poolMap := mazeequipaffixrandpoolv8.GetEquipPoolWeightCfg(poolId)
	// 		if len(poolMap) == 0 {
	// 			logger.ErrorWF("doll_equip_affix_pool_v8【人偶-装备-词条随机池】.xlsx 随机词条库缺失 ", zap.Any("装备id:", cfg.Equipment_id), zap.Any("池子id:", poolId))
	// 			continue
	// 		}
	// 		if len(poolMap) < int(maxRandNum) {
	// 			logger.ErrorWF("doll_equip_affix_pool_v8【人偶-装备-词条随机池】.xlsx 随机词条库条数不足", zap.Any("装备id:", cfg.Equipment_id), zap.Any("池子id:", poolId),
	// 				zap.Any("当前词条数量:", len(poolMap)), zap.Any("最大词条数量", maxRandNum))
	// 			continue
	// 		}
	// 		poolGroupMap := mazeequipaffixrandpoolv8.GetPoolGroupCfg(poolId)
	// 		if len(poolGroupMap) < int(maxRandNum) {
	// 			logger.ErrorWF("doll_equip_affix_pool_v8【人偶-装备-词条随机池】.xlsx 随机词条库条去重组不足", zap.Any("装备id:", cfg.Equipment_id), zap.Any("池子id:", poolId),
	// 				zap.Any("当前词条去重组数量:", len(poolGroupMap)), zap.Any("最大词条数量", maxRandNum))
	// 			continue
	// 		}
	// 	}

	// 	//校验天赋属性
	// 	var maxModNum int32
	// 	for modNum := range cfg.Affix_mod_num {
	// 		if modNum > maxModNum {
	// 			maxModNum = modNum
	// 		}
	// 	}
	// 	for poolId := range cfg.Affix_mod_pool {
	// 		if poolId == 0 {
	// 			continue
	// 		}
	// 		poolWeightMap := GDollEquipAffixModPoolV8CfgEx.GetEnchantPoolWeightCfg(poolId)
	// 		if len(poolWeightMap) == 0 {
	// 			logger.ErrorWF("doll_equip_affix_mod_pool_v8【人偶-装备-附魔词条随机池】.xlsx 随机词条库缺失 ", zap.Any("装备id:", cfg.Equipment_id), zap.Any("池子id:", poolId))
	// 			continue
	// 		}
	// 		if len(poolWeightMap) < int(maxModNum) {
	// 			logger.ErrorWF("doll_equip_affix_mod_pool_v8【人偶-装备-附魔词条随机池】.xlsx 随机词条库条数不足", zap.Any("装备id:", cfg.Equipment_id), zap.Any("池子id:", poolId),
	// 				zap.Any("当前词条数量:", len(poolWeightMap)), zap.Any("最大词条数量", maxModNum))
	// 			continue
	// 		}
	// 		poolGroupMap := GDollEquipAffixModPoolV8CfgEx.GetEnchantPoolGroupCfg(poolId)
	// 		if len(poolGroupMap) < int(maxModNum) {
	// 			logger.ErrorWF("doll_equip_affix_mod_pool_v8【人偶-装备-附魔词条随机池】.xlsx 随机词条库条去重组不足", zap.Any("装备id:", cfg.Equipment_id), zap.Any("池子id:", poolId),
	// 				zap.Any("当前词条去重组数量:", len(poolGroupMap)), zap.Any("最大词条数量", maxModNum))
	// 			continue
	// 		}
	// 	}

	// 	//校验特殊属性
	// 	for poolId := range cfg.Affix_mod_pool {
	// 		if poolId == 0 {
	// 			continue
	// 		}
	// 		poolMap := GDollEquipAffixSpPoolV8CfgEx.GetSpPoolWeightCfg(poolId)
	// 		if len(poolMap) == 0 {
	// 			logger.ErrorWF("doll_equip_affix_sp_pool_v8【人偶-装备-特殊词条随机库】.xlsx 随机词条库缺失 ", zap.Any("装备id:", cfg.Equipment_id), zap.Any("池子id:", poolId))
	// 			continue
	// 		}
	// 	}

	// }

	// for _, cfg := range GDollEquipAffixSpRuleV8Cfg.GetAll() {
	// 	//校验基础属性
	// 	var maxBaseNum int32
	// 	for baseNum := range cfg.Affix_base_num {
	// 		if baseNum > maxBaseNum {
	// 			maxBaseNum = baseNum
	// 		}
	// 	}
	// 	groupBasePoolMap := make(map[int32]map[int32]int32, 0)
	// 	for poolId, index := range cfg.Affix_base_pool {
	// 		if index == 0 {
	// 			continue
	// 		}
	// 		if groupBasePoolMap[index] == nil {
	// 			groupBasePoolMap[index] = make(map[int32]int32, 0)
	// 		}
	// 		groupBasePoolMap[index][poolId] = 1000
	// 	}
	// 	for index := int32(1); index <= maxBaseNum; index++ {
	// 		basePoolMap, ok := groupBasePoolMap[index]
	// 		if !ok {
	// 			logger.ErrorWF("doll_equip_affix_sp_rule_v8【人偶-装备-生成特殊词条规则】.xlsx 基础属性条数缺失 ", zap.Any("装备规则id：", cfg.Rule_id), zap.Any("缺失条数位置：", index))
	// 			continue
	// 		}
	// 		for poolId := range basePoolMap {
	// 			poolMap := mazeequipaffixrandpoolv8.GetEquipPoolWeightCfg(poolId)
	// 			if len(poolMap) == 0 {
	// 				logger.ErrorWF("doll_equip_affix_pool_v8【人偶-装备-词条随机池】.xlsx 基础词条库缺失 ", zap.Any("装备规则id:", cfg.Rule_id), zap.Any("池子id:", poolId))
	// 				continue
	// 			}
	// 		}
	// 	}

	// 	//校验随机属性
	// 	var maxRandNum int32
	// 	for randNum := range cfg.Affix_rand_num {
	// 		if randNum > maxRandNum {
	// 			maxRandNum = randNum
	// 		}
	// 	}
	// 	for poolId := range cfg.Affix_rand_pool {
	// 		if poolId == 0 {
	// 			continue
	// 		}
	// 		poolMap := mazeequipaffixrandpoolv8.GetEquipPoolWeightCfg(poolId)
	// 		if len(poolMap) == 0 {
	// 			logger.ErrorWF("doll_equip_affix_pool_v8【人偶-装备-词条随机池】.xlsx 随机词条库缺失 ", zap.Any("装备规则id:", cfg.Rule_id), zap.Any("池子id:", poolId))
	// 			continue
	// 		}
	// 		if len(poolMap) < int(maxRandNum) {
	// 			logger.ErrorWF("doll_equip_affix_pool_v8【人偶-装备-词条随机池】.xlsx 随机词条库条数不足", zap.Any("装备规则id:", cfg.Rule_id), zap.Any("池子id:", poolId),
	// 				zap.Any("当前词条数量:", len(poolMap)), zap.Any("最大词条数量", maxRandNum))
	// 			continue
	// 		}
	// 		poolGroupMap := mazeequipaffixrandpoolv8.GetPoolGroupCfg(poolId)
	// 		if len(poolGroupMap) < int(maxRandNum) {
	// 			logger.ErrorWF("doll_equip_affix_pool_v8【人偶-装备-词条随机池】.xlsx 随机词条库条去重组不足", zap.Any("装备规则id:", cfg.Rule_id), zap.Any("池子id:", poolId),
	// 				zap.Any("当前词条去重组数量:", len(poolGroupMap)), zap.Any("最大词条数量", maxRandNum))
	// 			continue
	// 		}
	// 	}

	// 	//校验天赋属性
	// 	var maxModNum int32
	// 	for modNum := range cfg.Affix_mod_num {
	// 		if modNum > maxModNum {
	// 			maxModNum = modNum
	// 		}
	// 	}
	// 	for poolId := range cfg.Affix_mod_pool {
	// 		if poolId == 0 {
	// 			continue
	// 		}
	// 		poolWeightMap := GDollEquipAffixModPoolV8CfgEx.GetEnchantPoolWeightCfg(poolId)
	// 		if len(poolWeightMap) == 0 {
	// 			logger.ErrorWF("doll_equip_affix_mod_pool_v8【人偶-装备-附魔词条随机池】.xlsx 随机词条库缺失 ", zap.Any("装备规则id:", cfg.Rule_id), zap.Any("池子id:", poolId))
	// 			continue
	// 		}
	// 		if len(poolWeightMap) < int(maxModNum) {
	// 			logger.ErrorWF("doll_equip_affix_mod_pool_v8【人偶-装备-附魔词条随机池】.xlsx 随机词条库条数不足", zap.Any("装备规则id:", cfg.Rule_id), zap.Any("池子id:", poolId),
	// 				zap.Any("当前词条数量:", len(poolWeightMap)), zap.Any("最大词条数量", maxModNum))
	// 			continue
	// 		}
	// 		poolGroupMap := GDollEquipAffixModPoolV8CfgEx.GetEnchantPoolGroupCfg(poolId)
	// 		if len(poolGroupMap) < int(maxModNum) {
	// 			logger.ErrorWF("doll_equip_affix_mod_pool_v8【人偶-装备-附魔词条随机池】.xlsx 随机词条库条去重组不足", zap.Any("装备规则id:", cfg.Rule_id), zap.Any("池子id:", poolId),
	// 				zap.Any("当前词条去重组数量:", len(poolGroupMap)), zap.Any("最大词条数量", maxModNum))
	// 			continue
	// 		}
	// 	}

	// 	//校验特殊属性
	// 	for poolId := range cfg.Affix_extra_pool {
	// 		if poolId == 0 {
	// 			continue
	// 		}
	// 		poolMap := GDollEquipAffixSpPoolV8CfgEx.GetSpPoolWeightCfg(poolId)
	// 		if len(poolMap) == 0 {
	// 			logger.ErrorWF("doll_equip_affix_sp_pool_v8【人偶-装备-特殊词条随机库】.xlsx 随机词条库缺失 ", zap.Any("装备规则id:", cfg.Rule_id), zap.Any("池子id:", poolId))
	// 			continue
	// 		}
	// 	}

	// }
	// 	writer.Write([]byte("ok"))

	// 	return
	// })

	// gm.SafeHttpRegister(logger, "/FixDelModExpEquipByMap", func(writer http.ResponseWriter, request *http.Request) {
	// 	request.ParseForm()
	// 	mapId := fkutil.ToUint64(request.Form.Get("mapId"))
	// 	run := fkutil.ToUint64(request.Form.Get("run"))
	// 	leagueIDMap, err := WorldLeagueRedis.GetWorldLeagueInfo(logger, mapId)
	// 	if err != nil {
	// 		logger.ErrorWF("FixDelModExpEquipByMap load world leagueInfo fail",
	// 			zap.Uint64("mapID", mapId),
	// 			zap.Error(err))
	// 		return
	// 	}
	// 	logger.InfoWF("FixDelModExpEquipByMap map info",
	// 		zap.Uint64("map", mapId),
	// 		zap.Int("leagueLen", len(leagueIDMap)),
	// 	)

	// 	var familyID uint64
	// 	var cnt int
	// 	for leagueID := range leagueIDMap {
	// 		// 取联盟下的散人家族
	// 		familyIDs, err := LeagueFamilyRedis.GetAllLeagueFamilyIDs(logger, leagueID)
	// 		if err != nil {
	// 			logger.ErrorWF("FixDelModExpEquipByMap get league familyIDs fail", zap.Any("leagueID", leagueID), zap.Error(err))
	// 			continue
	// 		}

	// 		for _, family := range familyIDs {
	// 			familyID = fkutil.ToUint64(family)
	// 			// 取家族下所有人
	// 			users, err := FamilyAllocUserRedis.GetAllFamilyUIDSliceFix(logger, familyID)
	// 			if err != nil {
	// 				logger.ErrorWF("FixDelModExpEquipByMap get family users fail", zap.Error(err))
	// 				continue
	// 			}

	// 			if len(users) == 0 {
	// 				continue
	// 			}

	// 			for _, uid := range users {
	// 				logger.SetUid(uid)
	// 				n, _ := GmDelExpModAttrBag(logger, uid, run == 1)
	// 				if n {
	// 					cnt++
	// 				}
	// 			}
	// 		}
	// 	}

	// 	writer.Write([]byte(fmt.Sprintf("需要修复数量:%d", cnt)))
	// })

	// gm.SafeHttpRegister(logger, "/FixDelModExpEquip", func(writer http.ResponseWriter, request *http.Request) {
	// 	userId := fkutil.ToUint64(request.Form.Get("userId"))
	// 	logger.SetUid(userId)
	// 	_, e := GmDelExpModAttrBag(logger, userId, true)
	// 	if e == nil {
	// 		writer.Write([]byte("ok"))
	// 	} else {
	// 		writer.Write([]byte("fail"))
	// 	}
	// })

	// gm.SafeHttpRegister(logger, "/SetEquipRollScore")

	// gm.SafeHttpRegister(logger, "/ClearBagNotAssemble")

	// gm.SafeHttpRegister(logger, "/BatchAddEquip")

	// gm.SafeHttpRegister(logger, "/fixAllEquipAttrLimit")
	// gm.SafeHttpRegister(logger, "/resetInsAllEquip", func(writer http.ResponseWriter, request *http.Request) {
	// 	userId := fkutil.ToUint64(request.Form.Get("userId"))
	// 	logger.SetUid(userId)
	// 	e := resetInsAllEquip(logger, userId)
	// 	if e == nil {
	// 		writer.Write([]byte("ok"))
	// 	} else {
	// 		writer.Write([]byte("fail"))
	// 	}
	// })
}

func ClearUserBag(ctx context.Context, userId uint64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	// todo 找装配的删除
	err = ClearDollAssembleInfo(ctx, userId)
	if err != nil {
		return
	}
	// todo 删除背包里的
	err = ClearEquipBag(logger, userId)
	if err != nil {
		return
	}

	return
}
