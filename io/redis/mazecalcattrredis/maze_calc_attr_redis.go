/*
 * @Author: majian
 * @Date: 2025-03-11 15:13:04
 * @Last Modified by: majian
 * @Last Modified time: 2025-04-01 19:19:59
 */
package mazecalcattrredis

import (
	"context"
	"errors"
	"fmt"

	"maze_game_server/common/constdef"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

var (
	gRedis = &fkredis.FkRedis{}
)

func init() {
	// 21630 maze:calc:attr:u:%llu 迷宫计算属性存储
	fkconfig.RegisterNameNode("mazecalcattrredis", 21630, gRedis)
}

func SaveMazeCalcAttr(ctx context.Context, userId uint64, attrs map[int32]int64) error {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:calc:attr:u:%d", userId)
	var args []interface{}
	args = append(args, key)
	for k, v := range attrs {
		args = append(args, k)
		args = append(args, v)
	}

	_, err := gRedis.Do(context.TODO(), "HMSET", args...)
	if err != nil {
		logger.CtxError(ctx, "SaveMazeCalcAttr fail",
			zap.Error(err),
			zap.String("key", key),
			zap.Any("attrs", attrs))
		return err
	}
	logger.CtxInfo(ctx, "SaveMazeCalcAttr succ",
		zap.String("key", key),
		zap.Any("attrs", attrs))
	return err
}

func GetAllMazeCalcAttr(ctx context.Context, userId uint64) (attrDbs map[int32]int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:calc:attr:u:%d", userId)
	res, err := redis.ByteSlices(gRedis.Do(context.TODO(), "hgetall", key))
	if err == redis.ErrNil {
		err = nil
		logger.CtxWarn(ctx, "GetAllMazeCalcAttr hvals nil", zap.String("key", key))
		return
	}

	if err != nil {
		logger.CtxError(ctx, "GetAllMazeCalcAttr hvals fail", zap.Error(err), zap.String("key", key))
		return
	}
	attrDbs = make(map[int32]int64)
	for i := 0; i < len(res); i += 2 {
		k, e := fkutil.Bytes2Int64(res[i])
		if e != nil {
			err = e
			return
		}
		v, e := fkutil.Bytes2Int64(res[i+1])
		if e != nil {
			err = e
			return
		}
		if k <= 0 {
			continue
		}
		attrDbs[int32(k)] = v
	}

	logger.CtxInfo(ctx, "GetAllMazeCalcAttr succ",
		zap.Any("attrDbs", attrDbs),
		zap.String("key", key))
	return attrDbs, err
}

func HScanMazeCalcAttr(ctx context.Context, userId uint64) (attrDbs map[int32]int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	var cursor int64 = 0
	key := fmt.Sprintf("maze:calc:attr:u:%d", userId)
	attrDbs = make(map[int32]int64)
	for {
		rs1, err1 := redis.Values(gRedis.Do(context.TODO(), "HSCAN", key, cursor, "match", "*", "count", 100))
		if err1 != nil {
			logger.CtxError(ctx, "HScanMazeCalcAttr error", zap.Error(err1),
				zap.String("key", key))
			return attrDbs, err1
		}
		if len(rs1) != 2 {
			logger.CtxError(ctx, "HScanMazeCalcAttr data rs error",
				zap.String("key", key))
			return attrDbs, errors.New("data rs error")
		}

		cursor, err = redis.Int64(rs1[0], err)
		if err != nil {
			logger.CtxError(ctx, "HScanMazeCalcAttr Int64 error",
				zap.Error(err),
				zap.String("key", key))
			return attrDbs, err
		}
		data, err := redis.ByteSlices(rs1[1], err)
		if err != nil {
			logger.CtxError(ctx, "HScanMazeCalcAttr ByteSlices error",
				zap.Error(err),
				zap.String("key", key))
			return attrDbs, err
		}
		for i := 0; i < len(data); i += 2 {
			id, e := fkutil.Bytes2Int64(data[i])
			if e != nil {
				return attrDbs, e
			}
			val, e := fkutil.Bytes2Int64(data[i+1])
			if e != nil {
				return attrDbs, e
			}
			attrDbs[int32(id)] = val
		}
		if cursor == 0 {
			break
		}
	}
	logger.CtxInfo(ctx, "HScanMazeCalcAttr succ", zap.String("key", key), zap.Int("len", len(attrDbs)))
	return
}

func DelMazeCalcAttr(ctx context.Context, userId uint64) error {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:calc:attr:u:%d", userId)

	_, err := gRedis.Do(context.TODO(), "DEL", key)
	if err != nil {
		logger.CtxError(ctx, "DelMazeCalcAttr fail",
			zap.Error(err),
			zap.String("key", key))
		return err
	}
	logger.CtxInfo(ctx, "DelMazeCalcAttr succ",
		zap.String("key", key))
	return err
}

func BatchGetMazeCalcAttr(ctx context.Context, userId uint64, attrIds []int32) (attrDbs map[int32]int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:calc:attr:u:%d", userId)
	var args []interface{}
	args = append(args, key)
	for _, attrId := range attrIds {
		args = append(args, attrId)
	}
	res, err := redis.ByteSlices(gRedis.Do(context.TODO(), "HMGET", args...))
	if err == redis.ErrNil {
		err = nil
		logger.CtxWarn(ctx, "BatchGetMazeCalcAttr HMGET nil", zap.String("key", key), zap.Any("attrIds", attrIds))
		return
	}

	if err != nil {
		logger.CtxError(ctx, "BatchGetMazeCalcAttr HMGET fail", zap.Error(err),
			zap.String("key", key), zap.Any("attrIds", attrIds))
		return
	}
	if len(attrIds) != len(res) {
		logger.CtxError(ctx, "BatchGetMazeCalcAttr len not match",
			zap.String("key", key), zap.Any("attrIds", attrIds),
			zap.Int("len", len(res)))
		return nil, errors.New("len not match")
	}

	attrDbs = make(map[int32]int64)
	for i := 0; i < len(attrIds); i++ {
		v, e := fkutil.Bytes2Int64(res[i])
		if e != nil {
			return attrDbs, e
		}
		attrDbs[attrIds[i]] = v
	}

	logger.CtxInfo(ctx, "BatchGetMazeCalcAttr succ", zap.Any("attrDbs", attrDbs), zap.String("key", key))
	return attrDbs, err
}

// 删除指定属性
func HDelMazeCalcAttr(ctx context.Context, userId uint64, attrIds []int32) error {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:calc:attr:u:%d", userId)
	args := make([]interface{}, 0, len(attrIds)+1)
	args = append(args, key)
	for _, id := range attrIds {
		args = append(args, id)
	}
	_, err := gRedis.Do(context.TODO(), "HDEL", args...)
	if err != nil {
		logger.CtxError(ctx, "HDelMazeCalcAttr fail",
			zap.Error(err),
			zap.Any("attrIds", attrIds),
			zap.String("key", key))
		return err
	}
	logger.CtxInfo(ctx, "HDelMazeCalcAttr succ",
		zap.Any("attrIds", attrIds),
		zap.String("key", key))
	return err
}

func GetMazeForce(ctx context.Context, userId uint64) (force int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:calc:attr:u:%d", userId)
	force, err = redis.Int64(gRedis.Do(context.TODO(), "HGET", key, constdef.MazeForce))
	var empty bool
	if err == redis.ErrNil {
		err = nil
		empty = true
	}
	if err != nil {
		logger.CtxError(ctx, "GetMazeForce redis err with", zap.Error(err), zap.String("key", key))
	} else {
		logger.CtxInfo(ctx, "GetMazeForce succ", zap.Int64("force", force),
			zap.Bool("empty", empty),
			zap.String("key", key))
	}
	return
}

// 取计算属性长度
func HlenMazeCalcAttr(ctx context.Context, userId uint64) (slen int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:calc:attr:u:%d", userId)
	slen, err = redis.Int64(gRedis.Do(context.TODO(), "HLEN", key))
	if err != nil {
		logger.CtxError(ctx, "HlenMazeCalcAttr fail",
			zap.Error(err),
			zap.String("key", key))
		return slen, err
	}
	logger.CtxInfo(ctx, "HlenMazeCalcAttr succ", zap.Int64("len", slen),
		zap.String("key", key))
	return slen, err
}
