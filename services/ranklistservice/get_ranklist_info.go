package ranklistservice

import (
	"context"
	"errors"
	"maze_game_server/io/redis/mazeranklistredis"
	"maze_game_server/model/ranklistmodel"
	"sort"
	"strconv"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 获取排行榜 从start到stop 不带分数 不带比较器
func (s *service) GetRankListLimit(ctx context.Context, r *ranklistmodel.RankListModel, start int64, stop int64) ([]ranklistmodel.RankItem, error) {
	logger := fklog.ContextAppLogger(ctx)
	if start > stop || start < 0 || stop < 0 {
		return nil, errors.New("start or stop is invalid")
	}

	var rankList []ranklistmodel.RankItem
	rankListKey := r.GetRankListKey()

	tmprankList, err := mazeranklistredis.GetNowRankListLimit(ctx, rankListKey, start, stop, r.Order)
	if err != nil {
		logger.CtxError(ctx, "GetRankListLimit fail",
			zap.String("rankListname", rankListKey),
			zap.Int64("start", start),
			zap.Int64("stop", stop),
			zap.Any("ranklistmodel", r),
			zap.Error(err))
		return nil, err
	}

	for i, v := range tmprankList {
		userID, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			logger.CtxError(ctx, "GetRankListLimit type assert fail",
				zap.String("name", rankListKey),
				zap.Int64("start", start),
				zap.Int64("stop", stop),
				zap.Any("ranklistmodel", r),
				zap.Error(err))
			return nil, err
		}
		rankList = append(rankList, ranklistmodel.RankItem{
			UserID: userID,
			Index:  int32(i),
		})
	}

	return rankList, nil
}

// 获取排行榜 从start到stop 不带分数 带比较器
func (s *service) GetRankListLimitWithComparator(ctx context.Context, r *ranklistmodel.RankListModel, start int64, stop int64, comparator func(a, b *ranklistmodel.RankItem) bool) ([]ranklistmodel.RankItem, error) {
	logger := fklog.ContextAppLogger(ctx)
	rankList, err := s.GetRankListLimit(ctx, r, start, stop)
	if err != nil {
		logger.CtxError(ctx, "GetRankListLimitWithComparator fail",
			zap.String("rankListname", r.GetRankListKey()),
			zap.Int64("start", start),
			zap.Int64("stop", stop),
			zap.Any("ranklistmodel", r),
			zap.Error(err))
		return nil, err
	}

	sort.Slice(rankList, func(i, j int) bool {
		if r.Order {
			if rankList[i].Score != rankList[j].Score {
				return rankList[i].Score < rankList[j].Score
			}
			return comparator(&rankList[i], &rankList[j])
		} else {
			if rankList[i].Score != rankList[j].Score {
				return rankList[i].Score > rankList[j].Score
			}
			return comparator(&rankList[i], &rankList[j])
		}
	})

	for i := range rankList {
		rankList[i].Index = int32(i)
	}

	return rankList, nil
}

// 获取排行榜 从start到stop 带分数 不带比较器
func (s *service) GetRankListLimitWithScore(ctx context.Context, r *ranklistmodel.RankListModel, start int64, stop int64) ([]ranklistmodel.RankItem, error) {
	logger := fklog.ContextAppLogger(ctx)
	if start > stop || start < 0 || stop < 0 {
		return nil, errors.New("start or stop is invalid")
	}

	var rankList []ranklistmodel.RankItem
	rankListKey := r.GetRankListKey()

	tmprankList, err := mazeranklistredis.GetNowRankListLimitWithScore(ctx, rankListKey, start, stop, r.Order)
	if err != nil {
		logger.CtxError(ctx, "GetRankListLimitWithScore fail",
			zap.String("rankListname", rankListKey),
			zap.Int64("start", start),
			zap.Int64("stop", stop),
			zap.Any("ranklistmodel", r),
			zap.Error(err))
		return nil, err
	}

	for i, v := range tmprankList {
		userIDString, ok := v.Member.(string)
		if !ok {
			logger.CtxError(ctx, "GetRankListLimitWithScore type assert fail",
				zap.String("name", rankListKey),
				zap.Int64("start", start),
				zap.Int64("stop", stop),
				zap.Any("ranklistmodel", r),
				zap.Any("member", v.Member),
			)
			return nil, err
		}

		userID, err := strconv.ParseUint(userIDString, 10, 64)
		if err != nil {
			logger.CtxError(ctx, "GetRankListLimitWithScore type assert fail",
				zap.String("name", rankListKey),
				zap.Int64("start", start),
				zap.Int64("stop", stop),
				zap.Any("ranklistmodel", r),
				zap.Error(err))
			return nil, err
		}

		rankList = append(rankList, ranklistmodel.RankItem{
			UserID: userID,
			Score:  int64(v.Score),
			Index:  int32(i),
		})
	}

	return rankList, nil
}

// 获取排行榜 从start到stop 带分数 带比较器
func (s *service) GetRankListLimitWithScoreAndComparator(ctx context.Context, r *ranklistmodel.RankListModel, start int64, stop int64, comparator func(a, b *ranklistmodel.RankItem) bool) ([]ranklistmodel.RankItem, error) {
	logger := fklog.ContextAppLogger(ctx)
	rankList, err := s.GetRankListLimitWithScore(ctx, r, start, stop)
	if err != nil {
		logger.CtxError(ctx, "GetRankListLimitWithScoreAndComparator fail",
			zap.String("rankListname", r.GetRankListKey()),
			zap.Int64("start", start),
			zap.Int64("stop", stop),
			zap.Any("ranklistmodel", r),
			zap.Error(err))
		return nil, err
	}

	sort.Slice(rankList, func(i, j int) bool {
		if r.Order {
			if rankList[i].Score != rankList[j].Score {
				return rankList[i].Score < rankList[j].Score
			}
			return comparator(&rankList[i], &rankList[j])
		} else {
			if rankList[i].Score != rankList[j].Score {
				return rankList[i].Score > rankList[j].Score
			}
			return comparator(&rankList[i], &rankList[j])
		}
	})

	for i := range rankList {
		rankList[i].Index = int32(i)
	}

	return rankList, nil
}
