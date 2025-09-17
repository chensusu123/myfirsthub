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
	blacklistModel, err := friendmodel.NewBlacklistModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "AddBlacklist NewBlacklistModel err", zap.Error(err), zap.Uint64("toId", toId))
		return
	}
	// 2.判断是否已经拉黑对方了
	for _, i := range blacklistModel.Blacklist {
		if i.UserId == toId {
			logger.CtxWarn(ctx, "AddBlacklist fail", zap.Any("err", "已经拉黑对方"))
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

	logger.CtxInfo(ctx, "AddBlacklist Successful", zap.Any("userID", userId), zap.Any("toID", toId))
	return nil
}

// 是否在黑名单里面
func (s *service) IsBlacklist(ctx context.Context, userId, toId uint64) (bool, error) {
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
