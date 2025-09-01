package tempbuffservice

import (
	"context"
	"fmt"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/common/tradeno"
	"maze_game_server/excel/mazebarriesv8config"
	"maze_game_server/excel/mazeenergyresetcostv8config"
	"maze_game_server/model/tempbuffmodel"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/services/itemservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) RefreshOptionalMazeTempBuffList(ctx context.Context, userId uint64, barrierId, level, areaId, attrMask int32, cost []*MazeCommon.MazeItem) (*OptionalBuffInfo, error) {
	logger := fklog.ContextAppLogger(ctx)
	buffInfo, err := tempbuffmodel.NewTempBuffInfoModel(ctx, userId, barrierId)
	if err != nil {
		logger.CtxError(ctx, "RefreshOptionalMazeTempBuffListRQ GetMazeTempBuff failed", zap.Error(err))
		return nil, fmt.Errorf("获取用户buff信息失败")
	}

	if buffInfo == nil {
		logger.CtxError(ctx, "RefreshOptionalMazeTempBuffListRQ buff is nil", zap.Error(err))
		return nil, fmt.Errorf("获取用户buff信息失败")
	}

	// 是否可以刷新
	config := mazeenergyresetcostv8config.GetEnergyResetCostConfig(buffInfo.BuffSequence.RefreshCount + 1)
	if config == nil {
		logger.CtxError(ctx, "RefreshOptionalMazeTempBuffListRQ refresh config is nil", zap.Error(err))
		return nil, fmt.Errorf("获取刷新配置失败")
	}

	if !s.checkCost(config.Cost, cost) {
		logger.CtxError(ctx, "RefreshOptionalMazeTempBuffListRQ checkCost failed", zap.Error(err))
		return nil, fmt.Errorf("刷新消耗异常")
	}

	var newCost []*MazeCommon.MazeItem
	for _, item := range cost {
		if item.GetItemId() == 46200001 {
			continue
		}

		newCost = append(newCost, item)
	}

	logger.CtxInfo(ctx, "RefreshOptionalMazeTempBuffListRQ DeductItems start", zap.Any("cost", newCost))
	if len(newCost) > 0 {
		items := itemutil.ItemPb2ItemInfo(newCost)
		errInfo := itemservice.GlobalItemService.SubItem(ctx, userId, itemservice.ItemOpTypeRefreshTempBuff, tradeno.GetTradeNum(), items...)
		if errInfo != nil {
			logger.CtxError(ctx, "RefreshOptionalMazeTempBuffListRQ DeductItems failed", zap.Error(err))
			return nil, fmt.Errorf("扣钱失败")
		}
	}

	// 刷新可选buff
	err = s.refreshOptionalBuff(ctx, userId, barrierId, level, areaId, attrMask, buffInfo)
	if err != nil {
		logger.CtxError(ctx, "RefreshOptionalMazeTempBuffListRQ refreshOptionalBuff failed", zap.Error(err))
		return nil, fmt.Errorf("刷新buff失败")
	}

	optionalBuffInfo := s.packOptionalInfo(ctx, buffInfo)
	if optionalBuffInfo != nil {
		return optionalBuffInfo, nil
	}

	if len(buffInfo.BuffSequence.OptionalBuffList) == 0 {
		// 无buff可选
		return nil, fmt.Errorf("buff已全部选择")
	}

	return nil, fmt.Errorf("buff配置异常")
}

func (s *service) checkCost(costMap map[int32]int64, costList []*MazeCommon.MazeItem) bool {
	if len(costList) != len(costMap) {
		return false
	}

	for _, item := range costList {
		if count, ok := costMap[item.GetItemId()]; !ok || count != item.GetCount() {
			return false
		}
	}

	return true
}

func (s *service) refreshOptionalBuff(ctx context.Context, userId uint64, stageId, level, areaId, attrMask int32,
	buffInfo *tempbuffmodel.TempBuffInfoModel) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	stageConfig := mazebarriesv8config.GetStageConfig(ctx, stageId)
	if stageConfig == nil {
		logger.CtxError(ctx, "getOptionalBuffList stage config unknown", zap.Int32("stageId", stageId))
		return errors.New("关卡配置异常")
	}

	buffInfo.BuffSequence.RefreshCount = buffInfo.BuffSequence.RefreshCount + 1
	buffInfo.BuffSequence.OptionalBuffList, err = s.createOptionalBuffList(ctx, buffInfo, level, areaId, attrMask, stageConfig)
	if err != nil {
		logger.CtxError(ctx, "refreshOptionalBuff createOptionalBuffList failed", zap.Error(err))
		return err
	}

	err = buffInfo.Save(ctx, userId, stageId)
	if err != nil {
		logger.CtxError(ctx, "refreshOptionalBuff SetMazeTempBuff failed",
			zap.Int32("stageId", stageId), zap.Any("info", buffInfo), zap.Error(err))
		return err
	}

	return nil
}
