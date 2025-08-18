package game

import (
	"github.com/gogo/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/services/barrierstagecounterservice"
)

func (g *Game) OnBarrierKillMonsterRQ_10620_10621(s *session.Session, req *MazeGame.BarrierKillMonsterRQ) (err error) {
	defer fkprometheus.InfoPMT("OnBarrierKillMonsterRQ")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeGame.BarrierKillMonsterRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.StageId = req.StageId
	res.BarrierId = req.BarrierId

	logger.InfoWF("OnBarrierKillMonsterRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnBarrierKillMonsterRQ end", zap.Any("res", res))
	}()

	userId := uint64(s.UID())

	if req.GetBarrierId() <= 0 || req.GetStageId() < 0 {
		logger.ErrorWF("OnBarrierKillMonsterRQ req barrier or areaId invalid", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡id未设置")
		return
	}
	barrierCfg := GMazeBarriesV8Cfg.Get(req.GetBarrierId())
	if barrierCfg == nil {
		logger.ErrorWF("OnBarrierKillMonsterRQ get barrier cfg fail", zap.Any("barrier", req.GetStageId()))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("找不到该关卡配置")
		return
	}

	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("OnBarrierKillMonsterRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	if userInfo.Barrier != req.GetBarrierId() {
		logger.ErrorWF("OnBarrierKillMonsterRQ barrier err", zap.Any("req", req), zap.Any("barrier", userInfo.Barrier))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡id错误")
		return
	}

	totalNum, _, err := barrierstagecounterservice.GlobalBarrierStageCounterService.AddKillMonsterNum(logger, userId, req.GetBarrierId(), req.GetStageId(), req.GetMonsterId(), 1, req.GetMonsterGuid())
	if err != nil {
		logger.ErrorWF("OnBarrierKillMonsterRQ AddKillMonsterNum fail", zap.Any("req", req), zap.Error(err))
		return err
	}

	res.TotalNum = proto.Int32(totalNum)
	logger.InfoWF("OnBarrierKillMonsterRQ AddKillMonsterNum success", zap.Any("res", res))
	return nil
}
