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
)

// OnAlterUserProfile 修改用户资料
func OnAlterUserProfile(ctx fklog.FKLogI, shardingID int64, rqMsg proto.Message, rsMsg proto.Message, opData string) (err error) {
	defer fkprometheus.DebugPMT("OnAlterUserProfile")()
	userCtx := fkserver.NewUserContext(context.TODO(), uint64(shardingID), ctx)
	req := rqMsg.(*UserProfile.AlterUserProfileRQ)
	res := rsMsg.(*UserProfile.AlterUserProfileRS)
	res.ErrInfo = errors.NO_ERROR
	userCtx.WarnWF("OnAlterUserProfile with", zap.Any("rq", req))

	addStartTime := time.Now()
	defer func() {
		costTime := time.Since(addStartTime).Seconds()
		userCtx.WarnWF("OnAlterUserProfile end ", zap.Any("req", req), zap.Any("res", res), zap.Float64("costTime", costTime))
	}()

	// 1. 检查权限
	if req.GetAlterProfile() == nil || req.GetAlterProfile().GetUserId() != uint64(shardingID) {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}

	userID := uint64(shardingID)

	// 2. 获取分布式锁
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

	// 3. 更新数据库
	dbProfile := &userprofilemysql.UserProfile{
		UserID:   userID,
		Nickname: req.GetAlterProfile().GetNickName(),
		IconUrl:  uint32(req.GetAlterProfile().GetIconUrl()),
		Sex:      uint8(req.GetAlterProfile().GetSex()),
	}
	if err := userprofilemysql.Update(dbProfile); err != nil {
		userCtx.ErrorWF("OnAlterUserProfile update db profile failed", zap.Error(err))
		return err
	}

	// 4. 更新本地缓存
	cache.Add(userID, req.GetAlterProfile())

	// 5. 删除Redis缓存(下次查询会自动重建缓存)
	if err := userprofileredis.DeleteCache(userID); err != nil {
		userCtx.ErrorWF("OnAlterUserProfile delete cache failed", zap.Error(err))
	}

	res.UserProfile = req.GetAlterProfile()

	return nil
}

// notifyProfileChange 通知用户资料变更
func notifyProfileChange(ctx fklog.FKLogI, userId uint64) {
	// TODO: 实现通知逻辑，可以通过消息队列或其他方式通知相关服务
}
