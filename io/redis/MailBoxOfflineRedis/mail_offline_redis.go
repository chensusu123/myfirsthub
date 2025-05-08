package MailBoxOfflineRedis

//import (
//	"context"
//	"fmt"
//	"time"
//
//	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
//	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
//	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
//	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
//	_ "gitlab.ifreetalk.com/plate/protodef/MazeMail"
//
//	"go.uber.org/zap"
//)
//
//func init() {
//	fkconfig.RegisterNameNode("MailBoxOfflineRedis", 16883, redisConn)
//}
//
//var redisConn = &fkredis.FkRedis{}
//
//func getOfflineKey(userID uint64) string {
//	return fmt.Sprintf("user:mail:offline:push::%d", userID)
//}
//
////设置用户上一次离线推送的时间信息
//func SetUserLastOfflinePush(ctx context.Context, logger fklog.FKLogI, userID uint64) (err error) {
//	key := getOfflineKey(userID)
//
//	now := time.Now().Unix()
//	var expire = 86400 * 2
//
//	_, err = redisConn.Do(ctx, "SET", key, now, "EX", expire)
//	if err != nil {
//		logger.ErrorWF("set user last offline push ", zap.Error(err), zap.String("key", key), zap.Int64("now", now), zap.Any("expire", expire))
//		return err
//	}
//
//	logger.InfoWF("success", zap.String("key", key), zap.Int64("now", now), zap.Any("expire", expire))
//	return err
//}
//
////获取用户上一次离线推送的时间信息
//func GetUserLastOfflinePush(ctx context.Context, logger fklog.FKLogI, userID uint64) (lastTime int64, err error) {
//	key := getOfflineKey(userID)
//
//	lastTime, err = redis.Int64(redisConn.Do(ctx, "GET", key))
//	if err == redis.ErrNil {
//		return 0, nil
//	}
//
//	if err != nil {
//		logger.ErrorWF("get user last offline push", zap.Error(err), zap.String("key", key))
//		return
//	}
//
//	logger.InfoWF("success", zap.String("key", key), zap.Int64("lastTime", lastTime))
//	return
//}
