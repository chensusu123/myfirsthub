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
	"maze_game_server/io/mysql/userprofilemysql"
	"maze_game_server/io/redis/userprofileredis"
	"maze_game_server/io/redis/userprofilelock"
	"gorm.io/gorm"
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
	dbProfile, err := userprofilemysql.GetByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("用户不存在")
		} else {
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("获取用户资料失败")
		}
		return
	}
	if dbProfile.UserID != userID {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("用户不存在")
		return
	}

	// 修改用户资料
	FixedProfile := alterProfile(req.GetAlterProfile(), dbProfile)

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

	mysqlDB := userprofilemysql.GetDB()
	if mysqlDB == nil || mysqlDB.Error != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("获取数据库失败")
		userCtx.ErrorWF("OnAlterUserProfile userprofilemysql.GetDB nil")
		return
	}
	err = mysqlDB.Transaction(func(tx *gorm.DB) error {
		dbProfile = convertToDbProfile(FixedProfile, dbProfile.CreatedAt, time.Now())
		if err = tx.Save(dbProfile).Error; err != nil {
			userCtx.ErrorWF("OnAlterUserProfile update db profile failed", zap.Error(err))
			return err
		}

		cache.Add(userID, req.GetAlterProfile())

		if err = userprofileredis.DeleteCache(userID); err != nil {
			cache.Remove(userID)
			tx.Rollback()
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("删除缓存失败")
			userCtx.ErrorWF("OnAlterUserProfile delete cache failed", zap.Error(err))
			return err
		}
		return nil
	})
	res.UserProfile = FixedProfile

	return nil
}

// notifyProfileChange 通知用户资料变更
func notifyProfileChange(ctx fklog.FKLogI, userId uint64) {
	// TODO: 实现通知逻辑，可以通过消息队列或其他方式通知相关服务
}

// 修改需要更改的资料
func alterProfile(alterProfile *UserProfile.UserProfile,
	dbProfile *userprofilemysql.UserProfile) *UserProfile.UserProfile {
	ret := convertToProfile(dbProfile)
	if alterProfile.GetNickName() != "" {
		ret.NickName = alterProfile.NickName
	}
	if alterProfile.GetIconUrl() != "" {
		ret.IconUrl = alterProfile.IconUrl
	}
	if alterProfile.GetSex() != 0 {
		ret.Sex = alterProfile.Sex
	}
	return ret
}
