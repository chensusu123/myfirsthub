package itemservice

import (
	"context"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/pb/common/MessageType"
)

func (s *service) QueryItems(ctx context.Context, userId uint64, items ...*ItemInfo) (queryItems []*ItemInfo, errInfo *MessageType.ErrorInfo) {
	logger := fklog.ContextAppLogger(ctx)

	defer logger.CtxInfo(ctx, "QueryItems end")

	if userId <= 0 {
		logger.CtxError(ctx, "QueryItems invalid args")
		return
	}
	queryItems = s.MergeQueryItems(items)
	if len(queryItems) == 0 {
		logger.CtxWarn(ctx, "QueryItems MergeQueryItems items nil")
		return
	}
	queryItems, errInfo = s.queryItems(ctx, userId, queryItems)
	if errInfo != nil {
		logger.CtxError(ctx, "QueryItems query items with err", zap.Any("errInfo", errInfo))
		return
	}

	return
}

func (s *service) MergeQueryItems(input []*ItemInfo) (output []*ItemInfo) {
	if len(input) == 0 {
		return
	}

	itemMap := make(map[int32]struct{}, len(input))
	for _, item := range input {
		if item.ItemId <= 0 {
			continue
		}

		itemMap[item.ItemId] = struct{}{}
	}

	for id := range itemMap {
		output = append(output, &ItemInfo{
			ItemId: id,
		})
	}
	return
}

func (s *service) queryItems(ctx context.Context, userId uint64, items []*ItemInfo) (resItems []*ItemInfo, errInfo *MessageType.ErrorInfo) {
	logger := fklog.ContextAppLogger(ctx)

	groupItems, err := GloRegIns.GroupItemsByProcessor(ctx, items...)
	if err != nil {
		return items, errors.NewCommonCodeError(err.Error())
	}

	for handleProcess, queryItems := range groupItems {
		if handleProcess == nil || len(queryItems) == 0 {
			logger.CtxError(ctx, "queryItems not fount processor", zap.Any("queryItems", queryItems))
			errInfo = errors.NewCommonCodeError("handle class nil")
			return
		}
		errInfo = handleProcess.GetItem(ctx, userId, queryItems)
		if errInfo != nil && errInfo.GetErrCode() != errors.NO_ERROR.GetErrCode() {
			logger.CtxError(ctx, "queryItems err", zap.Any("items", items), zap.Any("errInfo", errInfo))
			return
		}

		if len(queryItems) > 0 {
			resItems = append(resItems, queryItems...)
		}
	}
	return
}
