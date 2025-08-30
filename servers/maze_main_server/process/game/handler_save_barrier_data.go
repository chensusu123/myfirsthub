package game

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/barriersavedatamodel"
	"maze_game_server/module/mazecommonvalue"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/services/barriersavedataservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

// 关卡存档
func (g *Game) OnSaveBarrierDataRQ_10624_10625(s *session.Session, req *MazeGame.SaveBarrierDataRQ) (err error) {
	defer fkprometheus.InfoPMT("OnSaveBarrierDataRQ")()

	res := &MazeGame.SaveBarrierDataRS{}
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "OnSaveBarrierDataRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnSaveBarrierDataRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	userId := uint64(s.UID())
	rescueItems := make([]*barriersavedatamodel.RescueItemInfo, 0, len(req.RescueItems))
	for _, i := range req.RescueItems {
		rescueItems = append(rescueItems, &barriersavedatamodel.RescueItemInfo{
			MapConfigId: i.GetMapConfigId(),
		})
	}
	err = barriersavedataservice.GlobalBarrierSaveDataService.SaveBarrierData(ctx, userId, req.GetBarrierId(),
		req.GetStageId(), req.GetRescueValue(), req.GetBossPower(), req.GetBossProgress(), rescueItems)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return nil
	}

	mazecommonvalue.SendPassValueIdPack(ctx, userId, req.GetBarrierId(), req.GetStageId())

	return
}
