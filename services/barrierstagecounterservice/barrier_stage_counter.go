package barrierstagecounterservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/config/GMazeBrushFoeV8Cfg"
	"maze_game_server/config/GMazeFoeV8Cfg"
	"maze_game_server/excel/dollmappuzzlenewcfgex"
	"maze_game_server/model/barrierstagecountermodel"
)

func (s service) GetBarrierStageCounter(logger fklog.FKLogI, userId uint64, barrierId int32) (killMonsterNum int32, totalDamage, totalExp int64, guidList []int64, err error) {
	model, err := barrierstagecountermodel.NewBarrierStageCounterModel(logger, userId, barrierId)
	if err != nil {
		logger.ErrorWF("GetBarrierAreaRecord NewBarrierStageCounterModel fail", zap.Error(err))
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
	for _, guid := range model.KillMonsterGuidMap {
		guidList = append(guidList, guid...)
	}

	barrierNum, barrierExp := getBarrierMonsterNum(logger, barrierId)
	if killMonsterNum > barrierNum {
		logger.InfoWF("GetBarrierStageCounter killMonsterNum > barrierNum", zap.Int32("barrierId", barrierId), zap.Int32("barrierNum", barrierNum), zap.Int32("killMonsterNum", killMonsterNum))
		killMonsterNum = barrierNum
	}
	if totalExp > barrierExp {
		logger.InfoWF("GetBarrierStageCounter totalExp > barrierExp", zap.Int32("barrierId", barrierId), zap.Int64("barrierExp", barrierExp), zap.Int64("totalExp", totalExp))
		totalExp = barrierExp
	}

	logger.InfoWF("GetBarrierStageCounter success", zap.Int32("barrierId", barrierId), zap.Int32("killMonsterNum", killMonsterNum), zap.Int64("totalDamage", totalDamage), zap.Int64("totalExp", totalExp))
	return
}

func (s service) AddKillMonsterNum(logger fklog.FKLogI, userId uint64, barrierId, stageId, monsterId, addVal int32, monsterGuid int64) (killMonsterNum int32, guidList []int64, err error) {

	recordModel, err := barrierstagecountermodel.NewBarrierStageCounterModel(logger, userId, barrierId)
	if err != nil {
		logger.ErrorWF("AddKillMonsterNum NewBarrierStageCounterModel fail", zap.Error(err))
		return 0, nil, err
	}

	barrierNum, barrierExp := getBarrierMonsterNum(logger, barrierId)

	killMonsterNum, ok := recordModel.KillMonsterRecordMap[stageId]
	logger.InfoWF("AddKillMonsterNum before kill num", zap.Int32("number", killMonsterNum), zap.Int32("stageId", stageId))
	if !ok {
		recordModel.KillMonsterRecordMap[stageId] = addVal
	} else {
		recordModel.KillMonsterRecordMap[stageId] += addVal
	}
	killMonsterNum += addVal
	if killMonsterNum > barrierNum {
		logger.InfoWF("AddKillMonsterNum killMonsterNum > barrierNum", zap.Int32("barrierId", barrierId), zap.Int32("barrierNum", barrierNum), zap.Int32("killMonsterNum", killMonsterNum))
		killMonsterNum = barrierNum
		recordModel.KillMonsterRecordMap[stageId] = barrierNum
	}
	logger.InfoWF("AddKillMonsterNum after kill num", zap.Int32("number", killMonsterNum), zap.Int32("stageId", stageId))

	guidList, ok = recordModel.KillMonsterGuidMap[stageId]
	if !ok {
		guidList = make([]int64, 0)
	}
	guidList = append(guidList, monsterGuid)
	recordModel.KillMonsterGuidMap[stageId] = guidList
	logger.InfoWF("AddKillMonsterNum add monsterGuid", zap.Int64("monsterGuid", monsterGuid), zap.Int32("stageId", stageId))

	addExp := getMonsterExp(logger, monsterId)
	if addExp > 0 {
		expNum, ok := recordModel.ExpMap[stageId]
		logger.InfoWF("AddKillMonsterNum before exp num", zap.Int64("expNum", expNum), zap.Int32("stageId", stageId))
		if !ok {
			recordModel.ExpMap[stageId] = addExp
		} else {
			recordModel.ExpMap[stageId] += addExp
		}
		expNum += addExp
		if expNum > barrierExp {
			logger.InfoWF("AddKillMonsterNum killMonsterNum > barrierNum", zap.Int32("barrierId", barrierId), zap.Int32("barrierNum", barrierNum), zap.Int32("killMonsterNum", killMonsterNum))
			expNum = barrierExp
			recordModel.ExpMap[stageId] = barrierExp
		}

		logger.InfoWF("AddKillMonsterNum after exp num", zap.Int64("expNum", expNum), zap.Int32("stageId", stageId))
	}

	err = recordModel.Save(logger, userId, barrierId)
	if err != nil {
		logger.ErrorWF("AddKillMonsterNum Save fail", zap.Error(err))
		return killMonsterNum, nil, err
	}

	return killMonsterNum, guidList, nil
}

func (s service) AddDamage(logger fklog.FKLogI, userId uint64, barrierId, stageId int32, addVal int64) (damage int64, err error) {

	recordModel, err := barrierstagecountermodel.NewBarrierStageCounterModel(logger, userId, barrierId)
	if err != nil {
		logger.ErrorWF("AddDamage NewBarrierStageCounterModel fail", zap.Error(err))
		return 0, err
	}

	damage, ok := recordModel.DamageRecordMap[stageId]
	logger.InfoWF("AddDamage  before", zap.Int64("damage", damage), zap.Int32("stageId", stageId))
	if !ok {
		recordModel.DamageRecordMap[stageId] = addVal
	} else {
		recordModel.DamageRecordMap[stageId] += addVal
	}
	damage += addVal
	logger.InfoWF("AddDamage  after", zap.Int64("damage", damage), zap.Int32("stageId", stageId))

	err = recordModel.Save(logger, userId, barrierId)
	if err != nil {
		logger.ErrorWF("AddDamage Save fail", zap.Error(err))
		return damage, err
	}

	return damage, nil
}

func (s service) DelBarrierStageCounter(logger fklog.FKLogI, userId uint64, barrierId, stageId int32) error {

	passArea := dollmappuzzlenewcfgex.GetPassAreaInfos(barrierId, stageId)
	recordModel, err := barrierstagecountermodel.NewBarrierStageCounterModel(logger, userId, barrierId)
	if err != nil {
		logger.ErrorWF("DelBarrierStageCounter NewBarrierStageCounterModel fail", zap.Error(err))
		return err
	}

	if len(passArea) == 0 {
		err = recordModel.Del(logger, userId, barrierId)
		if err != nil {
			logger.ErrorWF("DelBarrierStageCounter DEL fail", zap.Error(err))
			return err
		}
		logger.InfoWF("DelBarrierStageCounter key", zap.Uint64("userId", userId), zap.Int32("barrierId", barrierId), zap.Int32("stageId", stageId))
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
		logger.InfoWF("DelBarrierStageCounter delExpList", zap.Any("delExpList", delExpList))
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
		logger.InfoWF("DelBarrierStageCounter delList", zap.Any("delList", delList))
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
		logger.InfoWF("DelBarrierStageCounter delArr", zap.Any("delArr", delArr))
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

		err = recordModel.Save(logger, userId, barrierId)
		if err != nil {
			logger.ErrorWF("DelBarrierStageCounter Save fail", zap.Error(err))
		}
		logger.InfoWF("DelBarrierStageCounter success", zap.Uint64("userId", userId), zap.Int32("barrierId", barrierId), zap.Int32("stageId", stageId))
	}

	return err
}

func (s service) DelBarrierStageCounterOnPass(logger fklog.FKLogI, userId uint64, barrierId int32) error {
	for i := int32(0); i <= barrierId; i++ {
		recordModel, err := barrierstagecountermodel.NewBarrierStageCounterModel(logger, userId, barrierId)
		if err != nil {
			logger.ErrorWF("DelBarrierStageCounter NewBarrierStageCounterModel fail", zap.Error(err))
			continue
		}
		err = recordModel.Del(logger, userId, barrierId)
		if err != nil {
			logger.ErrorWF("DelBarrierStageCounter DEL fail", zap.Error(err))
			continue
		}
		logger.InfoWF("DelBarrierStageCounterOnPass key", zap.Uint64("userId", userId), zap.Int32("barrierId", barrierId))
	}
	return nil
}

// 是否通过关卡区域
func isPassBarrierArea(stageId int32, passArea []*dollmappuzzlenewcfgex.AreaInfo) bool {
	for _, info := range passArea {
		if info.StageId == 0 || info.StageId > stageId {
			return true
		}
	}
	return false
}

func getMonsterExp(logger fklog.FKLogI, monsterId int32) int64 {
	foeCfg := GMazeFoeV8Cfg.Get(monsterId)
	if foeCfg == nil {
		logger.ErrorWF("getMonsterExp get foe cfg fail", zap.Any("foeId", monsterId))
		return 0
	}
	return int64(foeCfg.Drop_exp_num)
}

func getBarrierMonsterNum(logger fklog.FKLogI, barrierId int32) (totalMonsterNum int32, totalExpNum int64) {
	barrierCfg := GMazeBarriesV8Cfg.Get(barrierId)
	if barrierCfg == nil {
		logger.ErrorWF("getBarrierMonsterNum get barrier cfg fail", zap.Any("barrierId", barrierId))
		return 0, 0
	}

	allFoe := GMazeBrushFoeV8Cfg.GetAll()
	if len(allFoe) == 0 {
		logger.ErrorWF("getBarrierMonsterNum barrier foe empty", zap.Any("barrierId", barrierId))
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
				totalExpNum += getMonsterExp(logger, foe)
			}
		}
		for foe, num := range cfg.Monsterslist_ids_and_nums {
			if foe > 0 {
				totalMonsterNum += num
				totalExpNum += getMonsterExp(logger, foe) * int64(num)
			}
		}
	}

	logger.InfoWF("getBarrierMonsterNum ", zap.Int32("barrierId", barrierId), zap.Int32("totalMonsterNum", totalMonsterNum), zap.Int64("totalExpNum", totalExpNum))
	return totalMonsterNum, totalExpNum
}
