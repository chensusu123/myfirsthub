package bagmodule

import (
	"gitlab.ifreetalk.com/plate/protodef/MazeEquipCache"
)

type MainBagEquip struct {
	BaseBagEquip
}

type BaseBagEquip struct {
	equips map[int64]*MazeEquipCache.MazeEquipInfoDb
}

func NewMainBagEquip() *MainBagEquip {
	return &MainBagEquip{BaseBagEquip: BaseBagEquip{equips: make(map[int64]*MazeEquipCache.MazeEquipInfoDb)}}
}

func (b *BaseBagEquip) LoadMainEquip(equip *MazeEquipCache.MazeEquipInfoDb) {
	b.equips[equip.GetEquipGuid()] = equip
}

func (b *BaseBagEquip) RemMainEquip(equip *MazeEquipCache.MazeEquipInfoDb) {
	delete(b.equips, equip.GetEquipGuid())
}

// 返回主背包数据
func (b *BaseBagEquip) GetMainEquip() map[int64]*MazeEquipCache.MazeEquipInfoDb {
	return b.equips
}
