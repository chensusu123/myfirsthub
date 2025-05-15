/*
 * @Author: majian
 * @Date: 2025-03-25 10:50:37
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-25 15:03:02
 * @Desc 游戏复活
 */
package game

import (
	"sort"
	"time"

	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeRebornCostV8Cfg"
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGame"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/errors"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/gentradeno"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/itemutil"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/uniqueid"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeconfigv8"

	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazerebornkafka"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeuserbarrierredis"
	"go.uber.org/zap"
)

func OnMazeBarrierRebornRQ(logger fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	fkprometheus.InfoPMT("OnMazeBarrierRebornRQ")()

	req := rqMsg.(*MazeGame.MazeBarrierRebornRQ)
	res := rsMsg.(*MazeGame.MazeBarrierRebornRS)

	logger.InfoWF("OnMazeBarrierRebornRQ start", zap.Any("req", req))
	defer func() {
		logger.InfoWF("OnMazeBarrierRebornRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.RebornAck = req.RebornAck

	if req.GetRebornAck() != 1 && req.GetRebornAck() != 2 {
		logger.ErrorWF("OnMazeBarrierRebornRQ req type invalid", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("请求类型无效")
		return
	}

	barrierId := req.GetBarrierId()
	if barrierId <= 0 {
		logger.ErrorWF("OnMazeBarrierRebornRQ req barrier invalid", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡id未设置")
		return
	}

	barrierInfo, err := mazeuserbarrierredis.GetUserBarrierInfo(logger, shardingID, barrierId)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierRebornRQ GetUserBarrierInfo fail",
			zap.Any("req", req), zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	if barrierInfo == nil || barrierInfo.GetBarrierId() == 0 {
		logger.ErrorWF("OnMazeBarrierRebornRQ no barrier data", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("未找到关卡数据")
		return
	}

	if barrierInfo.GetBarrierId() != barrierId {
		logger.ErrorWF("OnMazeBarrierRebornRQ barrier data not match",
			zap.Int32("cliBarrierId", barrierId),
			zap.Int32("svrBarrireId", barrierInfo.GetBarrierId()))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡数据不匹配")
		return
	}
	if barrierInfo.GetEndTime() > 0 {
		logger.ErrorWF("OnMazeBarrierRebornRQ barrier challenge already end",
			zap.Int32("barrierId", barrierId), zap.Int64("endTime", barrierInfo.GetEndTime()))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("本次挑战已结束")
		return
	}

	if barrierInfo.GetStartTime() == 0 {
		logger.ErrorWF("OnMazeBarrierRebornRQ barrier challenge no start",
			zap.Int32("barrierId", barrierId))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("挑战未开始")
		return
	}

	svrCost, rebornMax, canReborn := GetReviveCost(barrierInfo.GetRebornCount() + 1)
	res.RebornMax = proto.Int32(rebornMax)
	res.RebornCount = proto.Int32(barrierInfo.GetRebornCount())
	if req.GetRebornAck() == 1 { // 查询
		if !canReborn {
			res.ErrInfo = errors.NewErrorInfo(ERROR_CODE_USER_CANT_REBORN, errMsg[ERROR_CODE_USER_CANT_REBORN])
			return
		}
		res.RebornAckTime = proto.Int64(mazeconfigv8.GetUserReviveTime() + time.Now().Unix())
		res.RebornCost = svrCost
	} else {
		if req.GetRebornAck() == 2 && len(svrCost) > 0 && len(req.GetRebornCost()) == 0 {
			logger.ErrorWF("OnMazeBarrierRebornRQ no cost param",
				zap.Any("rqCost", req.GetRebornCost()),
				zap.Any("svrCost", svrCost))
			res.ErrInfo = errors.ARGS_NOT_MATCH.ToInfo()
			return
		}

		nextCost, _, nextCanReborn := GetReviveCost(barrierInfo.GetRebornCount() + 2)
		res.RebornCost = nextCost
		if !nextCanReborn {
			res.RebornCost = []*MazeCommon.MazeItem{}
		}

		if !itemutil.CheckItemMatch(req.GetRebornCost(), svrCost) {
			logger.ErrorWF("OnMazeBarrierRebornRQ cost check fail",
				zap.Any("rqCost", req.GetRebornCost()),
				zap.Any("svrCost", svrCost))
			res.ErrInfo = errors.NewErrorInfo(ERROR_CODE_REBORN_COST_NOT_MATCH, "消耗不匹配")

			res.RebornAckTime = proto.Int64(mazeconfigv8.GetUserReviveTime() + time.Now().Unix())
			return
		}

		// 扣除消耗
		tid := uniqueid.GenUniqueIdUInt64()
		if len(svrCost) > 0 {
			// 通用	698	UN_CGK_COMMON_BILL_TYPE_698	迷宫挑战复活		否	马健	2025-03-25 13:48:10
			errInfo := gentradeno.DeductItemsEx(logger, shardingID, 698, tid, svrCost...)
			if errInfo != nil {
				logger.ErrorWF("OnMazeBarrierRebornRQ DeductItemsEx",
					zap.Any("svrCost", svrCost),
					zap.Uint64("tid", tid),
					zap.Any("errInfo", errInfo))
				if errInfo.GetErrCode() == 50049 {
					res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("物品不足")
				} else {
					res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("扣物品失败")
				}
				return
			}
		}

		// 更新复活次数
		barrierInfo.RebornCount = proto.Int32(barrierInfo.GetRebornCount() + 1)
		err = mazeuserbarrierredis.SetUserBarrierInfo(logger, shardingID, barrierInfo.GetBarrierId(), barrierInfo)
		if err != nil {
			logger.ErrorWF("OnMazeBarrierRebornRQ SetUserBarrierInfo fail", zap.Error(err),
				zap.Any("barrierInfo", barrierInfo))
			res.ErrInfo = errors.DB_SAVE_ERROR.ToInfo()
			return
		}

		res.RebornCount = proto.Int32(barrierInfo.GetRebornCount())

		record := &mazerebornkafka.MazeRebornRecord{
			UserId:      shardingID,
			Barrier:     barrierInfo.GetBarrierId(),
			RebornCount: int64(barrierInfo.GetRebornCount()),
			RebornCost:  itemutil.CommonItemsToString(svrCost),
		}
		mazerebornkafka.PushMazeRebornRecord(logger, record)
	}
	return nil
}

func GetReviveCost(reviveCnt int32) ([]*MazeCommon.MazeItem, int32, bool) {
	allRows := GMazeRebornCostV8Cfg.GetAllMazeRebornCostV8Config()
	if len(allRows) == 0 {
		return nil, 0, false
	}
	sort.Slice(allRows, func(i, j int) bool {
		return allRows[i].Order > allRows[j].Order
	})
	ret := make([]*MazeCommon.MazeItem, 0)
	var maxReborn int32

	for _, row := range allRows {
		if maxReborn < row.Reborn_max {
			maxReborn = row.Reborn_max
		}
		if reviveCnt <= row.Reborn_max && reviveCnt >= row.Reborn_min {
			for k, v := range row.Cost {
				if k <= 0 || v <= 0 {
					continue
				}
				ret = append(ret, &MazeCommon.MazeItem{
					ItemId: proto.Int32(k),
					Count:  proto.Int64(v)})
			}
			// break
		}
	}

	if len(ret) == 0 {
		return ret, maxReborn, false
	}

	return ret, maxReborn, true
}
