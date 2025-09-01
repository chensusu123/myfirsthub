package bagservice

import (
	"context"
	"time"

	"maze_game_server/common/errors"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/model/bagmodel"
	"maze_game_server/pb/common/MazeBag"
	"maze_game_server/pb/common/MessageType"
	"maze_game_server/services/itemservice"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// GatherItem 加道具
func (s *service) GatherItem(ctx context.Context, userId uint64, opType int32, tradeNo uint64, items ...*itemservice.ItemInfo) (addRes *itemservice.AddItemRes, err error) {
	addRes = &itemservice.AddItemRes{}
	logger := fklog.ContextAppLogger(ctx)
	bagModel, err := bagmodel.NewBagModel(ctx, userId)
	if err != nil {
		return nil, err
	}
	var sucItems []*itemservice.ItemInfo
	for idx, item := range items {
		currItemCount, addErr := bagModel.IncrValue(ctx, userId, item.ItemId, item.Count)
		if addErr != nil {
			logger.CtxError(ctx, "bag GatherItem IncrValue err",
				zap.Int32("itemID", item.ItemId),
				zap.Int64("count", item.Count),
				zap.Error(addErr),
			)
			err = addErr
			addRes.FailItem = append(addRes.FailItem, items[idx])
			continue
		}

		addRes.SucItem = append(addRes.SucItem, items[idx])

		sucItems = append(sucItems, &itemservice.ItemInfo{
			ItemId: item.ItemId,
			Count:  currItemCount,
		})
	}

	// 变化通知ID包
	s.sendBagItemChgID(ctx, userId, sucItems)

	// 流水 TODO
	return
}

// DeductItem 扣道具
func (s *service) DeductItem(ctx context.Context, userId uint64, opType int32, tradeNo uint64, items []*itemservice.ItemInfo) (deductRes *itemservice.AddItemRes, errInfo *MessageType.ErrorInfo) {
	deductRes = &itemservice.AddItemRes{}
	var sucItems []*itemservice.ItemInfo // 推id包使用

	bagModel, err := bagmodel.NewBagModel(ctx, userId)
	if err != nil {
		errInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	logger := fklog.ContextAppLogger(ctx)
	for _, item := range items {
		curCount, err := bagModel.IncrValue(ctx, userId, item.ItemId, -item.Count)
		if err != nil {
			logger.CtxError(ctx, "DeductItem IncrValue err", zap.Any("item", item), zap.Error(err))
			errInfo = errors.MODULE_ERROR.ToInfo()
			deductRes.FailItem = append(deductRes.FailItem, item)
			continue
		}
		sucItems = append(sucItems, &itemservice.ItemInfo{
			ItemId: item.ItemId,
			Count:  curCount,
		})
		deductRes.SucItem = append(deductRes.SucItem, item)
	}

	// 变化通知ID包
	s.sendBagItemChgID(ctx, userId, sucItems)

	// 流水 TODO
	return
}

// DeductItemCheck 扣道具检查
func (s *service) DeductItemCheck(ctx context.Context, userId uint64, opType int32, tradeNo uint64, items []*itemservice.ItemInfo) (checkRes *itemservice.AddItemRes, errInfo *MessageType.ErrorInfo) {
	checkRes = &itemservice.AddItemRes{}

	queryIds := make([]int32, 0, len(items))
	for _, item := range items {
		queryIds = append(queryIds, item.ItemId)
	}
	logger := fklog.ContextAppLogger(ctx)
	bagModel, err := bagmodel.NewBagModel(ctx, userId)
	if err != nil {
		errInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	dbCountMap, err := bagModel.BatchLoad(ctx, userId, queryIds)
	if err != nil {
		logger.CtxError(ctx, "bag DeductItemCheck BatchLoad err", zap.Int32s("queryIds", queryIds), zap.Error(err))
		errInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	var bagCount int64
	for _, item := range items {
		bagCount = dbCountMap[item.ItemId]
		if bagCount < item.Count {
			logger.CtxWarn(ctx, "bag DeductItemCheck item not enough", zap.Int64("bagCount", bagCount), zap.Any("item", item))
			errInfo = errors.ITEM_CHECK_ITEM_NOT_ENOUGH.ToInfo()
			return
		}
	}
	return
}

// GetItem 查询数据
func (s *service) GetItem(ctx context.Context, userId uint64, items []*itemservice.ItemInfo) (errInfo *MessageType.ErrorInfo) {
	queryIds := make([]int32, 0, len(items))
	for _, item := range items {
		queryIds = append(queryIds, item.ItemId)
	}
	logger := fklog.ContextAppLogger(ctx)

	bagModel, err := bagmodel.NewBagModel(ctx, userId)
	if err != nil {
		errInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	dbCountMap, err := bagModel.BatchLoad(ctx, userId, queryIds)
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

func (s *service) sendBagItemChgID(ctx context.Context, userId uint64, items []*itemservice.ItemInfo) {
	if len(items) == 0 {
		return
	}
	// logger := fklog.ContextAppLogger(ctx)
	idPack := &MazeBag.MazeBagChgID{
		Token: proto.Int64(time.Now().UnixMilli()),
		Items: make([]*MazeBag.MazeBagItem, 0, len(items)),
	}
	for _, item := range items {
		idPack.Items = append(idPack.Items, itemutil.BuildMazeBagItem(ctx, item.ItemId, item.Count))
	}
	_ = online.PushWithContext(ctx, userId, 10404, idPack)
}
