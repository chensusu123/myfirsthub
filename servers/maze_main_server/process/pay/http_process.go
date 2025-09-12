package pay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"maze_game_server/common/constdef"
	"maze_game_server/common/jwt"
	"maze_game_server/common/tradeno"
	"maze_game_server/config/GMazeChargeV8Cfg"
	"maze_game_server/pb/common/MazePay"
	"maze_game_server/services/itemservice"
	"maze_game_server/usecase/online"

	"github.com/google/uuid"
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/logidutil"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
)

const (
	headerKey     = "X-Trace-Id"
	USER_ID_FIELD = "user_id"
)

type HeaderStrKey string

type payDeliveryResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SafeGETRegister(logger fklog.FKLogI, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	appConfig := appconfig.GlobalConfig()
	tracerPattern := http.MethodGet + pattern

	handlerFactor := func() http.Handler {
		return otelhttp.NewHandler(
			http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				defer fkutil.CaptureException()

				uid := fkutil.ToUint64(request.Form.Get(USER_ID_FIELD))

				logidInt := logidutil.GenerateLogID()
				logID := strconv.Itoa(int(logidInt))
				shardingID := fmt.Sprintf("%d", appConfig.Global.ShardingID)
				writer.Header().Set("X-App-Namespace", appConfig.Global.Namespace)
				writer.Header().Set("X-App-Section", appConfig.Global.SectionID)
				writer.Header().Set("X-App-Name", appConfig.Server.Server)
				writer.Header().Set("X-App-Sharding", shardingID)
				writer.Header().Set("X-Log-ID", logID)

				callLogger := fklog.AppLogger().Clone("")
				callLogger.SetLogId(logidInt)
				callLogger.SetUid(uid)

				ctx, span := otel.Tracer("gm-handler").Start(request.Context(), tracerPattern)
				defer span.End()

				rid := request.Header.Get(headerKey)
				if rid == "" {
					if span.SpanContext().TraceID().IsValid() {
						rid = span.SpanContext().TraceID().String()
					}
					if rid == "" {
						rid = uuid.New().String()
					}
				}
				writer.Header().Set(headerKey, rid)
				// 将callLogger添加添加到ctx中
				ctx = fklog.ContextWithLogger(ctx, callLogger)
				ctx = context.WithValue(ctx, HeaderStrKey(headerKey), rid)

				request.ParseForm()

				span.SetAttributes(attribute.Int64("enduser.id", int64(uid)))
				//if !CheckGM.CheckGMOnline(context.TODO(), logger, uid, pattern, request.RemoteAddr) {
				//	return
				//}

				callLogger.CtxWarn(ctx, "execute gm", zap.Uint64("userId", uid),
					zap.String("pattern", pattern), zap.Any("header", request.Header),
					zap.Any("host", request.Host), zap.Any("remoteAddr", request.RemoteAddr))
				handler(writer, request.WithContext(ctx))
			}), tracerPattern)
	}
	http.Handle(pattern, handlerFactor())
}

// func safeHttpRegister(logger fklog.FKLogI, pattern string, handler func(http.ResponseWriter, *http.Request)) {
// 	http.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
// 		defer fkutil.CaptureException()

// 		request.ParseForm()
// 		uid := fkutil.ToUint64(request.Form.Get("userId"))

// 		logger.WarnWF("execute http request", zap.Uint64("userId", uid),
// 			zap.String("uri", request.RequestURI), zap.Any("header", request.Header),
// 			zap.Any("host", request.Host), zap.Any("remoteAddr", request.RemoteAddr))
// 		handler(writer, request)
// 	})
// }

func RegPayDelivery(loggerx fklog.FKLogI) {
	SafeGETRegister(loggerx, "/v1/pay/delivery", func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		logger := fklog.ContextAppLogger(ctx)
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
			tradeNo := tradeno.GetTradeNum()
			item := &itemservice.ItemInfo{
				ItemId: constdef.MazeCommonItemDiamond,
				Count:  int64(chargeCfg.Currency_num),
			}
			errInfo := itemservice.GlobalItemService.AddItem(ctx, deliveryClaim.UserId, itemservice.ItemOpTypePay, tradeNo, item)
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
		PushPaySuccess(ctx, logger, int64(deliveryClaim.UserId), strconv.FormatInt(deliveryClaim.TradeNo, 10))
	})
}

// 推送发货成功， 推送失败也不处理，因为已经先发货成功了
func PushPaySuccess(ctx context.Context, logger fklog.FKLogI, userId int64, tradeNo string) {
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
		err := online.ClusterPush(ctx, uint64(userId), 10509, push)
		if err != nil {
			logger.ErrorWF("PushPay error", zap.Error(err), zap.Int64("userId", userId), zap.String("tradeNo", tradeNo))
		} else {
			logger.InfoWF("PushPay success", zap.Int64("userId", userId), zap.String("tradeNo", tradeNo))
		}
	}()
}
