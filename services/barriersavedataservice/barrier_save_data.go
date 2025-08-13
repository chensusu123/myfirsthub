package barriersavedataservice

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/model/barriersavedatamodel"
	"maze_game_server/module/mazecommonvalue"
)

func (s *service) SaveBarrierData(logger fklog.FKLogI, userId uint64, barrier, stageId, rescueValue, bossPower int32) error {
	model, err := barriersavedatamodel.NewBarrierSaveDataModel(logger, userId, barrier, false)
	if err != nil {
		logger.ErrorWF("GetBarrierSaveData NewBarrierSaveDataModel err", zap.Error(err))
		return fmt.Errorf("获取关卡存档失败")
	}
	model.StageId = stageId
	model.RescueValue = rescueValue
	model.BossPower = bossPower
	err = model.Save(logger, userId, barrier)
	if err != nil {
		logger.ErrorWF("SaveBarrierData save err", zap.Error(err))
	}
	return nil
}

func (s *service) GetBarrierSaveData(logger fklog.FKLogI, userId uint64, barrier int32) (*barriersavedatamodel.BarrierSaveDataModel, error) {
	model, err := barriersavedatamodel.NewBarrierSaveDataModel(logger, userId, barrier, true)
	if err != nil {
		logger.ErrorWF("GetBarrierSaveData NewBarrierSaveDataModel err", zap.Error(err))
		return nil, fmt.Errorf("获取关卡存档失败")
	}

	return model, nil
}

func (s *service) DelBarrierSaveData(logger fklog.FKLogI, userId uint64, barrier int32) error {
	model, err := barriersavedatamodel.NewBarrierSaveDataModel(logger, userId, barrier, false)
	if err != nil {
		logger.ErrorWF("DelBarrierSaveData NewBarrierSaveDataModel err", zap.Error(err))
		return nil
	}
	err = model.Del(logger, userId, barrier)
	if err != nil {
		logger.ErrorWF("DelBarrierSaveData del err", zap.Error(err))
		return nil
	}

	return nil
}

func (s *service) GetPassValue(logger fklog.FKLogI, userId uint64, barrier int32) (int64, error) {
	model, err := barriersavedatamodel.NewBarrierSaveDataModel(logger, userId, barrier, true)
	if err != nil {
		logger.ErrorWF("GetPassValue NewBarrierSaveDataModel err", zap.Error(err))
		return 0, nil
	}
	passValue, err := mazecommonvalue.CalcPassValue(logger, barrier, model.StageId)
	if err != nil {
		logger.ErrorWF("GetPassValue calcPassValue err", zap.Error(err))
		return 0, err
	}
	return passValue, nil
}
