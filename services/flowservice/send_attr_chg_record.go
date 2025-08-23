package flowservice

import (
	"context"
	flowmodel "maze_game_server/model/flowmodel/flow_model"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) SendFlowData(ctx context.Context, record interface{}) {
	logger := fklog.ContextAppLogger(ctx)
	defer func() {
		logger.CtxInfo(ctx, "flowservice SendFlowData End",
			zap.Any("data", record))
	}()
	logger.CtxInfo(ctx, "flowservice SendFlowData Start",
		zap.Any("data", record),
	)
	data := flowmodel.NewFlowData(ctx, record)
	s.ch <- data
}
