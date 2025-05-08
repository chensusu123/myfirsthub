package mail_send_queue

//import (
//	"context"
//
//	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
//	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
//	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
//	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
//	"gitlab.ifreetalk.com/plate/protodef/MailBoxSvr"
//	"gitlab.ifreetalk.com/plate/protodef/MazeMail"
//	"go.uber.org/zap"
//)
//
//func init() {
//	//fkcore.RegisterIO(redisQueue, 16252)
//	fkconfig.RegisterNode(16252, redisQueue)
//}
//
//const redisKey = "ppworld:mail:send:queue"
//
//var redisQueue = &fkredis.QueueRedis{MsgKey: redisKey}
//
//func PushNewInfo(ctx context.Context, logger fklog.FKLogI, uid uint64, mailInfo *MazeMail.MailInfo) error {
//	data := &MailBoxSvr.SvrSendMailRQ{
//		ToUser:   proto.Uint64(uid),
//		MailInfo: mailInfo,
//	}
//	cnt, err := proto.Marshal(data)
//	if err != nil {
//		logger.ErrorWF("PushInfo error", zap.Error(err))
//		return err // 240204增加,接口返回Marshal错误
//	}
//	_, err = redisQueue.Do(ctx, "LPUSH", redisQueue.GetKey(), cnt)
//	if err != nil {
//		logger.ErrorWF("PushNewInfo push mail error", zap.Error(err))
//	}
//	return err
//}
