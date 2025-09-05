/*
 * @Author: majian
 * @Date: 2024-07-03 20:18:16
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-13 20:32:16
 */
package dollassemblesuitredis

import (
	"context"
	"errors"
	"fmt"

	"maze_game_server/common/function/assemble"
	"maze_game_server/pb/server/MazeEquipCache"

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
	// 21638 maze:dressed:equip:u:%llu 迷宫游戏穿戴装备存储
	fkconfig.RegisterNameNode("dollassemblesuitredis", 21638, gRedis)
}

// 查询指定装备套
func GetDollAssembleSuit(ctx context.Context, userId uint64, index int32, posCnt int) (equipSuit []*MazeEquipCache.MazeEquipPosDb, err error) {
	key := fmt.Sprintf("maze:dressed:equip:u:%d", userId)
	logger := fklog.ContextAppLogger(ctx)
	args := make([]interface{}, 0, 1+posCnt)
	args = append(args, key)
	for i := 1; i <= posCnt; i++ {
		field := assemble.EnCodeAssembleEquipField(index, int32(i))
		args = append(args, field)
	}

	res, err := redis.ByteSlices(gRedis.Do(ctx, "hmget", args...))
	if err == redis.ErrNil {
		err = nil
		logger.CtxInfo(ctx, "GetDollAssembleSuit hmget nil", zap.String("key", key), zap.Int32("index", index))
		return
	}
	if err != nil {
		logger.CtxError(ctx, "GetDollAssembleSuit hmget fail", zap.Error(err), zap.String("key", key), zap.Int32("index", index))
		return
	}

	sLen := len(res)
	for i := 0; i < sLen; i++ {
		assembleDb := &MazeEquipCache.MazeEquipPosDb{}
		e := proto.Unmarshal(res[i], assembleDb)
		if e != nil {
			logger.CtxError(ctx, "GetDollAssembleSuit Unmarshal fail", zap.Error(e),
				zap.String("key", key), zap.Int("i", i), zap.Int32("index", index))
			return nil, e
		}
		equipSuit = append(equipSuit, assembleDb)
	}
	logger.CtxInfo(ctx, "GetDollAssembleSuit hgetall succ", zap.Any("res", equipSuit),
		zap.String("key", key), zap.Int32("index", index))
	return equipSuit, err
}

// 查询指定装备套的指定位置的装备信息
func GetDollAssembleByPos(ctx context.Context, userId uint64, index int32, pos int32) (assembleDb *MazeEquipCache.MazeEquipPosDb, err error) {
	key := fmt.Sprintf("maze:dressed:equip:u:%d", userId)
	logger := fklog.ContextAppLogger(ctx)
	field := assemble.EnCodeAssembleEquipField(index, pos)

	res, err := redis.Bytes(gRedis.Do(ctx, "hget", key, field))
	if err == redis.ErrNil {
		err = nil
		logger.CtxInfo(ctx, "GetDollAssembleByPos hmget nil", zap.String("key", key), zap.Int32("index", index))
		return
	}
	if err != nil {
		logger.CtxError(ctx, "GetDollAssembleByPos hmget fail", zap.Error(err), zap.String("key", key), zap.Int32("index", index))
		return
	}

	assembleDb = &MazeEquipCache.MazeEquipPosDb{}
	e := proto.Unmarshal(res, assembleDb)
	if e != nil {
		logger.CtxError(ctx, "GetDollAssembleByPos Unmarshal fail", zap.Error(e),
			zap.String("key", key), zap.Int32("pos", pos), zap.Int32("index", index))
		return nil, e
	}
	logger.CtxInfo(ctx, "GetDollAssembleByPos hget succ", zap.Any("assembleDb", assembleDb),
		zap.String("key", key), zap.Int32("index", index))
	return assembleDb, err
}

// 查询所有装备套
func GetAllDollAssembleSuit(ctx context.Context, userId uint64) (allEquipSuit map[int32][]*MazeEquipCache.MazeEquipPosDb, err error) {
	key := fmt.Sprintf("maze:dressed:equip:u:%d", userId)
	logger := fklog.ContextAppLogger(ctx)
	res, err := redis.ByteSlices(gRedis.Do(ctx, "hgetall", key))
	if err == redis.ErrNil {
		err = nil
		logger.CtxInfo(ctx, "GetAllDollAssembleSuit hgetall nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.CtxError(ctx, "GetAllDollAssembleSuit hgetall fail", zap.Error(err), zap.String("key", key))
		return
	}
	allEquipSuit = make(map[int32][]*MazeEquipCache.MazeEquipPosDb)
	sLen := len(res)
	for i := 0; i < sLen; i += 2 {
		ks := fkutil.ByteSliceToString(res[i])
		suitIndex, pos := assemble.DecodeAssembleEquipField(ks)
		if suitIndex <= 0 || pos <= 0 {
			continue
		}
		if !assemble.IsValidEquipPos(uint32(pos)) {
			return nil, errors.New("equip pos invalid")
		}
		assembleDb := &MazeEquipCache.MazeEquipPosDb{}
		e := proto.Unmarshal(res[i+1], assembleDb)
		if e != nil {
			return nil, e
		}
		allEquipSuit[suitIndex] = append(allEquipSuit[suitIndex], assembleDb)
		continue
	}
	logger.CtxInfo(ctx, "GetAllDollAssembleSuit hgetall succ", zap.Any("res", allEquipSuit), zap.String("key", key))
	return allEquipSuit, err
}

// 保存人偶装备数据
func SetDollAssembleSuit(ctx context.Context, userId uint64, index int32, equips []*MazeEquipCache.MazeEquipPosDb) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	args := make([]interface{}, 0, 1+2*len(equips))
	if len(equips) == 0 {
		return
	}
	key := fmt.Sprintf("maze:dressed:equip:u:%d", userId)
	args = append(args, key)
	for _, equip := range equips {
		args = append(args, assemble.EnCodeAssembleEquipField(index, equip.GetPos()))
		equipPosPb, e := proto.Marshal(equip)
		if e != nil {
			return e
		}
		args = append(args, equipPosPb)
	}
	if len(args) == 1 {
		logger.CtxWarn(ctx, "equip pos nil")
		return nil
	}
	// redis操作
	_, err = gRedis.Do(ctx, "HMSET", args...)
	if err != nil {
		logger.CtxError(ctx, "SetDollAssembleSuit redis with fail",
			zap.Error(err),
			zap.Any("equips", equips),
			zap.String("key", key))
		return err
	}
	logger.CtxInfo(ctx, "SetDollAssembleSuit succ",
		zap.Any("equips", equips),
		zap.String("key", key))
	return nil
}

func SaveEquipAssembleInfo(ctx context.Context, userId uint64, index int32, equips []*MazeEquipCache.MazeEquipPosInfo) error {
	var equipDbs []*MazeEquipCache.MazeEquipPosDb
	for _, equip := range equips {
		if equip.GetEquipLoadInfo() != nil {
			equipDbs = append(equipDbs, equip.GetEquipLoadInfo())
		}
	}
	return SetDollAssembleSuit(ctx, userId, index, equipDbs)
}

func SaveEquipAssembleInfoV2(ctx context.Context, userId uint64, index int32, equips []*MazeEquipCache.MazeEquipPosInfo) error {
	logger := fklog.ContextAppLogger(ctx)
	var equipDbs []*MazeEquipCache.MazeEquipPosDb
	for _, equip := range equips {
		if equip.GetEquipLoadInfo() != nil {
			savePb := equip.GetEquipLoadInfo()
			if equip.GetEquipLoadInfo().GetActivateMask() > 0 {
				savePb = proto.Clone(equip.GetEquipLoadInfo()).(*MazeEquipCache.MazeEquipPosDb)
				if savePb == nil {
					logger.CtxError(ctx, "SaveEquipAssembleInfoV2 clone fail", zap.Any("equip", equip))
					return errors.New("clone equip fail")
				}
				savePb.ActivateMask = proto.Int32(0)
			}
			equipDbs = append(equipDbs, savePb)
		}
	}
	return SetDollAssembleSuit(ctx, userId, index, equipDbs)
}

// 删除套装信息
func DelEquipSuitInfo(ctx context.Context, userId uint64) error {
	key := fmt.Sprintf("maze:dressed:equip:u:%d", userId)
	logger := fklog.ContextAppLogger(ctx)
	_, err := gRedis.Do(ctx, "DEL", key)
	if err != nil {
		logger.CtxError(ctx, "DelEquipSuitInfo fail",
			zap.Error(err),
			zap.String("key", key))
		return err
	}
	logger.CtxInfo(ctx, "DelEquipSuitInfo succ",
		zap.String("key", key))
	return err
}

// 批量保存装配数据
func BatchSaveDollAssembleSuit(ctx context.Context, userId uint64, equipsMap map[string]*MazeEquipCache.MazeEquipPosDb) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	args := make([]interface{}, 0, 1+2*len(equipsMap))
	if len(equipsMap) == 0 {
		return
	}
	key := fmt.Sprintf("maze:dressed:equip:u:%d", userId)
	args = append(args, key)
	for key, equip := range equipsMap {
		equipPosPb, e := proto.Marshal(equip)
		if e != nil {
			return e
		}
		args = append(args, key, equipPosPb)
	}

	// redis操作
	_, err = gRedis.Do(ctx, "HMSET", args...)
	if err != nil {
		logger.CtxError(ctx, "BatchSaveDollAssembleSuit redis with fail",
			zap.Error(err),
			zap.Any("equips", equipsMap),
			zap.String("key", key))
		return err
	}
	logger.CtxInfo(ctx, "BatchSaveDollAssembleSuit succ",
		zap.Any("equips", equipsMap),
		zap.String("key", key))
	return nil
}
