package game

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/services/barrierservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

func (g *Game) OnMazeBarrierListRQ_10457_10458(s *session.Session, req *MazeGame.MazeBarrierListRQ) (err error) {
	defer fkprometheus.InfoPMT("OnMazeBarrierListRQ")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeGame.MazeBarrierListRS{}

	logger.InfoWF("OnMazeBarrierListRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnMazeBarrierListRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userID := uint64(s.UID())

	barriers, errinfo := barrierservice.Global.GetBarrierInfos(logger, userID)
	if errinfo.GetErrCode() != errors.NO_ERROR_CODE {
		res.ErrInfo = errinfo
	} else {
		res.MazeBarrierList = barriers
	}
	return
}
