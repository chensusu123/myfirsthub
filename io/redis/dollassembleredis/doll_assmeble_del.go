/*
 * @Author: majian
 * @Date: 2024-07-05 20:54:40
 * @Last Modified by: majian
 * @Last Modified time: 2024-07-05 20:56:40
 */
package dollassembleredis

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 删除装配信息
func DelAssmebleInfo(ctx context.Context, userId uint64) error {
	key := fmt.Sprintf("maze:assemble:info:u:%d", userId)
	logger := fklog.ContextAppLogger(ctx)
	_, err := gRedis.Do(context.TODO(), "DEL", key)
	if err != nil {
		logger.CtxError(ctx, "DelAssmebleInfo fail",
			zap.Error(err),
			zap.String("key", key))
		return err
	}
	logger.CtxInfo(ctx, "DelAssmebleInfo succ",
		zap.String("key", key))
	return err
}

func BatchDelAssmebleInfo(ctx context.Context, userId uint64, fields ...string) error {
	key := fmt.Sprintf("maze:assemble:info:u:%d", userId)
	logger := fklog.ContextAppLogger(ctx)
	var args []interface{}
	args = append(args, key)
	for _, field := range fields {
		args = append(args, field)
	}
	_, err := gRedis.Do(context.TODO(), "HDEL", args...)
	if err != nil {
		logger.CtxError(ctx, "BatchDelAssmebleInfo fail",
			zap.Error(err),
			zap.String("key", key), zap.Any("fields", fields))
		return err
	}
	logger.CtxInfo(ctx, "BatchDelAssmebleInfo succ",
		zap.String("key", key), zap.Any("fields", fields))
	return err
}
