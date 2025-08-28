package gmservice

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/logidutil"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	headerKey     = "X-Trace-Id"
	USER_ID_FIELD = "user_id"
)

type HeaderStrKey string

func (s *service) SafeGETRegister(logger fklog.FKLogI, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	appConfig := appconfig.GlobalConfig()
	tracerPattern := http.MethodGet + pattern

	oiginPattern := pattern
	// /s4/AddExp
	pattern = "/s" + appConfig.Global.SectionID + pattern
	// 域名调用适配 带服前缀

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

	// 域名适配 带服前缀
	http.Handle(pattern, handlerFactor())
	// 后台适配 不带服前缀
	http.Handle(oiginPattern, handlerFactor())
}

func (s *service) SafePOSTRegister(logger fklog.FKLogI, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	appConfig := appconfig.GlobalConfig()
	tracerPattern := http.MethodPost + pattern

	oiginPattern := pattern
	// /s4/AddExp
	pattern = "/s" + appConfig.Global.SectionID + pattern
	// 域名调用适配 带服前缀

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

	// 域名适配 带服前缀
	http.Handle(pattern, handlerFactor())
	// 后台适配 不带服前缀
	http.Handle(oiginPattern, handlerFactor())
}
