package flowservice

import (
	"encoding/json"
	"fmt"
	"maze_game_server/io/kafka"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) ProcessFlowData() {
	for {
		select {
		case data := <-s.ch:
			// 打到kafka 中
			ctx, record := data.Ctx, data.Data
			logger := fklog.ContextAppLogger(ctx)

			jsonData, err := json.Marshal(record)
			if err != nil {
				logger.CtxError(ctx, "flowservice ProcessFlowData Marshal Fail",
					zap.Any("record", record),
					zap.Error(err),
				)
				return
			}
			err = kafka.GflowKafka.SendMsg(ctx, fmt.Sprintf("%v", time.Now().UnixNano()), jsonData)
			if err != nil {
				logger.CtxError(ctx, "flowservice ProcessFlowData SendMsg Fail",
					zap.Any("record", record),
					zap.Error(err),
				)
				return
			}
		}
	}
}
