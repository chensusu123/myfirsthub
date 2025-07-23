package ranklistservice

import (
	"context"
	"maze_game_server/io/redis/mazeranklistredis"
	"maze_game_server/model/ranklistmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) DelUserRank(logger fklog.FKLogI, r *ranklistmodel.RankListModel, userID uint64) error {
	rankListKey := r.GetRankListKey()
	err := mazeranklistredis.DelUserRank(context.TODO(), rankListKey, userID)
	if err != nil {
		logger.ErrorWF("DelUserRank fail",
			zap.String("rankListname", rankListKey),
			zap.Uint64("userID", userID),
			zap.Error(err))
		return err
	}
	return nil
}
