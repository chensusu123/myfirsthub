package friendservice

import (
	"context"
	"maze_game_server/common/errors"
	"maze_game_server/model/friendmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) Blacklist(ctx context.Context, userId uint64, page, pageSize int32) ([]*friendmodel.BlacklistInfo, *errors.CodeError) {
	logger := fklog.ContextAppLogger(ctx)
	start, end := (page-1)*pageSize, page*pageSize-1
	blacklistModel, err := friendmodel.NewBlacklistModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "Blacklist err", zap.Error(err), zap.Int32("page", page), zap.Int32("pageSize", pageSize))
		return nil, errors.MODULE_ERROR
	}
	count := int32(len(blacklistModel.Blacklist))
	if count < start {
		return []*friendmodel.BlacklistInfo{}, nil
	}
	if count < end {
		end = int32(len(blacklistModel.Blacklist))
	}
	res := blacklistModel.Blacklist[start:end]

	return res, nil
}
