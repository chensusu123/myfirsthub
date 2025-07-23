// 目前model内部用来过滤排行榜相关参数数据
package ranklistmodel

import (
	"fmt"
	"maze_game_server/io/redis/mazeranklistredis"
	"time"
)

const (
	RankListLimitMax = 100
)

// 排行榜数据
type RankItem struct {
	UserID uint64 `json:"user_id"`
	Score  int64  `json:"score"`
	Index  int32  `json:"index"` // 排名 返回数组的时候和下标一致 主要给查询个人排名使用
}

// 排行榜
type RankListModel struct {
	Name  string
	Order bool  // true 升序 false 降序
	Cycle int32 // 周期 0 表示不循环
}

// option 来设置排行榜的配置
type Option func(*RankListModel)

func WithName(name string) Option {
	return func(r *RankListModel) {
		r.Name = name
	}
}

// 周期 0 表示不循环 1 表示每天 2 表示每周 3 表示每月 4 表示每年
func WithCycle(cycle int32) Option {
	return func(r *RankListModel) {
		r.Cycle = cycle
	}
}

func NewRankListModel(opts ...Option) *RankListModel {
	rankList := &RankListModel{}
	for _, opt := range opts {
		opt(rankList)
	}
	return rankList
}

func (r *RankListModel) GetRankListKey() string {
	if r.Cycle == 0 {
		return mazeranklistredis.GetRankListKey(r.Name)
	}
	switch r.Cycle {
	case 1:
		// 天榜
		now := time.Now()
		nowStr := now.Format("20060102")
		return mazeranklistredis.GetRankListKey(fmt.Sprintf("%s:%s", r.Name, nowStr))
	case 2:
		// 周榜
		// 额外处理一下 周一到周日的都统一格式化周一
		now := time.Now()

		weekday := now.Weekday()
		daysToSubtract := int(weekday - time.Monday)
		if daysToSubtract < 0 {
			// 如果是周日，需要减去6天
			daysToSubtract += 7
		}
		monday := now.AddDate(0, 0, -daysToSubtract)
		mondayStr := monday.Format("20060102")
		return mazeranklistredis.GetRankListKey(fmt.Sprintf("%s:%s", r.Name, mondayStr))
	case 3:
		// 月榜
		now := time.Now()
		nowStr := now.Format("200601")
		return mazeranklistredis.GetRankListKey(fmt.Sprintf("%s:%s", r.Name, nowStr))
	case 4:
		// 年榜
		now := time.Now()
		nowStr := now.Format("2006")
		return mazeranklistredis.GetRankListKey(fmt.Sprintf("%s:%s", r.Name, nowStr))
	default:
		return mazeranklistredis.GetRankListKey(r.Name)
	}

}
