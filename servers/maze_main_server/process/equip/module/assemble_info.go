package module

import (
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"maze_game_server/common/function/assemble"
	"maze_game_server/pb/server/MazeEquipCache"
)

func GetEquipPosInfo(assembleInfo *MazeEquipCache.MazeAssembleDb, pos int32) *MazeEquipCache.MazeEquipPosInfo {
	for _, equipPos := range assembleInfo.MazeEquips {
		if equipPos.GetEquipPos().GetPos() == pos {
			return equipPos
		}
	}
	return nil
}

func FindEquipPosInfo(assembleInfo *MazeEquipCache.MazeAssembleDb, guid int64) *MazeEquipCache.MazeEquipPosInfo {
	if guid == 0 {
		return nil
	}
	for _, equipPos := range assembleInfo.MazeEquips {
		if equipPos.GetEquipLoadInfo().GetEquipGuid() == guid {
			return equipPos
		}
	}
	return nil
}

func GetDressedGuid(equip *MazeEquipCache.MazeEquipPosInfo) int64 {
	if equip != nil && equip.GetEquipLoadInfo().GetEquipGuid() > 0 {
		return equip.GetEquipLoadInfo().GetEquipGuid()
	}
	return 0
}

func ReplaceEquip(logger fklog.FKLogI, equipPos *MazeEquipCache.MazeEquipPosInfo, equip *MazeEquipCache.MazeEquipInfoDb) {
	equipPos.EquipLoadInfo = &MazeEquipCache.MazeEquipPosDb{}
	equipPos.EquipLoadInfo.EquipGuid = proto.Int64(equip.GetEquipGuid())
	equipPos.EquipLoadInfo.EquipId = proto.Int32(equip.GetEquipId())
	equipPos.EquipLoadInfo.Pos = proto.Int32(equipPos.GetEquipPos().GetPos())
	// if equip.EquipSubType != nil {
	// 	equipPos.EquipLoadInfo.EquipSubType = equip.EquipSubType
	// }
	// if equip.GetSuitId() > 0 {
	// 	equipPos.EquipLoadInfo.SuitId = equip.SuitId
	// }
	// var resId int32
	// if equipPos.GetEquipPos().GetPos() == 1 {
	// 	resId = pbutil.GetEquipResId(logger, equip)
	// }

	// if resId > 0 {
	// 	equipPos.EquipLoadInfo.ResId = proto.Int32(resId)
	// }
	equipPos.EquipInfo = equip
}

func ResetEquipPos(equipPos *MazeEquipCache.MazeEquipPosInfo) {
	equipPos.Force = proto.Int64(0)
	equipPos.EquipLoadInfo = &MazeEquipCache.MazeEquipPosDb{}
	equipPos.EquipLoadInfo.Pos = proto.Int32(equipPos.GetEquipPos().GetPos())
	equipPos.EquipInfo = nil
}

// 是否穿戴装备
func IsDressEquip(equips []*MazeEquipCache.MazeEquipPosInfo) bool {
	for _, equip := range equips {
		if assemble.IsAssembleEquip(equip) {
			return true
		}
	}
	return false
}
