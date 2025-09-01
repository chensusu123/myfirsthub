package mazecollectredis

import (
	"context"
	"fmt"
	"maze_game_server/pb/server/MazeCollectCache"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var gRedis = &fkredis.FkRedis{}
var collectKsy = "maze:collect:info:%d"

const field = 1

func init() {
	// 21726 maze:collect:info:%d
	fkconfig.RegisterNameNode("babycollectredis", 21726, gRedis)
}

func SetCollectInfo(ctx context.Context, uid uint64, info *MazeCollectCache.MazeCollectInfo) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	bts, err := proto.Marshal(info)
	if err != nil {
		return
	}
	_, err = gRedis.Do(context.TODO(), "set", fmt.Sprintf(collectKsy, uid), bts)
	logger.CtxInfo(ctx, "SetCollectInfo", zap.Uint64("userId", uid), zap.Any("info", info))
	return
}

func GetCollectInfo(ctx context.Context, uid uint64) (info *MazeCollectCache.MazeCollectInfo, err error) {
	logger := fklog.ContextAppLogger(ctx)
	bts, err := redis.Bytes(gRedis.Do(context.TODO(), "get", fmt.Sprintf(collectKsy, uid)))
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil
		}
		return nil, err
	}
	info = &MazeCollectCache.MazeCollectInfo{}
	err = proto.Unmarshal(bts, info)
	logger.CtxInfo(ctx, "GetCollectInfo", zap.Uint64("userId", uid), zap.Any("info", info))
	return
}

func GMDel(ctx context.Context, userId uint64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf(collectKsy, userId)
	_, err = redis.Int64(gRedis.Do(context.TODO(), "del", key))
	if err != nil {
		logger.CtxError(ctx, "GMDel fail", zap.String("key", key), zap.Error(err))
		return
	}
	logger.CtxInfo(ctx, "GMDel succ", zap.Any("key", key))
	return
}
