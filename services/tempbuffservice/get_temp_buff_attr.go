package tempbuffservice

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/model/tempbuffmodel"
)

func (s *service) GetTempBuffAttr(logger fklog.FKLogI, userID uint64, stageId int32) (map[int32]int64, error) {
	buffInfo, err := tempbuffmodel.NewTempBuffInfoModel(logger, userID, stageId)
	if err != nil {
		logger.ErrorWF("SelectMazeTempBuffRQ GetMazeTempBuff", zap.Error(err))
		return nil, fmt.Errorf("获取用户buff信息失败")
	}
	if buffInfo == nil {
		logger.WarnWF("SelectMazeTempBuffRQ buff is nil", zap.Error(err))
		return nil, fmt.Errorf("获取用户buff信息失败")
	}
	attr := make(map[int32]int64)
	for _, i := range buffInfo.TotalBuff {
		attr[i.BuffId] += i.BuffValue
	}

	return attr, nil
}

func (s *service) GetTempBuffInfo(logger fklog.FKLogI, userID uint64, stageId int32) (*tempbuffmodel.TempBuffInfoModel, error) {
	buffInfo, err := tempbuffmodel.NewTempBuffInfoModel(logger, userID, stageId)
	if err != nil {
		logger.ErrorWF("SelectMazeTempBuffRQ GetMazeTempBuff", zap.Error(err))
		return nil, fmt.Errorf("获取用户buff信息失败")
	}
	if buffInfo == nil {
		logger.WarnWF("SelectMazeTempBuffRQ buff is nil", zap.Error(err))
		return nil, fmt.Errorf("获取用户buff信息失败")
	}

	return buffInfo, nil
}
