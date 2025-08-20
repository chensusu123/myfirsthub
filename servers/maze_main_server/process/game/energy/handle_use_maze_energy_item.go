package energy

import (
	"context"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/uniqueid"
	"maze_game_server/io/kafka/mazeenergyrecord"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeEnergy"
	"maze_game_server/services/barrierenergyservice"
	"maze_game_server/services/itemservice"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (e *Energy) OnUseMazeEnergyItemRQ_10611_10612(s *session.Session, req *MazeEnergy.UseMazeEnergyItemRQ) (err error) {
	defer fkprometheus.DebugPMT("OnUseMazeEnergyItemRQ")()

	logger := log.Clone("Energy", uint64(s.UID()), 0)
	res := &MazeEnergy.UseMazeEnergyItemRS{}

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header

	userId := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnUseMazeEnergyItemRQ end", zap.Any("res", res))
	}()

	logger.InfoWF("OnUseMazeEnergyItemRQ with", zap.Any("req", req))

	uInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("OnUseMazeEnergyItemRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	oldEnergy := uInfo.Energy
	maxVal := barrierenergyservice.GlobalBarrierEnergyService.GetEnergyMaxValue() // 体力最大值

	if uInfo.Energy >= maxVal {
		// 补发一个id包
		recoverTime := barrierenergyservice.GlobalBarrierEnergyService.GetEnergyRecoverCfg() // 迷宫体力回复间隔时间（秒）
		nextTime := uInfo.EnergyLastTime + recoverTime - time.Now().Unix()
		err = barrierenergyservice.GlobalBarrierEnergyService.SendEnergyChgPack(logger, userId, uInfo.Energy, nextTime)
		if err != nil {
			logger.ErrorWF("OnUseMazeEnergyItemRQ SendEnergyChgPack failed", zap.Error(err))
			//return err
			err = nil
		}

		logger.WarnWF("OnUseMazeEnergyItemRQ energy already full", zap.Int32("has", uInfo.Energy), zap.Int32("maxVal", maxVal))
		res.EnergyInfo.CurVal = proto.Int32(uInfo.Energy)
		res.EnergyInfo.MaxVal = proto.Int32(maxVal)
		res.EnergyInfo.NextRecoveryTime = proto.Int64(nextTime)
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("energy already full")
		return err
	}

	// 扣物品
	itemId, recoverNum := barrierenergyservice.GlobalBarrierEnergyService.GetEnergyItemCfg()
	//careCost := make([]*MazeCommon.MazeItem, 1)
	//careCost = append(careCost, &MazeCommon.MazeItem{
	//	ItemId: proto.Int32(itemId),
	//	Count:  proto.Int64(1),
	//})

	tid := uniqueid.GenUniqueIdUInt64()
	careCost := &itemservice.ItemInfo{
		ItemId: itemId,
		Count:  1,
	}
	errInfo := itemservice.GlobalItemService.SubItem(context.TODO(), userId, itemservice.ItemOpTypeUseEnergy, tid, careCost)
	if errInfo != nil {
		logger.ErrorWF("OnUseMazeEnergyItemRQ DeductItemsEx", zap.Any("careCost", careCost), zap.Uint64("tid", tid), zap.Any("errInfo", errInfo))

		if errInfo.GetErrCode() == 50049 {
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("物品不足")
		} else {
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("扣物品失败")
		}
		return
	}

	logger.InfoWF("OnUseMazeEnergyItemRQ DeductItemsEx succ", zap.Any("careCost", careCost), zap.Uint64("tid", tid))

	energy, nextTime, err := barrierenergyservice.GlobalBarrierEnergyService.AddEnergy(logger, userId, recoverNum)
	if err != nil {
		logger.ErrorWF("OnUseMazeEnergyItemRQ AddEnergy fail", zap.Error(err))
		return err
	}

	res.EnergyInfo = &MazeEnergy.EnergyInfo{
		CurVal:           proto.Int32(energy),
		MaxVal:           proto.Int32(maxVal),
		NextRecoveryTime: proto.Int64(nextTime)}

	defer func() {
		barrierenergyservice.GlobalBarrierEnergyService.PushEnergyRecord(logger, userId, oldEnergy, energy, mazeenergyrecord.ItemEnergy, uInfo.EnergyLastTime)
	}()

	return
}
