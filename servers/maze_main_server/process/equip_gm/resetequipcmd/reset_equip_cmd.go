/*
 * @Author: majian
 * @Date: 2025-02-21 19:26:40
 * @Last Modified by: majian
 * @Last Modified time: 2025-02-21 19:57:25
 */
package resetequipcmd

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip_gm/equipbaggm"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip"
)

func RunCmd1001(logger fklog.FKLogI, userID uint64, session string, param string) error {
	e := equipbaggm.ClearUserBag(logger, userID)
	if e != nil {
		return e
	}
	e = equip.ChkEquipPosUnlock(logger, userID, "gm", true)
	if e != nil {
		return e
	}
	// 初始装备套检查
	e = equip.InitDollEquipSuitSeq(logger, userID)
	if e != nil {
		return e
	}
	// 处理初始化装备
	e = equip.HandleDollEquipInit(logger, userID, true)
	if e != nil {
		return e
	}
	return nil
}
