package game

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazecommonvalue"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/services/barriersavedataservice"
)

// 关卡存档
func (g *Game) OnSaveBarrierDataRQ_10624_10625(s *session.Session, req *MazeGame.SaveBarrierDataRQ) (err error) {
	defer fkprometheus.InfoPMT("OnSaveBarrierDataRQ")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeGame.SaveBarrierDataRS{}

	logger.InfoWF("OnSaveBarrierDataRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnSaveBarrierDataRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	userId := uint64(s.UID())

	err = barriersavedataservice.GlobalBarrierSaveDataService.SaveBarrierData(logger, userId, req.GetBarrierId(), req.GetStageId(), req.GetRescueValue(), req.GetBossPower())
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return nil
	}

	mazecommonvalue.SendPassValueIdPack(logger, userId, req.GetBarrierId(), req.GetStageId())

	return
}
