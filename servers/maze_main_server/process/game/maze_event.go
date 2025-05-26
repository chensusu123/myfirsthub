package game

import (
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/pb/common/MazeGame"
)

// OnMazeReportBattleEventRQ 关卡事件上报
func OnMazeReportBattleEventRQ(logger fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	defer fkprometheus.InfoPMT("OnMazeReportBattleEventRQ")()

	req := rqMsg.(*MazeGame.ReportBattleEventRQ)
	res := rsMsg.(*MazeGame.ReportBattleEventRS)

	logger.InfoWF("OnMazeReportBattleEventRQ start", zap.Any("req", req))
	defer func() {
		logger.InfoWF("OnMazeReportBattleEventRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	return
}
