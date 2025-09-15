package friendservice

import (
	"context"
	"maze_game_server/model/friendmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) FriendList(ctx context.Context, userId uint64, page, pageSize int32) ([]*friendmodel.FriendInfo, bool, error) {
	logger := fklog.ContextAppLogger(ctx)

	friendModel, err := friendmodel.NewFriendModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "FriendList err", zap.Error(err), zap.Int32("page", page), zap.Int32("pageSize", pageSize))
		return nil, true, err
	}

	// 客户端0是开始
	realyPage := page + 1
	finsh := false

	start, end := (realyPage-1)*pageSize, realyPage*pageSize-1

	count := int32(len(friendModel.FriendList))
	if count < start || start > end || start < 0 || end < 0 {
		return []*friendmodel.FriendInfo{}, true, nil
	}
	if count <= end {
		end = int32(len(friendModel.FriendList))
		finsh = true
	}
	res := friendModel.FriendList[start:end]

	return res, finsh, nil
}
