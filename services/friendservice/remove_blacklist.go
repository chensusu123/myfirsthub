package friendservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/model/friendmodel"
)

func (s *service) RemoveBlacklist(logger fklog.FKLogI, userID, toId uint64) *errors.CodeError {
	blacklistModel, err := friendmodel.NewBlacklistModel(logger, userID)
	if err != nil {
		logger.ErrorWF("RemoveBlacklist GetBlacklist err", zap.Error(err), zap.Uint64("toId", toId))
		return errors.MODULE_ERROR
	}
	index := -1
	for i, b := range blacklistModel.Blacklist {
		if b.UserId == toId {
			index = i
			break
		}
	}
	// 判断是否是在黑名单中
	if index == -1 {
		return errors.COMMON_ERROR_TIPS.WrapMsg("不在黑名单中")
	}
	blacklistModel.Blacklist = append(blacklistModel.Blacklist[:index], blacklistModel.Blacklist[index+1:]...)
	// 删除黑名单
	if err = blacklistModel.Save(logger, userID); err != nil {
		logger.ErrorWF("RemoveBlacklist SetBlacklist err", zap.Error(err), zap.Uint64("toId", toId))
		return errors.MODULE_ERROR
	}
	return nil
}
