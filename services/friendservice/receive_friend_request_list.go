package friendservice

import (
	"context"
	"maze_game_server/common/errors"
	"maze_game_server/model/friendmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) ReceiveFriendRequestList(ctx context.Context, userId uint64, page, pageSize int32) ([]*friendmodel.ReceiveFriendRequestInfo, *errors.CodeError) {
	logger := fklog.ContextAppLogger(ctx)
	start, end := (page-1)*pageSize, page*pageSize-1
	receiveModel, err := friendmodel.NewReceiveFriendRequestModel(ctx, userId)
	if err != nil {
		logger.ErrorWF("ReceiveFriendRequestList NewReceiveFriendRequestModel err", zap.Error(err), zap.Int32("page", page), zap.Int32("pageSize", pageSize))
		return nil, errors.MODULE_ERROR
	}
	count := int32(len(receiveModel.ReceiveList))
	if count < start {
		return []*friendmodel.ReceiveFriendRequestInfo{}, nil
	}
	if count < end {
		end = int32(len(receiveModel.ReceiveList))
	}
	res := receiveModel.ReceiveList[start:end]

	return res, nil
}
