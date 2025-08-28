/*
 * @Author: majian
 * @Date: 2024-07-10 11:36:40
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-14 19:39:03
 */
package attr_calc

import (
	"context"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"

	"maze_game_server/common/structsdef"
	"maze_game_server/servers/maze_main_server/process/attr_calc/mazeattrlogic"

	"go.uber.org/zap"
)

func OnMazeAttrCalcMsg(ctx context.Context, logger fklog.FKLogI, index int, obj interface{}) error {
	defer fkprometheus.DebugPMT("OnMazeAttrCalcMsg")()
	msg, ok := obj.(*structsdef.MazeCalcAttrNotifyMsg)
	if !ok {
		logger.ErrorWF("OnMazeAttrCalcMsg invalid msg")
		return nil
	}
	nLogger := logger.Clone("OnMazeAttrCalcMsg")
	nLogger.SetUid(msg.UserId)
	nLogger.SetLogId(time.Now().UnixNano())
	nLogger.WarnWF("OnMazeAttrCalcMsg pop", zap.Any("msg", msg))

	if msg.UserId <= 0 {
		return nil
	}
	var retry bool
	needRetry, e := mazeattrlogic.RunDacFromQue(ctx, msg)
	if e != nil && needRetry {
		retry = true
	}
	if retry {
		doMazeAttrCalcRetry(ctx, msg)
	}
	return nil
}
