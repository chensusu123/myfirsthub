package game

import (
	"github.com/gogo/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeGame"
)

func (g *Game) OnReportScoreAwardItemRQ_10618_10619(s *session.Session, req *MazeGame.ReportScoreAwardItemRQ) (err error) {
	defer fkprometheus.InfoPMT("OnReportScoreAwardItemRQ")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeGame.ReportScoreAwardItemRS{}

	logger.InfoWF("OnReportScoreAwardItemRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnReportScoreAwardItemRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	// 根据类型找显示的道具
	if req.GetRewardType() == MazeGame.ReportScoreAwardItemType_REWARD_GOLD {
		res.ItemId = proto.Int32(constdef.GoldPileItemCfgId)
	} else if req.GetRewardType() == MazeGame.ReportScoreAwardItemType_REWARD_Strengthen_Stone {
		res.ItemId = proto.Int32(constdef.StrengthenStonePileItemCfgId)
	}
	res.UserLevel = req.UserLevel
	res.Count = req.Count

	return
}
