/*
 * @Author: majian
 * @Date: 2024-08-19 17:24:05
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-22 18:04:49
 */
package packequipostopb

import (
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/assemble"
	"maze_game_server/common/function/excelutil"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/common/function/packtopb"
	"maze_game_server/config/GMazeEquipPosLvV8Cfg"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/server/MazeEquipCache"
)

// 打包装备位信息
func PackEquipPosPb(logger fklog.FKLogI, asEquip *MazeEquipCache.MazeEquipPosInfo, mask int32) (*MazeGameEquip.EquipPosInfo, error) {
	equipPos := &MazeGameEquip.EquipPosInfo{}
	equipPos.EquipPos = proto.Int32(asEquip.GetEquipPos().GetPos())
	if mask == -1 {
		mask = int32(MazeGameEquip.ENUM_MAZE_EQUIP_POS_MASK_UNLOCK_INFO) |
			int32(MazeGameEquip.ENUM_MAZE_EQUIP_POS_MASK_STRENGTHEN_INFO) |
			int32(MazeGameEquip.ENUM_MAZE_EQUIP_POS_MASK_LOAD_EQUIP_INFO)
	}
	equipPos.Mask = proto.Int32(mask)
	if mask&int32(MazeGameEquip.ENUM_MAZE_EQUIP_POS_MASK_UNLOCK_INFO) > 0 {
		equipPos.IsUnlock = proto.Int32(1)
	}

	if mask&int32(MazeGameEquip.ENUM_MAZE_EQUIP_POS_MASK_LOAD_EQUIP_INFO) > 0 {
		// 打包装备
		if asEquip.EquipInfo != nil {
			var e error
			equipPos.EquipInfo, e = packtopb.EquipInfoToCliPB(logger, asEquip.EquipInfo)
			if e != nil {
				logger.ErrorWF("PackEquipPosPb EquipInfoToCliPB fail", zap.Error(e), zap.Any("equip", asEquip.EquipInfo))
				return nil, e
			}
			// 设置武力值,人偶不需要展示武力或者对比武力
			//	equipPos.EquipInfo.ForceValue = proto.Int64(asEquip.GetForce())
			equipPos.EquipInfo.ForceValue = proto.Int64(asEquip.GetForce())
			// if assemble.IsFiveElemActivate(asEquip.GetEquipLoadInfo().GetActivateMask()) {
			// 	// fe := equipPos.GetEquipInfo().GetFiveElemInfo()
			// 	// if fe != nil {
			// 	// 	fe.IsActive = proto.Bool(true)
			// 	// }
			// }
			// 检查装备套装是否激活
			if assemble.IsHurtSuitActivate(asEquip.GetEquipLoadInfo().GetActivateMask()) {
				if equipPos.GetEquipInfo().GetSuitInfo() != nil {
					equipPos.GetEquipInfo().GetSuitInfo().IsActive = proto.Int32(1)
				}
			}
		}
	}

	if mask&int32(MazeGameEquip.ENUM_MAZE_EQUIP_POS_MASK_STRENGTHEN_INFO) > 0 {
		// 打包装备位强化信息
		stInfo := &MazeGameEquip.EquipPosStrengthenInfo{}
		stInfo.Level = proto.Int32(asEquip.GetEquipPos().GetLevel())
		key := excelutil.GetEquipPosEnLevelKey(asEquip.GetEquipPos().GetPos(), asEquip.GetEquipPos().GetLevel())

		lpCfg := GMazeEquipPosLvV8Cfg.GetMazeEquipPosLvV8Config(key)
		if lpCfg == nil {
			e := errors.New(fmt.Sprintf("no found pos pvp level cfg:%d", key))
			logger.ErrorWF("PackEquipPosPb GetMazeEquipPosLvV8Config fail", zap.Error(e),
				zap.Any("equip", asEquip), zap.Int32("mask", mask))

			return nil, errors.New(fmt.Sprintf("no found pos pvp level cfg:%d", key))
		}
		if lpCfg.Next_order > 0 {
			stInfo.Cond = &MazeGameEquip.EquipPosStCond{}
			stInfo.Cond.NeedOtherLv = proto.Int32(lpCfg.Need_other)
			stInfo.Cond.CostItems = itemutil.Map2Common(lpCfg.Cost)
			stInfo.Cond.NeedMazeLv = proto.Int32(lpCfg.Need_maze_level)
		}
		equipPos.StInfo = stInfo
	}

	return equipPos, nil
}
