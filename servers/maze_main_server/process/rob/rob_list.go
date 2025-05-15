// @Author: ZhaoXiming 2025/3/24 15:04
// @Desc:

package rob

import (
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/common/errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeRobGuaJi"
	"go.uber.org/zap"
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

func OnMazeRobGuaJiListRQ(ctx fknet.TCPContext, uid uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	req := rqMsg.(*MazeRobGuaJi.MazeRobGuaJiListRQ)
	res := rsMsg.(*MazeRobGuaJi.MazeRobGuaJiListRS)
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header

	defer fkprometheus.DebugPMT("OnMazeRobGuaJiListRQ")()
	defer func() {
		ctx.InfoWF("OnMazeRobGuaJiListRQ end",
			zap.Any("req", req),
			zap.Any("res", res),
		)
	}()

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
