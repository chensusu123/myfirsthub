package bagservice

import (
	"context"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/model/bagmodel"
)

func (s *service) GetAllBagItem(ctx context.Context, userID uint64) (map[int32]int64, error) {
	logger := fklog.ContextAppLogger(ctx)
	bagModel, err := bagmodel.NewBagModel(ctx, userID)
	if err != nil {
		return nil, err
	}
	bagItemMap, err := bagModel.LoadAll(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "GetAllBagItem LoadAll err", zap.Error(err))
		return nil, err
	}
	return bagItemMap, nil
}

func (s *service) DelAllBagItem(ctx context.Context, userID uint64) error {
	logger := fklog.ContextAppLogger(ctx)
	bagModel, err := bagmodel.NewBagModel(ctx, userID)
	if err != nil {
		return err
	}
	err = bagModel.Del(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "DelAllBagItem del err", zap.Error(err))
		return err
	}
	return nil
}
