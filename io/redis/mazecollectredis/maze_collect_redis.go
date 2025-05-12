package mazecollectredis

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/plate/protodef/MazeCollectCache"

	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
)

var gRedis = &fkredis.FkRedis{}

const field = 1

func init() {
	// 21726 maze:collect:info:%d
	fkconfig.RegisterNameNode("babycollectredis", 21726, gRedis)
}

func SetCollectInfo(logger fklog.FKLogI, uid uint64, info *MazeCollectCache.MazeCollectInfo) (err error) {
	bts, err := proto.Marshal(info)
	if err != nil {
		return
	}
	_, err = gRedis.Do(context.TODO(), "set", fmt.Sprintf("maze:collect:info:%d", uid), bts)
	logger.InfoWF("SetCollectInfo", zap.Uint64("userId", uid), zap.Any("info", info))
	return
}

func GetCollectInfo(logger fklog.FKLogI, uid uint64) (info *MazeCollectCache.MazeCollectInfo, err error) {
	bts, err := redis.Bytes(gRedis.Do(context.TODO(), "get", fmt.Sprintf("maze:collect:info:%d", uid)))
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil
		}
		return nil, err
	}
	info = &MazeCollectCache.MazeCollectInfo{}
	err = proto.Unmarshal(bts, info)
	logger.InfoWF("GetCollectInfo", zap.Uint64("userId", uid), zap.Any("info", info))
	return
}

func GMDel(logger fklog.FKLogI, userId uint64) (err error) {
	key := fmt.Sprintf("maze:collect:info:%d", userId)
	_, err = redis.Int64(gRedis.Do(context.TODO(), "del", key))
	if err != nil {
		logger.ErrorWF("GMDel fail", zap.String("key", key), zap.Error(err))
		return
	}
	logger.InfoWF("GMDel succ", zap.Any("key", key))
	return
}
