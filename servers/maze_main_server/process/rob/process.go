// @Author: ZhaoXiming 2025/3/24 14:12
// @Desc:

package rob

import (
	"maze_game_server/lib/nano/component"
	"net/http"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

type Rob struct {
	component.Base
}

func NewRob() *Rob {
	return &Rob{}
}

func RegTcpHandler() {

	// // 迷宫掠夺列表
	// websocket_service.RegProcSimple(10488, &MazeRobGuaJi.MazeRobGuaJiListRQ{},
	// 	10489, &MazeRobGuaJi.MazeRobGuaJiListRS{}, OnMazeRobGuaJiListRQ)

	// // 迷宫掠夺
	// websocket_service.RegProcSimple(10490, &MazeRobGuaJi.MazeRobGuaJiRQ{},
	// 	10491, &MazeRobGuaJi.MazeRobGuaJiRS{}, OnMazeRobGuaJiRQ)

}

func SafeHttpRegister(logger fklog.FKLogI, pattern string, handler func(fklog.FKLogI, http.ResponseWriter, *http.Request)) {
	appConfig := appconfig.GlobalConfig()
	// /s4/AddExp
	pattern = "/s" + appConfig.Global.SectionID + pattern
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
