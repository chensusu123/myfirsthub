/*
 * @Author: majian
 * @Date: 2024-08-14 15:56:58
 * @Last Modified by: majian
 * @Last Modified time: 2024-12-16 16:04:52
 */
package dollassembleredis

import (
	"context"
	"fmt"

	"maze_game_server/common/errors"
	"maze_game_server/pb/server/MazeEquipCache"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

// 按指定字段查询装配数据
func GetAssembleInfoByFields(ctx context.Context, uid uint64, fields []string) (assembleInfo *MazeEquipCache.MazeAssembleDb, err error) {
	rsMap, e := hmgetAssembleData(ctx, uid, fields)
	if e != nil {
		err = e
		return
	}
	assembleInfo = &MazeEquipCache.MazeAssembleDb{}
	err = unpackFieldsToPb(rsMap, assembleInfo)
	return assembleInfo, err
}

// 查询全量装配数据
func GetAllAssembleInfo(ctx context.Context, uid uint64) (assembleInfo *MazeEquipCache.MazeAssembleDb, err error) {
	rsMap, e := hgetAllAssembleData(ctx, uid)
	if e != nil {
		err = e
		return
	}
	assembleInfo = &MazeEquipCache.MazeAssembleDb{}
	err = unpackFieldsToPb(rsMap, assembleInfo)
	return assembleInfo, err
}

func hmgetAssembleData(ctx context.Context, uid uint64, fields []string) (rs map[string][]byte, err error) {
	rs = make(map[string][]byte)
	if len(fields) == 0 {
		return
	}

	key := fmt.Sprintf("maze:assemble:info:u:%d", uid)
	args := make([]interface{}, 0, len(fields)+1)
	args = append(args, key)
	for _, field := range fields {
		args = append(args, field)
	}
	rs1, err1 := redis.ByteSlices(gRedis.Do(ctx, "HMGET", args...))
	if err1 != nil {
		return rs, err1
	}

	if len(rs1) != len(fields) {
		err = errors.New("fields data less")
		return rs, err
	}

	for i := 0; i < len(rs1); i++ {
		rs[fields[i]] = rs1[i]
	}
	return
}

func hgetAllAssembleData(ctx context.Context, uid uint64) (rs map[string][]byte, err error) {
	logger := fklog.ContextAppLogger(ctx)
	rs = make(map[string][]byte)
	key := fmt.Sprintf("maze:assemble:info:u:%d", uid)

	rs1, e := redis.ByteSlices(gRedis.Do(ctx, "HGETALL", key))
	if e == redis.ErrNil {
		err = nil
		return
	}
	if e != nil {
		logger.CtxError(ctx, "hgetAllAssembleData with err", zap.Error(e), zap.String("key", key))
		err = e
		return
	}

	for i := 0; i < len(rs1); i += 2 {
		k := fkutil.ByteSliceToString(rs1[i])
		rs[k] = rs1[i+1]
	}
	return
}

func hmsetAssembleData(ctx context.Context, uid uint64, fields map[string]interface{}) (err error) {
	if len(fields) == 0 {
		return
	}
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:assemble:info:u:%d", uid)
	args := make([]interface{}, 0, len(fields)*2+1)
	args = append(args, key)
	for k, v := range fields {
		args = append(args, k)
		args = append(args, v)
	}
	_, err = gRedis.Do(ctx, "HMSET", args...)
	if err != nil {
		logger.CtxError(ctx, "hmsetAssembleData save fail", zap.Error(err), zap.String("key", key))
		return err
	}
	return err
}

// 保存装配数据
func SetAssembleInfoByFields(ctx context.Context, uid uint64, fields []string,
	assembleInfo *MazeEquipCache.MazeAssembleDb) (err error) {
	rs, e := packFieldsFromPb(assembleInfo, fields)
	if e != nil {
		err = e
		return err
	}
	return hmsetAssembleData(ctx, uid, rs)
}
