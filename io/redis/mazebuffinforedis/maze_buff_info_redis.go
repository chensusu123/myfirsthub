/*
 * @Author: majian
 * @Date: 2025-03-15 11:07:31
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-17 21:03:16
 */
package mazebuffinforedis

import (
	"context"
	"fmt"
	"maze_game_server/common/constdef"
	"maze_game_server/pb/server/MazeBuffData"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var (
	gRedis = &fkredis.FkRedis{}
)

func init() {
	// 21639 maze:buff:center:u:%llu 迷宫游戏buff数据存储
	fkconfig.RegisterNameNode("mazebuffinforedis", 21639, gRedis)
}

// 保存非武力值属性
func SaveMazeEquipBuff(ctx context.Context, userId uint64, attrs *MazeBuffData.MazeBuffDb) error {
	return SaveMazeBuffInfo(ctx, userId, constdef.MazeBuffSrcEquip, attrs)
}

func SaveMazeLvBuff(ctx context.Context, userId uint64, attrs *MazeBuffData.MazeBuffDb) error {
	return SaveMazeBuffInfo(ctx, userId, constdef.MazeBuffSrcLv, attrs)
}

func SaveMazeEquipPosBuff(ctx context.Context, userId uint64, attrs *MazeBuffData.MazeBuffDb) error {
	return SaveMazeBuffInfo(ctx, userId, constdef.MazeBuffSrcEquipPos, attrs)
}

func SaveMazeBuffInfo(ctx context.Context, userId uint64, field int32, attrs *MazeBuffData.MazeBuffDb) error {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:buff:center:u:%d", userId)
	data, err := proto.Marshal(attrs)
	if err != nil {
		logger.CtxError(ctx, "SaveMazeBuffInfo Marshal pb fail",
			zap.Error(err),
			zap.String("key", key),
			zap.Any("attrs", attrs))
		return err
	}
	_, err = gRedis.Do(ctx, "HSET", key, field, data)
	if err != nil {
		logger.CtxError(ctx, "SaveMazeBuffInfo fail",
			zap.Error(err),
			zap.String("key", key),
			zap.Int32("field", field),
			zap.Any("attrs", attrs))
		return err
	}
	logger.CtxInfo(ctx, "SaveMazeBuffInfo succ",
		zap.String("key", key),
		zap.Int32("field", field),
		zap.Any("attrs", attrs))
	return err
}

func GetMazeBuffBySrc(ctx context.Context, userId uint64, field int32) (attrDb *MazeBuffData.MazeBuffDb, err error) {
	key := fmt.Sprintf("maze:buff:center:u:%d", userId)
	res, err := redis.Bytes(gRedis.Do(ctx, "hget", key, field))
	logger := fklog.ContextAppLogger(ctx)
	if err == redis.ErrNil {
		err = nil
		logger.CtxWarn(ctx, "GetMazeBuffBySrc hget nil", zap.String("key", key), zap.Int32("field", field))
		return
	}

	if err != nil {
		logger.CtxError(ctx, "GetMazeBuffBySrc hget fail", zap.Error(err), zap.String("key", key), zap.Int32("field", field))
		return
	}

	attrDb = &MazeBuffData.MazeBuffDb{}
	err = proto.Unmarshal(res, attrDb)
	if err != nil {
		logger.CtxError(ctx, "GetMazeBuffBySrc Unmarshal pb fail", zap.Error(err),
			zap.String("key", key), zap.Int32("field", field))
		return nil, err
	}

	logger.CtxInfo(ctx, "GetMazeBuffBySrc succ", zap.Any("attrDb", attrDb), zap.String("key", key),
		zap.Int32("field", field))
	return attrDb, err
}

func GetMazeEquipBuff(ctx context.Context, userId uint64) (attrDb *MazeBuffData.MazeBuffDb, err error) {
	return GetMazeBuffBySrc(ctx, userId, constdef.MazeBuffSrcEquip)
}

// 1=展示属性 2=实际属性
func GetMazeBuffsV2(ctx context.Context, userId uint64, mask int32) (attrDbs map[int32][]*MazeBuffData.MazeBuffAttr, err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:buff:center:u:%d", userId)
	res, err := redis.ByteSlices(gRedis.Do(ctx, "HGETALL", key))
	if err == redis.ErrNil {
		err = nil
		logger.CtxWarn(ctx, "GetMazeBuffsV2 HGETALL nil", zap.String("key", key), zap.Int32("mask", mask))
		return
	}

	if err != nil {
		logger.CtxError(ctx, "GetMazeBuffsV2 HGETALL fail", zap.Error(err), zap.String("key", key), zap.Int32("mask", mask))
		return
	}
	attrDbs = make(map[int32][]*MazeBuffData.MazeBuffAttr)
	for i := 0; i < len(res); i += 2 {
		src, e := fkutil.Bytes2Int64(res[i])
		if e != nil {
			err = e
			return nil, err
		}
		attrDb := &MazeBuffData.MazeBuffDb{}
		err = proto.Unmarshal(res[i+1], attrDb)
		if err != nil {
			logger.CtxError(ctx, "GetMazeBuffsV2 Unmarshal pb fail", zap.Error(err),
				zap.Int32("mask", mask),
				zap.String("key", key))
			return nil, err
		}
		var needAttrs []*MazeBuffData.MazeBuffAttr
		if mask&1 > 0 {
			needAttrs = append(needAttrs, attrDb.MazeShowBuffs...)
		}
		if mask&2 > 0 {
			needAttrs = append(needAttrs, attrDb.MazeRealBuffs...)
		}
		attrDbs[int32(src)] = needAttrs
	}

	logger.CtxInfo(ctx, "GetMazeBuffsV2 succ", zap.Int32("mask", mask), zap.Any("attrDbs", attrDbs), zap.String("key", key))
	return attrDbs, err
}

func GetAllMazeBuffs(ctx context.Context, userId uint64) (attrDbs map[int32]*MazeBuffData.MazeBuffDb, err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:buff:center:u:%d", userId)
	res, err := redis.ByteSlices(gRedis.Do(ctx, "HGETALL", key))
	if err == redis.ErrNil {
		err = nil
		logger.CtxWarn(ctx, "GetAllMazeBuffs HGETALL nil", zap.String("key", key))
		return
	}

	if err != nil {
		logger.CtxError(ctx, "GetAllMazeBuffs HGETALL fail", zap.Error(err), zap.String("key", key))
		return
	}
	attrDbs = make(map[int32]*MazeBuffData.MazeBuffDb)
	for i := 0; i < len(res); i += 2 {
		src, e := fkutil.Bytes2Int64(res[i])
		if e != nil {
			err = e
			return nil, err
		}
		attrDb := &MazeBuffData.MazeBuffDb{}
		err = proto.Unmarshal(res[i+1], attrDb)
		if err != nil {
			logger.CtxError(ctx, "GetAllMazeBuffs Unmarshal pb fail", zap.Error(err),
				zap.String("key", key))
			return nil, err
		}
		attrDbs[int32(src)] = attrDb
	}

	logger.CtxInfo(ctx, "GetAllMazeBuffs succ", zap.Any("attrDbs", attrDbs), zap.String("key", key))
	return attrDbs, err
}

func DelMazeBuff(ctx context.Context, userId uint64) error {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:buff:center:u:%d", userId)

	_, err := gRedis.Do(ctx, "DEL", key)
	if err != nil {
		logger.CtxError(ctx, "DelMazeBuff fail",
			zap.Error(err),
			zap.String("key", key))
		return err
	}
	logger.CtxInfo(ctx, "DelMazeBuff succ",
		zap.String("key", key))
	return err
}

func DelMazeBuffBySrc(ctx context.Context, userId uint64, field int32) error {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:buff:center:u:%d", userId)

	_, err := gRedis.Do(ctx, "HDEL", key, field)
	if err != nil {
		logger.CtxError(ctx, "DelMazeBuffBySrc fail",
			zap.Error(err),
			zap.Int32("field", field),
			zap.String("key", key))
		return err
	}
	logger.CtxInfo(ctx, "DelMazeBuffBySrc succ",
		zap.Int32("field", field),
		zap.String("key", key))
	return err
}
