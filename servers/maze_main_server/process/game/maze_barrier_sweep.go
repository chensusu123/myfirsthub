package game

import (
	"time"

	"maze_game_server/common/errors"
	"maze_game_server/lib/codec"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeEnergy"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/services/barrierservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// OnStartMazeSweepRQ start sweep
func (g *Game) OnStartMazeSweepRQ_10471_10472(s *session.Session, req *MazeGame.StartMazeSweepRQ) (err error) {
	defer fkprometheus.InfoPMT("OnStartMazeSweepRQ")()

	logger := log.Clone("Sweep", uint64(s.UID()), 0)
	res := &MazeGame.StartMazeSweepRS{}
	energyID := &MazeEnergy.EnergyChangeID{} // defer时多补一个体力ID包

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userID := uint64(s.UID())

	logger.InfoWF("OnStartMazeSweepRQ start", zap.Any("req", req))
	ctx := s.Context()
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnStartMazeSweepRQ end", zap.Any("res", res), zap.Any("errMsg", string(res.GetErrInfo().GetErrMsg())))
		err = s.ResponseMID(ctx, codec.ToMessageID(uint32(time.Now().Unix()), 0, 10610), energyID)
		logger.InfoWF("OnStartMazeSweepRQ end send EnergyChangeID", zap.Any("energyID", energyID))
	}()

	barrierId := req.GetBarrierId()
	res.BarrierId = req.BarrierId
	if userID <= 0 {
		logger.ErrorWF("OnStartMazeSweepRQ userId invalid", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效的用户ID")
		return
	}
	// check barrier
	if req.GetBarrierId() <= 0 {
		logger.ErrorWF("OnStartMazeSweepRQ req barrier invalid", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("未设置关卡id")
		return
	}

	// 扫荡关卡
	energyInfo, remainVal, gameID, awardItem, rareItem, errinfo := barrierservice.Global.SweepBarrier(logger, req.GetHeader(), userID, barrierId)
	if errinfo != nil {
		res.ErrInfo = errinfo
	} else {
		res.RemainEnergy = proto.Int32(remainVal)
		res.Awards = awardItem
		res.RareAward = rareItem
		res.GameId = proto.Uint64(gameID)
		res.ErrInfo = errors.NO_ERROR
	}
	energyID.EnergyInfo = energyInfo
	return nil
}
