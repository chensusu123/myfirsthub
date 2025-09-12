package friendservice

import (
	"context"
	"maze_game_server/common/errors"
	"maze_game_server/model/friendmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) SendFriendRequestList(ctx context.Context, userId uint64, page, pageSize int32) ([]*friendmodel.SendFriendRequestInfo, error) {
	logger := fklog.ContextAppLogger(ctx)
	start, end := (page-1)*pageSize, page*pageSize-1
	sendModel, err := friendmodel.NewSendFriendRequestModel(ctx, userId)
	if err != nil {
		logger.ErrorWF("SendFriendRequestList GetSendFriendRequest err", zap.Error(err), zap.Int32("page", page), zap.Int32("pageSize", pageSize))
		return nil, errors.MODULE_ERROR
	}
	count := int32(len(sendModel.SendList))
	if count < start {
		return []*friendmodel.SendFriendRequestInfo{}, nil
	}
	if count < end {
		end = int32(len(sendModel.SendList))
	}
	res := sendModel.SendList[start:end]

	return res, nil
}
