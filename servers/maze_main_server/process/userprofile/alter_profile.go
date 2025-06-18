// @Author pangchenyang 2025/6/9 21:34:00
// @Desc: 
package userprofile

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"maze_game_server/common/errors"
	"go.uber.org/zap"
	"time"
	"context"
	"google.golang.org/protobuf/proto"
	"maze_game_server/pb/common/UserProfile"
	"maze_game_server/io/redis/userprofileredis"
	"maze_game_server/io/redis/userprofilelock"
	"maze_game_server/io/mysql/flowrecord"
)

// OnAlterUserProfile 修改用户资料
func (p *Profile) OnAlterUserProfile_10483_10484(ctx fklog.FKLogI, shardingID int64, rqMsg proto.Message, rsMsg proto.Message, opData string) (err error) {
	defer fkprometheus.DebugPMT("OnAlterUserProfile")()
	userCtx := fkserver.NewUserContext(context.TODO(), uint64(shardingID), ctx)
	req := rqMsg.(*UserProfile.AlterUserProfileRQ)
	res := rsMsg.(*UserProfile.AlterUserProfileRS)
	res.ErrInfo = errors.NO_ERROR
	userCtx.WarnWF("OnAlterUserProfile with", zap.Any("rq", req))

	addStartTime := time.Now()
	defer func() {
		costTime := time.Since(addStartTime).Seconds()
		userCtx.WarnWF("OnAlterUserProfile end ", zap.Any("req", req), zap.Any("res", res), zap.String("errMsg", string(res.GetErrInfo().GetErrMsg())), zap.Float64("costTime", costTime))
	}()

	if req.GetAlterProfile() == nil || req.GetAlterProfile().GetUserId() != uint64(shardingID) || shardingID <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}

	userID := uint64(shardingID)

	// lock
	locked, err := userprofilelock.Lock(userID)
	if err != nil {
		userCtx.ErrorWF("OnAlterUserProfile get lock failed", zap.Error(err))
		return err
	}
	if !locked {
		userCtx.WarnWF("OnAlterUserProfile lock already exists")
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("服务器忙")
		return nil
	}
	defer userprofilelock.Unlock(userID)

	dbProfile, err := userprofileredis.GetProfile(userID)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("获取用户资料失败")
		return
	}

	var (
		chgDesc, newVal, oldVal string
	)

	// 修改用户资料
	FixedProfile := alterProfile(req.GetAlterProfile(), dbProfile, &chgDesc, &newVal, &oldVal)

	// 修改资料流水
	flowrecord.SaveAlterProfileRecord(userCtx, flowrecord.AlterProfileRecord{
		UserId:     userID,
		NewVal:     newVal,
		OldVal:     oldVal,
		ChgDesc:    chgDesc,
		CreateTime: time.Now().UnixMilli(), // ms
	})
	res.UserProfile = FixedProfile

	return nil
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

// notifyProfileChange 通知用户资料变更
func notifyProfileChange(ctx fklog.FKLogI, userId uint64) {
	// TODO: 实现通知逻辑，可以通过消息队列或其他方式通知相关服务
}
