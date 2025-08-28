package ranklistservice

import (
	"context"
	"maze_game_server/model/ranklistmodel"
)

type RankListService interface {
	// 更新排行榜
	UpdateRankList(ctx context.Context, r *ranklistmodel.RankListModel, userID uint64, score int64) error
	// 获取排行榜 从start到stop 不带分数 不带比较器
	GetRankListLimit(ctx context.Context, r *ranklistmodel.RankListModel, start int64, stop int64) ([]ranklistmodel.RankItem, error)
	// 获取排行榜 从start到stop 不带分数 带比较器
	GetRankListLimitWithComparator(ctx context.Context, r *ranklistmodel.RankListModel, start int64, stop int64, comparator func(a, b *ranklistmodel.RankItem) bool) ([]ranklistmodel.RankItem, error)
	// 获取排行榜 从start到stop 带分数 不带比较器
	GetRankListLimitWithScore(ctx context.Context, r *ranklistmodel.RankListModel, start int64, stop int64) ([]ranklistmodel.RankItem, error)
	// 获取排行榜 从start到stop 带分数 带比较器
	GetRankListLimitWithScoreAndComparator(ctx context.Context, r *ranklistmodel.RankListModel, start int64, stop int64, comparator func(a, b *ranklistmodel.RankItem) bool) ([]ranklistmodel.RankItem, error)
	// 获取玩家排名 不带分数 不带比较器
	GetUserRank(ctx context.Context, r *ranklistmodel.RankListModel, userID uint64) (int32, int64, error)
	// 获取玩家排名 不带分数 带比较器
	GetUserRankWithComparator(ctx context.Context, r *ranklistmodel.RankListModel, userID uint64, comparator func(a, b *ranklistmodel.RankItem) bool) (int32, int64, error)
	// 获取玩家排名 带分数 不带比较器
	GetUserRankWithScore(ctx context.Context, r *ranklistmodel.RankListModel, userID uint64) (int32, int64, error)
	// 获取玩家排名 带分数 带比较器
	GetUserRankWithScoreAndComparator(ctx context.Context, r *ranklistmodel.RankListModel, userID uint64, comparator func(a, b *ranklistmodel.RankItem) bool) (int32, int64, error)
	// 删除玩家排名
	DelUserRank(ctx context.Context, r *ranklistmodel.RankListModel, userID uint64) error
}

type service struct {
}

var GRankListService RankListService

func init() {
	GRankListService = NewRankListService()
}

func NewRankListService() RankListService {
	return &service{}
}
