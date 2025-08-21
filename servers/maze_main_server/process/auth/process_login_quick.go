// 代理服务使用，快速重登功能，主要用于对Nano的session进行初始化
package auth

import (
	"time"

	"maze_game_server/common/errors"
	"maze_game_server/io/redis/UnionIDBindRedis"
	"maze_game_server/io/redis/useridredis"
	"maze_game_server/io/redis/usersection"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/UserLogin"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (a *Auth) OnLoginQuickRQ_10550_10551(s *session.Session, req *UserLogin.UserLoginRq) (err error) {
	defer fkprometheus.InfoPMT("OnLoginQuickRQ")()
	ctx := s.Context()
	logger := log.Clone("Auth", uint64(req.GetAuthId()), 0)
	res := &UserLogin.UserLoginRs{}

	// res.Session = req.Session
	// res.ClientTime = req.ClientTime
	// res.ServerTime = proto.Int64(time.Now().UnixMilli())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnLoginQuickRQ end", zap.Any("req", req), zap.Any("res", res))
	}()

	// 认证
	// 失败直接返回
	if req.GetAuthId() == 0 {
		res.Error = errors.COMMON_ERROR_TIPS.Wrap("id is 0")
		return nil
	}

	users, err := UnionIDBindRedis.GetUsersWithUnionID(ctx, logger, uint64(req.GetAuthId()))
	if err != nil {
		res.Error = errors.COMMON_ERROR_TIPS.Wrap("get users with unionID fail")
		return nil
	}

	userID := uint64(0)
	if len(users) == 0 {
		newUserID := useridredis.Generate(ctx, logger)
		if newUserID == 0 {
			res.Error = errors.COMMON_ERROR_TIPS.Wrap("generate userID fail")
			return nil
		}
		err = UnionIDBindRedis.AddUnionID2UserID(ctx, logger, uint64(req.GetAuthId()), newUserID)
		if err != nil {
			res.Error = errors.COMMON_ERROR_TIPS.Wrap("add unionID to userID fail")
			return nil
		}
		err = UnionIDBindRedis.AddUserID2UnionID(ctx, logger, newUserID, uint64(req.GetAuthId()))
		if err != nil {
			res.Error = errors.COMMON_ERROR_TIPS.Wrap("add userID to unionID fail")
			return nil
		}
		err = usersection.Set(ctx, newUserID, appconfig.GlobalConfig().Global.SectionID)
		if err != nil {
			logger.ErrorWF("usersection.Set fail",
				zap.Uint64("userID", userID),
				zap.Error(err))
		}
		userID = newUserID
	} else {
		userID = users[0]
	}
	res.Error = errors.NO_ERROR
	res.UserId = proto.Uint64(userID)
	// 认证成功设置用户ID, 底层会处理

	// ctx.SetTag("userID", userID)
	s.Bind(int64(userID))
	online.Bind(logger, s, userID)

	res.ServerTime = proto.Int64(time.Now().UnixMilli())

	// time.AfterFunc(time.Second*2, func() {
	// 	SendArrivePacketWithContext(ctx, logger, int64(userID), 111, &UserLogin.UserLiveRs{
	// 		ClientTime: proto.Int64(time.Now().UnixMilli()),
	// 	})
	// })
	return nil
}
