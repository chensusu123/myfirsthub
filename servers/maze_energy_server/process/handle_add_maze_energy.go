/*
 * @Author: majian
 * @Date: 2025-03-27 14:04:26
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-27 14:14:55
 */
package process

import (
	"time"

	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkrpc"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/protodef/MazeEnergySvr"
	"gitlab.ifreetalk.com/servers/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/servers/maze_game_server/excel/mazeconfigv8"
	"gitlab.ifreetalk.com/servers/maze_game_server/module/mazeuserinfo"
	"go.uber.org/zap"
)

func OnAddMazeEnergyRQ(ctx fkrpc.RPCContext, shardingID int64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	defer fkprometheus.DebugPMT("OnAddMazeEnergyRQ")()
	userCtx := fkserver.NewUserContext(ctx.Context, uint64(shardingID), ctx.FKLogI)
	req := rqMsg.(*MazeEnergySvr.AddMazeEnergyRQ)
	res := rsMsg.(*MazeEnergySvr.AddMazeEnergyRS)
	res.ErrInfo = errors.NO_ERROR
	userCtx.InfoWF("OnAddMazeEnergyRQ with", zap.Any("rq", req))

	defer func() {
		userCtx.InfoWF("OnAddMazeEnergyRQ end ", zap.Any("req", req), zap.Any("res", res))
	}()
	if req.GetUserId() == 0 {
		userCtx.WarnWF("OnAddMazeEnergyRQ invalid userId ", zap.Any("rq", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("invalid userId")
		return nil
	}
	if req.GetTradeNumber() == 0 {
		userCtx.WarnWF("OnAddMazeEnergyRQ invalid tardeNo ", zap.Any("rq", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("invalid tardeNo")
		return nil
	}

	if req.GetOpType() == 0 {
		userCtx.WarnWF("OnAddMazeEnergyRQ invalid optype", zap.Any("rq", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("invalid optype")
		return nil
	}
	if req.GetAddVal() <= 0 {
		userCtx.WarnWF("OnAddMazeEnergyRQ invalid val", zap.Any("rq", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("invalid val")
		return nil
	}

	uInfo, err := mazeuserinfo.GetUserInfoV2(userCtx, req.GetUserId())
	if err != nil {
		userCtx.ErrorWF("OnAddMazeEnergyRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	maxVal := mazeconfigv8.GetEnergyMax() // 体力最大值
	if uInfo.Energy >= maxVal {
		userCtx.WarnWF("OnAddMazeEnergyRQ energy already full",
			zap.Int32("has", uInfo.Energy),
			zap.Int32("maxVal", maxVal))
		res.RemainVal = proto.Int32(uInfo.Energy)
		res.ErrInfo = errors.NewErrorInfo(constdef.MAZE_ERR_ENERGY_FULL, "energy already full")
		return
	}

	record := BeginRecord(req.GetUserId(), req.GetOpType(), uInfo)
	remain := uInfo.Energy + req.GetAddVal()
	if remain >= maxVal { // 如果加到满值，更新上次恢复时间
		remain = maxVal
		uInfo.SetEnergyLastTime(time.Now().Unix())
	}
	uInfo.SetEnergy(remain)
	err = mazeuserinfo.SetUserInfoV2(userCtx, req.GetUserId(), uInfo)
	if err != nil {
		userCtx.ErrorWF("OnAddMazeEnergyRQ SetUserInfoV2 fail", zap.Error(err), zap.Any("uInfo", uInfo))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	res.RemainVal = proto.Int32(remain)
	EndRecord(userCtx, record, req.GetAddVal(), uInfo)
	return nil
}
