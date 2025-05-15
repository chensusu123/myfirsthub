/*
 * @Author: majian
 * @Date: 2024-09-19 19:35:53
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-17 19:19:58
 */
package equipaassemblegm

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipCache"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/dollassemblesuitredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/effectequip"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/assemble"
)

// 修复装配信息里装备信息
func FixAssembleEquipInfo(logger fklog.FKLogI, userId uint64) (fix int32, err error) {
	// 获取到所有已装配的装备
	allEquips, err := dollassemblesuitredis.GetAllDollAssembleSuit(logger, userId)
	if err != nil {
		return 0, err
	}
	var needGuids []int64
	for _, equips := range allEquips {
		for _, equip := range equips {
			if equip.GetEquipGuid() > 0 {
				needGuids = append(needGuids, equip.GetEquipGuid())
			}
		}
	}
	// 获取到所需要的装备详情
	realEquipMap, err := effectequip.BatchGetEffectEquipInfo(logger, userId, needGuids...)
	if err != nil {
		return 0, err
	}

	// 检查装备的数据是否一致
	needFixEquipMap := make(map[string]*MazeEquipCache.MazeEquipPosDb)
	var needFixAttr bool
	for index, equips := range allEquips {
		for _, equip := range equips {
			if equip.GetEquipGuid() <= 0 {
				continue
			}
			equipInfo := realEquipMap[equip.GetEquipGuid()]
			var chg bool
			// 如果装备已经没有，需要清空装配信息
			if equipInfo == nil {
				equip.EquipGuid = nil
				equip.EquipId = nil
				//	equip.EquipSubType = nil
				chg = true
				needFixAttr = true
			} else {
				if equip.GetEquipId() != equipInfo.GetEquipId() {
					equip.EquipId = equipInfo.EquipId
					chg = true
				}
				// if equip.GetEquipSubType() != equipInfo.GetEquipSubType() {
				// 	equip.EquipSubType = equipInfo.EquipSubType
				// 	chg = true
				// }
			}
			if chg {
				needFixEquipMap[assemble.EnCodeAssembleEquipField(index, equip.GetPos())] = equip
			}
		}
	}
	if len(needFixEquipMap) > 0 {
		err = dollassemblesuitredis.BatchSaveDollAssembleSuit(logger, userId, needFixEquipMap)
		if needFixAttr {
			ReCalcDollEquipAttr(logger, userId, 3)
		}
		logger.InfoWF("FixAssembleEquipInfo fix ok", zap.Uint64("uid", userId),
			zap.Bool("needFixAttr", needFixAttr), zap.Any("fixEquips", needFixEquipMap))
	} else {
		logger.InfoWF("FixAssembleEquipInfo no need fix", zap.Uint64("uid", userId))
	}
	return int32(len(needFixEquipMap)), err
}
