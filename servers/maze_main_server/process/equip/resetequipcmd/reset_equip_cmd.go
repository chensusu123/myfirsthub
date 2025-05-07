/*
 * @Author: majian
 * @Date: 2025-02-21 19:26:40
 * @Last Modified by: majian
 * @Last Modified time: 2025-02-21 19:57:25
 */
package resetequipcmd

import (
	"gitlab.ifreetalk.com/maze/maze_equip_server/servers/maze_equip_gm_server/process/equipbaggm"
	"gitlab.ifreetalk.com/maze/maze_equip_server/servers/maze_equip_main_server/process"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
)

func RunCmd1001(logger fklog.FKLogI, userID uint64, session string, param string) error {
	e := equipbaggm.ClearUserBag(logger, userID)
	if e != nil {
		return e
	}
	e = process.ChkEquipPosUnlock(logger, userID, "gm", true)
	if e != nil {
		return e
	}
	// 初始装备套检查
	e = process.InitDollEquipSuitSeq(logger, userID)
	if e != nil {
		return e
	}
	// 处理初始化装备
	e = process.HandleDollEquipInit(logger, userID, true)
	if e != nil {
		return e
	}
	return nil
}
