package friendservice

import (
	"context"
	"fmt"
	"maze_game_server/model/friendmodel"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) AddBlacklist(ctx context.Context, userId, toId uint64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	friendModel, err := friendmodel.NewFriendModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "AddBlacklist GetFriends error", zap.Error(err))
		return
	}
	// 1.判断是否是好友
	if isFriend := s.IsFriend(friendModel, toId); isFriend {
		// 先通知对方删除 todo
		if _, err = s.RemoveFriendEvent(ctx, toId, userId); err != nil {
			logger.CtxError(ctx, "RemoveFriendEvent err", zap.Error(err))
			return
		}
		// 删除好友
		if err = s.delFriendAndFriendRequest(ctx, friendModel, userId, toId); err != nil {
			logger.CtxError(ctx, "AddBlacklist delFriendAndFriendRequest err", zap.Error(err), zap.Uint64("toId", toId))
			return
		}
	}
	blacklistModel, err := friendmodel.NewBlacklistModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "AddBlacklist NewBlacklistModel err", zap.Error(err), zap.Uint64("toId", toId))
		return
	}
	// 2.判断是否已经拉黑对方了
	for _, i := range blacklistModel.Blacklist {
		if i.UserId == toId {
			return fmt.Errorf("已经拉黑对方")
		}
	}
	blacklistModel.Blacklist = append(blacklistModel.Blacklist, &friendmodel.BlacklistInfo{
		UserId:   toId,
		CreateAt: time.Now().Unix(),
	})
	// 3.设置黑名单
	if err = blacklistModel.Save(ctx, userId); err != nil {
		logger.CtxError(ctx, "AddBlacklist Save err", zap.Error(err), zap.Uint64("toId", toId))
		return err
	}
	return nil
}

// 是否在黑名单里面
func (s *service) isBlacklist(ctx context.Context, userId, toId uint64) (bool, error) {
	logger := fklog.ContextAppLogger(ctx)
	blacklistModel, err := friendmodel.NewBlacklistModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "IsBlacklist NewBlacklistModel err", zap.Error(err))
		return false, err
	}
	for _, b := range blacklistModel.Blacklist {
		if b.UserId == toId {
			return true, nil
		}
	}
	return false, nil
}
