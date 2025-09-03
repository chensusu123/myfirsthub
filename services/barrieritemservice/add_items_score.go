package barrieritemservice

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/excel/mazebarriesv8config"
	"maze_game_server/excel/mazeconfigv8"
	"maze_game_server/model/barrieritemsmodel"
	"maze_game_server/servers/maze_main_server/process/item"
	"maze_game_server/services/itemservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) AddItemScore(ctx context.Context, userID uint64, barrierID int32, itemType int32, score int32, monsterGuid int64, monsterPos string) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "AddItemScore Start",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Int32("score", score),
		zap.Int64("monsterGuid", monsterGuid),
		zap.String("monsterPos", monsterPos),
	)

	defer func() {
		logger.CtxInfo(ctx, "AddItemScore End",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("score", score),
			zap.Int64("monsterGuid", monsterGuid),
			zap.String("monsterPos", monsterPos),
		)
	}()

	data, err := barrieritemsmodel.NewBarrierItems(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "AddItemScore NewBarrierItems Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("score", score),
			zap.Int64("monsterGuid", monsterGuid),
			zap.String("monsterPos", monsterPos),
			zap.Error(err),
		)
		return err
	}

	barrierCfg := mazebarriesv8config.GetStageConfig(ctx, barrierID)
	if barrierCfg.Need_item1_score == 0 || barrierCfg.Need_item2_score == 0 {
		logger.CtxError(ctx, "AddItemScore barrierCfg Need_item1_score or Need_item2_score Equal zero",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("score", score),
			zap.Int64("monsterGuid", monsterGuid),
			zap.String("monsterPos", monsterPos),
			zap.Any("data", data),
		)
		return nil
	}

	dropItems := make([]*itemservice.ItemInfo, 0)

	var realyItemID int32
	if itemType == constdef.MazeCfgId901 {
		itemMap := mazeconfigv8.GetReallyItemMap(ctx, itemType)
		for k := range itemMap {
			realyItemID = k
		}

		nowScore := data.ItemsScore[realyItemID] + score
		data.ItemsScore[realyItemID] = nowScore % barrierCfg.Need_item1_score

		itemNum := nowScore / barrierCfg.Need_item1_score
		data.Items[int64(realyItemID)] += int64(itemNum)
		dropItems = append(dropItems, &itemservice.ItemInfo{
			ItemId: realyItemID,
			Count:  int64(itemNum),
		})
	} else if itemType == constdef.MazeCfgId902 {
		itemMap := mazeconfigv8.GetReallyItemMap(ctx, itemType)
		for k := range itemMap {
			realyItemID = k
		}

		nowScore := data.ItemsScore[realyItemID] + score
		data.ItemsScore[realyItemID] = nowScore % barrierCfg.Need_item2_score

		itemNum := nowScore / barrierCfg.Need_item2_score
		data.Items[int64(realyItemID)] += int64(itemNum)
		dropItems = append(dropItems, &itemservice.ItemInfo{
			ItemId: realyItemID,
			Count:  int64(itemNum),
		})
	} else {
		data.Items[int64(itemType)] += int64(score)
		dropItems = append(dropItems, &itemservice.ItemInfo{
			ItemId: realyItemID,
			Count:  int64(score),
		})
	}

	err = data.Save(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "AddItemScore Save Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("score", score),
			zap.Int64("monsterGuid", monsterGuid),
			zap.String("monsterPos", monsterPos),
			zap.Any("data", data),
		)
		return
	}

	// 推包
	err = item.OnSendItemsPack(ctx, userID, dropItems, make([]*itemservice.ItemInfo, 0), monsterGuid, monsterPos)
	if err != nil {
		logger.CtxWarn(ctx, "AddEquipScore OnSendItemsPack Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("score", score),
			zap.Int64("monsterGuid", monsterGuid),
			zap.String("monsterPos", monsterPos),
			zap.Any("dropItems", dropItems),
			zap.Error(err),
		)
		return
	}

	return err
}
