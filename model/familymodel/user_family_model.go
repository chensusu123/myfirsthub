package familymodel

import (
	"maze_game_server/io/redis/familyredis"
	"time"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
)

type UserFamilyModel struct {
	UserID uint64 `json:"user_id,omitempty"`
}

func NewUserFamilyModel(logger fklog.FKLogI, userID uint64) *UserFamilyModel {
	return &UserFamilyModel{
		UserID: userID,
	}
}

// 获取用户所在家族
func (r *UserFamilyModel) GetUserFamily(logger fklog.FKLogI) (int32, error) {
	familyID, err := familyredis.GetUserFamilyID(r.UserID)
	if err != nil {
		logger.ErrorWF("GetUserFamily GetUserFamilyID err",
			zap.Uint64("userID", r.UserID), zap.Error(err))
		return 0, err
	}
	return familyID, nil
}

// 设置用户上次退出家族时间
func (r *UserFamilyModel) SetUserLastLeaveFamilyTime(logger fklog.FKLogI) error {
	lastLeaveFamilyTime := time.Now().Unix()
	err := familyredis.SetUserLastLeaveFamilyTime(r.UserID, lastLeaveFamilyTime)
	if err != nil {
		logger.ErrorWF("SetUserLastLeaveFamilyTime SetUserLastLeaveFamilyTime err",
			zap.Uint64("userID", r.UserID), zap.Int64("lastLeaveFamilyTime", lastLeaveFamilyTime), zap.Error(err))
		return err
	}
	return nil
}

// 获取用户上次退出家族时间
func (r *UserFamilyModel) GetUserLastLeaveFamilyTime(logger fklog.FKLogI) (int64, error) {
	lastLeaveFamilyTime, err := familyredis.GetUserLastLeaveFamilyTime(r.UserID)
	if err != nil {
		logger.ErrorWF("GetUserLastLeaveFamilyTime GetUserLastLeaveFamilyTime err",
			zap.Uint64("userID", r.UserID), zap.Error(err))
		return 0, err
	}
	return lastLeaveFamilyTime, nil
}

// 设置玩家对应的家族
func (r *UserFamilyModel) SetUserFamily(logger fklog.FKLogI, familyID int32) error {
	err := familyredis.SetUserFamilyID(r.UserID, familyID)
	if err != nil {
		logger.ErrorWF("SetUserFamilyID SetUserFamilyID err",
			zap.Uint64("userID", r.UserID), zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	return nil
}

// 删除玩家对应的家族
func (r *UserFamilyModel) DelUserFamily(logger fklog.FKLogI) error {
	err := familyredis.DelUserFamilyID(r.UserID)
	if err != nil {
		logger.ErrorWF("DelUserFamilyID DelUserFamilyID err",
			zap.Uint64("userID", r.UserID), zap.Error(err))
		return err
	}
	return nil
}
