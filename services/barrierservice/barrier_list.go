package barrierservice

import (
	"context"
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeActionCountV8Cfg"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/config/GMazeConfigV8Cfg"
	"maze_game_server/io/redis/mazefixedbarrierredis"
	"maze_game_server/model/userbarriermodel"
	"maze_game_server/model/userinfomodel"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/pb/common/MessageType"
	"sort"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// GetBarrierInfos implements BarrierService.
func (b *barrier) GetBarrierInfos(ctx context.Context, userID uint64) (barrierInfos []*MazeGame.MazeBarrierInfo, errinfo *MessageType.ErrorInfo) {
	logger := fklog.ContextAppLogger(ctx)
	var err error
	var maxNum, curNumToday int32
	maxNumCfg := GMazeActionCountV8Cfg.GetWithCtx(ctx, 101)
	if maxNumCfg == nil {
		logger.CtxError(ctx, "GetBarrierInfos GMazeActionCountV8Cfg fail, config not found", zap.Int32("configID", 101))
		return nil, errors.CONFIG_NOT_FOUND.ToInfo()
	}
	maxNum = maxNumCfg.Day_count_v8

	userInfo, err := userinfomodel.NewUserInfoModel(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "GetBarrierInfos NewUserInfoModel fail", zap.Error(err))
		return nil, errors.MODULE_ERROR.ToInfo()
	}

	userBarrier, err := userbarriermodel.NewUserBarrierModel(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "GetBarrierInfos GetBarrierReport fail", zap.Error(err))
		return nil, errors.MODULE_ERROR.ToInfo()
	}

	// challengeNum, err := challengenummodel.NewChallengeNumModel(logger, userID)
	// if err != nil {
	// 	logger.CtxError(ctx,"GetBarrierInfos GetBarrierReport fail", zap.Error(err))
	// 	return nil, errors.MODULE_ERROR.ToInfo()
	// }
	// curNumToday = maxNum - int32(challengeNum.Num)

	// challengeInfo = &MazeCommon.MazeCount{
	// 	CurCount:   proto.Int32(curNumToday),
	// 	TotalCount: proto.Int32(maxNum),
	// }
	// challengeRefresh = timeutil.GetTodayZeroTimestamp() + 86400

	var passBarrierOffset, newBarrierOffset int32 = 1, 3
	barrierListCfg := GMazeConfigV8Cfg.GetWithCtx(ctx, 501)
	if barrierListCfg != nil {
		for k, v := range barrierListCfg.Value_map {
			//由于配置的是以当前尚未通关的最小关卡id为基准，但是用户未进入新关卡id时 存储的新关卡id不更新 所以这里以最大通关id为基准查找
			passBarrierOffset = k
			newBarrierOffset = int32(v)
			break
		}
	}

	currBarrier := userInfo.Barrier
	if userInfo.Barrier == 0 {
		currBarrier = 1
	}
	if userInfo.Barrier == userInfo.PassBarrier && userInfo.Barrier > 0 {
		cfg := GMazeBarriesV8Cfg.GetWithCtx(ctx, userInfo.Barrier)
		if cfg != nil {
			currBarrier = cfg.Next_id
		}
	}

	// 检查固定关卡
	fixedBarrierId, err := mazefixedbarrierredis.GetUserFixedBarrierID(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "GetBarrierInfos GetUserFixedBarrierID failed", zap.Error(err), zap.Uint64("userId", userID))
	} else if fixedBarrierId > 0 && currBarrier == 1 {
		currBarrier = fixedBarrierId
		logger.CtxWarn(ctx, "Fix current barrier", zap.Uint64("userId", userID), zap.Int32("currBarrier", currBarrier))
	}

	passBarrier := userInfo.PassBarrier
	if currBarrier > 1 {
		passBarrier = currBarrier - 1
	}

	totalBarrierNum := passBarrierOffset + newBarrierOffset + 1

	if passBarrier > 0 {
		i := passBarrier
		for {
			if i <= 0 || len(barrierInfos) == int(passBarrierOffset) {
				break
			}
			sweepCfg := GMazeBarriesV8Cfg.GetWithCtx(ctx, i)
			if sweepCfg == nil {
				logger.CtxError(ctx, "GetBarrierInfos get pass barrier cfg fail", zap.Any("barrierId", i))
				return nil, errors.CONFIG_NOT_FOUND.Wrap("获取扫荡关卡数值失败")
			}
			sweepBarrier := &MazeGame.MazeBarrierInfo{
				BarrierId:        proto.Int32(i),
				BarrierName:      proto.String(sweepCfg.Name),
				MopCost:          proto.Int32(sweepCfg.Mop_cost),
				BarrierStatus:    proto.Int32(3), //3 已通关 可以扫荡
				ChallengeTimes:   proto.Int32(maxNum),
				ReChallengeTimes: proto.Int32(curNumToday),
			}
			barrierInfos = append(barrierInfos, sweepBarrier)
			i = sweepCfg.Last_id
		}
	}

	index := currBarrier

	for {
		if index == -1 {
			break
		}
		if len(barrierInfos) >= int(totalBarrierNum) {
			break
		}
		barrierCfg := GMazeBarriesV8Cfg.GetWithCtx(ctx, index)
		if barrierCfg == nil {
			logger.CtxError(ctx, "GetBarrierInfos get barrier cfg fail", zap.Any("barrier", index))
			return nil, errors.CONFIG_NOT_FOUND.Wrap("获取新关卡数值失败")
		}

		var newBarrier *MazeGame.MazeBarrierInfo
		if barrierCfg.Order == currBarrier {
			newBarrier = &MazeGame.MazeBarrierInfo{
				BarrierId:        proto.Int32(barrierCfg.Order),
				BarrierName:      proto.String(barrierCfg.Name),
				MopCost:          proto.Int32(barrierCfg.Mop_cost),
				BarrierStatus:    proto.Int32(1),
				ChallengeTimes:   proto.Int32(maxNum),
				ReChallengeTimes: proto.Int32(curNumToday),
			}
			if userBarrier.BarrierStatus == 2 {
				newBarrier.BarrierStatus = proto.Int32(2)
			}
		} else {
			newBarrier = &MazeGame.MazeBarrierInfo{
				BarrierId:        proto.Int32(barrierCfg.Order),
				BarrierName:      proto.String(barrierCfg.Name),
				MopCost:          proto.Int32(barrierCfg.Mop_cost),
				BarrierStatus:    proto.Int32(0),
				ChallengeTimes:   proto.Int32(maxNum),
				ReChallengeTimes: proto.Int32(curNumToday),
			}
		}

		barrierInfos = append(barrierInfos, newBarrier)
		index = barrierCfg.Next_id
	}

	sort.Slice(barrierInfos, func(i, j int) bool {
		return barrierInfos[i].GetBarrierId() < barrierInfos[j].GetBarrierId()
	})

	return barrierInfos, errors.NO_ERROR
}
