/*
@Author: xiaobo
@Date: 2024/2/2 10:46
@Description: 信封seq 列表
*/

package mailboxsetdb

import (
	"context"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/plate/protodef/MazeMail"
	"go.uber.org/zap"
	"strconv"
)

func init() {
	_ = fkconfig.RegisterNameNode("mailboxsetdb", 20155, redisConn) //TODO 用新类型
}

var redisConn = &fkredis.FkRedis{}

// u:%d:mail:%d:box
func getSetKey(userId uint64, classType int32) string {
	return redisConn.GetKey(userId, classType)
}

func AddMailBox(logger fklog.FKLogI, userID uint64, classType int32, mailID, seq uint64) (err error) {
	ctx := context.Background()
	key := getSetKey(userID, classType)
	_, err = redisConn.Do(ctx, "ZADD", key, seq, mailID)
	if err != nil {
		logger.ErrorWF("call AddMailBox do", zap.Error(err))
	}
	logger.InfoWF("success", zap.String("key", key), zap.Uint64("mailID", mailID), zap.Any("seq", seq))
	return err
}

func DelMailBox(logger fklog.FKLogI, userID uint64, classType int32, mailIDS []uint64) (err error) {
	_, err = delMailBox(logger, userID, classType, mailIDS)
	return
}

func DelMailBox2(logger fklog.FKLogI, userID uint64, classType int32, mailIDS []uint64) (num int, err error) {
	return delMailBox(logger, userID, classType, mailIDS)
}

func delMailBox(logger fklog.FKLogI, userID uint64, classType int32, mailIDS []uint64) (num int, err error) {
	ctx := context.Background()
	key := getSetKey(userID, classType)
	logger.DebugWF("DelMailBox begin", zap.String("key", key), zap.Any("mailIDS", mailIDS))
	fields := make([]interface{}, 0, len(mailIDS)+1)
	fields = append(fields, key)
	for _, mailID := range mailIDS {
		fields = append(fields, mailID)
	}
	num, err = redis.Int(redisConn.Do(ctx, "ZREM", fields...))
	if err != nil {
		logger.ErrorWF("call DelMailBox do", zap.Error(err))
	}
	return num, err
}

//判断信封是有已经存在
func IsMailBoxExist(logger fklog.FKLogI, userID uint64, classType int32, mailID uint64) (bool, error) {
	ctx := context.Background()
	key := getSetKey(userID, classType)
	_, err := redis.Int(redisConn.Do(ctx, "ZSCORE", key, mailID))
	if err == redis.ErrNil {
		logger.DebugWF("IsMailBoxExist not exist", zap.String("key", key), zap.Uint64("mailID", mailID))
		return false, nil
	}

	if err != nil {
		logger.ErrorWF("IsMailBoxExist do", zap.Error(err))
		return false, err
	}
	logger.InfoWF("success", zap.String("key", key), zap.Uint64("mailID", mailID))
	return true, err
}

func GetMailIndexes(logger fklog.FKLogI, userID uint64, classType int32, maxSeq uint64, limit int32) (mailIndexes []*MazeMail.ReceiptItem, err error) {
	ctx := context.Background()
	key := getSetKey(userID, classType)
	reply, err := redisConn.Do(ctx, "ZREVRANGEBYSCORE", key, maxSeq, "-inf", "WITHSCORES", "LIMIT", 0, limit)
	if err != nil {
		logger.ErrorWF("call GetMailIndexes do", zap.Error(err))
	}

	res, err := redis.ByteSlices(reply, err)
	if err != nil {
		return nil, err
	}

	ret := make([]*MazeMail.ReceiptItem, 0, len(res)/2)

	for i := 0; i < len(res); i += 2 {
		id, err := strconv.ParseUint(string(res[i]), 10, 64)
		if err != nil {
			continue
		}
		seq, err := strconv.ParseUint(string(res[i+1]), 10, 64)
		if err != nil {
			continue
		}
		ret = append(ret, &MazeMail.ReceiptItem{Id: proto.Uint64(id), Sequence: proto.Uint64(seq)})
	}
	return ret, nil
}
