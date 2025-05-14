package game

import (
	"sort"
	"time"

	"github.com/lonng/nano/session"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/timeutil"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazechallengenumredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeuserbarrierredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserinfo"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeActionCountV8Cfg"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeBarriesV8Cfg"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeConfigV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGame"

	"go.uber.org/zap"
)

func (g *Game) OnMazeBarrierListRQ(s *session.Session, req *MazeGame.MazeBarrierListRQ) (err error) {
	fkprometheus.InfoPMT("OnMazeBarrierListRQ")()
	logger := fklog.AppLogger().Clone("game")

	res := &MazeGame.MazeBarrierListRS{}

	logger.InfoWF("OnMazeBarrierListRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnMazeBarrierListRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userId := uint64(s.UID())

	var maxNum, curNumToday int32
	maxNumCfg := GMazeActionCountV8Cfg.Get(101)
	if maxNumCfg == nil {
		logger.ErrorWF("OnMazeBarrierListRQ GMazeActionCountV8Cfg fail", zap.Error(err))
		res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
		return
	}
	maxNum = maxNumCfg.Day_count_v8

	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierListRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	userBarrier, err := mazeuserbarrierredis.GetUserBarrierInfo(logger, userId, 0)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierListRQ GetBarrierReport fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	now := time.Now()
	today := now.Year()*10000 + int(now.Month())*100 + now.Day()

	useNumToday, err := mazechallengenumredis.GetUserChallengeNum(logger, userId, today)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierListRQ GetBarrierReport fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	curNumToday = maxNum - int32(useNumToday)

	res.ChallengeInfo = &MazeCommon.MazeCount{
		CurCount:   proto.Int32(curNumToday),
		TotalCount: proto.Int32(maxNum),
	}
	res.ChallengeRefresh = proto.Int64(timeutil.GetTodayZeroTimestamp() + 86400)

	var passBarrierOffset, newBarrierOffset int32 = 1, 3
	barrierListCfg := GMazeConfigV8Cfg.Get(501)
	if barrierListCfg != nil {
		for k, v := range barrierListCfg.Value_map {
			//由于配置的是以当前尚未通关的最小关卡id为基准，但是用户未进入新关卡id时 存储的新关卡id不更新 所以这里以最大通关id为基准查找
			passBarrierOffset = k
			newBarrierOffset = int32(v)
			break
		}
	}

	totalBarrierNum := passBarrierOffset + newBarrierOffset + 1

	if userInfo.PassBarrier > 0 {
		i := userInfo.PassBarrier
		for {
			if i <= 0 || len(res.GetMazeBarrierList()) == int(passBarrierOffset) {
				break
			}
			sweepCfg := GMazeBarriesV8Cfg.Get(i)
			if sweepCfg == nil {
				logger.ErrorWF("OnMazeBarrierListRQ get pass barrier cfg fail", zap.Any("barrierId", i))
				res.ErrInfo = errors.CONFIG_NOT_FOUND.Wrap("获取扫荡关卡数值失败")
				return
			}
			sweepBarrier := &MazeGame.MazeBarrierInfo{
				BarrierId:        proto.Int32(i),
				BarrierName:      proto.String(sweepCfg.Name),
				MopCost:          proto.Int32(sweepCfg.Mop_cost),
				BarrierStatus:    proto.Int32(3), //3 已通关 可以扫荡
				ChallengeTimes:   proto.Int32(maxNum),
				ReChallengeTimes: proto.Int32(curNumToday),
			}
			res.MazeBarrierList = append(res.MazeBarrierList, sweepBarrier)
			i = sweepCfg.Last_id
		}
	}

	currBarrier := userInfo.Barrier
	if userInfo.Barrier == 0 {
		currBarrier = 1
	}
	if userInfo.Barrier == userInfo.PassBarrier && userInfo.Barrier > 0 {
		cfg := GMazeBarriesV8Cfg.Get(userInfo.Barrier)
		if cfg != nil {
			currBarrier = cfg.Next_id
		}
	}

	index := currBarrier

	for {
		if index == -1 {
			break
		}
		if len(res.GetMazeBarrierList()) >= int(totalBarrierNum) {
			break
		}
		barrierCfg := GMazeBarriesV8Cfg.Get(index)
		if barrierCfg == nil {
			logger.ErrorWF("OnMazeBarrierListRQ get barrier cfg fail", zap.Any("barrier", index))
			res.ErrInfo = errors.CONFIG_NOT_FOUND.Wrap("获取新关卡数值失败")
			return
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
			if userBarrier.GetBarrierStatus() == 2 {
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

		res.MazeBarrierList = append(res.MazeBarrierList, newBarrier)
		index = barrierCfg.Next_id
	}

	sort.Slice(res.GetMazeBarrierList(), func(i, j int) bool {
		return res.GetMazeBarrierList()[i].GetBarrierId() < res.GetMazeBarrierList()[j].GetBarrierId()
	})

	return
}
