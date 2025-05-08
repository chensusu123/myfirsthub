package maillimitredis

//import (
//	"context"
//
//	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
//	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
//	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
//	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
//	"go.uber.org/zap"
//)
//
//var client = &fkredis.FkRedis{}
//
//func init() {
//	// 17488	UN_CGK_REDIS_TYPE_USER_AWARD_MAIL_LIMIT_REDIS
//	fkconfig.RegisterNameNode("MailLimitRedis", 17488, client)
//}
//
//// 时间戳
//func PushMailTime(logger fklog.FKLogI, userId uint64, opType int32, time int64) error {
//	key := client.GetKey(userId, opType)
//	_, err := client.Do(context.Background(), "LPUSH", key, time)
//	if err != nil {
//		logger.ErrorWF("push mail limit time fail", zap.String("key", key), zap.Int64("time", time))
//	}
//	return err
//}
//
//func GetMailSendTime(userId uint64, opType int32, num int32) ([]int64, error) {
//	key := client.GetKey(userId, opType)
//	list, err := redis.Int64s(client.Do(context.Background(), "LRANGE", key, 0, num-1))
//	return list, err
//}
//
//// 删除过期发送信封时间列表
//func TrimMailSendTime(logger fklog.FKLogI, userId uint64, opType int32, num int32) error {
//	key := client.GetKey(userId, opType)
//	logger.DebugWF("TrimMailSendTime begin", zap.String("key", key), zap.Int32("num", num))
//	_, err := client.Do(context.Background(), "LTRIM", key, 0, num)
//	if err != nil {
//		logger.ErrorWF("TrimMailSendTime fail", zap.String("key", key), zap.Int32("num", num), zap.Error(err))
//	}
//	return err
//}
