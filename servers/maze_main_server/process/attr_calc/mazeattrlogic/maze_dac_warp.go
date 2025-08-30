/*
 * @Author: majian
 * @Date: 2025-03-11 16:48:57
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-14 21:59:12
 */
package mazeattrlogic

import (
	"context"
	"fmt"

	"maze_game_server/common/constdef"
	"maze_game_server/common/structsdef"
)

func RunDacFromQue(ctx context.Context, msg *structsdef.MazeCalcAttrNotifyMsg) (needRetry bool, err error) {
	dacp := NewDACParam()
	dacp.ChgType = msg.ChgType
	dacp.ChgDesc = msg.ChgDesc
	dacp.Session = msg.Session
	return RunDac(ctx, msg.UserId, dacp)
}

func RunDacFromBC(ctx context.Context, userId uint64, subType int32) (needRetry bool, err error) {
	dacp := NewDACParam()
	dacp.ChgType = constdef.MazeBuffCenter
	dacp.ChgSubType = subType
	dacp.ChgDesc = fmt.Sprintf("oldBc_%d", subType)

	return RunDac(ctx, userId, dacp)
}

func RunDac(ctx context.Context, userId uint64, dacParam *DACParam) (needRetry bool, err error) {
	// logger := fklog.ContextAppLogger(ctx)
	dac := NewDAC(ctx, userId)
	defer func() {
		needRetry = dac.ErrNeedRetry()
	}()
	err = dac.InitData(ctx, dacParam)
	if err != nil {
		return
	}
	dac.Prepare()

	err = dac.Calc()
	if err != nil {
		return
	}
	err = dac.End(ctx)
	return
}
