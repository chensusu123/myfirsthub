/*
 * @Author: majian
 * @Date: 2024-04-26 10:51:44
 * @Last Modified by: majian
 * @Last Modified time: 2024-07-04 18:01:57
 */
package dollassembleredis

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
)

// 查询人偶装备初始化装备
func GetDollEquipInitState(logger fklog.FKLogI, userId uint64) (state int64, err error) {
	state, err = redis.Int64(gRedis.Do(context.TODO(), "hget", fmt.Sprintf("maze:assemble:info:u:%d", userId),
		constdef.AssemblePrefixInitEquip))
	if err == redis.ErrNil {
		err = nil
		logger.InfoWF("GetDollEquipInitState no init")
		return
	}
	if err != nil {
		logger.ErrorWF("GetDollEquipInitState fail", zap.Error(err), zap.Int64("state", state))
		return
	}

	logger.InfoWF("GetDollEquipInitState succ", zap.Int64("state", state))
	return
}

// 设置初始装备标记
func SetDollEquipInitState(logger fklog.FKLogI, userId uint64, state int64) (err error) {
	_, err = gRedis.Do(context.TODO(), "hset", fmt.Sprintf("maze:assemble:info:u:%d", userId), constdef.AssemblePrefixInitEquip, state)
	if err != nil {
		logger.ErrorWF("SetDollEquipInitState hset with err", zap.Error(err), zap.Int64("state", state))
		return
	}

	logger.InfoWF("SetDollEquipInitState succ", zap.Int64("state", state))
	return
}
