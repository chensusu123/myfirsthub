package mazeproduceredis

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
)

//doll:maze:produce:u:%d
var gRedis = &fkredis.FkRedis{}

func init() {
	fkconfig.RegisterNameNode("mazeproduceredis", 21687, gRedis)
}

// // 查询生产
// func GetUserProduce(ctx context.Context, userId uint64) (producePb *DollMazeBarrierCache.DollMazeMoneyDb, err error) {
// 	key := fmt.Sprintf("maze:produce:u:%d",userId)

// 	res, err := redis.Bytes(gRedis.Do(ctx, "get", key))
// 	if err == redis.ErrNil {
// 		err = nil
// 		logger.CtxInfo(ctx,"GetUserProduce nil", zap.String("key", key))
// 		return
// 	}
// 	if err != nil {
// 		logger.CtxError(ctx,"GetUserProduce get fail", zap.String("key", key), zap.Error(err))
// 		return
// 	}

// 	producePb = &DollMazeBarrierCache.DollMazeMoneyDb{}
// 	err = proto.Unmarshal(res, producePb)
// 	if err != nil {
// 		logger.CtxError(ctx,"GetUserProduce unmarshal fail", zap.Error(err), zap.Any("moneyData", res))
// 		return
// 	}

// 	logger.CtxInfo(ctx,"GetUserProduce succ", zap.String("key", key), zap.Any("producePb", producePb))
// 	return
// }

// // 设置生产
// func SetUserProduce(ctx context.Context, userId uint64, producePb *DollMazeBarrierCache.DollMazeMoneyDb) (err error) {

// 	data, err := proto.Marshal(producePb)
// 	if err != nil {
// 		logger.CtxError(ctx,"SetUserProduce Marshal fail", zap.Error(err), zap.Any("producePb", producePb))
// 		return
// 	}
// 	key := fmt.Sprintf("maze:produce:u:%d",userId)
// 	_, err = gRedis.Do(ctx, "set", key, data)
// 	if err != nil {
// 		logger.CtxError(ctx,"SetUserProduce set fail", zap.String("key", key), zap.Any("producePb", producePb), zap.Error(err))
// 		return
// 	}

// 	logger.CtxInfo(ctx,"SetUserProduce succ", zap.Any("producePb", producePb))
// 	return
// }

func ClearProduce(ctx context.Context, userId uint64, barrierId int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:produce:u:%d", userId)
	_, err = redis.Int(gRedis.Do(ctx, "DEL", key))
	if err != nil {
		logger.CtxError(ctx, "ClearProduce fail", zap.String("key", key), zap.Error(err))
		return
	}
	return
}

// gm删除
func GMDel(ctx context.Context, userId uint64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("maze:produce:u:%d", userId)
	_, err = redis.Int(gRedis.Do(ctx, "DEL", key))
	if err != nil {
		logger.CtxError(ctx, "GMDel fail", zap.String("key", key), zap.Error(err))
		return
	}
	return
}
