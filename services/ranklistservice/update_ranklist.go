package ranklistservice

import (
	"context"
	"maze_game_server/io/redis/mazeranklistredis"
	"maze_game_server/model/ranklistmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 覆盖
func (s *service) UpdateRankList(logger fklog.FKLogI, r *ranklistmodel.RankListModel, userID uint64, score int64) error {
	rankListKey := r.GetRankListKey()
	err := mazeranklistredis.AddRankList(context.TODO(), rankListKey, score, userID)
	if err != nil {
		logger.ErrorWF("UpdateRankList fail",
			zap.String("rankListname", rankListKey),
			zap.Uint64("userID", userID),
			zap.Int64("score", score),
			zap.Any("ranklistmodel", r),
			zap.Error(err))
		return err
	}
	return nil
}
