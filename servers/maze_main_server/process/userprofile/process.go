// @Author pangchenyang 2025/6/9 21:32:00
// @Desc: 
package userprofile

import (
	"maze_game_server/lib/nano/component"
	"maze_game_server/module/profilemodule"
	"maze_game_server/lib/nano/session"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"maze_game_server/pb/common/UserProfile"
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"go.uber.org/zap"
	"time"
	"google.golang.org/protobuf/proto"
)

type Profile struct {
	component.Base
}

func NewUserProfile() *Profile {
	return &Profile{}
}

// OnQueryUserProfile 查询用户资料
func (p *Profile) OnQueryUserProfile_10481_10482(s *session.Session, req *UserProfile.QueryUserProfileRQ) (err error) {
	defer fkprometheus.DebugPMT("OnQueryUserProfile")()
	res := &UserProfile.QueryUserProfileRS{}
	res.ErrInfo = errors.NO_ERROR
	userID := uint64(s.UID())
	logger := log.Clone("profile", userID, 0)
	logger.WarnWF("OnQueryUserProfile with", zap.Any("rq", req))

	addStartTime := time.Now()
	defer func() {
		err = s.Response(res)
		if err != nil {
			logger.ErrorWF("OnQueryUserProfile Response failed", zap.Error(err))
		}
		costTime := time.Since(addStartTime).Seconds()
		logger.WarnWF("OnQueryUserProfile end ", zap.Any("req", req), zap.Any("res", res),
			zap.String("errMsg", string(res.GetErrInfo().GetErrMsg())),
			zap.Float64("costTime", costTime))
	}()

	// 检查rq
	if len(req.GetUserList()) == 0 || userID <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}
	res.UserProfile, err = profilemodule.MGetUserProfile(logger, userID, req.GetUserList())
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return
	}
	return
}

// OnAlterUserProfile 修改用户资料
func (p *Profile) OnAlterUserProfile_10483_10484(s *session.Session, req *UserProfile.AlterUserProfileRQ) (err error) {
	defer fkprometheus.DebugPMT("OnAlterUserProfile")()
	res := &UserProfile.AlterUserProfileRS{}
	res.ErrInfo = errors.NO_ERROR
	userID := uint64(s.UID())
	logger := log.Clone("profile", userID, 0)
	logger.WarnWF("OnAlterUserProfile with", zap.Any("rq", req))

	addStartTime := time.Now()
	defer func() {
		err = s.Response(res)
		if err != nil {
			logger.ErrorWF("OnAlterUserProfile Response failed", zap.Error(err))
		}
		costTime := time.Since(addStartTime).Seconds()
		logger.WarnWF("OnAlterUserProfile end ", zap.Any("req", req), zap.Any("res", res), zap.String("errMsg", string(res.GetErrInfo().GetErrMsg())), zap.Float64("costTime", costTime))
	}()

	if req.GetAlterProfile() == nil || req.GetAlterProfile().GetUserId() != userID || userID <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}

	res.UserProfile, err = profilemodule.UpdateUserProfile(logger, req.GetAlterProfile())
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return
	}

	return nil
}

// OnQueryAvatarToken 查询用户操作头像jwt token
func (p *Profile) OnQueryAvatarToken_10511_10512(s *session.Session, req *UserProfile.QueryAvatarTokenRQ) (err error) {
	defer fkprometheus.DebugPMT("OnQueryAvatarToken")()
	res := &UserProfile.QueryAvatarTokenRS{}
	res.ErrInfo = errors.NO_ERROR
	userID := uint64(s.UID())
	logger := log.Clone("profile", userID, 0)
	logger.WarnWF("OnQueryAvatarToken with", zap.Any("rq", req))

	defer func() {
		err = s.Response(res)
		if err != nil {
			logger.ErrorWF("OnQueryAvatarToken Response failed", zap.Error(err))
		}
		logger.WarnWF("OnQueryAvatarToken end ", zap.Any("req", req), zap.Any("res", res),
			zap.String("errMsg", string(res.GetErrInfo().GetErrMsg())))
	}()

	// 检查rq
	if userID <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}
	// 生成token
	token, err := profilemodule.CreateAvatarToken(logger)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return
	}
	res.Token = proto.String(token)
	return
}
