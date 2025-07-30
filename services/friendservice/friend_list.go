package friendservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/model/friendmodel"
)

func (s *service) FriendList(logger fklog.FKLogI, userId uint64, page, pageSize int32) ([]*friendmodel.FriendInfo, *errors.CodeError) {
	start, end := (page-1)*pageSize, page*pageSize-1
	friendModel, err := friendmodel.NewFriendModel(logger, userId)
	if err != nil {
		logger.ErrorWF("FriendList err", zap.Error(err), zap.Int32("page", page), zap.Int32("pageSize", pageSize))
		return nil, errors.MODULE_ERROR
	}
	count := int32(len(friendModel.FriendList))
	if count < start {
		return []*friendmodel.FriendInfo{}, nil
	}
	if count < end {
		end = int32(len(friendModel.FriendList))
	}
	res := friendModel.FriendList[start:end]

	return res, nil
}
