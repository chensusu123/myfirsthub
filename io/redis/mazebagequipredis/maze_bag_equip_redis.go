package mazebagequipredis

import (
	"context"
	"errors"

	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/plate/protodef/MazeEquipCache"
	"go.uber.org/zap"
)

var gRedis = &fkredis.FkRedis{}

func init() {
	//21640 maze:bag:equip:info:%llu 迷宫游戏装备背包存储
	fkconfig.RegisterNameNode("mazebagequipredis", 21640, gRedis)
}

func GetEquipInfo(logger fklog.FKLogI, userId uint64, equipGuid int64) (equipInfo *MazeEquipCache.MazeEquipInfoDb, err error) {
	key := gRedis.GetKey(userId)
	ret, err := redis.Bytes(gRedis.Do(context.TODO(), "hget", key, equipGuid))
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil
		}
		logger.ErrorWF("GetEquipInfo get equip failed", zap.Error(err), zap.String("key", key), zap.Any("equipGuid", equipGuid))
		return nil, err
	}
	equipInfo = &MazeEquipCache.MazeEquipInfoDb{}
	err = proto.Unmarshal(ret, equipInfo)
	if err != nil {
		logger.ErrorWF("GetEquipInfo Unmarshal ret error", zap.Any("ret", ret), zap.Error(err))
		return nil, err
	}
	logger.DebugWF("GetEquipInfo succ", zap.Any("equipInfo", equipInfo))
	return equipInfo, nil
}

func GetBatchEquipInfo(logger fklog.FKLogI, userId uint64, equipGuids ...int64) (equipMap map[int64]*MazeEquipCache.MazeEquipInfoDb, err error) {
	key := gRedis.GetKey(userId)
	var args []interface{}
	args = append(args, key)
	for _, equipGuid := range equipGuids {
		args = append(args, equipGuid)
	}
	ret, err := redis.ByteSlices(gRedis.Do(context.TODO(), "HMGET", args...))
	if err != nil {
		logger.ErrorWF("GetBatchEquipInfo get equip failed", zap.Error(err), zap.String("key", key), zap.Any("args", args))
		return nil, err
	}
	equipMap = make(map[int64]*MazeEquipCache.MazeEquipInfoDb, 0)
	for i := 0; i < len(ret); i += 1 {
		info := &MazeEquipCache.MazeEquipInfoDb{}
		err := proto.Unmarshal(ret[i], info)
		if err != nil {
			logger.ErrorWF("GetBatchEquipInfo equipMap Unmarshal fail", zap.Error(err))
			return equipMap, err
		}
		if info.GetEquipGuid() == 0 {
			continue
		}
		equipMap[info.GetEquipGuid()] = info
	}
	return
}

func SaveEquipInfo(logger fklog.FKLogI, userId uint64, equipInfo *MazeEquipCache.MazeEquipInfoDb) (err error) {
	key := gRedis.GetKey(userId)

	data, err := proto.Marshal(equipInfo)
	if err != nil {
		logger.ErrorWF("SaveEquipInfo Marshal error", zap.Int64("equipGuid", equipInfo.GetEquipGuid()),
			zap.Any("equipInfo", equipInfo), zap.Error(err))
		return
	}

	_, err = redis.Int(gRedis.Do(context.TODO(), "hset", key, equipInfo.GetEquipGuid(), data))
	if err != nil {
		logger.ErrorWF("SaveOneAuction hset error", zap.Int64("equipGuid", equipInfo.GetEquipGuid()),
			zap.Any("equipInfo", equipInfo), zap.Error(err))
		return
	}
	return
}

// 保存装备的装配信息
func BatchSaveEquipInfo(logger fklog.FKLogI, userId uint64, equipList []*MazeEquipCache.MazeEquipInfoDb) error {
	args := make([]interface{}, 0, 1+2*len(equipList))

	key := gRedis.GetKey(userId)
	args = append(args, key)
	for _, equip := range equipList {
		args = append(args, equip.GetEquipGuid())
		equipPb, e := proto.Marshal(equip)
		if e != nil {
			return e
		}
		args = append(args, equipPb)
	}
	if len(args) == 1 {
		logger.WarnWF("equip nil")
		return nil
	}
	// redis操作
	_, err := gRedis.Do(context.TODO(), "HMSET", args...)
	if err != nil {
		logger.ErrorWF("BatchSaveEquipInfo redis with fail",
			zap.Error(err),
			zap.Any("equipList", equipList),
			zap.String("key", key))
		return err
	}
	logger.DebugWF("BatchSaveEquipInfo succ",
		zap.Any("equipList", equipList),
		zap.String("key", key))
	return nil
}

// GM DEL背包
func GMDelEquip(logger fklog.FKLogI, userId uint64) (err error) {
	key := gRedis.GetKey(userId)

	_, err = gRedis.Do(context.TODO(), "DEL", key)
	if err != nil {
		logger.ErrorWF("del equip info fail", zap.Error(err), zap.String("key", key))
		return err
	}
	logger.InfoWF("del equip succ ", zap.String("key", key))
	return nil
}

func BatchDelEquip(logger fklog.FKLogI, userId uint64, equipGuids ...int64) (err error) {
	key := gRedis.GetKey(userId)
	args := make([]interface{}, 0)
	args = append(args, key)
	for _, guid := range equipGuids {
		args = append(args, guid) //field
	}
	if len(args) <= 1 {
		logger.WarnWF("no has equip to save", zap.Any("args", args), zap.String("key", key))
		return errors.New("配置参数错误")
	}
	_, err = gRedis.Do(context.TODO(), "HDEL", args...)
	if err != nil {
		logger.ErrorWF("del equip info fail", zap.Error(err), zap.String("key", key), zap.Any("args", args))
		return err
	}
	logger.InfoWF("del equip succ ", zap.String("key", key), zap.Any("equipGuids", equipGuids))
	return nil
}

func GetAllEquipInfo(logger fklog.FKLogI, userId uint64) (equipMap map[int64]*MazeEquipCache.MazeEquipInfoDb, err error) {
	key := gRedis.GetKey(userId)
	batchCount := 300
	cursor := 0 // 初始hscan游标
	equipMap = make(map[int64]*MazeEquipCache.MazeEquipInfoDb, 0)
	for {
		r, err := gRedis.Do(context.TODO(), "HSCAN", key, cursor, "match", "*", "count", batchCount)
		if err != nil {
			return equipMap, err
		}
		var values []interface{}
		ret, err := redis.Values(r, err) // 因hscan返回由新游标和元素数组组成的两个元素的数组，需将reply接口转换成切片接口
		if err != nil {
			return equipMap, err
		}
		_, err = redis.Scan(ret, &cursor, &values) // 提取新游标和元素数组
		res, err := redis.ByteSlices(values, err)
		if err != nil {
			return equipMap, err
		}
		for i := 0; i < len(res); i += 2 {
			info := &MazeEquipCache.MazeEquipInfoDb{}
			err := proto.Unmarshal(res[i+1], info)
			if err != nil {
				logger.ErrorWF("GetAllEquipInfo Unmarshal error", zap.Uint64("userId", userId), zap.String("guid", string(res[i])), zap.String("info", string(res[i+1])), zap.Error(err))
				return equipMap, err
			}
			if info.GetEquipGuid() == 0 {
				logger.WarnWF("GetAllEquipInfo Unmarshal error", zap.String("key", key), zap.String("guid", string(res[i+1])), zap.Any("info", info))
				continue
			}
			equipMap[info.GetEquipGuid()] = info
		}
		// 当新游标为0时，hscan结束
		if cursor == 0 {
			break
		}
	}
	if len(equipMap) > batchCount {
		logger.WarnWF("GetAllEquipInfo equip count too many", zap.String("key", key), zap.Any("equipMap", len(equipMap)))
	}
	//logger.InfoWF("GetAllEquipInfo end", zap.String("key", key), zap.Any("equipMap", equipMap))
	return equipMap, nil
}
