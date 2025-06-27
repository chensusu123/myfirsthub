package userprofilemodel

import (
	"maze_game_server/io/redis/userprofileredis"
	"maze_game_server/lib/serialize"
	"maze_game_server/pb/common/UserProfile"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type UserProfileModel struct {
	UserID   uint64 `json:"user_id,omitempty"`
	NickName string `json:"nick_name,omitempty"`
	Sex      int32  `json:"sex,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

type UserProfileModels []*UserProfileModel

func LoadUserProfileModel(logger fklog.FKLogI, userID uint64) (um *UserProfileModel, err error) {
	um = &UserProfileModel{}
	if err = um.load(logger, userID); err != nil {
		logger.ErrorWF("LoadUserProfileModel err",
			zap.Uint64("userID", userID), zap.Error(err))
		return nil, err
	}
	return
}

// LoadBatchUserProfileModel 获取批量用户数据模型
func LoadBatchUserProfileModel(logger fklog.FKLogI, userIDs []uint64) (um UserProfileModels, err error) {
	um = make(UserProfileModels, len(userIDs))
	if len(userIDs) > 0 {
		values, err := userprofileredis.GlobalUserProfileRedis.BatchGetProfile(userIDs)
		if err != nil {
			logger.ErrorWF("LoadBatchUserProfileModel err",
				zap.Uint64s("userIDs", userIDs), zap.Error(err))
			return nil, err
		}
		for i, value := range values {
			if value == nil {
				continue
			}
			var profile *UserProfileModel
			if err := serialize.Unmarshal(value, &profile); err != nil {
				logger.ErrorWF("LoadBatchUserProfileModel err",
					zap.Uint64s("userIDs", userIDs), zap.Error(err))
				return nil, err
			}
			um[i] = profile
		}
	}
	return
}

func (up *UserProfileModel) load(logger fklog.FKLogI, userID uint64) (err error) {
	value, err := userprofileredis.GlobalUserProfileRedis.GetProfile(userID)
	if err != nil {
		logger.ErrorWF("LoadUserProfileModel err",
			zap.Uint64("userID", userID), zap.Error(err))
		return err
	}
	err = serialize.Unmarshal(value, up)
	if err != nil {
		logger.ErrorWF("LoadUserProfileModel err",
			zap.Uint64("userID", userID), zap.Error(err))
		return err
	}
	return
}

func (up *UserProfileModel) Save(logger fklog.FKLogI, userID uint64) (err error) {
	value, err := serialize.Marshal(up)
	if err != nil {
		logger.ErrorWF("save Marshal err",
			zap.Uint64("userID", userID), zap.Error(err))
		return err
	}
	err = userprofileredis.GlobalUserProfileRedis.SetProfile(userID, value)
	if err != nil {
		logger.ErrorWF("save SetProfile err",
			zap.Uint64("userID", userID), zap.Error(err))
		return err
	}
	return
}

func (up *UserProfileModel) Delete(logger fklog.FKLogI, userID uint64) (err error) {
	err = userprofileredis.GlobalUserProfileRedis.DelProfile(userID)
	if err != nil {
		logger.ErrorWF("delete DelProfile err",
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

func (up *UserProfileModel) FillModelData(profile *UserProfile.UserProfile) {
	up.UserID = profile.GetUserId()
	up.NickName = profile.GetNickName()
	up.Sex = profile.GetSex()
	up.Avatar = profile.GetAvatar()
}
