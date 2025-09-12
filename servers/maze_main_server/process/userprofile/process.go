// @Author pangchenyang 2025/6/9 21:32:00
// @Desc:
package userprofile

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/component"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/userprofilemodel"
	"maze_game_server/pb/common/UserProfile"
	"maze_game_server/services/userprofileservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
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
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &UserProfile.QueryUserProfileRS{}
	res.ErrInfo = errors.NO_ERROR
	userID := uint64(s.UID())
	logger.CtxInfo(ctx, "OnQueryUserProfile with", zap.Any("rq", req))

	defer func() {
		err = s.Response(res)
		if err != nil {
			logger.CtxError(ctx, "OnQueryUserProfile Response failed", zap.Error(err))
		}
		logger.CtxInfo(ctx, "OnQueryUserProfile end ", zap.Any("req", req), zap.Any("res", res),
			zap.String("errMsg", string(res.GetErrInfo().GetErrMsg())))
	}()

	// 检查rq
	if len(req.GetUserList()) == 0 || userID <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}
	userProfiles, err := userprofileservice.GlobalUserProfileService.GetBatchUserProfile(ctx, req.GetUserList())
	if err != nil {
		logger.CtxError(ctx, "OnQueryUserProfile GetBatchUserProfile fail",
			zap.Any("UserList", req.GetUserList()),
			zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return
	}

	for _, userProfile := range userProfiles {
		res.UserProfile = append(res.UserProfile, userProfile.ModelDataToPb())
	}

	return
}

// OnAlterUserProfile 修改用户资料
func (p *Profile) OnAlterUserProfile_10483_10484(s *session.Session, req *UserProfile.AlterUserProfileRQ) (err error) {
	defer fkprometheus.DebugPMT("OnAlterUserProfile")()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &UserProfile.AlterUserProfileRS{}
	res.ErrInfo = errors.NO_ERROR
	userID := uint64(s.UID())

	logger.CtxInfo(ctx, "OnAlterUserProfile with", zap.Any("rq", req))

	defer func() {
		err = s.Response(res)
		if err != nil {
			logger.CtxError(ctx, "OnAlterUserProfile Response failed", zap.Error(err))
		}
		logger.CtxInfo(ctx, "OnAlterUserProfile end ", zap.Any("req", req), zap.Any("res", res),
			zap.String("errMsg", string(res.GetErrInfo().GetErrMsg())))
	}()

	if req.GetAlterProfile() == nil || req.GetAlterProfile().GetUserId() != userID || userID <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}

	curProfile, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "OnAlterUserProfile get user profile fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	var (
		chgDesc, newVal, oldVal string
	)

	// 组装修改信息(只修改需要更改的字段)
	alterProfile := userprofilemodel.AlterProfile(req.GetAlterProfile(), curProfile, &chgDesc, &newVal, &oldVal)

	err = userprofileservice.GlobalUserProfileService.AlterUserProfile(ctx, userID, alterProfile)
	if err != nil {
		logger.ErrorWF("OnAlterUserProfile alter user profile fail", zap.Error(err))
		res.ErrInfo = errors.DB_SAVE_ERROR.ToInfo()
		return
	}

	// // 修改资料流水
	// err = flowrecord.GlobalAlterProfileRecordMysql.SaveAlterProfileRecord(flowrecord.AlterProfileRecord{
	// 	UserId:     userID,
	// 	NewVal:     newVal,
	// 	OldVal:     oldVal,
	// 	ChgDesc:    chgDesc,
	// 	CreateTime: time.Now().UnixMilli(), // ms
	// })
	// if err != nil {
	// 	logger.ErrorWF("OnAlterUserProfile save alter profile record fail", zap.Error(err))
	// 	res.ErrInfo = errors.DB_SAVE_ERROR.ToInfo()
	// 	return
	// }

	res.UserProfile = alterProfile.ModelDataToPb()

	return nil
}

// OnQueryAvatarToken 查询用户操作头像jwt token
func (p *Profile) OnQueryAvatarToken_10511_10512(s *session.Session, req *UserProfile.QueryAvatarTokenRQ) (err error) {
	defer fkprometheus.DebugPMT("OnQueryAvatarToken")()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &UserProfile.QueryAvatarTokenRS{}
	res.ErrInfo = errors.NO_ERROR
	userID := uint64(s.UID())

	logger.CtxInfo(ctx, "OnQueryAvatarToken with", zap.Any("rq", req))

	defer func() {
		err = s.Response(res)
		if err != nil {
			logger.CtxError(ctx, "OnQueryAvatarToken Response failed", zap.Error(err))
		}
		logger.CtxInfo(ctx, "OnQueryAvatarToken end ", zap.Any("req", req), zap.Any("res", res),
			zap.String("errMsg", string(res.GetErrInfo().GetErrMsg())))
	}()

	// 检查rq
	if userID <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}
	// 生成token
	token, err := CreateAvatarToken(ctx)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return
	}
	res.Token = proto.String(token)
	return
}
