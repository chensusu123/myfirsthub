package moneyservice

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/io/kafka/mazemoneykafka"
	"maze_game_server/model/moneymodel"
	"maze_game_server/module/mazecommonvalue"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/pb/common/MessageType"
	"maze_game_server/services/itemservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// GatherItem 加道具
func (s *service) GatherItem(ctx context.Context, userId uint64, opType int32, tradeNo uint64, items ...*itemservice.ItemInfo) (addRes *itemservice.AddItemRes, err error) {
	addRes = &itemservice.AddItemRes{}
	logger := fklog.ContextAppLogger(ctx)

	logger.CtxInfo(ctx, "GatherItem with", zap.Any("opType", opType), zap.Any("items", items))

	model, err := moneymodel.NewMoneyModel(ctx, userId)
	if err != nil {
		return nil, err
	}
	oldMap := make(map[int32]int64)
	moneyMap := make(map[int32]int64)
	for idx, item := range items {
		currItemCount, addErr := model.IncrValue(ctx, userId, item.ItemId, item.Count)
		if addErr != nil {
			logger.CtxError(ctx, "money GatherItem IncrValue err",
				zap.Int32("itemID", item.ItemId),
				zap.Int64("count", item.Count),
				zap.Error(addErr),
			)
			err = addErr
			addRes.FailItem = append(addRes.FailItem, items[idx])
			continue
		}
		moneyMap[item.ItemId] = currItemCount
		oldMap[item.ItemId] = currItemCount - item.Count
		addRes.SucItem = append(addRes.SucItem, items[idx])
	}

	// 货币变化推包
	s.sendMoneyItemChgID(logger, userId, moneyMap)
	// 流水
	s.sendFlow(ctx, userId, items, oldMap, moneyMap, tradeNo, opType)

	return
}

// DeductItem 扣道具
func (s *service) DeductItem(ctx context.Context, userId uint64, opType int32, tradeNo uint64, items []*itemservice.ItemInfo) (deductRes *itemservice.AddItemRes, errInfo *MessageType.ErrorInfo) {
	deductRes = &itemservice.AddItemRes{}
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "DeductItem with", zap.Any("opType", opType), zap.Any("items", items))

	model, err := moneymodel.NewMoneyModel(ctx, userId)
	if err != nil {
		return nil, errors.MODULE_ERROR.ToInfo()
	}
	oldMap, err := model.LoadAll(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "DeductItem LoadAll err", zap.Error(err))
		return nil, errors.MODULE_ERROR.ToInfo()
	}
	moneyMap := make(map[int32]int64)
	for _, item := range items {
		if oldMap[item.ItemId] < item.Count {
			logger.CtxWarn(ctx, "money DeductItem item not enough", zap.Int64("bagCount", oldMap[item.ItemId]), zap.Any("item", item))
			return nil, errors.ITEM_CHECK_ITEM_NOT_ENOUGH.ToInfo()
		}
		moneyMap[item.ItemId] = oldMap[item.ItemId] - item.Count
	}
	err = model.BatchSet(ctx, userId, moneyMap)
	if err != nil {
		logger.CtxError(ctx, "DeductItem BatchSet err", zap.Error(err))
		return nil, errors.MODULE_ERROR.ToInfo()
	}

	// 货币变化推包
	s.sendMoneyItemChgID(logger, userId, moneyMap)
	// 流水
	s.sendFlow(ctx, userId, items, oldMap, moneyMap, tradeNo, opType)

	deductRes.SucItem = items
	return
}

// DeductItemCheck 扣道具检查
func (s *service) DeductItemCheck(ctx context.Context, userId uint64, opType int32, tradeNo uint64, items []*itemservice.ItemInfo) (checkRes *itemservice.AddItemRes, errInfo *MessageType.ErrorInfo) {
	checkRes = &itemservice.AddItemRes{}
	logger := fklog.ContextAppLogger(ctx)
	queryIds := make([]int32, 0, len(items))
	for _, item := range items {
		queryIds = append(queryIds, item.ItemId)
	}
	model, err := moneymodel.NewMoneyModel(ctx, userId)
	if err != nil {
		errInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	dbCountMap, err := model.BatchLoad(ctx, userId, queryIds)
	if err != nil {
		logger.CtxError(ctx, "DeductItemCheck BatchLoad err", zap.Int32s("queryIds", queryIds), zap.Error(err))
		errInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	var bagCount int64
	for _, item := range items {
		bagCount = dbCountMap[item.ItemId]
		if bagCount < item.Count {
			logger.CtxWarn(ctx, "DeductItemCheck item not enough", zap.Int64("bagCount", bagCount), zap.Any("item", item))
			errInfo = errors.ITEM_CHECK_ITEM_NOT_ENOUGH.ToInfo()
			return
		}
	}
	return
}

// GetItem 查询数据
func (s *service) GetItem(ctx context.Context, userId uint64, items []*itemservice.ItemInfo) (errInfo *MessageType.ErrorInfo) {
	logger := fklog.ContextAppLogger(ctx)
	queryIds := make([]int32, 0, len(items))
	for _, item := range items {
		queryIds = append(queryIds, item.ItemId)
	}

	model, err := moneymodel.NewMoneyModel(ctx, userId)
	if err != nil {
		errInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	dbCountMap, err := model.BatchLoad(ctx, userId, queryIds)
	if err != nil {
		logger.CtxError(ctx, "GetItem BatchLoad err", zap.Int32s("queryIds", queryIds), zap.Error(err))
		errInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	for _, item := range items {
		item.Count = dbCountMap[item.ItemId]
	}
	return
}

// 通知货币变化
func (s *service) sendMoneyItemChgID(logger fklog.FKLogI, userId uint64, moneyMap map[int32]int64) {
	commonList := make([]*mazecommonvalue.CommonValueStruct, 0)
	commonValueMap := make(map[int32]int64)
	for k, v := range moneyMap {
		if k == constdef.MazeCommonItemCoin {
			commonValueMap[int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_MONEY)] = v
		} else if k == constdef.MazeCommonItemDiamond {
			commonValueMap[int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_DIAMOND)] = v
		}
	}
	moneyCommon := mazecommonvalue.MakeCommonValueList(logger, commonValueMap, map[int32]int32{}, map[int32]string{})
	commonList = append(commonList, moneyCommon...)
	mazecommonvalue.SendCommonValueIdPack(logger, userId, commonList)
}

// 流水
func (s *service) sendFlow(ctx context.Context, userId uint64, items []*itemservice.ItemInfo, oldMap, moneyMap map[int32]int64, tradeNo uint64, opType int32) {
	for _, i := range items {
		record := &mazemoneykafka.MazeMoneyRecord{
			UserId:        userId,
			OldMoneyId:    i.ItemId,
			OldMoneyCount: oldMap[i.ItemId],
			NewMoneyId:    i.ItemId,
			NewMoneyCount: moneyMap[i.ItemId],
			TradeNo:       tradeNo,
			ChgReason:     opType,
		}
		mazemoneykafka.PushMazeMoneyRecord(ctx, record)
	}
}
