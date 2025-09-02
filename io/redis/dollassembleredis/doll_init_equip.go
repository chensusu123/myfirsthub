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

	"maze_game_server/common/constdef"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
)

// 查询人偶装备初始化装备
func GetDollEquipInitState(ctx context.Context, userId uint64) (state int64, err error) {
	state, err = redis.Int64(gRedis.Do(ctx, "hget", fmt.Sprintf("maze:assemble:info:u:%d", userId),
		constdef.AssemblePrefixInitEquip))
	logger := fklog.ContextAppLogger(ctx)
	if err == redis.ErrNil {
		err = nil
		logger.CtxInfo(ctx, "GetDollEquipInitState no init")
		return
	}
	if err != nil {
		logger.CtxError(ctx, "GetDollEquipInitState fail", zap.Error(err), zap.Int64("state", state))
		return
	}

	logger.CtxInfo(ctx, "GetDollEquipInitState succ", zap.Int64("state", state))
	return
}

// 设置初始装备标记
func SetDollEquipInitState(ctx context.Context, userId uint64, state int64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	_, err = gRedis.Do(ctx, "hset", fmt.Sprintf("maze:assemble:info:u:%d", userId), constdef.AssemblePrefixInitEquip, state)
	if err != nil {
		logger.CtxError(ctx, "SetDollEquipInitState hset with err", zap.Error(err), zap.Int64("state", state))
		return
	}

	logger.CtxInfo(ctx, "SetDollEquipInitState succ", zap.Int64("state", state))
	return
}
