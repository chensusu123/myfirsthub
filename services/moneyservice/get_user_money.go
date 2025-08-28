package moneyservice

import (
	"context"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/common/constdef"
	"maze_game_server/model/moneymodel"
)

func (s *service) GetUserMoney(ctx context.Context, userId uint64) (coin, diamond int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	model, err := moneymodel.NewMoneyModel(ctx, userId)
	if err != nil {
		return
	}

	moneyMap, err := model.LoadAll(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "GetUserMoney GetMoneyCount err", zap.Error(err))
		return
	}
	for k, v := range moneyMap {
		if k == constdef.MazeCommonItemCoin {
			coin = v
		} else if k == constdef.MazeCommonItemDiamond {
			diamond = v
		}
	}
	return
}

func (s *service) SetMoney(ctx context.Context, userId uint64, itemId int32, value int64) error {
	logger := fklog.ContextAppLogger(ctx)
	model, err := moneymodel.NewMoneyModel(ctx, userId)
	if err != nil {
		return err
	}

	err = model.SetValue(ctx, userId, itemId, value)
	if err != nil {
		logger.ErrorWF("GetUserMoney SetMoney err", zap.Error(err))
		return err
	}
	return nil
}
