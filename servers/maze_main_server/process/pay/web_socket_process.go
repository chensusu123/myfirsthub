package pay

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"maze_game_server/common/errors"
	"maze_game_server/common/jwt"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/component"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazePay"
	"strconv"
)

type Pay struct {
	component.Base
}

func NewPay() *Pay {
	return &Pay{}
}

func (i *Pay) OnGetPayTokenRQ_10507_10508(s *session.Session, req *MazePay.MazePayTokenRQ) (err error) {
	logger := log.Clone("Pay", uint64(s.UID()), 0)
	res := &MazePay.MazePayTokenRS{}

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	uid := uint64(s.UID())

	defer fkprometheus.DebugPMT("OnGetPayTokenRQ")()
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnGetPayTokenRQ end", zap.Any("req", req), zap.Any("res", res))
	}()
	// todo 检查unique_id是否可以购买
	// todo 获取用户信息

	serverIdStr := appconfig.GlobalConfig().Global.SectionID
	serverId, err := strconv.ParseUint(serverIdStr, 10, 32)
	if err != nil {
		logger.ErrorWF("ServerId ParseUint failed", zap.Error(err), zap.Any("req", req), zap.Any("serverIdStr", serverIdStr))
	}
	payJwt, err := jwt.GeneratePayJWT(uid, req.GetUniqueId(), "wx123", uint32(serverId))
	if err != nil {
		logger.ErrorWF("OnGetPayTokenRQ GeneratePayJWT err", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	res.PayToken = proto.String(payJwt)
	return
}
