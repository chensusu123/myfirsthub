// @Author: ZhaoXiming 2025/3/24 14:12
// @Desc:

package rob

import (
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"
	"gitlab.ifreetalk.com/plate/freetk/fkutil"
	"gitlab.ifreetalk.com/plate/protodef/MazeRobGuaJi"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func RegTcpHandler() {

	// 迷宫掠夺列表
	tcp_service.RegProcSimple(16275, &MazeRobGuaJi.MazeRobGuaJiListRQ{},
		16276, &MazeRobGuaJi.MazeRobGuaJiListRS{}, OnMazeRobGuaJiListRQ)

	// 迷宫掠夺
	tcp_service.RegProcSimple(16277, &MazeRobGuaJi.MazeRobGuaJiRQ{},
		16278, &MazeRobGuaJi.MazeRobGuaJiRS{}, OnMazeRobGuaJiRQ)

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
