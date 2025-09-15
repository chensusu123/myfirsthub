package friendservice

import (
	"context"
	"maze_game_server/model/friendmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) SendFriendRequestList(ctx context.Context, userId uint64, page, pageSize int32) ([]*friendmodel.SendFriendRequestInfo, bool, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 客户端默认传0 避免越界 + 1处理
	nowPage := page + 1
	isFinish := false

	start, end := (nowPage-1)*pageSize, nowPage*pageSize-1
	sendModel, err := friendmodel.NewSendFriendRequestModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "SendFriendRequestList GetSendFriendRequest err", zap.Error(err), zap.Int32("nowPage", nowPage), zap.Int32("pageSize", pageSize))
		return nil, true, err
	}
	count := int32(len(sendModel.SendList))
	if count < start {
		return []*friendmodel.SendFriendRequestInfo{}, true, nil
	}
	if count <= end {
		end = int32(len(sendModel.SendList))
		isFinish = true
	}
	res := sendModel.SendList[start:end]

	return res, isFinish, nil
}
