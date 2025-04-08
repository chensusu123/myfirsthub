package mazebarriertempbuffredis

import (
	"context"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/protodef/MazeTempBuffSvr"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
)

var gRedis = &fkredis.FkRedis{}

func init() {
	fkconfig.RegisterNameNode("mazebarriertempbuffredis", 21724, gRedis)
}

// 清除buff
func ClearBarrierTempBuff(logger fklog.FKLogI, userId uint64, barrierId int32) (err error) {
	key := gRedis.GetKey(userId, barrierId)

	_, err = redis.Int(gRedis.Do(context.TODO(), "del", key))
	if err != nil {
		logger.ErrorWF("ClearBarrierTempBuff del fail", zap.String("key", key), zap.Error(err))
		return
	}

	logger.InfoWF("ClearBarrierTempBuff succ", zap.String("key", key), zap.Any("barrierId", barrierId))
	return
}


// 获取临时buff
func GetBarrierTempBuff(logger fklog.FKLogI, userId uint64, barrierId int32) (data *MazeTempBuffSvr.TempBuffInfo, err error) {
	data = &MazeTempBuffSvr.TempBuffInfo{}
	key := gRedis.GetKey(userId, barrierId)
	res, err := redis.Bytes(gRedis.Do(context.TODO(), "get", key))
	if err == redis.ErrNil {
		err = nil
		logger.InfoWF("GetBarrierTempBuff nil", zap.String("key", key))
		return
	}
	if err != nil {
		logger.ErrorWF("GetBarrierTempBuff get fail", zap.String("key", key), zap.Error(err))
		return
	}
	err = proto.Unmarshal(res, data)
	if err != nil {
		logger.ErrorWF("GetBarrierTempBuff Unmarshal fail", zap.String("key", key), zap.Error(err))
		return
	}
	logger.InfoWF("GetBarrierTempBuff succ", zap.String("key", key), zap.Any("data", data))
	return
}