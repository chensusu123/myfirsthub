// @Author pangchenyang 2025/6/9 21:32:00
// @Desc:
package userprofile

import (
	"fmt"
	"maze_game_server/common/errors"
	"maze_game_server/io/mysql/flowrecord"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/component"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/UserProfile"
	"maze_game_server/service/userprofileservice"
	"time"

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
	res := &UserProfile.QueryUserProfileRS{}
	res.ErrInfo = errors.NO_ERROR
	userID := uint64(s.UID())
	logger := log.Clone("profile", userID, 0)
	logger.WarnWF("OnQueryUserProfile with", zap.Any("rq", req))

	defer func() {
		err = s.Response(res)
		if err != nil {
			logger.ErrorWF("OnQueryUserProfile Response failed", zap.Error(err))
		}
		logger.WarnWF("OnQueryUserProfile end ", zap.Any("req", req), zap.Any("res", res),
			zap.String("errMsg", string(res.GetErrInfo().GetErrMsg())))
	}()

	// 检查rq
	if len(req.GetUserList()) == 0 || userID <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}
	res.UserProfile, err = userprofileservice.GlobalUserProfileService.GetBatchUserProfile(logger, req.GetUserList())
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

	defer func() {
		err = s.Response(res)
		if err != nil {
			logger.ErrorWF("OnAlterUserProfile Response failed", zap.Error(err))
		}
		logger.WarnWF("OnAlterUserProfile end ", zap.Any("req", req), zap.Any("res", res),
			zap.String("errMsg", string(res.GetErrInfo().GetErrMsg())))
	}()

	if req.GetAlterProfile() == nil || req.GetAlterProfile().GetUserId() != userID || userID <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}

	curProfile, err := userprofileservice.GlobalUserProfileService.GetUserProfile(logger, userID)
	if err != nil {
		logger.ErrorWF("OnAlterUserProfile get user profile fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	var (
		chgDesc, newVal, oldVal string
	)

	// 组装修改信息(只修改需要更改的字段)
	alterProfile := alterProfile(req.GetAlterProfile(), curProfile, &chgDesc, &newVal, &oldVal)

	err = userprofileservice.GlobalUserProfileService.AlterUserProfile(logger, userID, alterProfile)
	if err != nil {
		logger.ErrorWF("OnAlterUserProfile alter user profile fail", zap.Error(err))
		res.ErrInfo = errors.DB_SAVE_ERROR.ToInfo()
		return
	}

	// 修改资料流水
	err = flowrecord.GlobalAlterProfileRecordMysql.SaveAlterProfileRecord(flowrecord.AlterProfileRecord{
		UserId:     userID,
		NewVal:     newVal,
		OldVal:     oldVal,
		ChgDesc:    chgDesc,
		CreateTime: time.Now().UnixMilli(), // ms
	})
	if err != nil {
		logger.ErrorWF("OnAlterUserProfile save alter profile record fail", zap.Error(err))
		res.ErrInfo = errors.DB_SAVE_ERROR.ToInfo()
		return
	}

	res.UserProfile = alterProfile

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
	token, err := CreateAvatarToken(logger)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return
	}
	res.Token = proto.String(token)
	return
}

// 修改需要更改的资料，支持用户资料为空时的初始化
func alterProfile(alterProfile, dbProfile *UserProfile.UserProfile,
	chgDesc *string, newVal *string, oldVal *string) *UserProfile.UserProfile {
	ret := proto.Clone(dbProfile).(*UserProfile.UserProfile)
	if alterProfile.GetNickName() != "" && alterProfile.GetNickName() != ret.GetNickName() {
		ret.NickName = alterProfile.NickName
		*chgDesc += chgDescMap[chgName]
		*newVal += getAlterVal(chgName, alterProfile.GetNickName())
		*oldVal += getAlterVal(chgName, ret.GetNickName())
	}
	if alterProfile.GetAvatar() != "" {
		ret.Avatar = alterProfile.Avatar
		*chgDesc += chgDescMap[chgAvatar]
		*newVal += getAlterVal(chgAvatar, alterProfile.GetAvatar())
		*oldVal += getAlterVal(chgAvatar, ret.GetAvatar())
	}
	if alterProfile.GetSex() != 0 {
		ret.Sex = alterProfile.Sex
		*chgDesc += chgDescMap[chgSex]
		*newVal += getAlterVal(chgSex, alterProfile.GetSex())
		*oldVal += getAlterVal(chgSex, ret.GetSex())
	}
	return ret
}

const (
	chgName = 1 + iota
	chgAvatar
	chgSex
)

var chgDescMap = map[int32]string{
	1: "修改昵称",
	2: "修改头像",
	3: "修改性别",
}

var descMap = map[int32]string{
	1: "昵称",
	2: "头像",
	3: "性别",
}

func getAlterVal(chgType int32, val interface{}) string {
	return fmt.Sprintf("%s:%s", descMap[chgType], val)
}
