/*
 * @Author: majian
 * @Date: 2025-03-19 20:49:08
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-20 10:50:34
 */
package kafka_dispatch

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkutil"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/tcp/dispatchtcp"
)

type SexChangeInfo struct {
	UserId     string `json:"user_id"`
	Sex        int32  `json:"sex"`
	CreateChg  int32  `json:"create_chg"`
	SourceType int32  `json:"source_type"`
	FromType   int32  `json:"from_type"`
}

func HandleSexMsg(c context.Context, logger fklog.FKLogI, index int, key, data []byte) error {
	msg := &SexChangeInfo{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		logger.ErrorWF("HandleSexMsg Unmarshal fail", zap.Error(err), zap.String("msg", string(data)))
		return err
	}
	userId := fkutil.ToUint64(msg.UserId)
	if userId == 0 {
		logger.WarnWF("HandleSexMsg has invalid param", zap.String("msg", fmt.Sprintf("%+v", msg)))
		return nil
	}

	logger.InfoWF("HandleSexMsg receive", zap.Any("msg", msg))

	return dispatchtcp.DispatchKafkaMsgTcp(logger, uint64(userId), constdef.KafkaMDTSexDesc, data)
}
