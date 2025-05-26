/*
 * @Author: majian
 * @Date: 2025-03-11 16:48:57
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-14 21:59:12
 */
package mazeattrlogic

import (
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"maze_game_server/common/constdef"
	"maze_game_server/common/structsdef"
)

func RunDacFromQue(logger fklog.FKLogI, msg *structsdef.MazeCalcAttrNotifyMsg) (needRetry bool, err error) {
	dacp := NewDACParam()
	dacp.ChgType = msg.ChgType
	dacp.ChgDesc = msg.ChgDesc
	dacp.Session = msg.Session
	return RunDac(logger, msg.UserId, dacp)
}

func RunDacFromBC(logger fklog.FKLogI, userId uint64, subType int32) (needRetry bool, err error) {
	dacp := NewDACParam()
	dacp.ChgType = constdef.MazeBuffCenter
	dacp.ChgSubType = subType
	dacp.ChgDesc = fmt.Sprintf("oldBc_%d", subType)

	return RunDac(logger, userId, dacp)
}

func RunDac(logger fklog.FKLogI, userId uint64, dacParam *DACParam) (needRetry bool, err error) {
	dac := NewDAC(logger, userId)
	defer func() {
		needRetry = dac.ErrNeedRetry()
	}()
	err = dac.InitData(dacParam)
	if err != nil {
		return
	}
	dac.Prepare()

	err = dac.Calc()
	if err != nil {
		return
	}
	err = dac.End()
	return
}
