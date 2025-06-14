package pay

import (
	"encoding/json"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"maze_game_server/pb/common/MazePay"
	"maze_game_server/usecase/mustarrive"
	"net/http"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
)

func safeHttpRegister(logger fklog.FKLogI, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		defer fkutil.CaptureException()

		request.ParseForm()
		uid := fkutil.ToUint64(request.Form.Get("userId"))

		logger.WarnWF("execute http request", zap.Uint64("userId", uid),
			zap.String("uri", request.RequestURI), zap.Any("header", request.Header),
			zap.Any("host", request.Host), zap.Any("remoteAddr", request.RemoteAddr))
		handler(writer, request)
	})
}

func RegPayDelivery(logger fklog.FKLogI) {
	safeHttpRegister(logger, "/v1/pay/delivery", func(writer http.ResponseWriter, request *http.Request) {
		userId := fkutil.ToUint64(request.Form.Get("userId"))
		uniqueId := fkutil.ToUint64(request.Form.Get("uniqueId"))
		tradeNo := request.Form.Get("tradeNo")
		payChannel := request.Form.Get("payChannel")
		logger.SetLogId(time.Now().UnixNano())
		logger.SetUid(userId)
		logger.InfoWF("pay delivery begin")
		_ = uniqueId
		_ = tradeNo
		_ = payChannel
		writer.Header().Set("Content-Type", "application/json")

		{
			// todo 发货逻辑
			// 要优化：发货逻辑和订单状态修改不是事务的，所以存在极限情况多发货  例如：发货后，服务挂掉，支付服务器没收到发货回复认为没有发货成功，将进行发货重试
			// todo 最好还是检查下礼包id是否合法
			var err error
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(writer).Encode(map[string]any{
					"code":    http.StatusInternalServerError,
					"message": err.Error(),
				})
				logger.ErrorWF("pay delivery failed", zap.Error(err))
			}
		}

		logger.InfoWF("pay delivery success")
		_ = json.NewEncoder(writer).Encode(map[string]any{
			"code":    http.StatusOK,
			"message": "success",
		})
		PushPaySuccess(logger, int64(userId), tradeNo)
	})
}

// 推送发货成功， 推送失败也不处理，因为已经先发货成功了
func PushPaySuccess(logger fklog.FKLogI, userId int64, tradeNo string) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.ErrorWF("push pay success panic", zap.Any("err", err))
			}
		}()
		push := &MazePay.PushPaySuccessID{
			TradeNo: proto.String(tradeNo),
		}
		// 通知用户发货成功
		err := mustarrive.SendArrivePacket(logger, userId, 10509, push)
		if err != nil {
			logger.ErrorWF("PushPay error", zap.Error(err), zap.Int64("userId", userId), zap.String("tradeNo", tradeNo))
		} else {
			logger.InfoWF("PushPay success", zap.Int64("userId", userId), zap.String("tradeNo", tradeNo))
		}
	}()
}
