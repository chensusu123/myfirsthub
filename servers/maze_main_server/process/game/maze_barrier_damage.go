package game

import (
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/services/barrierstagecounterservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

func (g *Game) OnBarrierDamageRQ_10622_10623(s *session.Session, req *MazeGame.BarrierDamageRQ) (err error) {
	defer fkprometheus.InfoPMT("OnBarrierDamageRQ")()
	ctx := s.Context()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeGame.BarrierDamageRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.StageId = req.StageId
	res.BarrierId = req.BarrierId

	logger.CtxInfo(ctx, "OnBarrierDamageRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnBarrierDamageRQ end", zap.Any("res", res))
	}()

	userId := uint64(s.UID())

	if req.GetStageId() < 0 || req.GetBarrierId() <= 0 {
		logger.CtxError(ctx, "OnBarrierDamageRQ req barrier or areaId invalid", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡id未设置")
		return
	}
	barrierCfg := GMazeBarriesV8Cfg.Get(req.GetBarrierId())
	if barrierCfg == nil {
		logger.CtxError(ctx, "OnBarrierDamageRQ get barrier cfg fail", zap.Any("barrier", req.GetStageId()))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("找不到该关卡配置")
		return
	}

	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnBarrierDamageRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	if userInfo.Barrier != req.GetBarrierId() {
		logger.CtxError(ctx, "OnBarrierDamageRQ barrier err", zap.Any("req", req), zap.Any("barrier", userInfo.Barrier))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡id错误")
		return
	}

	totalDamage, err := barrierstagecounterservice.GlobalBarrierStageCounterService.AddDamage(ctx, userId, req.GetBarrierId(), req.GetStageId(), req.GetDamage())
	if err != nil {
		logger.CtxError(ctx, "OnBarrierDamageRQ AddDamage fail", zap.Any("req", req), zap.Error(err))
		return err
	}

	logger.CtxInfo(ctx, "OnBarrierDamageRQ addDamage success", zap.Any("req", req), zap.Int64("totalDamage", totalDamage))
	return nil
}
