package auth

import (
	"time"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/protodef/UserLogin"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/errors"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/UnionIDBindRedis"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/useridredis"

	"go.uber.org/zap"
)

func OnLoginRQ(ctx fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	fkprometheus.InfoPMT("OnLoginRQ")()

	req := rqMsg.(*UserLogin.UserLoginRq)
	res := rsMsg.(*UserLogin.UserLoginRs)

	logger := ctx
	res.Session = req.Session
	defer func() {
		logger.InfoWF("OnLoginRQ end", zap.Any("req", req), zap.Any("res", res))
	}()

	// 认证
	// 失败直接返回
	if req.GetAuthId() == 0 {
		res.Error = errors.COMMON_ERROR_TIPS.Wrap("id is 0")
		return nil
	}

	users, err := UnionIDBindRedis.GetUsersWithUnionID(logger, uint64(req.GetAuthId()))
	if err != nil {
		res.Error = errors.COMMON_ERROR_TIPS.Wrap("get users with unionID fail")
		return nil
	}

	userID := uint64(0)
	if len(users) == 0 {
		newUserID := useridredis.Generate(logger)
		if newUserID == 0 {
			res.Error = errors.COMMON_ERROR_TIPS.Wrap("generate userID fail")
			return nil
		}
		err = UnionIDBindRedis.AddUnionID2UserID(logger, uint64(req.GetAuthId()), newUserID)
		if err != nil {
			res.Error = errors.COMMON_ERROR_TIPS.Wrap("add unionID to userID fail")
			return nil
		}
		err = UnionIDBindRedis.AddUserID2UnionID(logger, newUserID, uint64(req.GetAuthId()))
		if err != nil {
			res.Error = errors.COMMON_ERROR_TIPS.Wrap("add userID to unionID fail")
			return nil
		}
		userID = newUserID
	} else {
		userID = users[0]
	}
	res.Error = errors.NO_ERROR
	res.UserId = proto.Uint64(userID)
	// 认证成功设置用户ID, 底层会处理

	ctx.SetTag("userID", userID)
	res.ServerTime = proto.Int64(time.Now().UnixMilli())
	return nil
}

func OnLiveRQ(ctx fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	fkprometheus.InfoPMT("OnLiveRQ")()

	req := rqMsg.(*UserLogin.UserLiveRq)
	res := rsMsg.(*UserLogin.UserLiveRs)

	logger := ctx
	res.Session = req.Session
	defer func() {
		logger.InfoWF("OnLiveRQ end", zap.Any("req", req), zap.Any("res", res))
	}()

	res.Error = errors.NO_ERROR
	res.ServerTime = proto.Int64(time.Now().UnixMilli())
	return nil
}
