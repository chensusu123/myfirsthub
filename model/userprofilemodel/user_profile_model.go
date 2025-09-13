package userprofilemodel

import (
	"context"
	"fmt"
	"maze_game_server/io"
	"maze_game_server/io/redis/userprofileredis"
	"maze_game_server/lib/serialize"
	"maze_game_server/pb/common/UserProfile"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

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

func GetAlterVal(chgType int32, val interface{}) string {
	return fmt.Sprintf("%s:%s", descMap[chgType], val)
}

type UserProfileModel struct {
	UserID   uint64 `json:"user_id,omitempty"`
	NickName string `json:"nick_name,omitempty"`
	Sex      int32  `json:"sex,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

type UserProfileModels []*UserProfileModel

func getKey(userID uint64) string {
	return fmt.Sprintf("user:profile:%d", userID)
}

func LoadUserProfileModel(ctx context.Context, userID uint64) (um *UserProfileModel, err error) {
	logger := fklog.ContextAppLogger(ctx)
	um = &UserProfileModel{}
	if err = um.load(ctx, userID); err != nil {
		logger.CtxError(ctx, "LoadUserProfileModel err",
			zap.Uint64("userID", userID), zap.Error(err))
		return nil, err
	}
	return
}

// LoadBatchUserProfileModel 获取批量用户数据模型
func LoadBatchUserProfileModel(ctx context.Context, userIDs []uint64) (um UserProfileModels, err error) {
	logger := fklog.ContextAppLogger(ctx)
	um = make(UserProfileModels, len(userIDs))
	if len(userIDs) > 0 {
		values, err := userprofileredis.GlobalUserProfileRedis.BatchGetProfile(ctx, userIDs)
		if err != nil {
			logger.CtxError(ctx, "LoadBatchUserProfileModel err",
				zap.Uint64s("userIDs", userIDs), zap.Error(err))
			return nil, err
		}
		for i, value := range values {
			if value == nil {
				continue
			}
			var profile *UserProfileModel
			if err := serialize.Unmarshal(value, &profile); err != nil {
				logger.CtxError(ctx, "LoadBatchUserProfileModel err",
					zap.Uint64s("userIDs", userIDs), zap.Error(err))
				return nil, err
			}
			um[i] = profile
		}
	}
	return
}

func (up *UserProfileModel) load(ctx context.Context, userID uint64) (err error) {
	return io.LoadSvrData(ctx, getKey(userID), up)
}

func (up *UserProfileModel) Save(ctx context.Context, userID uint64) (err error) {
	return io.SaveSvrData(ctx, getKey(userID), up)
}

func (up *UserProfileModel) Delete(ctx context.Context, userID uint64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	err = userprofileredis.GlobalUserProfileRedis.DelProfile(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "delete DelProfile err",
			zap.Uint64("userID", userID), zap.Error(err))
		return err
	}
	return
}

func (up *UserProfileModel) ModelDataToPb() *UserProfile.UserProfile {
	return &UserProfile.UserProfile{
		UserId:   proto.Uint64(up.UserID),
		NickName: proto.String(up.NickName),
		Sex:      proto.Int32(up.Sex),
		Avatar:   proto.String(up.Avatar),
	}
}

func (up *UserProfileModel) FillModelData(profile *UserProfileModel) {
	up.UserID = profile.UserID
	up.NickName = profile.NickName
	up.Sex = profile.Sex
	up.Avatar = profile.Avatar
}

// 修改需要更改的资料，支持用户资料为空时的初始化
func AlterProfile(alterProfile *UserProfile.UserProfile, dbProfile *UserProfileModel,
	chgDesc *string, newVal *string, oldVal *string) *UserProfileModel {
	ret := copyProfile(dbProfile)
	if alterProfile.GetNickName() != "" && alterProfile.GetNickName() != ret.NickName {
		ret.NickName = alterProfile.GetNickName()
		*chgDesc += chgDescMap[chgName]
		*newVal += GetAlterVal(chgName, alterProfile.GetNickName())
		*oldVal += GetAlterVal(chgName, ret.NickName)
	}
	if alterProfile.GetAvatar() != "" {
		ret.Avatar = alterProfile.GetAvatar()
		*chgDesc += chgDescMap[chgAvatar]
		*newVal += GetAlterVal(chgAvatar, alterProfile.GetAvatar())
		*oldVal += GetAlterVal(chgAvatar, ret.Avatar)
	}
	if alterProfile.GetSex() != 0 {
		ret.Sex = alterProfile.GetSex()
		*chgDesc += chgDescMap[chgSex]
		*newVal += GetAlterVal(chgSex, alterProfile.GetSex())
		*oldVal += GetAlterVal(chgSex, ret.Sex)
	}
	return ret
}

func copyProfile(profile *UserProfileModel) *UserProfileModel {
	ret := &UserProfileModel{}
	if profile == nil {
		return ret
	}

	ret.UserID = profile.UserID
	ret.Avatar = profile.Avatar
	ret.Sex = profile.Sex
	ret.NickName = profile.NickName
	return ret
}
