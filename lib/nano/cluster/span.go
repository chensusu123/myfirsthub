package cluster

import (
	"context"

	packCodec "maze_game_server/lib/codec"

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
