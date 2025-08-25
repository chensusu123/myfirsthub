package barriersavedataservice

import (
	"context"
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/model/barriersavedatamodel"
	"maze_game_server/module/mazecommonvalue"
)

func (s *service) SaveBarrierData(ctx context.Context, userId uint64, barrier, stageId, rescueValue, bossPower int32,
	bossProgress float32, rescueItems []*barriersavedatamodel.RescueItemInfo) error {

	logger := fklog.ContextAppLogger(ctx)
	model, err := barriersavedatamodel.NewBarrierSaveDataModel(ctx, userId, barrier, false)
	if err != nil {
		logger.CtxError(ctx, "GetBarrierSaveData NewBarrierSaveDataModel err", zap.Error(err))
		return fmt.Errorf("获取关卡存档失败")
	}
	model.StageId = stageId
	model.RescueValue = rescueValue
	model.BossPower = bossPower
	model.BossProgress = bossProgress
	model.RescueItems = rescueItems
	err = model.Save(ctx, userId, barrier)
	if err != nil {
		logger.CtxError(ctx, "SaveBarrierData save err", zap.Error(err))
	}
	return nil
}

func (s *service) GetBarrierSaveData(ctx context.Context, userId uint64, barrier int32) (*barriersavedatamodel.BarrierSaveDataModel, error) {
	logger := fklog.ContextAppLogger(ctx)
	model, err := barriersavedatamodel.NewBarrierSaveDataModel(ctx, userId, barrier, true)
	if err != nil {
		logger.CtxError(ctx, "GetBarrierSaveData NewBarrierSaveDataModel err", zap.Error(err))
		return nil, fmt.Errorf("获取关卡存档失败")
	}

	return model, nil
}

func (s *service) DelBarrierSaveData(ctx context.Context, userId uint64, barrier int32) error {
	logger := fklog.ContextAppLogger(ctx)
	model, err := barriersavedatamodel.NewBarrierSaveDataModel(ctx, userId, barrier, false)
	if err != nil {
		logger.CtxError(ctx, "DelBarrierSaveData NewBarrierSaveDataModel err", zap.Error(err))
		return nil
	}
	err = model.Del(ctx, userId, barrier)
	if err != nil {
		logger.CtxError(ctx, "DelBarrierSaveData del err", zap.Error(err))
		return nil
	}

	return nil
}

func (s *service) GetPassValue(ctx context.Context, userId uint64, barrier int32) (int64, error) {
	logger := fklog.ContextAppLogger(ctx)
	model, err := barriersavedatamodel.NewBarrierSaveDataModel(ctx, userId, barrier, true)
	if err != nil {
		logger.CtxError(ctx, "GetPassValue NewBarrierSaveDataModel err", zap.Error(err))
		return 0, nil
	}
	passValue, err := mazecommonvalue.CalcPassValue(ctx, logger, barrier, model.StageId)
	if err != nil {
		logger.CtxError(ctx, "GetPassValue calcPassValue err", zap.Error(err))
		return 0, err
	}
	return passValue, nil
}
