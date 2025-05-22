// @Author pangchenyang 2025/3/21 17:49:00
// @Desc:
package sweep

import (
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEnergySvr"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGame"
	"gitlab.ifreetalk.com/maze-plate/protodef/MessageType"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/errors"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/uniqueid"
	"gitlab.ifreetalk.com/maze/maze_game_server/config/GMazeBarriesV8Cfg"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/calsweepbarrier"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserinfo"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/game/energy"
	"go.uber.org/zap"
)

// OnStartMazeSweepRQ start sweep
func OnStartMazeSweepRQ(ctx fknet.TCPContext, userID uint64, rqMsg, rsMsg proto.Message) (err error) {
	defer fkprometheus.InfoPMT("OnStartMazeSweepRQ")()

	req := rqMsg.(*MazeGame.StartMazeSweepRQ)
	res := rsMsg.(*MazeGame.StartMazeSweepRS)

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	ctx.InfoWF("OnStartMazeSweepRQ start", zap.Any("req", req))
	defer func() {
		ctx.InfoWF("OnStartMazeSweepRQ end", zap.Any("res", res), zap.Any("errMsg", string(res.GetErrInfo().GetErrMsg())))
	}()
	barrierId := req.GetBarrierId()
	res.BarrierId = req.BarrierId
	if userID <= 0 {
		ctx.ErrorWF("OnStartMazeSweepRQ userId invalid", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效的用户ID")
		return
	}
	// check barrier
	if req.GetBarrierId() <= 0 {
		ctx.ErrorWF("OnStartMazeSweepRQ req barrier invalid", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("未设置关卡id")
		return
	}

	barrierCfg := GMazeBarriesV8Cfg.Get(req.GetBarrierId())
	if barrierCfg == nil {
		ctx.ErrorWF("OnStartMazeSweepRQ get barrier cfg fail", zap.Int32("barrierId", barrierId))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("找不到该关卡配置")
		return
	}

	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userID)
	if err != nil {
		ctx.ErrorWF("OnStartMazeSweepRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	if barrierId > userInfo.PassBarrier {
		ctx.ErrorWF("OnStartMazeSweepRQ exceed maxUserBarrierID", zap.Any("req", req), zap.Int32("save", userInfo.Barrier))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("不能扫荡未通关的关卡")
		return
	}
	// check and cost energy
	remainVal, errInfo := SubSweepEnergy(ctx, userID, barrierId, barrierCfg.Mop_cost)
	res.RemainEnergy = proto.Int32(remainVal)
	if errInfo != nil && errInfo.GetErrCode() != errors.NO_ERROR_CODE {
		res.ErrInfo = errInfo
		return
	}
	// make gameID
	gameID := uniqueid.GenUniqueIdUInt64()
	res.GameId = proto.Uint64(gameID)

	// query sweep award
	// res.Awards, err = GetSweepAward(ctx, userID, barrierId)
	res.Awards, res.RareAward, err = calsweepbarrier.CalUserSweepBarrierAward(ctx, userID, barrierId, req.GetHeader())
	if err != nil {
		ctx.ErrorWF("OnStartMazeSweepRQ CalUserSweepBarrierAward fail", zap.Int32("barrierId", barrierId), zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("获取扫荡奖励失败")
		return
	}

	return nil
}

// SubSweepEnergy sub energy
func SubSweepEnergy(logger fklog.FKLogI, userID uint64, barrierId int32, subVal int32) (remainVal int32, errInfo *MessageType.ErrorInfo) {
	errInfo = errors.NO_ERROR
	subEnergyRq := &MazeEnergySvr.SubMazeEnergyRQ{
		UserId:      proto.Uint64(userID),
		OpType:      proto.Int32(int32(MazeEnergySvr.ENUM_MAZE_ENERGY_OP_TYPE_SWEEP)),
		SubVal:      proto.Int32(subVal),
		TradeNumber: proto.Uint64(uniqueid.GenUniqueIdUInt64()),
		OpDesc:      proto.String(fkconfig.GetServerConfig().ServerName),
	}
	subEnergyRs := &MazeEnergySvr.SubMazeEnergyRS{}
	// 合并服务，内聚接口
	// err := mazeenergyrpc.SubMazeEnergyRQ(logger, subEnergyRq, subEnergyRs)
	err := energy.SubMazeEnergyRQ(logger, userID, subEnergyRq, subEnergyRs)
	if err != nil {
		errInfo = errors.COMMON_ERROR_TIPS.Wrap("扣体力失败")
		logger.ErrorWF("SubSweepEnergy SubMazeEnergyRQ fail", zap.Int32("barrierId", barrierId),
			zap.Error(err))
		return subEnergyRs.GetRemainVal(), errInfo
	}
	if subEnergyRs.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
		if subEnergyRs.GetErrInfo().GetErrCode() == constdef.MAZE_ERR_ENERGY_LESS {
			errInfo = errors.COMMON_ERROR_TIPS.Wrap("体力不足")
			logger.WarnWF("SubSweepEnergy SubMazeEnergyRQ less energy", zap.Int32("barrierId", barrierId))
		} else {
			logger.ErrorWF("SubSweepEnergy SubMazeEnergyRQ fail", zap.Int32("barrierId", barrierId),
				zap.Any("err", errInfo), zap.String("errMsg", string(errInfo.GetErrMsg())))
			errInfo = subEnergyRs.GetErrInfo()
		}
		return subEnergyRs.GetRemainVal(), errInfo
	}
	return subEnergyRs.GetRemainVal(), errInfo
}

// func GetSweepAward(logger fklog.FKLogI, userID uint64, barrierId int32) ([]*Common.Item, error) {
// 	addExp, addMoney, awardItems, err := calsweepbarrier.CalUserSweepBarrierAward(logger, userID, barrierId)
// 	if err != nil {
// 		logger.ErrorWF("GetSweepAward CalUserSweepBarrierAward fail", zap.Int32("barrierId", barrierId), zap.Error(err))
// 		return nil, err
// 	}
// 	award := ItemsMapToList(awardItems)
// 	award = append(award, &Common.Item{
// 		ItemId: proto.Int32(constdef.MazeCommonItemExp), // todo modify expItemID
// 		Count:  proto.Int64(int64(addExp)),
// 	}, &Common.Item{
// 		ItemId: proto.Int32(constdef.MazeCommonItemCoin), // todo modify moneyItemID
// 		Count:  proto.Int64(addMoney),
// 	})
// 	return award, nil
// }
