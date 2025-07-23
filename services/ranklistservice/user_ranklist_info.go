package ranklistservice

import (
	"context"
	"maze_game_server/io/redis/mazeranklistredis"
	"maze_game_server/model/ranklistmodel"
	"sort"
	"strconv"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 获取玩家排名 不带分数 不带比较器
func (s *service) GetUserRank(logger fklog.FKLogI, r *ranklistmodel.RankListModel, userID uint64) (int32, int64, error) {
	rank, _, err := s.GetUserRankWithScore(logger, r, userID)
	if err != nil {
		return 0, 0, err
	}

	return rank, 0, nil
}

// 获取玩家排名 带分数 不带比较器
func (s *service) GetUserRankWithScore(logger fklog.FKLogI, r *ranklistmodel.RankListModel, userID uint64) (int32, int64, error) {
	var rank int32
	var score int64
	var err error
	rankListKey := r.GetRankListKey()

	rank, score, err = mazeranklistredis.GetUserRank(context.TODO(), rankListKey, userID, r.Order, true)
	if err != nil {
		logger.ErrorWF("GetUserRank fail",
			zap.String("rankListname", rankListKey),
			zap.Uint64("userID", userID),
			zap.Error(err))
		return 0, 0, err
	}

	return rank, score, nil
}

// 获取玩家排名 不带分数 带比较器
func (s *service) GetUserRankWithComparator(logger fklog.FKLogI, r *ranklistmodel.RankListModel, userID uint64, comparator func(a, b *ranklistmodel.RankItem) bool) (int32, int64, error) {
	var rank int32
	var err error
	rank, _, err = s.GetUserRankWithScoreAndComparator(logger, r, userID, comparator)
	if err != nil {
		return 0, 0, err
	}

	return rank, 0, nil
}

// 获取玩家排名 带分数 带比较器
func (s *service) GetUserRankWithScoreAndComparator(logger fklog.FKLogI, r *ranklistmodel.RankListModel, userID uint64, comparator func(a, b *ranklistmodel.RankItem) bool) (int32, int64, error) {
	var rank int32
	var score int64
	var err error
	rankListKey := r.GetRankListKey()

	rankList, err := mazeranklistredis.GetUserAboveSameScore(context.TODO(), rankListKey, userID, r.Order)
	if err != nil {
		logger.ErrorWF("GetUserRank fail",
			zap.String("rankListname", rankListKey),
			zap.Uint64("userID", userID),
			zap.Error(err))
		return 0, 0, err
	}

	var tmpRankList []ranklistmodel.RankItem
	for _, v := range rankList {
		var nowRank ranklistmodel.RankItem
		userIDString, ok := v.Member.(string)
		if !ok {
			logger.ErrorWF("GetUserRank type assert fail",
				zap.String("name", rankListKey),
				zap.Uint64("userID", userID),
				zap.Any("member", v.Member))
			return 0, 0, err
		}
		userID, err := strconv.ParseUint(userIDString, 10, 64)
		if err != nil {
			logger.ErrorWF("GetUserRank type assert fail",
				zap.String("name", rankListKey),
				zap.Uint64("userID", userID),
				zap.Error(err))
			return 0, 0, err
		}
		nowRank.UserID = userID
		nowRank.Score = int64(v.Score)
		tmpRankList = append(tmpRankList, nowRank)
	}

	sort.Slice(tmpRankList, func(i, j int) bool {
		if r.Order {
			if tmpRankList[i].Score != tmpRankList[j].Score {
				return tmpRankList[i].Score < tmpRankList[j].Score
			}
			return comparator(&tmpRankList[i], &tmpRankList[j])
		} else {
			if tmpRankList[i].Score != tmpRankList[j].Score {
				return tmpRankList[i].Score > tmpRankList[j].Score
			}
			return comparator(&tmpRankList[i], &tmpRankList[j])
		}
	})

	for i, v := range tmpRankList {
		if v.UserID == userID {
			rank = int32(i)
			score = v.Score
			break
		}
	}

	logger.InfoWF("GetUserRank LoadData",
		zap.String("rankListname", rankListKey),
		zap.Uint64("userID", userID),
		zap.Any("tmpRankList", tmpRankList),
	)

	return rank, score, nil
}
