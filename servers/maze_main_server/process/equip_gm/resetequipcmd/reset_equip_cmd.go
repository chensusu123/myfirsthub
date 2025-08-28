/*
 * @Author: majian
 * @Date: 2025-02-21 19:26:40
 * @Last Modified by: majian
 * @Last Modified time: 2025-02-21 19:57:25
 */
package resetequipcmd

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/servers/maze_main_server/process/equip"
	"maze_game_server/servers/maze_main_server/process/equip_gm/equipbaggm"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func RunCmd1001(ctx context.Context, userID uint64, session string, param string) error {
	logger := fklog.ContextAppLogger(ctx)
	e := equipbaggm.ClearUserBag(ctx, userID)
	if e != nil {
		return e
	}
	e = equip.ChkEquipPosUnlock(ctx, userID, "gm", true)
	if e != nil {
		return e
	}
	// 初始装备套检查
	e = equip.InitDollEquipSuitSeq(logger, userID)
	if e != nil {
		return e
	}
	// 处理初始化装备
	e = equip.HandleDollEquipInit(ctx, userID, true)
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
