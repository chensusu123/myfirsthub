package itemservice

import (
	"context"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/pb/common/MessageType"
)

func (s *service) SubItem(ctx context.Context, userId uint64, opType ItemOpType, tradeNo uint64, items ...*ItemInfo) (errInfo *MessageType.ErrorInfo) {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "SubItem start", zap.Any("opType", opType), zap.Uint64("tradeNo", tradeNo), zap.Any("items", items))
	if userId <= 0 || opType <= 0 || tradeNo <= 0 {
		logger.CtxWarn(ctx, "SubItem invalid args", zap.Any("opType", opType))
		return
	}

	defer logger.CtxInfo(ctx, "SubItem end")
	defer fkprometheus.DebugPMT("SubItem")()

	// 检查并合并物品
	realSubItems := s.checkAndMergeItem(items)
	if len(realSubItems) == 0 {
		logger.CtxWarn(ctx, "SubItem sub items nil")
		return
	}

	subRes, errInfo := s.deductItems(ctx, userId, int32(opType), tradeNo, realSubItems...)
	if errInfo != nil && errInfo.GetErrCode() != errors.NO_ERROR_CODE {
		logger.CtxError(ctx, "SubItem deductItems err", zap.Any("errInfo", errInfo), zap.Any("realSubItems", realSubItems), zap.Any("subRes", subRes))
		return
	}
	logger.CtxInfo(ctx, "SubItem end", zap.Any("subRes", subRes))
	return
}

// deductItems 扣物品
func (s *service) deductItems(ctx context.Context, userId uint64, opType int32, tradeNo uint64, items ...*ItemInfo) (deductRes *AddItemRes, errInfo *MessageType.ErrorInfo) {
	logger := fklog.ContextAppLogger(ctx)
	deductRes = &AddItemRes{}

	groupItems, err := GloRegIns.GroupItemsByProcessor(ctx, items...)
	if err != nil {
		errInfo = errors.NewCommonCodeError(err.Error())
		deductRes.FailItem = append(deductRes.FailItem, items...)
		return
	}

	// 检查道具够不够
	var checkDuctRes *AddItemRes
	for handleProcess, opItems := range groupItems {
		if handleProcess == nil {
			logger.CtxError(ctx, "deductItems cannot find deduct check func", zap.Any("opItems", opItems), zap.Any("opType", opType))
			errInfo = errors.ITEM_CHECK_ERROR.ToInfo()
			return
		}

		checkDuctRes, errInfo = handleProcess.DeductItemCheck(ctx, userId, opType, tradeNo, opItems)
		if errInfo != nil && errInfo.GetErrCode() != errors.NO_ERROR_CODE {
			logger.CtxError(ctx, "deductItems check deduce item fail", zap.Any("opType", opType), zap.Any("opItems", opItems), zap.Any("errInfo", errInfo))
			deductRes.Merge(checkDuctRes)
			return
		}
	}

	// 扣道具
	var deductItems *AddItemRes
	for handProcess, opItems := range groupItems {
		deductItems, errInfo = handProcess.DeductItem(ctx, userId, opType, tradeNo, opItems)
		if errInfo != nil && errInfo.GetErrCode() != errors.NO_ERROR.GetErrCode() {
			logger.CtxError(ctx, "deductItems cost fail", zap.Any("items", opItems), zap.Any("errInfo", errInfo))
			deductRes.FailItem = append(deductRes.FailItem, deductItems.FailItem...)
			return
		}

		// 扣除成功的道具
		deductRes.SucItem = append(deductRes.SucItem, opItems...)
	}
	return
}
