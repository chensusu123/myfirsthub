package friendservice

import (
	"context"
	"maze_game_server/common/errors"
	"maze_game_server/model/friendmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) ReceiveFriendRequestList(ctx context.Context, userId uint64, page, pageSize int32) ([]*friendmodel.ReceiveFriendRequestInfo, bool, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 客户端默认从0开始
	nowPage := page + 1

	start, end := (nowPage-1)*pageSize, nowPage*pageSize-1
	receiveModel, err := friendmodel.NewReceiveFriendRequestModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "ReceiveFriendRequestList NewReceiveFriendRequestModel err", zap.Error(err), zap.Int32("page", page), zap.Int32("pageSize", pageSize))
		return nil, false, errors.MODULE_ERROR
	}
	count := int32(len(receiveModel.ReceiveList))
	if count < start {
		return []*friendmodel.ReceiveFriendRequestInfo{}, false, nil
	}
	if count < end {
		end = int32(len(receiveModel.ReceiveList))
	}
	res := receiveModel.ReceiveList[start:end]

	return res, true, nil
}
