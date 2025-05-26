package mazetempbuffredis

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
	"maze_game_server/pb/server/MazeTempBuffSvr"
)

/**
* @Description: 迷宫临时buff
* @Author: wangyongliang
* @Date: 2025/3/21 15:38
**/

var gRedis = fkredis.FkRedis{}

func init() {
	fkconfig.RegisterNameNode("MazeTempBuffRedis", 21724, &gRedis)
}

func getKey(userId uint64, stateId int32) string {
	return fmt.Sprintf("u:%d:stage:%d:temp:buff", userId, stateId)
}

func SetMazeTempBuff(logger fklog.FKLogI, userId uint64, stateId int32, buffInfo *MazeTempBuffSvr.TempBuffInfo) error {
	key := getKey(userId, stateId)
	data, err := proto.Marshal(buffInfo)
	if err != nil {
		return err
	}
	_, err = gRedis.Do(context.Background(), "SET", key, data)
	if err != nil {
		logger.ErrorWF("SetMazeTempBuff", zap.String("key", key), zap.Any("buffInfo", buffInfo),
			zap.Error(err))
		return err
	}
	logger.InfoWF("SetMazeTempBuff end", zap.String("key", key), zap.Any("buffInfo", buffInfo))
	return nil
}

func GetMazeTempBuff(logger fklog.FKLogI, userId uint64, stateId int32) (*MazeTempBuffSvr.TempBuffInfo, error) {
	key := getKey(userId, stateId)
	res, err := redis.Bytes(gRedis.Do(context.Background(), "GET", key))
	if err == redis.ErrNil {
		return nil, nil
	}

	if err != nil {
		logger.ErrorWF("GetMazeTempBuff GET", zap.String("key", key), zap.Error(err))
		return nil, err
	}

	buffInfo := &MazeTempBuffSvr.TempBuffInfo{}
	err = proto.Unmarshal(res, buffInfo)
	if err != nil {
		logger.ErrorWF("GetMazeTempBuff Unmarshal", zap.String("key", key), zap.Error(err))
		return nil, err
	}
	logger.DebugWF("GetMazeTempBuff end", zap.String("key", key), zap.Any("buffInfo", buffInfo))
	return buffInfo, nil
}

func DelMazeTempBuff(logger fklog.FKLogI, userId uint64, stateId int32) error {
	key := getKey(userId, stateId)
	_, err := gRedis.Do(context.Background(), "DEL", key)
	if err != nil {
		logger.ErrorWF("DelMazeTempBuff", zap.String("key", key), zap.Error(err))
		return err
	}
	logger.InfoWF("DelMazeTempBuff end", zap.String("key", key))
	return nil
}
