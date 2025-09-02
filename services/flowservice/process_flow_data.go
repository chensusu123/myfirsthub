package flowservice

import (
	"context"
	"encoding/json"
	"fmt"
	"maze_game_server/io/kafka"
	flowmodel "maze_game_server/model/flowmodel/flow_model"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) ProcessFlowData() {
	defer func() {
		if err := recover(); err != nil {
			logger.CtxErrorf(context.TODO(), "flowservice ProcessFlowData Restart")
			// 重启该携程 防止panic退出
			go s.ProcessFlowData()
		}
	}()

	for {
		select {
		case data := <-s.ch:
			// 检测容量告警
			if len(s.ch) >= flowmodel.ERRORFLOWCHANSIZE {
				logger.CtxErrorf(context.TODO(), "flowservice ProcessFlowData Size Greater than ERRORFLOWCHANSIZE",
					zap.Any("ERRORFLOWCHANSIZE", flowmodel.ERRORFLOWCHANSIZE))
			} else if len(s.ch) >= flowmodel.WARNFLOWCHANSIZE {
				logger.CtxWarnf(context.TODO(), "flowservice ProcessFlowData Size Greater than WARNFLOWCHANSIZE",
					zap.Any("ERRORFLOWCHANSIZE", flowmodel.WARNFLOWCHANSIZE))
			}

			// 打到kafka 中
			ctx, record := data.Ctx, data.Data
			logger := fklog.ContextAppLogger(ctx)

			jsonData, err := json.Marshal(record)
			if err != nil {
				logger.CtxError(ctx, "flowservice ProcessFlowData Marshal Fail",
					zap.Any("record", record),
					zap.Error(err),
				)
				continue
			}

			err = kafka.GflowKafka.SendMsg(ctx, fmt.Sprintf("%v", time.Now().UnixNano()), jsonData)
			if err != nil {
				logger.CtxError(ctx, "flowservice ProcessFlowData SendMsg Fail",
					zap.Any("record", record),
					zap.Error(err),
				)
			}
		}
	}
}
