package game

import (
	"context"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/addequip"
	"maze_game_server/common/tradeno"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/pb/server/MazeEquipSvr"
	"maze_game_server/services/barrierscorerewardservice"
	"maze_game_server/services/equipdropservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

func (g *Game) OnReportAwardFoeEquipRQ_10455_10456(s *session.Session, req *MazeGame.ReportAwardFoeEquipRQ) (err error) {
	defer fkprometheus.InfoPMT("OnReportAwardFoeEquipRQ")()

	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeGame.ReportAwardFoeEquipRS{}

	logger.InfoWF("OnReportAwardFoeEquipRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnReportAwardFoeEquipRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userId := uint64(s.UID())

	equipNum := req.GetEquipNum()
	if equipNum <= 0 {
		logger.ErrorWF("OnReportAwardFoeEquipRQ equipNum fail", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("请求的加装备数量为0")
		return
	}

	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.ErrorWF("OnReportAwardFoeEquipRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	level := userInfo.Level

	if req.GetBarrierId() != userInfo.Barrier {
		logger.ErrorWF("OnReportAwardFoeEquipRQ check barrier fail", zap.Error(err), zap.Any("req", req.GetBarrierId()), zap.Any("save", userInfo.Barrier))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("请求的关卡id和上报存储的关卡不一致")
		return
	}

	//addEquipMap, err := calequipsequence.GetNewEquip(logger, userId, req.GetBarrierId(), int32(level), equipNum)
	//if err != nil {
	//	logger.ErrorWF("OnReportAwardFoeEquipRQ GetMazeShopInfo fail", zap.Error(err), zap.Any("barrier", req.GetBarrierId()), zap.Any("level", level))
	//	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	//	return
	//}
	addEquipMap, err := equipdropservice.GlobalEquipDropService.GetNewEquip(ctx, userId, int32(level), req.GetBarrierId(), equipNum)
	if err != nil {
		logger.ErrorWF("OnReportAwardFoeEquipRQ GetMazeShopInfo fail", zap.Error(err), zap.Any("barrier", req.GetBarrierId()), zap.Any("level", level))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	tradeNo := tradeno.GetTradeNum()
	//rs, err2 := addequip.InstanceEquip(logger, userId, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_FOE), tradeNo, equipNumPerCycle, addEquipMap)
	rs, err2 := addequip.AddEquipToBag(ctx, userId, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_FOE), tradeNo, addEquipMap)
	if err2 != nil {
		logger.ErrorWF("OnReportAwardFoeEquipRQ addEquipToBag fail", zap.Error(err2), zap.Any("optype", MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_FOE),
			zap.Any("tradeNo", tradeNo), zap.Any("addEquip", addEquipMap), zap.Any("rs", rs))
	}
	// PushDollMazeShopInfoLog(logger, userId, int32(level), shopInfo, rs.EquipList, 19, tradeNo, 0)

	if err = barrierscorerewardservice.GlobalScoreRewardService.SaveBarrierScoreRewardEquip(context.TODO(), userId, req.GetBarrierId(), addEquipMap); err != nil {
		logger.ErrorWF("OnReportAwardFoeEquipRQ SaveBarrierScoreRewardEquip err", zap.Error(err), zap.Any("barrier", req.GetBarrierId()))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	}

	return
}
