package tempbuffservice

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/excel/mazebarriesv8config"
	"maze_game_server/excel/mazeenergyresetcostv8config"
	"maze_game_server/model/tempbuffmodel"
	"maze_game_server/module/itemmodule"
	"maze_game_server/pb/common/MazeCommon"
)

func (s *service) RefreshOptionalMazeTempBuffList(logger fklog.FKLogI, userId uint64, stageId, level, areaId int32, cost []*MazeCommon.MazeItem) (*OptionalBuffInfo, error) {
	buffInfo, err := tempbuffmodel.NewTempBuffInfoModel(logger, userId, stageId)
	if err != nil {
		logger.ErrorWF("RefreshOptionalMazeTempBuffListRQ GetMazeTempBuff failed", zap.Error(err))
		return nil, fmt.Errorf("获取用户buff信息失败")
	}

	if buffInfo == nil {
		logger.WarnWF("RefreshOptionalMazeTempBuffListRQ buff is nil", zap.Error(err))
		return nil, fmt.Errorf("获取用户buff信息失败")
	}

	// 是否可以刷新
	config := mazeenergyresetcostv8config.GetEnergyResetCostConfig(buffInfo.BuffSequence.RefreshCount + 1)
	if config == nil {
		logger.WarnWF("RefreshOptionalMazeTempBuffListRQ refresh config is nil", zap.Error(err))
		return nil, fmt.Errorf("获取刷新配置失败")
	}

	if !s.checkCost(config.Cost, cost) {
		logger.WarnWF("RefreshOptionalMazeTempBuffListRQ checkCost failed", zap.Error(err))
		return nil, fmt.Errorf("刷新消耗异常")
	}

	var newCost []*MazeCommon.MazeItem
	for _, item := range cost {
		if item.GetItemId() == 46200001 {
			continue
		}

		newCost = append(newCost, item)
	}

	logger.InfoWF("RefreshOptionalMazeTempBuffListRQ DeductItems start", zap.Any("cost", newCost))
	if len(newCost) > 0 {
		err = itemmodule.DeductItems(logger, userId, itemmodule.CostRefreshType, newCost)
		if err != nil {
			logger.ErrorWF("RefreshOptionalMazeTempBuffListRQ DeductItems failed", zap.Error(err))
			return nil, fmt.Errorf("扣钱失败")
		}
	}

	// 刷新可选buff
	err = s.refreshOptionalBuff(logger, userId, stageId, level, areaId, buffInfo)
	if err != nil {
		logger.ErrorWF("RefreshOptionalMazeTempBuffListRQ refreshOptionalBuff failed", zap.Error(err))
		return nil, fmt.Errorf("刷新buff失败")
	}

	optionalBuffInfo := s.packOptionalInfo(logger, buffInfo)
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

func (s *service) refreshOptionalBuff(logger fklog.FKLogI, userId uint64, stageId, level, areaId int32,
	buffInfo *tempbuffmodel.TempBuffInfoModel) (err error) {
	stageConfig := mazebarriesv8config.GetStageConfig(stageId)
	if stageConfig == nil {
		logger.WarnWF("getOptionalBuffList stage config unknown", zap.Int32("stageId", stageId))
		return errors.New("关卡配置异常")
	}

	buffInfo.BuffSequence.RefreshCount = buffInfo.BuffSequence.RefreshCount + 1
	buffInfo.BuffSequence.OptionalBuffList, err = s.createOptionalBuffList(logger, buffInfo, level, areaId, stageConfig)
	if err != nil {
		logger.ErrorWF("refreshOptionalBuff createOptionalBuffList failed", zap.Error(err))
		return err
	}

	err = buffInfo.Save(logger, userId, stageId)
	if err != nil {
		logger.ErrorWF("refreshOptionalBuff SetMazeTempBuff failed",
			zap.Int32("stageId", stageId), zap.Any("info", buffInfo), zap.Error(err))
		return err
	}

	return nil
}
