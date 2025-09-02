package mazeuserbarrierredis

import (
	"context"
	"fmt"
	"maze_game_server/pb/server/MazeBarrierCache"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

//doll:maze:u:%d:barrier:%d
var gRedis = &fkredis.FkRedis{}

func init() {
	fkconfig.RegisterNameNode("mazeuserbarrierredis", 21688, gRedis)
}

func GMDel(ctx context.Context, userId uint64, barrierId int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:u:%d:barrier:%d", userId, 0)
	_, err = redis.Int(gRedis.Do(ctx, "DEL", key))
	if err != nil {
		logger.CtxError(ctx, "GMDel fail", zap.String("key", key), zap.Error(err))
		return
	}
	return
}

func GetUserBarrierInfo(ctx context.Context, userId uint64, barrierId int32) (data *MazeBarrierCache.MazeBarrierCache, err error) {
	logger := fklog.ContextAppLogger(ctx)
	data = &MazeBarrierCache.MazeBarrierCache{}
	key := fmt.Sprintf("maze:u:%d:barrier:%d", userId, 0)

	res, err := redis.Bytes(gRedis.Do(ctx, "get", key))
	if err == redis.ErrNil {
		err = nil
		logger.CtxInfo(ctx, "GetUserBarrierInfo nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.CtxError(ctx, "GetUserBarrierInfo get fail", zap.String("key", key), zap.Error(err))
		return
	}

	err = proto.Unmarshal(res, data)
	if err != nil {
		logger.CtxError(ctx, "GetUserBarrierInfo Unmarshal fail", zap.String("key", key), zap.Error(err))
		return
	}
	logger.CtxInfo(ctx, "GetUserBarrierInfo succ", zap.String("key", key), zap.Any("data", data))
	return
}

func SetUserBarrierInfo(ctx context.Context, userId uint64, barrierId int32, data *MazeBarrierCache.MazeBarrierCache) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	res, err := proto.Marshal(data)
	if err != nil {
		logger.CtxError(ctx, "SetUserBarrierInfo marshal fail", zap.Any("data", data), zap.Error(err))
		return
	}

	key := fmt.Sprintf("maze:u:%d:barrier:%d", userId, 0)
	_, err = gRedis.Do(ctx, "set", key, res)
	if err != nil {
		logger.CtxError(ctx, "SetUserBarrierInfo set fail", zap.String("key", key), zap.Error(err))
		return
	}
	logger.CtxInfo(ctx, "SetUserBarrierInfo succ", zap.String("key", key), zap.Any("data", data))
	return
}
