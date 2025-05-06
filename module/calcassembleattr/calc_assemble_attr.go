package calcassembleattr

import (
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/protodef/MazeBuffData"
	"gitlab.ifreetalk.com/plate/protodef/MazeEquipCache"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/assemble"
)

// 计算装备属性加成
func CalcEquipAttrs(logger fklog.FKLogI, equips []*MazeEquipCache.MazeEquipPosInfo) (attrs map[int32]int64, other *MazeBuffData.MazeBuffDb, err error) {
	effect, e := CalcEquipEffectAll(logger, equips, EffectCalcInParam{IsLog: false, IsForce: true})
	if e != nil {
		return nil, nil, e
	}
	return effect.ForceAttrs, effect.Other, nil
}

// 计算不包含某件装备的属性加成
func CalcEquipAttrsWithoutOne(logger fklog.FKLogI, equips []*MazeEquipCache.MazeEquipPosInfo, pos int32) (attrs map[int32]int64, err error) {
	var tmpEquips []*MazeEquipCache.MazeEquipPosInfo
	for _, equipPos := range equips {
		lPos := equipPos.GetEquipPos().GetPos()
		if lPos == 0 {
			continue
		}
		if pos == lPos {
			continue
		}
		if !assemble.IsAssembleEquip(equipPos) {
			continue
		}
		tmpEquips = append(tmpEquips, equipPos)
	}
	attrs, _, err = CalcEquipAttrs(logger, tmpEquips)
	return
}

// 替换某件装备计算属性加成
func CalcEquipAttrsReplaceOne(logger fklog.FKLogI, equips []*MazeEquipCache.MazeEquipPosInfo, pos int32, replaceEquip *MazeEquipCache.MazeEquipInfoDb) (attrs map[int32]int64, err error) {
	var tmpEquips []*MazeEquipCache.MazeEquipPosInfo
	var unLock bool
	for _, equipPos := range equips {
		lpos := equipPos.GetEquipPos().GetPos()
		if lpos == 0 {
			continue
		}

		if pos == lpos {
			equipTmp := &MazeEquipCache.MazeEquipPosInfo{}
			equipTmp.EquipPos = equipPos.EquipPos
			equipTmp.EquipInfo = replaceEquip
			equipTmp.EquipLoadInfo = &MazeEquipCache.MazeEquipPosDb{}
			equipTmp.EquipLoadInfo.EquipId = proto.Int32(replaceEquip.GetEquipId())
			equipTmp.EquipLoadInfo.EquipGuid = proto.Int64(replaceEquip.GetEquipGuid())
			equipTmp.EquipLoadInfo.Pos = proto.Int32(lpos)
			tmpEquips = append(tmpEquips, equipTmp)
			unLock = true
		} else {
			if !assemble.IsAssembleEquip(equipPos) {
				continue
			}
			tmpEquips = append(tmpEquips, equipPos)
		}
	}
	if !unLock {
		equipTmp := &MazeEquipCache.MazeEquipPosInfo{}
		equipTmp.EquipPos = &MazeEquipCache.MazeEquipSlotDb{Pos: proto.Int32(pos)}
		equipTmp.EquipInfo = replaceEquip
		equipTmp.EquipLoadInfo = &MazeEquipCache.MazeEquipPosDb{}
		equipTmp.EquipLoadInfo.EquipId = proto.Int32(replaceEquip.GetEquipId())
		equipTmp.EquipLoadInfo.EquipGuid = proto.Int64(replaceEquip.GetEquipGuid())
		equipTmp.EquipLoadInfo.Pos = proto.Int32(pos)
		tmpEquips = append(tmpEquips, equipTmp)
	}
	attrs, _, err = CalcEquipAttrs(logger, tmpEquips)
	return
}

func CalcEquipEffect(logger fklog.FKLogI, equips []*MazeEquipCache.MazeEquipPosInfo) (effect *EquipmentEffectInfo, err error) {
	return CalcEquipEffectAll(logger, equips, EffectCalcInParam{IsLog: true, IsForce: false})
}
