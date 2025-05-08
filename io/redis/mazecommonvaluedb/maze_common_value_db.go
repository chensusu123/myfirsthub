/*
@Author: xiaobo
@Date: 2025/3/24 11:46
@Description:
*/

package mazecommonvaluedb

import (
	"context"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
)

var db = &fkredis.FkRedis{}

func init() {
	_ = fkconfig.RegisterNameNode("mazecommonvaluedb", 21685, db)
}

func getKey(uid uint64) string {
	return db.GetKey(uid)
}

func BatchGetItem(_ fklog.FKLogI, uid uint64, itemIds []int32) (countMap map[int32]int64, err error) {
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
