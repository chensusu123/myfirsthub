// @Author: ZhaoXiming 2025/3/24 14:12
// @Desc:

package rob

import (
	"net/http"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeRobGuaJi"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
	"go.uber.org/zap"
)

func RegTcpHandler() {

	// 迷宫掠夺列表
	websocket_service.RegProcSimple(16275, &MazeRobGuaJi.MazeRobGuaJiListRQ{},
		16276, &MazeRobGuaJi.MazeRobGuaJiListRS{}, OnMazeRobGuaJiListRQ)

	// 迷宫掠夺
	websocket_service.RegProcSimple(16277, &MazeRobGuaJi.MazeRobGuaJiRQ{},
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
