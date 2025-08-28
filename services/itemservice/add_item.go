package itemservice

import (
	"context"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/pb/common/MessageType"
)

type ItemInfo struct {
	ItemId int32 `json:"item_id,omitempty"` //物品Id
	Count  int64 `json:"count,omitempty"`   //物品数量
}

func (s *service) AddItem(ctx context.Context, userId uint64, opType ItemOpType, tradeNo uint64, items ...*ItemInfo) (errInfo *MessageType.ErrorInfo) {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "AddItem start", zap.Any("opType", opType), zap.Uint64("tradeNo", tradeNo), zap.Any("items", items))
	if userId <= 0 || opType <= 0 || tradeNo <= 0 {
		logger.CtxError(ctx, "AddItem invalid args", zap.Any("opType", opType))
		return
	}

	// 检查并合并物品
	realAddItems := s.checkAndMergeItem(items)
	if len(realAddItems) <= 0 {
		logger.CtxWarn(ctx, "AddItem add items nil", zap.Any("addItems", items))
		return
	}

	addRes, err := s.addItems(ctx, userId, int32(opType), tradeNo, realAddItems...)
	if err != nil {
		logger.CtxError(ctx, "AddItem addItems err", zap.Any("gatherItem", realAddItems),
			zap.Any("addRes", addRes), zap.Error(err),
		)
		return
	}
	logger.CtxInfo(ctx, "AddItem end", zap.Any("addRes", addRes))
	return
}

func (s *service) addItems(ctx context.Context, userId uint64, opType int32, tradeNo uint64, items ...*ItemInfo) (addRes *AddItemRes, err error) {
	logger := fklog.ContextAppLogger(ctx)
	addRes = &AddItemRes{}

	groupItems, err := GloRegIns.GroupItemsByProcessor(ctx, items...)
	if err != nil {
		addRes.FailItem = append(addRes.FailItem, items...)
		return
	}

	var (
		gatherItems *AddItemRes
		addErr      error
	)
	for handleProcess, opItems := range groupItems {
		gatherItems, addErr = handleProcess.GatherItem(ctx, userId, opType, tradeNo, opItems...)
		if addErr != nil {
			logger.CtxError(ctx, "addItems GatherItem error", zap.Any("items", opItems), zap.Any("gatherRes", gatherItems), zap.Error(addErr))
			err = addErr
		}
		addRes.Merge(gatherItems)
	}

	return
}
