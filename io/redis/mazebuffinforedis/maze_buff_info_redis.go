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

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
	"maze_game_server/common/constdef"
	"maze_game_server/pb/server/MazeBuffData"
)

var (
	gRedis = &fkredis.FkRedis{}
)

func init() {
	// 21639 maze:buff:center:u:%llu 迷宫游戏buff数据存储
	fkconfig.RegisterNameNode("mazebuffinforedis", 21639, gRedis)
}

// 保存非武力值属性
func SaveMazeEquipBuff(logger fklog.FKLogI, userId uint64, attrs *MazeBuffData.MazeBuffDb) error {
	return SaveMazeBuffInfo(logger, userId, constdef.MazeBuffSrcEquip, attrs)
}

func SaveMazeLvBuff(logger fklog.FKLogI, userId uint64, attrs *MazeBuffData.MazeBuffDb) error {
	return SaveMazeBuffInfo(logger, userId, constdef.MazeBuffSrcLv, attrs)
}

func SaveMazeEquipPosBuff(logger fklog.FKLogI, userId uint64, attrs *MazeBuffData.MazeBuffDb) error {
	return SaveMazeBuffInfo(logger, userId, constdef.MazeBuffSrcEquipPos, attrs)
}

func SaveMazeBuffInfo(logger fklog.FKLogI, userId uint64, field int32, attrs *MazeBuffData.MazeBuffDb) error {
	key := fmt.Sprintf("maze:buff:center:u:%d", userId)
	data, err := proto.Marshal(attrs)
	if err != nil {
		logger.ErrorWF("SaveMazeBuffInfo Marshal pb fail",
			zap.Error(err),
			zap.String("key", key),
			zap.Any("attrs", attrs))
		return err
	}
	_, err = gRedis.Do(context.TODO(), "HSET", key, field, data)
	if err != nil {
		logger.ErrorWF("SaveMazeBuffInfo fail",
			zap.Error(err),
			zap.String("key", key),
			zap.Int32("field", field),
			zap.Any("attrs", attrs))
		return err
	}
	logger.InfoWF("SaveMazeBuffInfo succ",
		zap.String("key", key),
		zap.Int32("field", field),
		zap.Any("attrs", attrs))
	return err
}

func GetMazeBuffBySrc(logger fklog.FKLogI, userId uint64, field int32) (attrDb *MazeBuffData.MazeBuffDb, err error) {
	key := fmt.Sprintf("maze:buff:center:u:%d", userId)
	res, err := redis.Bytes(gRedis.Do(context.TODO(), "hget", key, field))
	if err == redis.ErrNil {
		err = nil
		logger.WarnWF("GetMazeBuffBySrc hget nil", zap.String("key", key), zap.Int32("field", field))
		return
	}

	if err != nil {
		logger.ErrorWF("GetMazeBuffBySrc hget fail", zap.Error(err), zap.String("key", key), zap.Int32("field", field))
		return
	}

	attrDb = &MazeBuffData.MazeBuffDb{}
	err = proto.Unmarshal(res, attrDb)
	if err != nil {
		logger.ErrorWF("GetMazeBuffBySrc Unmarshal pb fail", zap.Error(err),
			zap.String("key", key), zap.Int32("field", field))
		return nil, err
	}

	logger.InfoWF("GetMazeBuffBySrc succ", zap.Any("attrDb", attrDb), zap.String("key", key),
		zap.Int32("field", field))
	return attrDb, err
}

func GetMazeEquipBuff(logger fklog.FKLogI, userId uint64) (attrDb *MazeBuffData.MazeBuffDb, err error) {
	return GetMazeBuffBySrc(logger, userId, constdef.MazeBuffSrcEquip)
}

// 1=展示属性 2=实际属性
func GetMazeBuffsV2(logger fklog.FKLogI, userId uint64, mask int32) (attrDbs map[int32][]*MazeBuffData.MazeBuffAttr, err error) {
	key := fmt.Sprintf("maze:buff:center:u:%d", userId)
	res, err := redis.ByteSlices(gRedis.Do(context.TODO(), "HGETALL", key))
	if err == redis.ErrNil {
		err = nil
		logger.WarnWF("GetMazeBuffsV2 HGETALL nil", zap.String("key", key), zap.Int32("mask", mask))
		return
	}

	if err != nil {
		logger.ErrorWF("GetMazeBuffsV2 HGETALL fail", zap.Error(err), zap.String("key", key), zap.Int32("mask", mask))
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
			logger.ErrorWF("GetMazeBuffsV2 Unmarshal pb fail", zap.Error(err),
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

	logger.InfoWF("GetMazeBuffsV2 succ", zap.Int32("mask", mask), zap.Any("attrDbs", attrDbs), zap.String("key", key))
	return attrDbs, err
}

func GetAllMazeBuffs(logger fklog.FKLogI, userId uint64) (attrDbs map[int32]*MazeBuffData.MazeBuffDb, err error) {
	key := fmt.Sprintf("maze:buff:center:u:%d", userId)
	res, err := redis.ByteSlices(gRedis.Do(context.TODO(), "HGETALL", key))
	if err == redis.ErrNil {
		err = nil
		logger.WarnWF("GetAllMazeBuffs HGETALL nil", zap.String("key", key))
		return
	}

	if err != nil {
		logger.ErrorWF("GetAllMazeBuffs HGETALL fail", zap.Error(err), zap.String("key", key))
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
			logger.ErrorWF("GetAllMazeBuffs Unmarshal pb fail", zap.Error(err),
				zap.String("key", key))
			return nil, err
		}
		attrDbs[int32(src)] = attrDb
	}

	logger.InfoWF("GetAllMazeBuffs succ", zap.Any("attrDbs", attrDbs), zap.String("key", key))
	return attrDbs, err
}

func DelMazeBuff(logger fklog.FKLogI, userId uint64) error {
	key := fmt.Sprintf("maze:buff:center:u:%d", userId)

	_, err := gRedis.Do(context.TODO(), "DEL", key)
	if err != nil {
		logger.ErrorWF("DelMazeBuff fail",
			zap.Error(err),
			zap.String("key", key))
		return err
	}
	logger.InfoWF("DelMazeBuff succ",
		zap.String("key", key))
	return err
}

func DelMazeBuffBySrc(logger fklog.FKLogI, userId uint64, field int32) error {
	key := fmt.Sprintf("maze:buff:center:u:%d", userId)

	_, err := gRedis.Do(context.TODO(), "HDEL", key, field)
	if err != nil {
		logger.ErrorWF("DelMazeBuffBySrc fail",
			zap.Error(err),
			zap.Int32("field", field),
			zap.String("key", key))
		return err
	}
	logger.InfoWF("DelMazeBuffBySrc succ",
		zap.Int32("field", field),
		zap.String("key", key))
	return err
}
