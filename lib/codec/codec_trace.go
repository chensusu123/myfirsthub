package codec

import (
	"context"

	"maze_game_server/lib/codec/raw_pkg"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func receiveSpan() (context.Context, trace.Span) {
	tracer := otel.Tracer("nano.codec")
	spanKind := trace.WithSpanKind(trace.SpanKindClient)
	ctx, span := tracer.Start(context.Background(), "recive_pack", spanKind)
	span.AddEvent("nano.recive.pack.init")
	return ctx, span
}

type BaseHeadHeader struct {
	pack *raw_pkg.StruSvrEsRawBaseHead
}

func NewBaseHeader(pack *raw_pkg.StruSvrEsRawBaseHead) BaseHeadHeader {
	return BaseHeadHeader{pack: pack}
}

// Set sets the header entries associated with key to the single
// element value. It is case-sensitive and replaces any existing
// values associated with key.
func (h BaseHeadHeader) Set(key, value string) {
	if h.pack == nil {
		return
	}
	h.pack.Header[key] = value
}

// Get gets the first value associated with the given key.
// It is case-sensitive.
func (h BaseHeadHeader) Get(key string) string {
	if h.pack == nil {
		return ""
	}
	value, ok := h.pack.Header[key]
	if !ok {
		return ""
	}
	return value
}

// Values returns all values associated with the given key.
// It is case-sensitive.
func (h BaseHeadHeader) Keys() []string {
	l := make([]string, 0)
	if h.pack == nil {
		return l
	}
	for k := range h.pack.Header {
		l = append(l, k)
	}
	return l
}
