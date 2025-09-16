// @Author pangchenyang 2025/6/9 21:32:00
// @Desc:
package userprofile

import (
	"maze_game_server/common/errors"
	"maze_game_server/common/function/packtopb/packequipostopb"
	"maze_game_server/lib/nano/component"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/userprofilemodel"
	"maze_game_server/module/dollassembleinfo"
	"maze_game_server/pb/common/Costume"
	"maze_game_server/pb/common/UserProfile"
	"maze_game_server/services/allianceservice"
	"maze_game_server/services/costumeservice"
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

// 用于当前用户查询其他人的详细信息
func (p *Profile) OnQueryUserDetailInfo_10719_10720(s *session.Session, req *UserProfile.QueryUserDetailInfoRQ) (err error) {
	defer fkprometheus.DebugPMT("OnQueryUserDetailInfo")()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &UserProfile.QueryUserDetailInfoRS{}
	res.ErrInfo = errors.NO_ERROR

	logger.CtxInfo(ctx, "OnQueryUserDetailInfo with", zap.Any("rq", req))

	defer func() {
		err = s.Response(res)
		if err != nil {
			logger.CtxError(ctx, "OnQueryUserDetailInfo Response failed", zap.Error(err))
		}
		logger.CtxInfo(ctx, "OnQueryUserDetailInfo end ", zap.Any("req", req), zap.Any("res", res),
			zap.String("errMsg", string(res.GetErrInfo().GetErrMsg())))
	}()

	// 检查rq
	if req.GetUserId() <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}
	// 查询用户资料
	userDetailProfile, err := userprofileservice.GlobalUserProfileService.GetUserDetailInfo(ctx, uint64(s.UID()), uint64(req.GetUserId()))
	if err != nil {
		logger.CtxError(ctx, "OnQueryUserDetailInfo get user profile fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	res.UserId = proto.Int64(int64(req.GetUserId()))
	res.NickName = proto.String(userDetailProfile.NickName)
	res.Sex = proto.Int32(int32(userDetailProfile.Sex))
	res.IconToken = proto.String(userDetailProfile.IconToken)
	res.IsBlack = proto.Bool(userDetailProfile.IsBlack)
	res.IsFriend = proto.Bool(userDetailProfile.IsFriend)
	res.ShowId = proto.Int32(int32(userDetailProfile.ShowID))

	// 查询用户联盟信息
	allianceID, err := allianceservice.GlobalAllianceService.QueryUserAlliance(ctx, uint64(req.GetUserId()))
	if err != nil {
		logger.CtxError(ctx, "OnQueryUserDetailInfo get user alliance fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	allianceInfo, err := allianceservice.GlobalAllianceService.QueryAllianceInfo(ctx, allianceID)
	if err != nil {
		logger.CtxError(ctx, "OnQueryUserDetailInfo get alliance info fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	res.AllianceInfo = allianceInfo.DataToAllianceInfoPb()

	//查询用户装扮列表
	costumeList, err := costumeservice.GlobalCostumeService.GetUserCostume(ctx, uint64(req.GetUserId()))
	if err != nil {
		logger.CtxError(ctx, "OnQueryUserDetailInfo get user costume fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	res.CostumeInfo = make([]*Costume.CostumeInfo, 0, len(costumeList))
	for k, v := range costumeList {
		res.CostumeInfo = append(res.CostumeInfo, &Costume.CostumeInfo{
			Pos:     proto.Int32(k),
			ModelId: proto.Int32(v),
		})
	}

	//装备位数据
	// 查询装配信息
	assembleInfo, err := dollassembleinfo.GetDollAssembleInfo(ctx, uint64(req.GetUserId()))
	logger.CtxInfo(ctx, "OnQueryUserDetailInfo get assemble info", zap.Int64("userId", req.GetUserId()), zap.Any("assembleInfo", assembleInfo))
	if err != nil {
		logger.CtxError(ctx, "GetDollAssembleInfo get fail", zap.Error(err))
		return
	}
	for _, equip := range assembleInfo.MazeEquips {
		cliEquip, err := packequipostopb.PackEquipPosPb(ctx, equip, -1)
		if err != nil {
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			logger.CtxError(ctx, "OnQueryUserDetailInfo PackEquipPosPb fail", zap.Error(err))
			return err
		}
		res.EquipPosList = append(res.EquipPosList, cliEquip)
	}
	return
}
