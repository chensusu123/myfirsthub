package mazeuserbarrierredis

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
	"maze_game_server/pb/server/MazeBarrierCache"
)

//doll:maze:u:%d:barrier:%d
var gRedis = &fkredis.FkRedis{}

func init() {
	fkconfig.RegisterNameNode("mazeuserbarrierredis", 21688, gRedis)
}

func GMDel(logger fklog.FKLogI, userId uint64, barrierId int32) (err error) {
	key := fmt.Sprintf("maze:u:%d:barrier:%d", userId, 0)
	_, err = redis.Int(gRedis.Do(context.TODO(), "DEL", key))
	if err != nil {
		logger.ErrorWF("GMDel fail", zap.String("key", key), zap.Error(err))
		return
	}
	return
}

func GetUserBarrierInfo(logger fklog.FKLogI, userId uint64, barrierId int32) (data *MazeBarrierCache.MazeBarrierCache, err error) {
	data = &MazeBarrierCache.MazeBarrierCache{}
	key := fmt.Sprintf("maze:u:%d:barrier:%d", userId, 0)

	res, err := redis.Bytes(gRedis.Do(context.TODO(), "get", key))
	if err == redis.ErrNil {
		err = nil
		logger.InfoWF("GetUserBarrierInfo nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.ErrorWF("GetUserBarrierInfo get fail", zap.String("key", key), zap.Error(err))
		return
	}

	err = proto.Unmarshal(res, data)
	if err != nil {
		logger.ErrorWF("GetUserBarrierInfo Unmarshal fail", zap.String("key", key), zap.Error(err))
		return
	}
	logger.InfoWF("GetUserBarrierInfo succ", zap.String("key", key), zap.Any("data", data))
	return
}

func SetUserBarrierInfo(logger fklog.FKLogI, userId uint64, barrierId int32, data *MazeBarrierCache.MazeBarrierCache) (err error) {
	res, err := proto.Marshal(data)
	if err != nil {
		logger.ErrorWF("SetUserBarrierInfo marshal fail", zap.Any("data", data), zap.Error(err))
		return
	}

	key := fmt.Sprintf("maze:u:%d:barrier:%d", userId, 0)
	_, err = gRedis.Do(context.TODO(), "set", key, res)
	if err != nil {
		logger.ErrorWF("SetUserBarrierInfo set fail", zap.String("key", key), zap.Error(err))
		return
	}
	logger.InfoWF("SetUserBarrierInfo succ", zap.String("key", key), zap.Any("data", data))
	return
}
