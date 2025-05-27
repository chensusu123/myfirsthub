package game

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/pb/common/MazeGame"

	"maze_game_server/lib/nano/session"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

// OnMazeReportBattleEventRQ 关卡事件上报
func (g *Game) OnMazeReportBattleEventRQ_10496_10497(s *session.Session, req *MazeGame.ReportBattleEventRQ) (err error) {
	defer fkprometheus.InfoPMT("OnMazeReportBattleEventRQ")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeGame.ReportBattleEventRS{}

	logger.InfoWF("OnMazeReportBattleEventRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnMazeReportBattleEventRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	return
}
