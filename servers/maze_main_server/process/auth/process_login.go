package auth

import (
	"context"
	"time"

	"maze_game_server/common/errors"
	"maze_game_server/io/redis/UnionIDBindRedis"
	"maze_game_server/io/redis/useridredis"
	"maze_game_server/io/redis/usersection"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/UserLogin"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func checkResResult(res *UserLogin.UserLoginRs) bool {
	if res.Error == nil {
		return false
	}
	if res.GetError().GetErrCode() != 0x80000000 {
		return false
	}
	return true
}

func (a *Auth) OnLoginRQ_10492_10493(s *session.Session, req *UserLogin.UserLoginRq) (err error) {
	defer fkprometheus.InfoPMT("OnLoginRQ")()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &UserLogin.UserLoginRs{}

	res.Session = req.Session
	res.ClientTime = req.ClientTime
	res.ServerTime = proto.Int64(time.Now().UnixMilli())

	defer func() {
		err = s.Response(res)
		if !checkResResult(res) {
			// logger.CtxError(ctx, "OnLoginRQ end", zap.Any("req", req), zap.Any("res", res), zap.String("ClientAddr", s.String("ClientAddr")))
		} else {
			logger.CtxInfo(ctx, "OnLoginRQ end", zap.Any("req", req), zap.Any("res", res), zap.String("ClientAddr", s.String("ClientAddr")))
		}
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
		err = usersection.Set(context.TODO(), newUserID, appconfig.GlobalConfig().Global.SectionID)
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

	time.AfterFunc(time.Second*2, func() {
		SendArrivePacketWithContext(ctx, logger, int64(userID), 111, &UserLogin.UserLiveRs{
			ClientTime: proto.Int64(time.Now().UnixMilli()),
		})
	})
	return nil
}

func (a *Auth) OnLiveRQ_10494_10495(s *session.Session, req *UserLogin.UserLiveRq) (err error) {
	defer fkprometheus.InfoPMT("OnLiveRQ")()

	logger := log.Clone("Auth", uint64(s.UID()), 0)
	res := &UserLogin.UserLiveRs{}

	res.Session = req.Session
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnLiveRQ end", zap.Any("req", req), zap.Any("res", res))
	}()

	res.Error = errors.NO_ERROR
	res.ServerTime = proto.Int64(time.Now().UnixMilli())
	return nil
}

func (a *Auth) OnConfigDataMd5Rq_10500_10501(s *session.Session, req *UserLogin.ConfigDataMd5Rq) (err error) {
	defer fkprometheus.InfoPMT("OnConfigDataMd5Rq")()

	logger := log.Clone("Auth", uint64(s.UID()), 0)
	res := &UserLogin.ConfigDataMd5Rs{}

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnConfigDataMd5Rq end", zap.Any("req", req), zap.Any("res", res))
	}()

	res.Error = errors.NO_ERROR
	cfg := config_manager.ShowSheet()
	var serverGitVersion string
	for _, v := range cfg {
		res.Items = append(res.Items, &UserLogin.ConfigDataItem{
			FileName:  proto.String(v.XlsxFile),
			Md5:       proto.String(v.Md5),
			SheetName: proto.String(v.XlsxSheet),
		})
		serverGitVersion = v.GitVersion
	}
	res.ConfigVersion = req.ConfigVersion
	res.ServerConfigVersion = proto.String(serverGitVersion)
	res.Result = proto.Bool(serverGitVersion == req.GetConfigVersion())
	return nil
}
