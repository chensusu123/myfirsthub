// @Author: ZhaoXiming 2025/3/24 15:04
// @Desc:

package rob

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeRobGuaJi"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

/*
9003200130205896
9003200130205890

46200001
*/

var (
	robItem = []*MazeCommon.MazeItem{
		{
			ItemId: proto.Int32(46200001),
			Count:  proto.Int64(10),
		},
	}
)

func (*Rob) OnMazeRobGuaJiListRQ_10488_10489(s *session.Session, req *MazeRobGuaJi.MazeRobGuaJiListRQ) (err error) {

	logger := log.Clone("Rob", uint64(s.UID()), 0)
	res := &MazeRobGuaJi.MazeRobGuaJiListRS{}
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header

	userId := uint64(s.UID())
	_ = userId

	defer fkprometheus.DebugPMT("OnMazeRobGuaJiListRQ")()
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(s.Context(), "OnMazeRobGuaJiListRQ end",
			zap.Any("req", req),
			zap.Any("res", res),
		)
	}()

	// TODO ID不对，接口如果继续用，则需要改
	startID := uint64(9003200130205890)
	count := uint64(5)
	robUid := uint64(0)

	res.RobList = make([]*MazeRobGuaJi.RobGuaJiMsg, 0, count)

	for i := uint64(0); i < count; i++ {
		robUid = startID + i
		res.RobList = append(res.RobList, &MazeRobGuaJi.RobGuaJiMsg{
			Uid:     proto.Uint64(robUid),
			RobItem: robItem,
		})
	}

	return
}
