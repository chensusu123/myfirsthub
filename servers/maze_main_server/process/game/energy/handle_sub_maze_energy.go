/*
 * @Author: majian
 * @Date: 2025-03-21 21:24:40
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-27 22:23:15
 */
package energy

import (
	"time"

	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeconfigv8"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserinfo"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/protodef/MazeEnergySvr"
	"go.uber.org/zap"
)

// func OnSubMazeEnergyRQ(ctx fkrpc.RPCContext, shardingID int64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
// 	defer fkprometheus.DebugPMT("OnSubMazeEnergyRQ")()
// 	userCtx := fkserver.NewUserContext(ctx.Context, uint64(shardingID), ctx.FKLogI)
// 	req := rqMsg.(*MazeEnergySvr.SubMazeEnergyRQ)
// 	res := rsMsg.(*MazeEnergySvr.SubMazeEnergyRS)
// 	res.ErrInfo = errors.NO_ERROR
// 	userCtx.InfoWF("OnSubMazeEnergyRQ with", zap.Any("rq", req))

// 	defer func() {
// 		userCtx.InfoWF("OnSubMazeEnergyRQ end ", zap.Any("req", req), zap.Any("res", res))
// 	}()

// 	return SubMazeEnergyRQ(ctx, shardingID, req, res)
// }

func SubMazeEnergyRQ(logger fklog.FKLogI, shardingID int64, req *MazeEnergySvr.SubMazeEnergyRQ, res *MazeEnergySvr.SubMazeEnergyRS) (err error) {
	res.ErrInfo = errors.NO_ERROR

	if req.GetUserId() == 0 {
		logger.WarnWF("SubMazeEnergyRQ invalid userId ", zap.Any("rq", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("invalid userId")
		return nil
	}
	if req.GetTradeNumber() == 0 {
		logger.WarnWF("SubMazeEnergyRQ invalid tardeNo ", zap.Any("rq", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("invalid tardeNo")
		return nil
	}

	if req.GetOpType() == 0 {
		logger.WarnWF("SubMazeEnergyRQ invalid optype", zap.Any("rq", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("invalid optype")
		return nil
	}
	if req.GetSubVal() <= 0 {
		logger.WarnWF("SubMazeEnergyRQ invalid val", zap.Any("rq", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("invalid val")
		return nil
	}

	uInfo, err := mazeuserinfo.GetUserInfoV2(logger, req.GetUserId())
	if err != nil {
		logger.ErrorWF("SubMazeEnergyRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	if uInfo.Energy < req.GetSubVal() {
		logger.WarnWF("SubMazeEnergyRQ energy less",
			zap.Int32("has", uInfo.Energy),
			zap.Int32("need", req.GetSubVal()))
		res.RemainVal = proto.Int32(uInfo.Energy)
		res.ErrInfo = errors.NewErrorInfo(constdef.MAZE_ERR_ENERGY_LESS, "energy not enough")
		return
	}
	record := BeginRecord(req.GetUserId(), req.GetOpType(), uInfo)
	remain := uInfo.Energy - req.GetSubVal()
	if uInfo.Energy >= mazeconfigv8.GetEnergyMax() { // 如果满值时扣除，更新上次恢复时间
		uInfo.SetEnergyLastTime(time.Now().Unix())
	}
	uInfo.SetEnergy(remain)
	err = mazeuserinfo.SetUserInfoV2(logger, req.GetUserId(), uInfo)
	if err != nil {
		logger.ErrorWF("SubMazeEnergyRQ SetUserInfoV2 fail", zap.Error(err), zap.Any("uInfo", uInfo))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	res.RemainVal = proto.Int32(remain)
	EndRecord(logger, record, 0-req.GetSubVal(), uInfo)
	return nil
}
