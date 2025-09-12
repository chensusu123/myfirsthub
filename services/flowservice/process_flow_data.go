package flowservice

import (
	"encoding/json"
	"fmt"
	"time"

	"maze_game_server/io/kafka"
	flowmodel "maze_game_server/model/flowmodel/flow_model"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func (s *service) ProcessFlowData() {
	defer func() {
		if err := recover(); err != nil {
			fklog.AppLogger().ErrorWF("flowservice ProcessFlowData Restart")
			// 重启该携程 防止panic退出
			go s.ProcessFlowData()
		}
	}()

	for {
		select {
		case data := <-s.ch:
			ctx := data.Ctx
			logger := fklog.ContextAppLogger(data.Ctx)
			// 检测容量告警
			if len(s.ch) >= flowmodel.ERRORFLOWCHANSIZE {
				logger.CtxError(ctx, "flowservice ProcessFlowData Size Greater than ERRORFLOWCHANSIZE",
					zap.Any("ERRORFLOWCHANSIZE", flowmodel.ERRORFLOWCHANSIZE))
			} else if len(s.ch) >= flowmodel.WARNFLOWCHANSIZE {
				logger.CtxWarn(ctx, "flowservice ProcessFlowData Size Greater than WARNFLOWCHANSIZE",
					zap.Any("ERRORFLOWCHANSIZE", flowmodel.WARNFLOWCHANSIZE))
			}
			s._processFlowData(data)
		}
	}
}

func (s *service) _processFlowData(data *flowmodel.FlowData) error {
	// 打到kafka 中
	ctx, record := data.Ctx, data.Data
	span := trace.SpanFromContext(ctx)

	defer func() {
		s.queueLen.Add(-1)
		span.AddEvent("process_flow_data_end")
		span.End()
	}()

	span.AddEvent("Marshal")

	logger := fklog.ContextAppLogger(ctx)

	jsonData, err := json.Marshal(record)
	if err != nil {
		logger.CtxError(ctx, "flowservice ProcessFlowData Marshal Fail",
			zap.Any("record", record),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.AddEvent("SendMsg")
	err = kafka.GflowKafka.SendMsg(ctx, fmt.Sprintf("%v", time.Now().UnixNano()), jsonData)
	if err != nil {
		logger.CtxError(ctx, "flowservice ProcessFlowData SendMsg Fail",
			zap.Any("record", record),
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return nil
}
