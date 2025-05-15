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
// func GetUserProduce(logger fklog.FKLogI, userId uint64) (producePb *DollMazeBarrierCache.DollMazeMoneyDb, err error) {
// 	key := fmt.Sprintf("maze:produce:u:%d",userId)

// 	res, err := redis.Bytes(gRedis.Do(context.TODO(), "get", key))
// 	if err == redis.ErrNil {
// 		err = nil
// 		logger.InfoWF("GetUserProduce nil", zap.String("key", key))
// 		return
// 	}
// 	if err != nil {
// 		logger.ErrorWF("GetUserProduce get fail", zap.String("key", key), zap.Error(err))
// 		return
// 	}

// 	producePb = &DollMazeBarrierCache.DollMazeMoneyDb{}
// 	err = proto.Unmarshal(res, producePb)
// 	if err != nil {
// 		logger.ErrorWF("GetUserProduce unmarshal fail", zap.Error(err), zap.Any("moneyData", res))
// 		return
// 	}

// 	logger.InfoWF("GetUserProduce succ", zap.String("key", key), zap.Any("producePb", producePb))
// 	return
// }

// // 设置生产
// func SetUserProduce(logger fklog.FKLogI, userId uint64, producePb *DollMazeBarrierCache.DollMazeMoneyDb) (err error) {

// 	data, err := proto.Marshal(producePb)
// 	if err != nil {
// 		logger.ErrorWF("SetUserProduce Marshal fail", zap.Error(err), zap.Any("producePb", producePb))
// 		return
// 	}
// 	key := fmt.Sprintf("maze:produce:u:%d",userId)
// 	_, err = gRedis.Do(context.TODO(), "set", key, data)
// 	if err != nil {
// 		logger.ErrorWF("SetUserProduce set fail", zap.String("key", key), zap.Any("producePb", producePb), zap.Error(err))
// 		return
// 	}

// 	logger.InfoWF("SetUserProduce succ", zap.Any("producePb", producePb))
// 	return
// }

func ClearProduce(logger fklog.FKLogI, userId uint64, barrierId int32) (err error) {
	key := fmt.Sprintf("maze:produce:u:%d", userId)
	_, err = redis.Int(gRedis.Do(context.TODO(), "DEL", key))
	if err != nil {
		logger.ErrorWF("ClearProduce fail", zap.String("key", key), zap.Error(err))
		return
	}
	return
}

// gm删除
func GMDel(logger fklog.FKLogI, userId uint64) (err error) {
	key := fmt.Sprintf("maze:produce:u:%d", userId)
	_, err = redis.Int(gRedis.Do(context.TODO(), "DEL", key))
	if err != nil {
		logger.ErrorWF("GMDel fail", zap.String("key", key), zap.Error(err))
		return
	}
	return
}
