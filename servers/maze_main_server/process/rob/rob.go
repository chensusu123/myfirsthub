// @Author: ZhaoXiming 2025/3/31 21:07
// @Desc:

package rob

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeRobGuaJi"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

func (*Rob) OnMazeRobGuaJiRQ_10490_10491(s *session.Session, req *MazeRobGuaJi.MazeRobGuaJiRQ) (err error) {

	logger := log.Clone("Rob", uint64(s.UID()), 0)
	res := &MazeRobGuaJi.MazeRobGuaJiRS{}
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header

	defer fkprometheus.DebugPMT("OnMazeRobGuaJiRQ")()
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(s.Context(), "OnMazeRobGuaJiRQ end",
			zap.Any("req", req),
			zap.Any("res", res),
		)
	}()

	return
}
