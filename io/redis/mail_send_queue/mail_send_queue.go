package mail_send_queue

// import (
//	"context"
//
//	"google.golang.org/protobuf/proto"
//	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
//	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
//	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
//	"maze_game_server/pb/common/MailBoxSvr"
//	"maze_game_server/pb/common/MazeMail"
//	"go.uber.org/zap"
// )
//
// func init() {
//	//fkcore.RegisterIO(redisQueue, 16252)
//	fkconfig.RegisterNode(16252, redisQueue)
// }
//
// const redisKey = "ppworld:mail:send:queue"
//
// var redisQueue = &fkredis.QueueRedis{MsgKey: redisKey}
//
// func PushNewInfo(ctx context.Context, logger fklog.FKLogI, uid uint64, mailInfo *MazeMail.MailInfo) error {
//	data := &MailBoxSvr.SvrSendMailRQ{
//		ToUser:   proto.Uint64(uid),
//		MailInfo: mailInfo,
//	}
//	cnt, err := proto.Marshal(data)
//	if err != nil {
//		logger.CtxError(ctx,"PushInfo error", zap.Error(err))
//		return err // 240204增加,接口返回Marshal错误
//	}
//	_, err = redisQueue.Do(ctx, "LPUSH", redisQueue.GetKey(), cnt)
//	if err != nil {
//		logger.CtxError(ctx,"PushNewInfo push mail error", zap.Error(err))
//	}
//	return err
// }
