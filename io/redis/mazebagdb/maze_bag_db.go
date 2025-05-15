/*
@Author: xiaobo
@Date: 2025/3/21 13:41
@Description:
*/

package mazebagdb

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
)

var db = &fkredis.FkRedis{}

func init() {
	_ = fkconfig.RegisterNameNode("mazebagdb", 21725, db)
}

func getKey(uid uint64) string {
	return fmt.Sprintf("maze:bag:%d", uid)
}

func IncrBagItem(_ fklog.FKLogI, uid uint64, itemId int32, count int64) (curCount int64, err error) {
	curCount, err = redis.Int64(db.Do(context.Background(), "HINCRBY", getKey(uid), itemId, count))
	return
}

func BatchGetBagItem(_ fklog.FKLogI, uid uint64, itemIds []int32) (countMap map[int32]int64, err error) {
	if len(itemIds) == 0 {
		return
	}

	args := make([]interface{}, 0, len(itemIds)+1)
	args = append(args, getKey(uid))
	for _, id := range itemIds {
		args = append(args, id)
	}

	reply, err := redis.Int64s(db.Do(context.Background(), "HMGET", args...))
	if err == redis.ErrNil {
		err = nil
		return
	}
	if err != nil {
		return
	}

	countMap = make(map[int32]int64)
	for idx, count := range reply {
		countMap[itemIds[idx]] = count
	}

	return
}

func BatchDelBagItem(_ fklog.FKLogI, uid uint64, itemIds []int32) (err error) {
	if len(itemIds) == 0 {
		return
	}

	args := make([]interface{}, 0, len(itemIds)+1)
	args = append(args, getKey(uid))
	for _, id := range itemIds {
		args = append(args, id)
	}

	_, err = redis.Int64s(db.Do(context.Background(), "HDEL", args...))
	return
}

func GetAllBagItem(_ fklog.FKLogI, uid uint64, batchCount int) (itemMap map[int32]int64, err error) {
	key := getKey(uid)

	var (
		cursor int64 = 0
		reply  interface{}
		rs     []interface{}
		res    []string
	)
	itemMap = make(map[int32]int64)

	for {
		reply, err = db.Do(context.Background(), "hscan", key, cursor, "match", "*", "count", batchCount)
		if err != nil {
			return
		}

		rs, err = redis.Values(reply, err)
		if err != nil {
			return
		}

		cursor, err = redis.Int64(rs[0], err)
		if err != nil {
			return
		}

		res, err = redis.Strings(rs[1], err)
		if err != nil {
			return
		}

		for i := 0; i < len(res); i += 2 {
			// i  id
			// i+1 count
			itemMap[fkutil.ToInt32(res[i])] = fkutil.ToInt64(res[i+1])
		}

		if cursor == 0 {
			break
		}
	}

	return
}

func DelKey(_ fklog.FKLogI, uid uint64) (err error) {
	_, err = db.Do(context.Background(), "DEL", getKey(uid))
	return
}
