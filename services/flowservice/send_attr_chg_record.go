package flowservice

import (
	"context"

	flowmodel "maze_game_server/model/flowmodel/flow_model"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func (s *service) SendFlowData(ctx context.Context, record interface{}) {
	ctx, span := sendFlowDataSpan(ctx)
	logger := fklog.ContextAppLogger(ctx)
	defer func() {
		logger.CtxInfo(ctx, "flowservice SendFlowData End",
			zap.Any("data", record))
	}()
	logger.CtxInfo(ctx, "flowservice SendFlowData Start",
		zap.Any("data", record),
	)
	data := flowmodel.NewFlowData(ctx, record)
	c := s.queueLen.Add(1)
	span.SetAttributes(attribute.Int64("flowdata.queue.len", c))
	span.AddEvent("push.channel")
	s.ch <- data
}

func sendFlowDataSpan(ctx context.Context) (context.Context, trace.Span) {
	tracer := otel.Tracer("flowservice")
	ctx, span := tracer.Start(ctx, "SendFlowData")
	span.AddEvent("SendFlowData.init")
	return ctx, span
}
