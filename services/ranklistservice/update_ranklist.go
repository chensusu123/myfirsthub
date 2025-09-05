package ranklistservice

import (
	"context"
	"maze_game_server/io/redis/mazeranklistredis"
	"maze_game_server/model/ranklistmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 覆盖
func (s *service) UpdateRankList(ctx context.Context, r *ranklistmodel.RankListModel, userID uint64, score int64) error {
	logger := fklog.ContextAppLogger(ctx)
	rankListKey := r.GetRankListKey()
	err := mazeranklistredis.AddRankList(ctx, rankListKey, score, userID)
	if err != nil {
		logger.CtxError(ctx, "UpdateRankList fail",
			zap.String("rankListname", rankListKey),
			zap.Uint64("userID", userID),
			zap.Int64("score", score),
			zap.Any("ranklistmodel", r),
			zap.Error(err))
		return err
	}
	return nil
}
