/*
 * @Author: majian
 * @Date: 2025-03-21 19:14:33
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-27 14:01:38
 */
package energy

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeEnergy"
	"maze_game_server/services/barrierenergyservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (e *Energy) OnQueryMazeEnergyRQ_10469_10470(s *session.Session, req *MazeEnergy.QueryMazeEnergyRQ) (err error) {
	defer fkprometheus.DebugPMT("OnQueryMazeEnergyRQ")()

	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeEnergy.QueryMazeEnergyRS{}

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header

	userId := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnQueryMazeEnergyRQ end", zap.Any("res", res))
	}()

	logger.CtxInfo(ctx, "OnQueryMazeEnergyRQ with", zap.Any("req", req))

	energy, nextTime, err := barrierenergyservice.GlobalBarrierEnergyService.GetBarrierEnergy(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnQueryMazeEnergyRQ SetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.DB_SAVE_ERROR.ToInfo()
		return
	}
	maxValue := barrierenergyservice.GlobalBarrierEnergyService.GetEnergyMaxValue()

	res.EnergyInfo = &MazeEnergy.EnergyInfo{
		CurVal:           proto.Int32(energy),
		MaxVal:           proto.Int32(maxValue),
		NextRecoveryTime: proto.Int64(nextTime)}
	return
}

//func BeginRecord(userId uint64, opType int32, energyInfo *mazeuserinfo.UserInfo) *mazeenergyrecord.MazeEnergyChgRecord {
//	rd := &mazeenergyrecord.MazeEnergyChgRecord{}
//	rd.UserId = userId
//	rd.OldVal = energyInfo.Energy
//	rd.OpType = opType
//	return rd
//}
//
//func EndRecord(logger fklog.FKLogI, record *mazeenergyrecord.MazeEnergyChgRecord, chgval int32, energyInfo *mazeuserinfo.UserInfo) error {
//	record.NewVal = energyInfo.Energy
//	record.ChgVal = chgval
//	record.LastTime = energyInfo.EnergyLastTime
//	return mazeenergyrecord.SendMazeEnergyChgRecord(logger, record)
//}
