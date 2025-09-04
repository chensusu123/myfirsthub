package cluster

import (
	"context"

	packCodec "maze_game_server/lib/codec"
	"maze_game_server/lib/codec/raw_pkg"
	"maze_game_server/lib/nano/session"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func agentSendSpan(ctx context.Context, sendType string, mid uint64, uid int64, agentSession int64) (context.Context, trace.Span) {
	if ctx == nil {
		ctx = context.Background()
	}
	tracer := otel.Tracer("agent.send.message")
	sessionID, rqTime, rsID := packCodec.SplitSessionAndPackType(mid)
	ctx, span := tracer.Start(ctx, "agent.send."+sendType)
	span.SetAttributes(
		attribute.Int64("packet.session", int64(sessionID)),
		attribute.Int64("packet.id", int64(rsID)),
		attribute.Int64("rq.time", int64(rqTime)),
		attribute.Int64("enduser.id", uid),
		attribute.Int64("agent.session", agentSession),
	)
	span.AddEvent("agent.send.init")
	return ctx, span
}

func receiveSpan(ctx context.Context, msgLen int) (context.Context, trace.Span) {
	if msgLen <= 0 {
		span := trace.SpanFromContext(ctx)
		return ctx, span
	}
	tracer := otel.Tracer("nano.recive")
	ctx, span := tracer.Start(ctx, "recive_data")
	span.AddEvent("nano.receive.span.init")
	return ctx, span
}

func closeHandleSpan(ctx context.Context, a *agent, closeName string) (context.Context, trace.Span) {
	tracer := otel.Tracer("nano.close")
	ctx, span := tracer.Start(ctx, "nano.close.handle")
	span.AddEvent("nano.close.span.init")
	span.SetAttributes(attribute.String("remote_addr", a.session.RemoteAddr().String()),
		attribute.Int64("agent.session", a.session.ID()),
		attribute.Int64("enduser.id", a.session.UID()),
		attribute.String("nano.close.name", closeName),
	)
	return ctx, span
}

func callSpan(ctx context.Context, s *session.Session, callName string) (context.Context, trace.Span) {
	tracer := otel.Tracer("nano.call")
	ctx, span := tracer.Start(ctx, "nano.call.handle")
	span.AddEvent("nano.call.span.init")
	span.SetAttributes(
		attribute.Int64("agent.session", s.ID()),
		attribute.Int64("enduser.id", s.UID()),
		attribute.String("nano.callName", callName),
	)
	return ctx, span
}

func packSpan(ctx context.Context, agentSession int64, userID int64, pack *raw_pkg.StruSvrEsRawBaseHead) (context.Context, trace.Span) {
	ctx = otel.GetTextMapPropagator().Extract(ctx, packCodec.NewBaseHeader(pack))
	// 创建自定义Tracer
	tracer := otel.Tracer("nano-net")
	// 消息处理函数中手动创建Span
	spanKind := trace.WithSpanKind(trace.SpanKindServer)

	ctx, span := tracer.Start(ctx, "nano.pack.process",
		spanKind)
	span.SetAttributes(
		attribute.Int64("agent.session", agentSession),
		attribute.Int64("enduser.id", userID),
		attribute.Int64("packet.session", int64(pack.SessionID)),
		attribute.Int64("packet.id", int64(pack.PackType)),
	)
	// span.SetAttributes(attribute.String("nats.subject", subject))
	span.AddEvent("begin")
	return ctx, span
}
