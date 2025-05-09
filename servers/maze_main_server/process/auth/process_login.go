package auth

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/usecase/loginauth"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/protodef/SysPackDef"

	"go.uber.org/zap"
)

func OnLoginRQ(ctx fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	fkprometheus.InfoPMT("OnLoginRQ")()

	req := rqMsg.(*SysPackDef.UserLoginRq)
	res := rsMsg.(*SysPackDef.UserLoginRs)

	logger := ctx

	defer func() {
		logger.InfoWF("OnLoginRQ end", zap.Any("req", req), zap.Any("res", res))
	}()

	// 认证
	// 失败直接返回
	if req.GetUserID() == 0 {
		res.Error = errors.COMMON_ERROR_TIPS.Wrap("id is 0")
		return nil
	}

	if req.GetUserAuthToken() == 0 {
		res.Error = errors.COMMON_ERROR_TIPS.Wrap("auth token is 0")
		return nil
	}

	loginOK, err := loginauth.AuthJudge(logger, req.GetUserID(), req.GetUserAuthToken())
	if err != nil {
		logger.ErrorWF("OnLoginRQ AuthJudge fail", zap.Error(err))
		res.Error = errors.COMMON_ERROR_TIPS.Wrap("AuthJudge fail")
		return
	}

	if !loginOK {
		res.Error = errors.COMMON_ERROR_TIPS.Wrap("auth token is 0")
		return
	}

	res.Error = errors.NO_ERROR

	// 认证成功设置用户ID, 底层会处理
	userId := req.GetUserID()
	ctx.SetTag("userID", userId)
	return nil
}
