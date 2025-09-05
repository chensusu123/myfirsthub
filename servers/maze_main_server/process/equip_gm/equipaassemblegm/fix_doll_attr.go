/*
 * @Author: majian
 * @Date: 2024-07-25 14:12:42
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-15 18:04:35
 */
package equipaassemblegm

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/common/structsdef"
	"maze_game_server/io/redis/mazeattrcalcnotifyqueue"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/module/dollassembleinfo"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 1 修复非武力值buff 2=修复武力值buff
func ReCalcDollEquipAttr(ctx context.Context, userId uint64, fixType int32) error {
	logger := fklog.ContextAppLogger(ctx)
	_, effect, err := dollassembleinfo.GetDollAssembleInfoEx(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "ReCalcDollEquipAttr Get Assemble info fail", zap.Error(err))
		return err
	}
	if effect == nil {
		logger.CtxInfo(ctx, "ReCalcDollEquipAttr effect nil")
		return nil
	}
	_, otherAttrs := effect.ForceAttrs, effect.Other
	if fixType&1 > 0 {
		// 更新buff中心
		e := mazebuffinforedis.SaveMazeEquipBuff(ctx, userId, otherAttrs)
		if e != nil {
			return e
		}

		calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{}
		calcAttrNotify.FromServer = "maze_equip_gm_server"
		calcAttrNotify.UserId = userId
		calcAttrNotify.ChgType = constdef.MazeBuffChgTypeEquipGm

		calcAttrNotify.Session = ""
		mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(ctx, calcAttrNotify)
	}
	// if fixType&2 > 0 {
	// 	e := dollassembleattrredis.SetDollAssembleAttr(logger, userId, constdef.AttrFieldEquip, forceAttrs)
	// 	if e != nil {
	// 		return e
	// 	}

	// 	senddollforceattrrecord.SendDollForceAttrRecord(logger, userId, constdef.AttrFieldEquip, int32(constdef.ForceChgTypeGm), nil, forceAttrs)
	// 	dollforcechgnotifyqueue.DollForceChgNotice(logger, userId, constdef.ForcePreviewEquip, int32(constdef.ForceChgTypeGm), "")
	// }
	return nil
}

func CalcDollAttrCalc(ctx context.Context, userId uint64) error {
	calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{}
	calcAttrNotify.FromServer = "maze_equip_gm_server"
	calcAttrNotify.UserId = userId
	calcAttrNotify.ChgType = constdef.MazeBuffChgTypeGm
	calcAttrNotify.Session = ""
	return mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(ctx, calcAttrNotify)
}
