/*
* @Author: majian
* @Date: 2021-03-05 01:08
 */
package gm

import (
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkutil"
	"go.uber.org/zap"
	"net/http"
)

const USER_ID_FIELD = "userId"

func SafeHttpRegister(logger fklog.FKLogI, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		defer fkutil.CaptureException()

		request.ParseForm()
		uid := fkutil.ToUint64(request.Form.Get(USER_ID_FIELD))

		//if !CheckGM.CheckGMOnline(context.TODO(), logger, uid, pattern, request.RemoteAddr) {
		//	return
		//}

		logger.WarnWF("execute gm", zap.Uint64("userId", uid),
			zap.String("pattern", pattern), zap.Any("header", request.Header),
			zap.Any("host", request.Host), zap.Any("remoteAddr", request.RemoteAddr))
		handler(writer, request)
	})
}
