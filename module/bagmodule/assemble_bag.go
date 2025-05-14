package bagmodule

import (
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipCache"
)

type AssembleEquip struct {
	BaseAssembleEquip
}

type BaseAssembleEquip struct {
	equips map[int64]*MazeEquipCache.MazeEquipInfoDb
}

func NewAssembleEquip() *AssembleEquip {
	return &AssembleEquip{BaseAssembleEquip: BaseAssembleEquip{equips: make(map[int64]*MazeEquipCache.MazeEquipInfoDb)}}
}

func (b *AssembleEquip) LoadAssembleEquip(equip *MazeEquipCache.MazeEquipInfoDb) {
	b.equips[equip.GetEquipGuid()] = equip
}

// 返回主背包数据
func (b *AssembleEquip) GetAssembleEquip() map[int64]*MazeEquipCache.MazeEquipInfoDb {
	return b.equips
}
