/*
 * @Author: majian
 * @Date: 2024-12-05 13:56:59
 * @Last Modified by: majian
 * @Last Modified time: 2024-12-05 14:26:11
 */
package equipaassemblegm

import (
	"context"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/buff/BuffManagerRedis"
)

func GetBcBuffBySrc(logger fklog.FKLogI, userID uint64, buffIDs []int32) (srcBuffsMap map[int32]map[int32]int64, e error) {
	buffMap, e := BuffManagerRedis.GetBuffs(context.TODO(), logger, userID, 0, buffIDs)
	if e != nil {
		return
	}
	srcBuffsMap = make(map[int32]map[int32]int64)
	for id, val := range buffMap {
		if id <= 0 || val <= 0 {
			continue
		}
		srcMap, e := BuffManagerRedis.GetBuffDetail(context.TODO(), logger, userID, 0, id)
		if e != nil {
			return srcBuffsMap, e
		}
		for src, value := range srcMap {
			if value <= 0 {
				continue
			}
			if _, ok := srcBuffsMap[src]; !ok {
				buffsMap := make(map[int32]int64)
				srcBuffsMap[src] = buffsMap
				buffsMap[id] = value
			} else {
				srcBuffsMap[src][id] = value
			}
		}
	}
	return srcBuffsMap, nil
}
