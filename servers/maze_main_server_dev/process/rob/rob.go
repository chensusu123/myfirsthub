// @Author: ZhaoXiming 2025/3/31 21:07
// @Desc:

package rob

import (
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/protodef/MazeRobGuaJi"
	"go.uber.org/zap"
)

func OnMazeRobGuaJiRQ(ctx fknet.TCPContext, uid uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	req := rqMsg.(*MazeRobGuaJi.MazeRobGuaJiRQ)
	res := rsMsg.(*MazeRobGuaJi.MazeRobGuaJiRS)
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header

	defer fkprometheus.DebugPMT("OnMazeRobGuaJiRQ")()
	defer func() {
		ctx.InfoWF("OnMazeRobGuaJiRQ end",
			zap.Any("req", req),
			zap.Any("res", res),
		)
	}()

	return
}
