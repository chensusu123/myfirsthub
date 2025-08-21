package tempbuffservice

import (
	"context"
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/model/tempbuffmodel"
)

func (s *service) GetTempBuffAttr(ctx context.Context, userID uint64, barrierId int32) (map[int32]int64, error) {
	logger := fklog.ContextAppLogger(ctx)
	buffInfo, err := tempbuffmodel.NewTempBuffInfoModel(ctx, userID, barrierId)
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

func (s *service) GetTempBuffInfo(ctx context.Context, userID uint64, barrierId int32) (*tempbuffmodel.TempBuffInfoModel, error) {
	logger := fklog.ContextAppLogger(ctx)
	buffInfo, err := tempbuffmodel.NewTempBuffInfoModel(ctx, userID, barrierId)
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
