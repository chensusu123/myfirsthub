/*
 * @Author: majian
 * @Date: 2024-08-10 16:39:14
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-17 17:07:18
 */
package dollassembleredis

import (
	"context"
	"fmt"

	"maze_game_server/common/errors"
	"maze_game_server/common/function/assemble"
	"maze_game_server/pb/server/MazeEquipCache"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 获取人偶装备位信息
func GetDollEquipPosInfo(ctx context.Context, userId uint64, posCnt int) (equipList map[int32]*MazeEquipCache.MazeEquipSlotDb, err error) {
	key := fmt.Sprintf("maze:assemble:info:u:%d", userId)
	logger := fklog.ContextAppLogger(ctx)
	args := make([]interface{}, 0, 1+posCnt)
	args = append(args, key)
	for i := 1; i <= posCnt; i++ {
		field := assemble.EnCodeAssemblePosField(int32(i))
		args = append(args, field)
	}

	res, err := redis.ByteSlices(gRedis.Do(context.TODO(), "hmget", args...))
	if err == redis.ErrNil {
		err = nil
		logger.CtxInfo(ctx, "GetDollEquipPosInfo hmget nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.CtxError(ctx, "GetDollEquipPosInfo hmget fail", zap.Error(err), zap.String("key", key))
		return
	}
	equipList = make(map[int32]*MazeEquipCache.MazeEquipSlotDb)
	sLen := len(res)
	for i := 0; i < sLen; i++ {
		if res[i] == nil {
			// 没有数据
			continue
		}
		assembleDb := &MazeEquipCache.MazeEquipSlotDb{}
		e := proto.Unmarshal(res[i], assembleDb)
		if e != nil {
			logger.CtxError(ctx, "GetDollEquipPosInfo Unmarshal fail", zap.Error(e),
				zap.String("key", key), zap.Int("i", i))
			return nil, e
		}
		if assembleDb.GetPos() <= 0 {
			logger.CtxError(ctx, "GetDollEquipPosInfo data err",
				zap.Any("assembleDb", assembleDb),
				zap.String("key", key), zap.Int("i", i))
			return nil, errors.New("data err")
		}
		equipList[assembleDb.GetPos()] = assembleDb
	}
	logger.CtxInfo(ctx, "GetDollEquipPosInfo hgetall succ", zap.Any("res", equipList),
		zap.String("key", key))
	return equipList, err
}

// 更新装备位信息
func SetDollEquipPosInfo(ctx context.Context, userId uint64, posList []*MazeEquipCache.MazeEquipSlotDb) (err error) {
	key := fmt.Sprintf("maze:assemble:info:u:%d", userId)
	logger := fklog.ContextAppLogger(ctx)
	args := make([]interface{}, 0, 1+len(posList))
	args = append(args, key)
	for _, pos := range posList {
		field := assemble.EnCodeAssemblePosField(int32(pos.GetPos()))
		args = append(args, field)
		data, e := proto.Marshal(pos)
		if e != nil {
			err = e
			logger.CtxError(ctx, "SetDollEquipPosInfo Marshal fail", zap.Error(err), zap.String("key", key),
				zap.Any("pos", pos))
			return
		}
		args = append(args, data)
	}
	if len(args) <= 1 {
		return
	}
	_, err = gRedis.Do(context.TODO(), "hmset", args...)
	if err != nil {
		logger.CtxError(ctx, "SetDollEquipPosInfo hmget fail", zap.Error(err), zap.String("key", key),
			zap.Any("posList", posList))
		return
	}
	logger.CtxInfo(ctx, "SetDollEquipPosInfo hgetall succ", zap.Any("posList", posList),
		zap.String("key", key))
	return err
}

// func SetDollEquipPosTotalInfo(ctx context.Context, userId uint64, assembleInfo *MazeEquipCache.MazeAssembleDb) error {
// 	var updateFields []string
// 	if assembleInfo.PkLevel != nil {
// 		updateFields = append(updateFields, constdef.AssemblePrefixPKLevel)
// 	}
// 	posList := assembleInfo.GetMazeEquips()
// 	for _, ePos := range posList {
// 		if ePos.GetEquipPos() != nil {
// 			updateFields = append(updateFields, assemble.EnCodeAssemblePosField(ePos.GetEquipPos().GetPos()))
// 		}
// 	}
// 	err := SetAssembleInfoByFields(logger, userId, updateFields, assembleInfo)
// 	if err != nil {
// 		logger.CtxError(ctx,"SetDollEquipPosTotalInfo fail", zap.Error(err), zap.Any("assembleInfo", assembleInfo))
// 	} else {
// 		logger.CtxInfo(ctx,"SetDollEquipPosTotalInfo succ", zap.Any("assembleInfo", assembleInfo))
// 	}
// 	return err
// }
