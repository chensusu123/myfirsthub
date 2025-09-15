package friendservice

import (
	"context"
	"maze_game_server/model/friendmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) Blacklist(ctx context.Context, userId uint64, page, pageSize int32) ([]*friendmodel.BlacklistInfo, bool, error) {
	logger := fklog.ContextAppLogger(ctx)
	nowPage := page + 1
	start, end := (nowPage-1)*pageSize, nowPage*pageSize-1
	blacklistModel, err := friendmodel.NewBlacklistModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "Blacklist err", zap.Error(err), zap.Int32("page", page), zap.Int32("pageSize", pageSize))
		return nil, true, err
	}

	isFinsh := false
	count := int32(len(blacklistModel.Blacklist))
	if count < start {
		return []*friendmodel.BlacklistInfo{}, true, nil
	}
	if count <= end {
		end = int32(len(blacklistModel.Blacklist))
		isFinsh = true
	}
	res := blacklistModel.Blacklist[start:end]

	return res, isFinsh, nil
}
