package barrierstagecounterservice

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/common/function/flowutil"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/config/GMazeBrushFoeV8Cfg"
	"maze_game_server/config/GMazeFoeV8Cfg"
	"maze_game_server/excel/dollmappuzzlenewcfgex"
	"maze_game_server/excel/mazefoev8"
	"maze_game_server/model/barrierstagecountermodel"
	"maze_game_server/model/flowmodel/mazemonstermodel"
	"maze_game_server/services/barrieritemservice"
	"maze_game_server/services/flowservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s service) GetBarrierStageCounter(ctx context.Context, userId uint64, barrierId int32) (killMonsterNum int32, totalDamage, totalExp int64, guidList []int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	model, err := barrierstagecountermodel.NewBarrierStageCounterModel(ctx, userId, barrierId)
	if err != nil {
		logger.CtxError(ctx, "GetBarrierAreaRecord NewBarrierStageCounterModel fail", zap.Error(err))
		return
	}

	killMonsterNum = 0
	totalDamage = 0
	totalExp = int64(0)
	guidList = make([]int64, 0, len(model.KillMonsterRecordMap))

	for _, damage := range model.DamageRecordMap {
		totalDamage += damage
	}
	for _, num := range model.KillMonsterRecordMap {
		killMonsterNum += num
	}
	for _, num := range model.ExpMap {
		totalExp += num
	}
	for _, areaGuids := range model.KillMonsterGuidMap {
		for guid := range areaGuids {
			guidList = append(guidList, guid)
		}
	}

	barrierNum, barrierExp := getBarrierMonsterNum(ctx, barrierId)
	if killMonsterNum > barrierNum {
		logger.CtxInfo(ctx, "GetBarrierStageCounter killMonsterNum > barrierNum", zap.Int32("barrierId", barrierId), zap.Int32("barrierNum", barrierNum), zap.Int32("killMonsterNum", killMonsterNum))
		killMonsterNum = barrierNum
	}
	if totalExp > barrierExp {
		logger.CtxInfo(ctx, "GetBarrierStageCounter totalExp > barrierExp", zap.Int32("barrierId", barrierId), zap.Int64("barrierExp", barrierExp), zap.Int64("totalExp", totalExp))
		totalExp = barrierExp
	}

	logger.CtxInfo(ctx, "GetBarrierStageCounter success", zap.Int32("barrierId", barrierId), zap.Int32("killMonsterNum", killMonsterNum), zap.Int64("totalDamage", totalDamage), zap.Int64("totalExp", totalExp))
	return
}

func (s service) AddKillMonsterNum(ctx context.Context, userId uint64, barrierId, stageId, areaID, areaIndex, monsterId int32, monsterGuid int64, curHp int64, maxHp int64, monsterPos string) (killMonsterNum int32, guidList []int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	recordModel, err := barrierstagecountermodel.NewBarrierStageCounterModel(ctx, userId, barrierId)
	if err != nil {
		logger.CtxError(ctx, "AddKillMonsterNum NewBarrierStageCounterModel fail", zap.Error(err))
		return
	}

	// 初始化
	var addVal int32
	guidList = make([]int64, 0)
	_, ok := recordModel.KillMonsterGuidMap[stageId]
	if !ok {
		recordModel.KillMonsterGuidMap[stageId] = make(map[int64]struct{})
	}

	// 去重
	if _, ok := recordModel.KillMonsterGuidMap[stageId][monsterGuid]; !ok {
		addVal = 1
		logger.CtxInfo(ctx, "AddKillMonsterNum add monsterGuid after", zap.Int64("monsterGuid", monsterGuid), zap.Int32("stageId", stageId))
	}
	recordModel.KillMonsterGuidMap[stageId][monsterGuid] = struct{}{}

	// 获取当前杀怪列表
	for nowGuid := range recordModel.KillMonsterGuidMap[stageId] {
		guidList = append(guidList, nowGuid)
	}
	logger.CtxInfo(ctx, "AddKillMonsterNum add monsterGuid after", zap.Int64("monsterGuid", monsterGuid), zap.Int32("stageId", stageId), zap.Any("guidList", guidList))

	logger.CtxInfo(ctx, "AddKillMonsterNum before kill num", zap.Int32("number", killMonsterNum), zap.Int32("stageId", stageId))
	recordModel.KillMonsterRecordMap[stageId] += addVal
	killMonsterNum = recordModel.KillMonsterRecordMap[stageId]

	// 校验杀怪最大上限和经验收益
	barrierNum, barrierExp := getBarrierMonsterNum(ctx, barrierId)
	if killMonsterNum > barrierNum {
		logger.CtxInfo(ctx, "AddKillMonsterNum killMonsterNum > barrierNum", zap.Int32("barrierId", barrierId), zap.Int32("barrierNum", barrierNum), zap.Int32("killMonsterNum", killMonsterNum))
		killMonsterNum = barrierNum
		recordModel.KillMonsterRecordMap[stageId] = barrierNum
	}
	logger.CtxInfo(ctx, "AddKillMonsterNum after kill num", zap.Int32("number", killMonsterNum), zap.Int32("stageId", stageId))

	addExp := getMonsterExp(ctx, monsterId)
	if addExp > 0 {
		expNum, ok := recordModel.ExpMap[stageId]
		logger.CtxInfo(ctx, "AddKillMonsterNum before exp num", zap.Int64("expNum", expNum), zap.Int32("stageId", stageId))
		if !ok {
			recordModel.ExpMap[stageId] = addExp
		} else {
			recordModel.ExpMap[stageId] += addExp
		}
		expNum += addExp
		if expNum > barrierExp {
			logger.CtxInfo(ctx, "AddKillMonsterNum killMonsterNum > barrierNum", zap.Int32("barrierId", barrierId), zap.Int32("barrierNum", barrierNum), zap.Int32("killMonsterNum", killMonsterNum))
			expNum = barrierExp
			recordModel.ExpMap[stageId] = barrierExp
		}

		logger.CtxInfo(ctx, "AddKillMonsterNum after exp num", zap.Int64("expNum", expNum), zap.Int32("stageId", stageId))
	}

	// 计算装备分数值 物品分数值
	foeCfg := mazefoev8.GetMazeFoeConfig(ctx, monsterId)
	nowEquipScore, dropEquips, err := barrieritemservice.GbarrierItemsService.AddEquipScore(ctx, userId, barrierId, foeCfg.Drop_equip_score_num, killMonsterNum, monsterGuid, monsterPos)
	if err != nil {
		logger.CtxError(ctx, "AddKillMonsterNum AddEquipScore Fail",
			zap.Uint64("userID", userId),
			zap.Int32("barrierId", barrierId),
			zap.Int32("monsterId", monsterId),
		)
		return
	}

	// 掉落物品
	nowItem1Score, dropItem1s, err := barrieritemservice.GbarrierItemsService.AddScoreItem(ctx, userId, barrierId, constdef.MazeCfgId901, foeCfg.Drop_item1_score_num, monsterGuid, monsterPos)
	if err != nil {
		logger.CtxError(ctx, "AddKillMonsterNum AddItemScore Fail",
			zap.Uint64("userID", userId),
			zap.Int32("barrierId", barrierId),
			zap.Int32("monsterId", monsterId),
		)
		return
	}
	nowItem2Score, dropItem2s, err := barrieritemservice.GbarrierItemsService.AddScoreItem(ctx, userId, barrierId, constdef.MazeCfgId902, foeCfg.Drop_item2_score_num, monsterGuid, monsterPos)
	if err != nil {
		logger.CtxError(ctx, "AddKillMonsterNum AddItemScore Fail",
			zap.Uint64("userID", userId),
			zap.Int32("barrierId", barrierId),
			zap.Int32("monsterId", monsterId),
		)
		return
	}

	// 技能物品掉落
	bloodBottle, err := barrieritemservice.GbarrierItemsService.FallOffSkillItems(ctx, userId, barrierId, killMonsterNum, curHp, maxHp, monsterGuid, monsterPos)
	if err != nil {
		logger.CtxError(ctx, "AddKillMonsterNum FallOffSkillItems Fail",
			zap.Uint64("userID", userId),
			zap.Int32("barrierId", barrierId),
			zap.Int32("monsterId", monsterId),
		)
		return
	}

	// 能量点数增加
	// err = tempbuffservice.GlobalTempBuffService.AddTmpBuffEnergy(ctx, userId, barrierId, areaID, areaIndex, foeCfg.Drop_energy_num)
	// if err != nil {
	// 	logger.CtxError(ctx, "AddKillMonsterNum AddTmpBuffEnergy Fail",
	// 		zap.Uint64("userID", userId),
	// 		zap.Int32("barrierId", barrierId),
	// 		zap.Int32("monsterId", monsterId),
	// 	)
	// 	return 0, nil, err
	// }

	err = recordModel.Save(ctx, userId, barrierId)
	if err != nil {
		logger.CtxError(ctx, "AddKillMonsterNum Save fail", zap.Error(err))
		return
	}

	// 发送流水
	monsterRecord := mazemonstermodel.NewMazeMonsterRecordModel(userId, uint32(barrierId), areaID, areaIndex, uint32(nowEquipScore), uint32(nowItem1Score), uint32(nowItem2Score),
		uint32(killMonsterNum), monsterGuid, monsterPos, flowutil.ItemInfo2String(dropEquips, dropItem1s, dropItem2s, bloodBottle))
	flowservice.GflowService.SendFlowData(ctx, monsterRecord)

	return
}

func (s service) AddDamage(ctx context.Context, userId uint64, barrierId, stageId int32, addVal int64) (damage int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	recordModel, err := barrierstagecountermodel.NewBarrierStageCounterModel(ctx, userId, barrierId)
	if err != nil {
		logger.CtxError(ctx, "AddDamage NewBarrierStageCounterModel fail", zap.Error(err))
		return 0, err
	}

	damage, ok := recordModel.DamageRecordMap[stageId]
	logger.CtxInfo(ctx, "AddDamage  before", zap.Int64("damage", damage), zap.Int32("stageId", stageId))
	if !ok {
		recordModel.DamageRecordMap[stageId] = addVal
	} else {
		recordModel.DamageRecordMap[stageId] += addVal
	}
	damage += addVal
	logger.CtxInfo(ctx, "AddDamage  after", zap.Int64("damage", damage), zap.Int32("stageId", stageId))

	err = recordModel.Save(ctx, userId, barrierId)
	if err != nil {
		logger.CtxError(ctx, "AddDamage Save fail", zap.Error(err))
		return damage, err
	}

	return damage, nil
}

func (s service) DelBarrierStageCounter(ctx context.Context, userId uint64, barrierId, stageId int32) error {
	logger := fklog.ContextAppLogger(ctx)
	passArea := dollmappuzzlenewcfgex.GetPassAreaInfos(barrierId, stageId)
	recordModel, err := barrierstagecountermodel.NewBarrierStageCounterModel(ctx, userId, barrierId)
	if err != nil {
		logger.CtxError(ctx, "DelBarrierStageCounter NewBarrierStageCounterModel fail", zap.Error(err))
		return err
	}

	if len(passArea) == 0 {
		err = recordModel.Del(ctx, userId, barrierId)
		if err != nil {
			logger.CtxError(ctx, "DelBarrierStageCounter DEL fail", zap.Error(err))
			return err
		}
		logger.CtxInfo(ctx, "DelBarrierStageCounter key", zap.Uint64("userId", userId), zap.Int32("barrierId", barrierId), zap.Int32("stageId", stageId))
	} else {
		//通过的区域不删除
		//删除未完成区域经验存档
		delExpList := make([]int32, 0)
		for k, _ := range recordModel.ExpMap {
			isPass := isPassBarrierArea(k, passArea)
			if !isPass {
				delExpList = append(delExpList, k)
			}
		}
		logger.CtxInfo(ctx, "DelBarrierStageCounter delExpList", zap.Any("delExpList", delExpList))
		for _, field := range delExpList {
			_, ok := recordModel.ExpMap[field]
			if ok {
				delete(recordModel.ExpMap, field)
			}
		}

		//删除未完成区域伤害值存档
		delList := make([]int32, 0)
		for k, _ := range recordModel.DamageRecordMap {
			isPass := isPassBarrierArea(k, passArea)
			if !isPass {
				delList = append(delList, k)
			}
		}
		logger.CtxInfo(ctx, "DelBarrierStageCounter delList", zap.Any("delList", delList))
		for _, field := range delList {
			_, ok := recordModel.DamageRecordMap[field]
			if ok {
				delete(recordModel.DamageRecordMap, field)
			}
		}

		//删除未完成区域杀怪数存档
		delArr := make([]int32, 0)
		for k, _ := range recordModel.KillMonsterRecordMap {
			isPass := isPassBarrierArea(k, passArea)
			if !isPass {
				delArr = append(delArr, k)
			}
		}
		logger.CtxInfo(ctx, "DelBarrierStageCounter delArr", zap.Any("delArr", delArr))
		for _, field := range delArr {
			_, ok := recordModel.KillMonsterRecordMap[field]
			if ok {
				delete(recordModel.KillMonsterRecordMap, field)
			}
			_, ok = recordModel.KillMonsterGuidMap[field]
			if ok {
				delete(recordModel.KillMonsterGuidMap, field)
			}
		}

		err = recordModel.Save(ctx, userId, barrierId)
		if err != nil {
			logger.CtxError(ctx, "DelBarrierStageCounter Save fail", zap.Error(err))
		}
		logger.CtxInfo(ctx, "DelBarrierStageCounter success", zap.Uint64("userId", userId), zap.Int32("barrierId", barrierId), zap.Int32("stageId", stageId))
	}

	return err
}

func (s service) DelBarrierStageCounterOnPass(ctx context.Context, userId uint64, barrierId int32) error {
	logger := fklog.ContextAppLogger(ctx)
	for i := int32(0); i <= barrierId; i++ {
		// recordModel, err := barrierstagecountermodel.NewBarrierStageCounterModel(ctx, userId, barrierId)
		// if err != nil {
		// 	logger.CtxError(ctx, "DelBarrierStageCounter NewBarrierStageCounterModel fail", zap.Error(err))
		// 	continue
		// }
		err := barrierstagecountermodel.GMDel(ctx, userId, barrierId)
		if err != nil {
			logger.CtxError(ctx, "DelBarrierStageCounter DEL fail", zap.Error(err))
			continue
		}
		logger.CtxInfo(ctx, "DelBarrierStageCounterOnPass key", zap.Uint64("userId", userId), zap.Int32("barrierId", barrierId))
	}
	return nil
}

// 是否通过关卡区域
func isPassBarrierArea(stageId int32, passArea []*dollmappuzzlenewcfgex.AreaInfo) bool {
	for _, info := range passArea {
		if stageId > 0 && info.StageId > stageId {
			return true
		}
	}
	return false
}

func getMonsterExp(ctx context.Context, monsterId int32) int64 {
	logger := fklog.ContextAppLogger(ctx)
	foeCfg := GMazeFoeV8Cfg.GetWithCtx(ctx, monsterId)
	if foeCfg == nil {
		logger.CtxError(ctx, "getMonsterExp get foe cfg fail", zap.Any("foeId", monsterId))
		return 0
	}
	return int64(foeCfg.Drop_exp_num)
}

func getBarrierMonsterNum(ctx context.Context, barrierId int32) (totalMonsterNum int32, totalExpNum int64) {
	logger := fklog.ContextAppLogger(ctx)
	barrierCfg := GMazeBarriesV8Cfg.GetWithCtx(ctx, barrierId)
	if barrierCfg == nil {
		logger.CtxError(ctx, "getBarrierMonsterNum get barrier cfg fail", zap.Any("barrierId", barrierId))
		return 0, 0
	}

	allFoe := GMazeBrushFoeV8Cfg.GetAll()
	if len(allFoe) == 0 {
		logger.CtxInfo(ctx, "getBarrierMonsterNum barrier foe empty", zap.Any("barrierId", barrierId))
		return 0, 0
	}
	totalMonsterNum = int32(0)
	totalExpNum = int64(0)
	for _, cfg := range allFoe {
		if cfg.Barries_id != barrierId {
			continue
		}
		for _, foe := range cfg.Monsters_id {
			if foe > 0 {
				totalMonsterNum += 1
				totalExpNum += getMonsterExp(ctx, foe)
			}
		}
		for foe, num := range cfg.Monsterslist_ids_and_nums {
			if foe > 0 {
				totalMonsterNum += num
				totalExpNum += getMonsterExp(ctx, foe) * int64(num)
			}
		}
	}

	logger.CtxInfo(ctx, "getBarrierMonsterNum ", zap.Int32("barrierId", barrierId), zap.Int32("totalMonsterNum", totalMonsterNum), zap.Int64("totalExpNum", totalExpNum))
	return totalMonsterNum, totalExpNum
}
