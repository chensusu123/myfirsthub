package barrieritemservice

import (
	"context"
	"maze_game_server/config/GMazeItemsV8Cfg"
	"maze_game_server/model/barrieritemsmodel"
	"maze_game_server/services/itemservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) TrySubBarrierItems(ctx context.Context, userID uint64, barrierID int32, items []*itemservice.ItemInfo, equips []*itemservice.ItemInfo) (bool, error) {
	logger := fklog.ContextAppLogger(ctx)
	defer func() {
		logger.CtxInfo(ctx, "TrySubBarrierItems End",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
		)
	}()
	logger.CtxInfo(ctx, "TrySubBarrierItems Start",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
	)

	data, err := barrieritemsmodel.NewBarrierItems(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "TrySubBarrierItems NewBarrierItems Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Error(err),
		)
		return false, err
	}

	for _, item := range items {
		if data.Items[int64(item.ItemId)] < item.Count {
			return false, nil
		}
	}
	for _, equip := range equips {
		if data.Equips[uint64(equip.ItemId)] < uint64(equip.Count) {
			return false, nil
		}
	}

	// 实际添加物品
	for _, item := range items {
		// 技能道具使用
		itemCfg := GMazeItemsV8Cfg.GetWithCtx(ctx, item.ItemId)
		if itemCfg.Type == int32(barrieritemsmodel.SpecialType) {
			data.SkillsCount[item.ItemId] += int32(item.Count)
		}
		data.Items[int64(item.ItemId)] -= item.Count
	}

	// 实际添加装备
	for _, equip := range equips {
		data.Equips[uint64(equip.ItemId)] -= uint64(equip.Count)
	}

	err = data.Save(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "TrySubBarrierItems Save Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Error(err),
		)
		return false, err
	}

	return true, nil
}
