package barrieritemservice

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/excel/mazebarriesv8config"
	"maze_game_server/excel/mazeconfigv8"
	"maze_game_server/io/redis/barrierguiditemredis"
	"maze_game_server/model/barrieritemsmodel"
	"maze_game_server/servers/maze_main_server/process/item"
	"maze_game_server/services/itemservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) AddScoreItem(ctx context.Context, userID uint64, barrierID int32, itemType int32, score int32, guid int64, pos string) (res int32, dropItems []*itemservice.ItemInfo, err error) {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "AddScoreItem Start",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Int32("score", score),
		zap.Int64("guid", guid),
		zap.String("pos", pos),
	)

	defer func() {
		logger.CtxInfo(ctx, "AddScoreItem End",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("score", score),
			zap.Int64("guid", guid),
			zap.String("pos", pos),
		)
	}()

	data, err := barrieritemsmodel.NewBarrierItems(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "AddScoreItem NewBarrierItems Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("score", score),
			zap.Int64("guid", guid),
			zap.String("pos", pos),
			zap.Error(err),
		)
		return
	}

	logger.CtxInfo(ctx, "AddScoreItem GetData Successful",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Any("nowdata", data),
	)

	barrierCfg := mazebarriesv8config.GetStageConfig(ctx, barrierID)
	if barrierCfg.Need_item1_score == 0 || barrierCfg.Need_item2_score == 0 {
		err = errors.New("配置不存在")
		logger.CtxError(ctx, "AddScoreItem barrierCfg Need_item1_score or Need_item2_score Equal zero",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("score", score),
			zap.Int64("guid", guid),
			zap.String("pos", pos),
			zap.Any("data", data),
		)
		return
	}

	var realyItemID int32
	var realyNum int32
	if itemType == constdef.MazeCfgId901 {
		itemMap := mazeconfigv8.GetReallyItemMap(ctx, itemType)
		for k := range itemMap {
			realyItemID = k
		}

		data.ItemsScore[realyItemID] += score
		realyNum = data.ItemsScore[realyItemID] / barrierCfg.Need_item1_score
		data.ItemsScore[realyItemID] %= barrierCfg.Need_item1_score
		res = data.ItemsScore[realyItemID]
	} else if itemType == constdef.MazeCfgId902 {
		itemMap := mazeconfigv8.GetReallyItemMap(ctx, itemType)
		for k := range itemMap {
			realyItemID = k
		}

		data.ItemsScore[realyItemID] += score
		realyNum = data.ItemsScore[realyItemID] / barrierCfg.Need_item2_score
		data.ItemsScore[realyItemID] %= barrierCfg.Need_item2_score
		res = data.ItemsScore[realyItemID]
	}

	for i := 1; i <= int(realyNum); i++ {
		itemGuid, err := barrierguiditemredis.IncrNowGuid(ctx, userID, barrierID)
		if err != nil {
			logger.CtxError(ctx, "AddScoreItem IncrNowGuid Fail",
				zap.Uint64("userID", userID),
				zap.Int32("barrierID", barrierID),
				zap.Int32("score", score),
				zap.Int64("guid", guid),
				zap.String("pos", pos),
				zap.Any("data", data),
				zap.Error(err),
			)
		}
		data.Items[itemGuid] = &itemservice.ItemInfo{
			ItemId: realyItemID,
			Count:  1,
			Guid:   itemGuid,
		}
		dropItems = append(dropItems, &itemservice.ItemInfo{
			ItemId: realyItemID,
			Count:  1,
			Guid:   itemGuid,
		})
	}

	err = data.Save(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "AddScoreItem Save Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("score", score),
			zap.Int64("guid", guid),
			zap.String("pos", pos),
			zap.Any("data", data),
		)
		return
	}

	logger.CtxInfo(ctx, "AddScoreItem Add Successful",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Any("nowdata", data),
		zap.Any("additems", dropItems),
	)

	// 推包
	if len(dropItems) > 0 {
		err = item.OnSendItemsPack(ctx, userID, dropItems, nil, guid, pos, 1)
		if err != nil {
			logger.CtxWarn(ctx, "AddScoreItem OnSendItemsPack Fail",
				zap.Uint64("userID", userID),
				zap.Int32("barrierID", barrierID),
				zap.Int32("score", score),
				zap.Int64("guid", guid),
				zap.String("pos", pos),
				zap.Any("dropItems", dropItems),
				zap.Error(err),
			)
			return
		}
	}

	return
}

func (s *service) AddItems(ctx context.Context, userID uint64, barrierID int32, items []*itemservice.ItemInfo, guid int64, pos string) error {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "AddItems Start",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Int64("guid", guid),
		zap.String("pos", pos),
		zap.Any("items", items),
	)

	defer func() {
		logger.CtxInfo(ctx, "AddItems End",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int64("guid", guid),
			zap.String("pos", pos),
			zap.Any("items", items),
		)
	}()

	data, err := barrieritemsmodel.NewBarrierItems(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "AddItems NewBarrierItems Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int64("guid", guid),
			zap.String("pos", pos),
			zap.Any("items", items),
			zap.Error(err),
		)
		return err
	}

	logger.CtxInfo(ctx, "AddItems GetData Successful",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Any("nowdata", data),
	)

	var dropItems []*itemservice.ItemInfo
	for _, item := range items {
		for i := 1; i <= int(item.Count); i++ {
			itemGuid, err := barrierguiditemredis.IncrNowGuid(ctx, userID, barrierID)
			if err != nil {
				logger.CtxError(ctx, "AddItems IncrNowGuid Fail",
					zap.Uint64("userID", userID),
					zap.Int32("barrierID", barrierID),
					zap.Int64("guid", guid),
					zap.String("pos", pos),
					zap.Any("data", data),
					zap.Error(err),
				)
			}
			data.Items[itemGuid] = &itemservice.ItemInfo{
				ItemId: item.ItemId,
				Count:  1,
				Guid:   itemGuid,
			}
			dropItems = append(dropItems, &itemservice.ItemInfo{
				ItemId: item.ItemId,
				Count:  1,
				Guid:   itemGuid,
			})
		}
	}

	err = data.Save(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "AddItems Save fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int64("guid", guid),
			zap.String("pos", pos),
			zap.Any("items", items),
		)
	}

	logger.CtxInfo(ctx, "AddItems Add Successful",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Any("nowdata", data),
		zap.Any("dropItems", dropItems),
	)

	// 推包
	if len(dropItems) > 0 {
		err = item.OnSendItemsPack(ctx, userID, dropItems, nil, guid, pos, 1)
		if err != nil {
			logger.CtxWarn(ctx, "AddScoreItem OnSendItemsPack Fail",
				zap.Uint64("userID", userID),
				zap.Int32("barrierID", barrierID),
				zap.Int64("guid", guid),
				zap.String("pos", pos),
				zap.Any("dropItems", dropItems),
				zap.Error(err),
			)
			return err
		}
	}

	return nil
}
