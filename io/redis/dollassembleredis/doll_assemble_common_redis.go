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

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/pb/server/MazeEquipCache"
)

// 按指定字段查询装配数据
func GetAssembleInfoByFields(agent fklog.FKLogI, uid uint64, fields []string) (assembleInfo *MazeEquipCache.MazeAssembleDb, err error) {
	rsMap, e := hmgetAssembleData(agent, uid, fields)
	if e != nil {
		err = e
		return
	}
	assembleInfo = &MazeEquipCache.MazeAssembleDb{}
	err = unpackFieldsToPb(rsMap, assembleInfo)
	return assembleInfo, err
}

// 查询全量装配数据
func GetAllAssembleInfo(agent fklog.FKLogI, uid uint64) (assembleInfo *MazeEquipCache.MazeAssembleDb, err error) {
	rsMap, e := hgetAllAssembleData(agent, uid)
	if e != nil {
		err = e
		return
	}
	assembleInfo = &MazeEquipCache.MazeAssembleDb{}
	err = unpackFieldsToPb(rsMap, assembleInfo)
	return assembleInfo, err
}

func hmgetAssembleData(logger fklog.FKLogI, uid uint64, fields []string) (rs map[string][]byte, err error) {
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
	rs1, err1 := redis.ByteSlices(gRedis.Do(context.TODO(), "HMGET", args...))
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

func hgetAllAssembleData(agent fklog.FKLogI, uid uint64) (rs map[string][]byte, err error) {
	rs = make(map[string][]byte)
	key := fmt.Sprintf("maze:assemble:info:u:%d", uid)

	rs1, e := redis.ByteSlices(gRedis.Do(context.TODO(), "HGETALL", key))
	if e == redis.ErrNil {
		err = nil
		return
	}
	if e != nil {
		agent.ErrorWF("hgetAllAssembleData with err", zap.Error(e), zap.String("key", key))
		err = e
		return
	}

	for i := 0; i < len(rs1); i += 2 {
		k := fkutil.ByteSliceToString(rs1[i])
		rs[k] = rs1[i+1]
	}
	return
}

func hmsetAssembleData(agent fklog.FKLogI, uid uint64, fields map[string]interface{}) (err error) {
	if len(fields) == 0 {
		return
	}

	key := fmt.Sprintf("maze:assemble:info:u:%d", uid)
	args := make([]interface{}, 0, len(fields)*2+1)
	args = append(args, key)
	for k, v := range fields {
		args = append(args, k)
		args = append(args, v)
	}
	_, err = gRedis.Do(context.TODO(), "HMSET", args...)
	if err != nil {
		agent.ErrorWF("hmsetAssembleData save fail", zap.Error(err), zap.String("key", key))
		return err
	}
	return err
}

// 保存装配数据
func SetAssembleInfoByFields(agent fklog.FKLogI, uid uint64, fields []string,
	assembleInfo *MazeEquipCache.MazeAssembleDb) (err error) {
	rs, e := packFieldsFromPb(assembleInfo, fields)
	if e != nil {
		err = e
		return err
	}
	return hmsetAssembleData(agent, uid, rs)
}
