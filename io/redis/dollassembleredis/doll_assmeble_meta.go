/*
 * @Author: majian
 * @Date: 2024-07-04 14:35:02
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-20 20:55:10
 */
package dollassembleredis

import (
	"context"
	"fmt"

	"maze_game_server/common/constdef"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

type DollAssembleMetaSt struct {
	CurAssmebleSuitIndex *int32 // 当前生效的装备套 从1开始
	LastSwitchSuitTime   *int64 //  上次切换套装时间
	InitState            *int64 // 初始化状态
}

func (m *DollAssembleMetaSt) GetCurSuitIndex() int32 {
	if m != nil && m.CurAssmebleSuitIndex != nil {
		return *m.CurAssmebleSuitIndex
	}
	return 0
}

func (m *DollAssembleMetaSt) GetSwitchTime() int64 {
	if m != nil && m.LastSwitchSuitTime != nil {
		return *m.LastSwitchSuitTime
	}
	return 0
}

func (m *DollAssembleMetaSt) GetInitState() int64 {
	if m != nil && m.InitState != nil {
		return *m.InitState
	}
	return 0
}

func GetDollAssembleMetaInfo(ctx context.Context, userId uint64, fields ...string) (result *DollAssembleMetaSt, err error) {
	key := fmt.Sprintf("maze:assemble:info:u:%d", userId)
	logger := fklog.ContextAppLogger(ctx)
	var args []interface{}
	args = append(args, key)
	for _, field := range fields {
		args = append(args, field)
	}
	if len(args) == 1 {
		return
	}
	res, err := redis.ByteSlices(gRedis.Do(context.TODO(), "hmget", args...))
	if err == redis.ErrNil {
		err = nil
		logger.CtxInfo(ctx, "GetDollAssembleMetaInfo hmget nil", zap.String("key", key), zap.Any("fields", fields))
		return
	}
	if err != nil {
		logger.CtxError(ctx, "GetDollAssembleMetaInfo hmget fail", zap.Error(err), zap.String("key", key))
		return
	}
	result = new(DollAssembleMetaSt)
	sLen := len(res)
	if len(fields) != sLen {
		logger.CtxError(ctx, "GetDollAssembleMetaInfo data len not match", zap.String("key", key),
			zap.Int("vLen", sLen), zap.Int("kLen", len(fields)))
		return
	}
	for i := 0; i < sLen; i++ {
		if fields[i] == constdef.AssemblePrefixCurAssembleSuitIndex {
			v, e := fkutil.Bytes2Int64(res[i])
			if e != nil {
				err = e
				return nil, err
			}
			v32 := int32(v)
			result.CurAssmebleSuitIndex = &v32
		}
		if fields[i] == constdef.AssemblePrefixSwitchSuitTime {
			v, e := fkutil.Bytes2Int64(res[i])
			if e != nil {
				err = e
				return nil, err
			}
			result.LastSwitchSuitTime = &v
		}
		if fields[i] == constdef.AssemblePrefixInitEquip {
			v, e := fkutil.Bytes2Int64(res[i])
			if e != nil {
				err = e
				return nil, err
			}
			result.InitState = &v
		}
	}
	logger.CtxInfo(ctx, "GetDollAssembleMetaInfo succ", zap.Any("result", result), zap.String("key", key))
	return result, err
}

func SetDollAssmebleMetaInfo(ctx context.Context, userId uint64, info *DollAssembleMetaSt) (err error) {
	key := fmt.Sprintf("maze:assemble:info:u:%d", userId)
	logger := fklog.ContextAppLogger(ctx)
	var args []interface{}
	args = append(args, key)

	if info.CurAssmebleSuitIndex != nil {
		args = append(args, constdef.AssemblePrefixCurAssembleSuitIndex)
		args = append(args, *info.CurAssmebleSuitIndex)
	}
	if info.LastSwitchSuitTime != nil {
		args = append(args, constdef.AssemblePrefixSwitchSuitTime)
		args = append(args, *info.LastSwitchSuitTime)
	}
	// if info.InitState != nil {
	// 	args = append(args, constdef.AssemblePrefixInitEquip)
	// 	args = append(args, *info.InitState)
	// }
	if len(args) == 1 {
		return
	}
	_, err = gRedis.Do(context.TODO(), "hmset", args...)
	if err != nil {
		logger.CtxError(ctx, "SetDollAssmebleMetaInfo hmset fail", zap.Error(err),
			zap.String("key", key), zap.Any("info", info))
		return
	}

	logger.CtxInfo(ctx, "SetDollAssmebleMetaInfo hmset succ",
		zap.String("key", key), zap.Any("info", info))
	return
}
