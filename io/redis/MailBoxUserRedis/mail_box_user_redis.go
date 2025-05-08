/*
@Author: xiaobo
@Date: 2024/2/2 10:42
@Description: seq id
*/

package MailBoxUserRedis

import (
	"context"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
)

func init() {
	_ = fkconfig.RegisterNameNode("MailBoxUserRedis", 16095, redisConn) //TODO 用新类型
}

var redisConn = &fkredis.FkRedis{}

// u:%d:mail:%d:seq
func getSequenceKey(userId uint64, classType int32) string {
	//return fmt.Sprintf("u:%d:mail:%d:seq", userId, classType)
	return redisConn.GetKey(userId, classType)
}

// 获取信封序列id
func GetMailBoxSequence(logger fklog.FKLogI, userID uint64, classType int32) (seq uint64, err error) {
	ctx := context.Background()
	key := getSequenceKey(userID, classType)
	seq, err = redis.Uint64(redisConn.Do(ctx, "INCR", key))
	if err != nil {
		logger.ErrorWF("call GetMailBoxSequence do", zap.String("key", key), zap.Error(err))
		return 0, err
	}

	logger.DebugWF("GetMailBoxSequence success", zap.String("key", key), zap.Uint64("seq", seq))
	return seq, err
}

// 读取
func ReadMailBoxSequence(logger fklog.FKLogI, userID uint64, classType int32) (seq uint64, err error) {
	key := getSequenceKey(userID, classType)
	ctx := context.Background()
	seq, err = redis.Uint64(redisConn.Do(ctx, "GET", key))
	if err == redis.ErrNil {
		return 0, nil
	}
	if err != nil {
		logger.ErrorWF("ReadMailBoxSequence fail", zap.String("key", key), zap.Error(err))
		return 0, err
	}

	logger.DebugWF("read mail sequence success", zap.String("key", key), zap.Uint64("seq", seq))
	return seq, err
}
