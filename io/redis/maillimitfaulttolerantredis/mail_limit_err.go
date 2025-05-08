package maillimitfaulttolerantredis

//import (
//	"context"
//
//	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
//	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
//	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
//	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
//	"gitlab.ifreetalk.com/plate/protodef/MailBoxSvr"
//	"go.uber.org/zap"
//)
//
//var client = &fkredis.FkRedis{}
//
//func init() {
//	// 17489	UN_CGK_REDIS_TYPE_LIMIT_MAIL_FAULT_TOLERANT_QUEUE
//	fkconfig.RegisterNameNode("MailLimitTolerantRedis", 17489, client)
//}
//
//func Push(logger fklog.FKLogI, info *MailBoxSvr.SvrSendMailRQ) error {
//	key := client.GetKey()
//	data, _ := proto.Marshal(info)
//	_, err := client.Do(context.Background(), "LPUSH", key, data)
//	if err != nil {
//		logger.ErrorWF("push mail limit tolerant fail", zap.ByteString("mail", data), zap.Error(err))
//	}
//	return err
//}
