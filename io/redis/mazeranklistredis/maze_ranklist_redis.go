package mazeranklistredis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	globalredis "maze_game_server/io/redis"

	"github.com/redis/go-redis/v9"
)

const (
	timeStamp   = 1e13
	RankListKey = "ranklist"
)

func GetRankListKey(name string) string {
	return fmt.Sprintf("%s:%s", RankListKey, name)
}

func getSectionKey(s string) (string, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return "", err
	}
	key := db.MakeSectionKey(s)
	return key, err
}

// 内部会处理分服
func AddRankList(ctx context.Context, name string, score int64, userID uint64) error {
	realScore := float64(score)

	db, err := globalredis.GCli.GetDB()
	key, err := getSectionKey(name)
	if err != nil {
		return err
	}

	return db.ZAdd(context.TODO(), key, redis.Z{
		Score:  realScore,
		Member: userID,
	}).Err()
}

// 批量获取
func GetNowRankListLimitWithScore(ctx context.Context, name string, start int64, stop int64, order bool) ([]redis.Z, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	key, err := getSectionKey(name)
	if err != nil {
		return nil, err
	}

	if order {
		return db.ZRangeWithScores(ctx, key, start, stop).Result()
	}

	return db.ZRevRangeWithScores(ctx, key, start, stop).Result()
}

func GetNowRankListLimit(ctx context.Context, name string, start int64, stop int64, order bool) ([]string, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	key, err := getSectionKey(name)
	if err != nil {
		return nil, err
	}
	if order {
		return db.ZRange(ctx, key, start, stop).Result()
	}

	return db.ZRevRange(ctx, key, start, stop).Result()
}

// 获取玩家排名
func GetUserRank(ctx context.Context, name string, userID uint64, order bool, withScore bool) (int32, int64, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return 0, 0, err
	}
	key, err := getSectionKey(name)
	if err != nil {
		return 0, 0, err
	}

	var rank redis.RankScore
	if withScore {
		if order {
			rank, err = db.ZRankWithScore(ctx, key, strconv.FormatUint(userID, 10)).Result()
		} else {
			rank, err = db.ZRevRankWithScore(ctx, key, strconv.FormatUint(userID, 10)).Result()
		}

		if err == redis.Nil {
			err = errors.New("user not found")
			return 0, 0, err
		}

		if err != nil {
			return 0, 0, err
		}

		return int32(rank.Rank), int64(rank.Score), nil
	}

	var nowRank int64
	if order {
		nowRank, err = db.ZRank(ctx, key, strconv.FormatUint(userID, 10)).Result()
		if err == redis.Nil {
			err = errors.New("user not found")
			return 0, 0, err
		}

		if err != nil {
			return 0, 0, err
		}

		return int32(nowRank), 0, nil
	}

	nowRank, err = db.ZRevRank(ctx, key, strconv.FormatUint(userID, 10)).Result()
	if err != nil {
		if err == redis.Nil {
			err = errors.New("user not found")
			return 0, 0, err
		}
		return 0, 0, err
	}
	return int32(nowRank), 0, nil
}

// 获取在玩家之上和和玩家相同score的所有元素
func GetUserAboveSameScore(ctx context.Context, name string, userID uint64, order bool) ([]redis.Z, error) {
	var rankList []redis.Z
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return nil, err
	}
	key, err := getSectionKey(name)
	if err != nil {
		return nil, err
	}

	// 先获取玩家之上的元素
	if order {
		tmp, err := db.ZRankWithScore(ctx, key, strconv.FormatUint(userID, 10)).Result()
		nowIndex := tmp.Rank
		nowScore := tmp.Score

		if err == redis.Nil {
			err = errors.New("user not found")
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		// 说明用户是当前第一名 防止越界
		if nowIndex == 0 {
			nowIndex = 1
		}

		lastRankList, err := db.ZRangeWithScores(ctx, key, 0, nowIndex-1).Result()
		if err != nil {
			return nil, err
		}
		rankList = append(rankList, lastRankList...)

		// 获取和玩家相同分数的元素
		sameScoreList, err := db.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
			Min: strconv.FormatFloat(float64(nowScore), 'f', -1, 64),
			Max: strconv.FormatFloat(float64(nowScore), 'f', -1, 64),
		}).Result()
		if err != nil {
			return nil, err
		}
		rankList = append(rankList, sameScoreList...)
	} else {
		tmp, err := db.ZRevRankWithScore(ctx, key, strconv.FormatUint(userID, 10)).Result()
		nowIndex := tmp.Rank
		nowScore := tmp.Score
		if err != nil {
			if err == redis.Nil {
				err = errors.New("user not found")
				return nil, err
			}
			return nil, err
		}

		if nowIndex == 0 {
			nowIndex = 1
		}

		lastRankList, err := db.ZRevRangeWithScores(ctx, key, 0, nowIndex-1).Result()
		if err != nil {
			return nil, err
		}
		rankList = append(rankList, lastRankList...)

		sameScoreList, err := db.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
			Min: strconv.FormatFloat(float64(nowScore), 'f', -1, 64),
			Max: strconv.FormatFloat(float64(nowScore), 'f', -1, 64),
		}).Result()
		if err != nil {
			return nil, err
		}
		rankList = append(rankList, sameScoreList...)
	}

	// 去重
	resRankList := make([]redis.Z, 0)
	userIDMap := make(map[uint64]struct{})
	for _, v := range rankList {
		tmpuserId := v.Member.(string)
		userId, err := strconv.ParseUint(tmpuserId, 10, 64)
		if err != nil {
			return nil, err
		}
		if _, ok := userIDMap[userId]; !ok {
			userIDMap[userId] = struct{}{}
			resRankList = append(resRankList, v)
		}
	}

	return resRankList, nil
}

// 删除排行榜指定元素
func DelUserRank(ctx context.Context, name string, userID uint64) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}
	key, err := getSectionKey(name)
	if err != nil {
		return err
	}
	return db.ZRem(ctx, key, strconv.FormatUint(userID, 10)).Err()
}

func SetTime(userID uint64) error {
	key := fmt.Sprintf("test:time:%d", userID)
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}

	value := strconv.FormatInt(time.Now().UnixNano(), 10)

	return db.Set(context.TODO(), key, value, 0).Err()
}

func GetTime(userID uint64) (int64, error) {
	key := fmt.Sprintf("test:time:%d", userID)
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return 0, err
	}
	return db.Get(context.TODO(), key).Int64()
}
