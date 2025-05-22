/*
 * @Author: majian
 * @Date: 2025-02-21 19:26:40
 * @Last Modified by: majian
 * @Last Modified time: 2025-02-21 19:57:25
 */
package resetequipcmd

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebuffinforedis"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip_gm/equipbaggm"
	"go.uber.org/zap"
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
	// 删除临时buff武力属性
	err := mazebuffinforedis.DelMazeBuffBySrc(logger, userID, constdef.MazeBuffSrcSelectBuffForce)
	if err != nil {
		logger.ErrorWF("MazeBarrierNotifyProcess DelMazeBuffBySrc failed", zap.Uint64("userId", userID), zap.Error(err))
		return err
	}
	return nil
}
