package pay

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"maze_game_server/common/constdef"
	"maze_game_server/common/function/gentradeno"
	"maze_game_server/common/jwt"
	"maze_game_server/config/GMazeChargeV8Cfg"
	"maze_game_server/pb/common/MazePay"
	"maze_game_server/services/itemservice"
	"maze_game_server/usecase/online"
	"net/http"
	"strconv"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
)

type payDeliveryResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

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
		logger.SetLogId(time.Now().UnixNano())
		httpCode := http.StatusInternalServerError
		res := &payDeliveryResponse{
			Code:    http.StatusInternalServerError,
			Message: "",
		}
		deliveryToken := request.Form.Get("deliveryToken")
		defer func() {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(httpCode)
			err := json.NewEncoder(writer).Encode(res)
			if err != nil {
				logger.ErrorWF("PayDelivery JsonEncode fail", zap.Error(err), zap.Any("req", deliveryToken), zap.Any("res", res))
				return
			}
		}()

		// 验证jwt
		deliveryClaim, err := jwt.ValidateDeliveryJWT(deliveryToken)
		if err != nil {
			logger.ErrorWF("pay delivery ValidateDeliveryJWT failed", zap.Error(err), zap.String("deliveryToken", deliveryToken))
			return
		}
		{
			// 发货
			// 查找礼包id
			var chargeCfg *GMazeChargeV8Cfg.MazeChargeV8ConfigRow
			chargeCfgAll := GMazeChargeV8Cfg.GetAll()
			for _, i := range chargeCfgAll {
				if i.Unique_id == deliveryClaim.UniqueId {
					chargeCfg = i
					break
				}
			}
			if chargeCfg == nil {
				logger.ErrorWF("pay delivery unique id not exist", zap.Error(err), zap.Any("deliveryClaim", deliveryClaim))
				return
			}
			// 要优化：发货逻辑和订单状态修改不是事务的，所以存在极限情况多发货  例如：发货后，服务挂掉，支付服务器没收到发货回复认为没有发货成功，将进行发货重试
			// todo 充值表要调整可能，目前没法通过maze_charge_v8找到具体的道具id,就临时用rmb的数量了
			tradeNo := gentradeno.GetTradeNum()
			item := &itemservice.ItemInfo{
				ItemId: constdef.MazeCommonItemDiamond,
				Count:  int64(chargeCfg.Currency_num),
			}
			errInfo := itemservice.GlobalItemService.AddItem(context.TODO(), deliveryClaim.UserId, itemservice.ItemOpTypePay, tradeNo, item)
			if errInfo != nil {
				res.Code = http.StatusInternalServerError
				res.Message = err.Error()
				logger.ErrorWF("pay delivery add item failed", zap.Error(err), zap.Any("deliveryClaim", deliveryClaim))
				return
			}
		}

		httpCode = http.StatusOK
		res.Code = httpCode
		logger.InfoWF("pay delivery success", zap.Any("deliveryClaim", deliveryClaim), zap.Any("res", res))
		PushPaySuccess(logger, int64(deliveryClaim.UserId), strconv.FormatInt(deliveryClaim.TradeNo, 10))
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
		err := online.Push(logger, uint64(userId), 10509, push)
		if err != nil {
			logger.ErrorWF("PushPay error", zap.Error(err), zap.Int64("userId", userId), zap.String("tradeNo", tradeNo))
		} else {
			logger.InfoWF("PushPay success", zap.Int64("userId", userId), zap.String("tradeNo", tradeNo))
		}
	}()
}
