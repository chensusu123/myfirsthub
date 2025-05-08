// @Author: ZhaoXiming 2025/3/21 16:06
// @Desc:

package mail

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/mail/MazeMailCli"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/redis_consumer"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/plate/freetk/fkutil"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func RegTcpHandler() {

	// 邮件获取
	tcp_service.RegProcSimple(000, &MazeMailCli.MailListQueryRQ{},
		000, &MazeMailCli.MailListQueryRS{}, OnMailListQueryRQ)

	// 打开附件
	tcp_service.RegProcSimple(000, &MazeMailCli.AnnexOpenRQ{},
		000, &MazeMailCli.AnnexOpenRS{}, OnAnnexOpenRQ)

	// 未读邮件数量
	tcp_service.RegProcSimple(000, &MazeMailCli.MailUnReadRQ{},
		000, &MazeMailCli.MailUnReadRS{}, OnMailUnReadRQ)

	// 读取信封内容
	tcp_service.RegProcSimple(000, &MazeMailCli.GetMailInfoRQ{},
		000, &MazeMailCli.GetMailInfoRS{}, OnGetMailInfoRQ)

	// 一键删除
	tcp_service.RegProcSimple(000, &MazeMailCli.TotalMailDelRQ{},
		000, &MazeMailCli.TotalMailDelRS{}, OnTotalMailDelRQ)

	// 一键领取
	tcp_service.RegProcSimple(000, &MazeMailCli.TotalMailRecvRQ{},
		000, &MazeMailCli.TotalMailRecvRS{}, OnTotalMailRecvRQ)

}

func InitRedis() {
	redis_consumer.PlugRedisConsumer("MailSendServer",
		000, //TODO
		redis_consumer.WithThreadCount(50),
		redis_consumer.WithListCustomContent(MailSendProcess),
	)

}

func CustomProc() {
	//custom.AddCustomProc("item_err_kafka", 1, itemrpc.DoItemErrRecord)
	//custom.AddCustomProc("expire_broadcast_phone", 1, mprocess.ExpireBroadcastPhone)
	//custom.AddCustomProc("mail_award_set_flag_range", 1, mprocess.RangeAwardMailFailFlag)
}

func SafeHttpRegister(logger fklog.FKLogI, pattern string, handler func(fklog.FKLogI, http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		defer fkutil.CaptureException()
		l := logger.Clone("")
		l.SetLogId(time.Now().UnixNano())

		l.DebugWF("execute gm", zap.String("pattern", pattern), zap.Any("header", request.Header),
			zap.Any("host", request.Host), zap.Any("remoteAddr", request.RemoteAddr))

		handler(l, writer, request)
	})
}

func InitHttp(logger fklog.FKLogI) {

}
