package barrierarearecordservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/model/barrierarearecordmodel"
	"maze_game_server/model/passareamodel"
)

func (s service) GetBarrierAreaRecord(logger fklog.FKLogI, userId uint64, stageId int32) (killMonsterNum int32, totalDamage int64, guidList []int64, err error) {
	model, err := barrierarearecordmodel.NewBarrierAreaNumRecordModel(logger, userId, stageId)
	if err != nil {
		logger.ErrorWF("GetBarrierAreaRecord NewBarrierAreaNumRecordModel fail", zap.Error(err))
		return 0, 0, nil, err
	}
	killMonsterNum = 0
	totalDamage = 0
	for _, damage := range model.DamageRecordMap {
		totalDamage += damage
	}
	for _, num := range model.KillMonsterRecordMap {
		killMonsterNum += num
	}

	guidList = make([]int64, 0, len(model.KillMonsterRecordMap))
	for _, guid := range model.KillMonsterGuidMap {
		guidList = append(guidList, guid...)
	}

	return
}

func (s service) AddKillMonsterNum(logger fklog.FKLogI, userId uint64, stageId, areaId, areaIndex, monsterId, addVal int32, monsterGuid int64) (killMonsterNum int32, guidList []int64, err error) {
	recordModel, err := barrierarearecordmodel.NewBarrierAreaNumRecordModel(logger, userId, stageId)
	if err != nil {
		logger.ErrorWF("AddKillMonsterNum NewBarrierAreaNumRecordModel fail", zap.Error(err))
		return 0, nil, err
	}

	field := EnCodeAreaField(areaId, areaIndex)
	number, ok := recordModel.KillMonsterRecordMap[field]
	logger.InfoWF("AddKillMonsterNum area before", zap.Int32("number", number), zap.Int32("field", field))
	if !ok {
		recordModel.KillMonsterRecordMap[field] = addVal
	} else {
		number += addVal
		recordModel.KillMonsterRecordMap[field] += addVal
	}
	logger.InfoWF("AddKillMonsterNum area after", zap.Int32("number", number), zap.Int32("field", field))

	isList, ok := recordModel.KillMonsterGuidMap[field]
	if !ok {
		isList = make([]int64, 0)
	}
	isList = append(guidList, monsterGuid)
	recordModel.KillMonsterGuidMap[field] = isList
	logger.InfoWF("AddKillMonsterNum area add monsterGuid", zap.Int64("monsterGuid", monsterGuid), zap.Int32("field", field))

	err = recordModel.Save(logger, userId, stageId)
	if err != nil {
		logger.ErrorWF("AddKillMonsterNum Save fail", zap.Error(err))
		return number, nil, err
	}

	killMonsterNum, _, guidList, err = s.GetBarrierAreaRecord(logger, userId, stageId)
	if err != nil {
		logger.ErrorWF("AddKillMonsterNum GetBarrierAreaRecord fail", zap.Error(err))
		return number, guidList, err
	}

	logger.InfoWF("AddKillMonsterNum success", zap.Int32("killMonsterNum", killMonsterNum))
	return killMonsterNum, guidList, nil
}

func (s service) AddDamage(logger fklog.FKLogI, userId uint64, stageId, areaId, areaIndex int32, addVal int64) (totalDamage int64, err error) {
	recordModel, err := barrierarearecordmodel.NewBarrierAreaNumRecordModel(logger, userId, stageId)
	if err != nil {
		logger.ErrorWF("AddDamage NewBarrierAreaNumRecordModel fail", zap.Error(err))
		return 0, err
	}

	field := EnCodeAreaField(areaId, areaIndex)

	damage, ok := recordModel.DamageRecordMap[field]
	logger.InfoWF("AddDamage area before", zap.Int64("damage", damage), zap.Int32("field", field))
	if !ok {
		recordModel.DamageRecordMap[field] = addVal
	} else {
		damage += addVal
		recordModel.DamageRecordMap[field] += addVal
	}
	logger.InfoWF("AddDamage area after", zap.Int64("damage", damage), zap.Int32("field", field))

	err = recordModel.Save(logger, userId, stageId)
	if err != nil {
		logger.ErrorWF("AddKillMonsterNum Save fail", zap.Error(err))
		return damage, err
	}

	_, totalDamage, _, err = s.GetBarrierAreaRecord(logger, userId, stageId)
	if err != nil {
		logger.ErrorWF("AddKillMonsterNum GetBarrierAreaRecord fail", zap.Error(err))
		return damage, err
	}

	return totalDamage, nil
}

func (s service) DelBarrierAreaRecord(logger fklog.FKLogI, userId uint64, stageId int32) error {

	passAreaModel, err := passareamodel.NewPassAreaModel(logger, userId, stageId)
	if err != nil {
		logger.ErrorWF("DelBarrierAreaRecord NewPassAreaModel fail", zap.Error(err))
		return err
	}

	recordModel, err := barrierarearecordmodel.NewBarrierAreaNumRecordModel(logger, userId, stageId)
	if err != nil {
		logger.ErrorWF("AddDamage NewBarrierAreaNumRecordModel fail", zap.Error(err))
		return err
	}

	if len(passAreaModel.PassAreaList) == 0 {
		err = recordModel.Del(logger, userId, stageId)
		if err != nil {
			logger.ErrorWF("DelBarrierAreaRecord DEL fail", zap.Error(err))
			return err
		}
	} else {
		//通过的区域不删除
		//删除未完成区域伤害值存档
		delList := make([]int32, 0)
		for k, _ := range recordModel.DamageRecordMap {
			isPass := isPassBarrierArea(logger, userId, stageId, k)
			if !isPass {
				delList = append(delList, k)
			}
		}
		logger.InfoWF("DelBarrierAreaRecord delList", zap.Any("delList", delList))
		for _, field := range delList {
			_, ok := recordModel.DamageRecordMap[field]
			if ok {
				delete(recordModel.DamageRecordMap, field)
			}
		}

		//删除未完成区域杀怪数存档
		delArr := make([]int32, 0)
		for k, _ := range recordModel.KillMonsterRecordMap {
			isPass := isPassBarrierArea(logger, userId, stageId, k)
			if !isPass {
				delArr = append(delArr, k)
			}
		}
		logger.InfoWF("DelBarrierAreaRecord delArr", zap.Any("delArr", delArr))
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
	}

	err = recordModel.Save(logger, userId, stageId)
	if err != nil {
		logger.ErrorWF("DelBarrierAreaRecord Save fail", zap.Error(err))
	}
	return err
}

const killMonsterPercent = 10000

// 编码field
func EnCodeAreaField(areaId, areaIndex int32) int32 {
	return areaId*killMonsterPercent + areaIndex
}

// 解码field
func DecodeAreaField(field int32) (areaId, areaIndex int32) {
	areaId = field / killMonsterPercent
	areaIndex = field % killMonsterPercent
	return
}

// 是否通过关卡区域
func isPassBarrierArea(logger fklog.FKLogI, userId uint64, stageId int32, field int32) bool {
	passAreaModel, err := passareamodel.NewPassAreaModel(logger, userId, stageId)
	if err != nil {
		logger.ErrorWF("isPassBarrierArea NewPassAreaModel fail", zap.Error(err))
		return false
	}

	areaId, areaIndex := DecodeAreaField(field)

	for _, info := range passAreaModel.PassAreaList {
		if info.AreaId == areaId && info.AreaIndex == areaIndex {
			return true
		}
	}
	return false
}
