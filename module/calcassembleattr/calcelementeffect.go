/*
 * @Author: majian
 * @Date: 2024-09-13 11:46:26
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-17 15:14:43
 */
package calcassembleattr

import (
	"google.golang.org/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/config/GMazeEquipSuiteInfoV8Cfg"
	"maze_game_server/excel/mazeequipconfigv8"
	"maze_game_server/pb/server/MazeBuffData"
	"maze_game_server/pb/server/MazeEquipCache"
)

// 10410效果规则
// 情况1:武器(非套装) 伤害类型与套装伤害类型一致
// 情况2:武器(非套装) 伤害类型与套装伤害类型不一致
// 情况3:武器(套装) 伤害类型与套装伤害类型不一致  第三种不会出现 和第四种一样
// 情况4:武器(套装) 伤害类型与套装伤害类型一致
// 属性值:情况类型*100+伤害类型枚举(1234) 例如: 情况4 + 火属性  402
// 情况4
func CalcElementEffect(logger fklog.FKLogI, equip *MazeEquipCache.MazeEquipPosInfo, suitMap map[int32]int32, isLog bool) (a1 *MazeBuffData.MazeBuffAttr) {
	wpSuitId := equip.GetEquipInfo().GetSuitId()
	var attrVal int32 // 存放结果  情况(n)*1000+伤害类型枚举(x:1 2 3 4)
	var elemId, effectType int32
	a1 = &MazeBuffData.MazeBuffAttr{}
	defer func() {
		if isLog {
			logger.InfoWF("CalcElementEffect end",
				zap.Int32("equipId", equip.GetEquipInfo().GetEquipId()),
				zap.Int64("guid", equip.GetEquipInfo().GetEquipGuid()),
				zap.Int32("suitId", wpSuitId),
				zap.Int32("elemId", elemId),
				zap.Any("a1", a1))
		}
	}()
	if wpSuitId > 0 {
		attrVal = 4
		var tmpSuitId int32
		elemId, tmpSuitId = findElemAttr(equip.GetEquipInfo())
		if tmpSuitId != wpSuitId {
			logger.ErrorWF("CalcElementEffect config not match",
				zap.Int32("equipId", equip.GetEquipInfo().GetEquipId()),
				zap.Int32("guid", int32(equip.GetEquipInfo().GetEquipGuid())),
				zap.Int32("suitId", wpSuitId), zap.Int32("tmpSuitId", tmpSuitId))
			return nil
		}
	} else {
		attrVal = 1
		// 找到跟武器伤害类型一致的套装
		elemId, wpSuitId = findElemAttr(equip.GetEquipInfo())
	}
	if wpSuitId == 0 || elemId == 0 {
		return nil
	}
	// 检查套装数量是否满足
	if !mazeequipconfigv8.IsShowSuitEffect(suitMap[wpSuitId]) {
		return nil
	}
	effectType = mazeequipconfigv8.GetSuitEffectType(elemId)
	if effectType > 0 {
		a1.AttrVal = proto.Int64(int64(effectType))
		if attrVal == 4 {
			a1.AttrId = proto.Int32(10410)
		} else if attrVal == 1 {
			a1.AttrId = proto.Int32(10411)
		}
		// logger.InfoWF("CalcElementEffect succ", zap.Any("effect", effectInfo.GetAttrVal()))
		return a1
	}
	return nil
}

// 找到武器的元素属性和对应套装Id
func findElemAttr(equip *MazeEquipCache.MazeEquipInfoDb) (elemId, suitId int32) {
	allRow := GMazeEquipSuiteInfoV8Cfg.GetAll()
	elemMap := make(map[int32]*GMazeEquipSuiteInfoV8Cfg.MazeEquipSuiteInfoV8ConfigRow)
	for _, row := range allRow {
		elemMap[row.Counter_weapons_damage_attr] = row
	}
	for _, entry := range equip.GetBaseAttrs() {
		for _, attr := range entry.GetRealAttrList() {
			if suitRow, ok := elemMap[attr.GetAttrId()]; ok {
				if suitRow != nil {
					return attr.GetAttrId(), suitRow.Suite_id
				}
			}
		}
	}
	return 0, 0
}
