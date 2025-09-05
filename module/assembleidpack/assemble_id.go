/*
 * @Author: majian
 * @Date: 2024-04-24 21:53:14
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-21 20:08:14
 */
package assembleidpack

import (
	"context"

	"maze_game_server/common/function/packtopb/asequipsuittopb"
	"maze_game_server/common/function/packtopb/packequipostopb"
	"maze_game_server/module/equippossuit"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/server/MazeEquipCache"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func SendAssembleChgID(ctx context.Context, userId uint64, assembleInfo *MazeEquipCache.MazeAssembleDb, wantMask, posMask, reason int32) error {
	logger := fklog.ContextAppLogger(ctx)
	var mask int32
	idp := &MazeGameEquip.MazeAssembleChgID{}
	idp.Reason = proto.Int32(reason)
	idp.MazeAssembleInfo = &MazeGameEquip.MazeAssembleInfo{}
	// 设置装备位信息
	if wantMask&int32(MazeGameEquip.ENUM_MAZE_ASSEMBLE_CHG_TYPE_MASK_EQUIP_POS_MASK) > 0 {
		for _, aEquip := range assembleInfo.MazeEquips {
			var e error
			posPb := &MazeGameEquip.EquipPosInfo{}
			posPb, e = packequipostopb.PackEquipPosPb(ctx, aEquip, posMask)

			if e != nil {
				logger.CtxError(ctx, "SendAssembleChgID PackEquipPosPb fail", zap.Error(e), zap.Any("assembleInfo", assembleInfo))
				return e
			}
			idp.MazeAssembleInfo.EquipPosList = append(idp.MazeAssembleInfo.EquipPosList, posPb)
		}
		if len(idp.MazeAssembleInfo.EquipPosList) > 0 {
			mask |= int32(MazeGameEquip.ENUM_MAZE_ASSEMBLE_CHG_TYPE_MASK_EQUIP_POS_MASK)
		}
	}
	// 设置装备位强化套装
	if wantMask&int32(MazeGameEquip.ENUM_MAZE_ASSEMBLE_CHG_TYPE_MASK_EQUIP_ST_SUIT_MASK) > 0 {
		if assembleInfo.EpEnSuitId != nil {
			// 打包装备强化套装信息
			var e error
			idp.MazeAssembleInfo.EquipPosStSuit,
				idp.MazeAssembleInfo.EquipPosNextStSuit, e = equippossuit.GetCurAndNextSuit(ctx, assembleInfo.GetEpEnSuitId())
			if e != nil {
				return e
			}
			mask |= int32(MazeGameEquip.ENUM_MAZE_ASSEMBLE_CHG_TYPE_MASK_EQUIP_ST_SUIT_MASK)
		}
	}
	// // 设置pk段位
	// if wantMask&int32(DollEquip.ENUM_ASSEMBLE_CHG_TYPE_MASK_ENUM_DOLL_PK_LEVEL_MASK) > 0 {
	// 	if assembleInfo.PkLevel != nil {
	// 		idp.DollAssembleInfo.PkLevel = proto.Int32(assembleInfo.GetPkLevel())
	// 		mask |= int32(DollEquip.ENUM_ASSEMBLE_CHG_TYPE_MASK_ENUM_DOLL_PK_LEVEL_MASK)
	// 	}
	// }
	// 设置装备套装

	if wantMask&int32(MazeGameEquip.ENUM_MAZE_ASSEMBLE_CHG_TYPE_MASK_EQUIP_SUIT_MASK) > 0 {
		if assembleInfo.EpSuitId != nil {
			idp.MazeAssembleInfo.AsEquipSuitInfo = asequipsuittopb.PackAsEquipSuitInfo(assembleInfo.GetEpSuitId())
			mask |= int32(MazeGameEquip.ENUM_MAZE_ASSEMBLE_CHG_TYPE_MASK_EQUIP_SUIT_MASK)
		}
	}
	if mask == 0 {
		logger.CtxWarn(ctx, "SendAssembleChgID no assemble data", zap.Any("assembleInfo", assembleInfo))
		return nil
	}
	idp.Mask = proto.Int32(mask)
	idp.Token = proto.Int64(GetAssembleToken())
	err := online.ClusterPush(context.TODO(), uint64(userId), 10422, idp)
	if err != nil {
		logger.CtxError(ctx, "SendAssembleChgID SendArrivePacket err", zap.Error(err))
		return err
	}
	logger.CtxInfo(ctx, "SendAssembleChgID SendArrivePacket succ", zap.Any("pack", idp))
	return nil
}
