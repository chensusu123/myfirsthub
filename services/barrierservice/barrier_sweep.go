package barrierservice

import (
	"context"
	"maze_game_server/common/errors"
	"maze_game_server/common/tradeno"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/io/kafka/mazeenergyrecord"
	"maze_game_server/module/calsweepbarrier"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/Common"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeEnergy"
	"maze_game_server/pb/common/MessageType"
	"maze_game_server/services/barrierenergyservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// SweepBarrier implements BarrierService.
func (b *barrier) SweepBarrier(ctx context.Context, header *Common.PacketHeader, userID uint64, barrierID int32) (energyInfo *MazeEnergy.EnergyInfo, remainVal int32, gameID uint64, awardItem, rareItem []*MazeCommon.MazeItem, errinfo *MessageType.ErrorInfo) {
	logger := fklog.ContextAppLogger(ctx)
	barrierCfg := GMazeBarriesV8Cfg.Get(barrierID)
	if barrierCfg == nil {
		logger.ErrorWF("OnStartMazeSweepRQ get barrier cfg fail", zap.Int32("barrierId", barrierID))
		errinfo = errors.COMMON_ERROR_TIPS.Wrap("找不到该关卡配置")
		return
	}

	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userID)
	if err != nil {
		logger.ErrorWF("SweepBarrier GetUserInfoV2 fail", zap.Error(err))
		errinfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	energy, _, err := barrierenergyservice.GlobalBarrierEnergyService.GetBarrierEnergy(ctx, userID)
	if err != nil {
		errinfo = errors.COMMON_ERROR_TIPS.Wrap("获取体力信息失败")
		logger.ErrorWF("SweepBarrier GetBarrierEnergy fail", zap.Error(err), zap.Uint64("userID", userID))
		return
	}
	energyInfo = &MazeEnergy.EnergyInfo{
		CurVal:           proto.Int32(energy),
		MaxVal:           proto.Int32(barrierenergyservice.GlobalBarrierEnergyService.GetEnergyMaxValue()),
		NextRecoveryTime: proto.Int64(userInfo.EnergyLastTime),
	}

	if barrierID > userInfo.PassBarrier {
		errinfo = errors.COMMON_ERROR_TIPS.Wrap("不能扫荡未通关的关卡")
		logger.ErrorWF("SweepBarrier exceed maxUserBarrierID", zap.Any("barrierID", barrierID), zap.Int32("save", userInfo.Barrier))
		return
	}
	// check and cost energy
	remainVal, err = barrierenergyservice.GlobalBarrierEnergyService.SubEnergy(ctx, userID, barrierCfg.Mop_cost)
	if err != nil {
		errinfo = errors.COMMON_ERROR_TIPS.Wrap("体力不足")
		logger.ErrorWF("SweepBarrier SubEnergy fail", zap.Error(err), zap.Uint64("userID", userID), zap.Any("Mop_cost", barrierCfg.Mop_cost))
		return
	}
	// make gameID
	gameID = tradeno.GetTradeNum()

	energyInfo = &MazeEnergy.EnergyInfo{
		CurVal:           proto.Int32(remainVal),
		MaxVal:           proto.Int32(barrierenergyservice.GlobalBarrierEnergyService.GetEnergyMaxValue()),
		NextRecoveryTime: proto.Int64(userInfo.EnergyLastTime),
	}

	defer func() {
		barrierenergyservice.GlobalBarrierEnergyService.PushEnergyRecord(ctx, userID, energy, remainVal, mazeenergyrecord.SweepBarrier, userInfo.EnergyLastTime)
	}()

	// query sweep award
	awardItem, rareItem, err = calsweepbarrier.CalUserSweepBarrierAward(ctx, userID, barrierID, header)
	if err != nil {
		errinfo = errors.COMMON_ERROR_TIPS.Wrap("获取扫荡奖励失败")
		logger.ErrorWF("SweepBarrier CalUserSweepBarrierAward fail", zap.Int32("barrierId", barrierID), zap.Error(err))
		return
	}

	return energyInfo, remainVal, gameID, awardItem, rareItem, nil
}
