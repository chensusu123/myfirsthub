package buff

import (
	"maze_game_server/common/errors"
	"maze_game_server/excel/mazebarriesv8config"
	"maze_game_server/excel/mazeenergyaffixrandrulev8config"
	"maze_game_server/excel/mazeenergyresetcostv8config"
	"maze_game_server/io/redis/mazetempbuffredis"
	"maze_game_server/lib/log"
	"maze_game_server/module/itemmodule"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeTempBuff"
	"maze_game_server/pb/server/MazeTempBuffSvr"
	"time"

	"maze_game_server/lib/nano/session"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/24 15:20
 * @Description: 刷新迷宫可选buff列表
 */

func (b *Buff) RefreshOptionalMazeTempBuffListRQ_10439_10440(s *session.Session, req *MazeTempBuff.RefreshOptionalMazeTempBuffListRQ) (err error) {
	defer fkprometheus.InfoPMT("RefreshOptionalMazeTempBuffListRQ")()

	start := time.Now()

	logger := log.Clone("Buff", uint64(s.UID()), 0)
	res := &MazeTempBuff.RefreshOptionalMazeTempBuffListRS{}
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.StageId = req.StageId
	res.Level = req.Level
	defer func() {
		err = s.Response(res)
		logger.InfoWF("RefreshOptionalMazeTempBuffListRQ end", zap.Any("req", req), zap.Any("res", res),
			zap.Duration("costTime", time.Now().Sub(start)))
	}()

	userId, stageId, level, cost := uint64(s.UID()), req.GetStageId(), req.GetLevel(), req.GetCost()
	if userId == 0 || stageId == 0 || level == 0 {
		logger.WarnWF("RefreshOptionalMazeTempBuffListRQ args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数错误")
		return nil
	}

	// todo 检查用户是不是小程序用户

	buffInfo, err := mazetempbuffredis.GetMazeTempBuff(logger, userId, stageId)
	if err != nil {
		logger.ErrorWF("RefreshOptionalMazeTempBuffListRQ GetMazeTempBuff failed", zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("获取用户buff信息失败")
		return err
	}

	if buffInfo == nil {
		logger.WarnWF("RefreshOptionalMazeTempBuffListRQ buff is nil", zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("获取用户buff信息失败")
		return nil
	}

	// 是否可以刷新
	config := mazeenergyresetcostv8config.GetEnergyResetCostConfig(buffInfo.GetBuffSequence().GetRefreshCount() + 1)
	if config == nil {
		logger.WarnWF("RefreshOptionalMazeTempBuffListRQ refresh config is nil", zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("获取刷新配置失败")
		return nil
	}

	if !checkCost(config.Cost, cost) {
		logger.WarnWF("RefreshOptionalMazeTempBuffListRQ checkCost failed", zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("刷新消耗异常")
		return nil
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
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
			return err
		}
	}

	// 刷新可选buff
	err = refreshOptionalBuff(logger, userId, stageId, level, buffInfo)
	if err != nil {
		logger.ErrorWF("RefreshOptionalMazeTempBuffListRQ refreshOptionalBuff failed", zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("刷新buff失败")
		return err
	}

	res.OptionalBuffInfo = packOptionalInfo(logger, buffInfo)
	if res.GetOptionalBuffInfo() != nil {
		return nil
	}

	if len(buffInfo.GetBuffSequence().GetSelectBuffList()) == 0 {
		// 无buff可选
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("buff已全部选择")
		return nil
	}

	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("buff配置异常")
	return nil
}

func checkCost(costMap map[int32]int64, costList []*MazeCommon.MazeItem) bool {
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

func refreshOptionalBuff(logger fklog.FKLogI, userId uint64, stageId, level int32,
	buffInfo *MazeTempBuffSvr.TempBuffInfo) (err error) {
	stageConfig := mazebarriesv8config.GetStageConfig(stageId)
	if stageConfig == nil {
		logger.WarnWF("getOptionalBuffList stage config unknown", zap.Int32("stageId", stageId))
		return errors.New("关卡配置异常")
	}

	buffInfo.BuffSequence.RefreshCount = proto.Int32(buffInfo.GetBuffSequence().GetRefreshCount() + 1)
	// 生成可选的buff列表
	configId := mazeenergyaffixrandrulev8config.GetKey(stageConfig.Energy_affix_rand_rule, level)
	buffInfo.BuffSequence.SelectBuffList, err = createOptionalBuffList(logger, buffInfo, configId)
	if err != nil {
		logger.ErrorWF("refreshOptionalBuff createOptionalBuffList failed", zap.Error(err))
		return err
	}

	err = mazetempbuffredis.SetMazeTempBuff(logger, userId, stageId, buffInfo)
	if err != nil {
		logger.ErrorWF("refreshOptionalBuff SetMazeTempBuff failed",
			zap.Int32("stageId", stageId), zap.Any("info", buffInfo), zap.Error(err))
		return err
	}

	return nil
}
