package friendservice

import (
	"maze_game_server/model/friendmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// func (s *service) AddBlacklist(ctx context.Context, userId, toId uint64) *errors.CodeError {
// 	logger := fklog.ContextAppLogger(ctx)
// 	friendModel, err := friendmodel.NewFriendModel(ctx, userId)
// 	if err != nil {
// 		logger.ErrorWF("AddBlacklist GetFriends error", zap.Error(err))
// 		return errors.MODULE_ERROR
// 	}
// 	// 1.判断是否是好友
// 	if isFriend := s.IsFriend(friendModel, toId); isFriend {
// 		// 先通知对方删除 todo
// 		if codeErr := s.RemoveFriendEvent(logger, toId, userId); codeErr != nil {
// 			logger.ErrorWF("RemoveFriendEvent err", zap.Error(codeErr))
// 			return codeErr
// 		}
// 		// 删除好友
// 		if err = s.delFriendAndFriendRequest(logger, friendModel, userId, toId); err != nil {
// 			logger.ErrorWF("AddBlacklist delFriendAndFriendRequest err", zap.Error(err), zap.Uint64("toId", toId))
// 			return errors.MODULE_ERROR
// 		}
// 	}
// 	blacklistModel, err := friendmodel.NewBlacklistModel(logger, userId)
// 	if err != nil {
// 		logger.ErrorWF("AddBlacklist NewBlacklistModel err", zap.Error(err), zap.Uint64("toId", toId))
// 		return errors.MODULE_ERROR
// 	}
// 	// 2.判断是否已经拉黑对方了
// 	for _, i := range blacklistModel.Blacklist {
// 		if i.UserId == toId {
// 			return errors.COMMON_ERROR_TIPS.WrapMsg("已经拉黑对方")
// 		}
// 	}
// 	blacklistModel.Blacklist = append(blacklistModel.Blacklist, &friendmodel.BlacklistInfo{
// 		UserId:   toId,
// 		CreateAt: time.Now().Unix(),
// 	})
// 	// 3.设置黑名单
// 	if err = blacklistModel.Save(logger, userId); err != nil {
// 		logger.ErrorWF("AddBlacklist Save err", zap.Error(err), zap.Uint64("toId", toId))
// 		return errors.MODULE_ERROR
// 	}
// 	return nil
// }

// 是否在黑名单里面
func (s *service) isBlacklist(logger fklog.FKLogI, userId, toId uint64) (bool, error) {
	blacklistModel, err := friendmodel.NewBlacklistModel(logger, userId)
	if err != nil {
		logger.ErrorWF("IsBlacklist NewBlacklistModel err", zap.Error(err))
		return false, err
	}
	for _, b := range blacklistModel.Blacklist {
		if b.UserId == toId {
			return true, nil
		}
	}
	return false, nil
}
