package bagmodule

import (
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipCache"
)

const (
	DB_OP_TYPE_ADD = iota
	DB_OP_TYPE_UPDATE
	DB_OP_TYPE_REM
)

type BagChangeInfo struct {
	fklog.FKLogI
	changeEquips []*BagChangeEquip
}

type BagChangeEquip struct {
	BagEquip *MazeEquipCache.MazeEquipInfoDb
	ChgType  int32
}

func NewBagChangeInfo(logger fklog.FKLogI) *BagChangeInfo {
	return &BagChangeInfo{FKLogI: logger, changeEquips: []*BagChangeEquip{}}
}

func (b *BagChangeInfo) addInfo(equip *MazeEquipCache.MazeEquipInfoDb, chgType int32) {
	change := &MazeEquipCache.MazeEquipInfoDb{}
	proto.Merge(change, equip)
	info := &BagChangeEquip{BagEquip: change, ChgType: chgType}
	b.changeEquips = append(b.changeEquips, info)
}

func (b *BagChangeInfo) AddBagChange(equip *MazeEquipCache.MazeEquipInfoDb) {
	b.addInfo(equip, DB_OP_TYPE_UPDATE)
}

func (b *BagChangeInfo) AddBagRem(equip *MazeEquipCache.MazeEquipInfoDb) {
	b.addInfo(equip, DB_OP_TYPE_REM)
}

func (b *BagChangeInfo) AddBagAdd(equip *MazeEquipCache.MazeEquipInfoDb) {
	b.addInfo(equip, DB_OP_TYPE_ADD)
}

//func (b *BagChangeInfo) GetAllChgBag() int {
//	return len(b.changeEquips)
//}

// 合并添加上又删除的装备信息
func (b *BagChangeInfo) MergeAddAndRemoveUpdateInfo() (rems []int64, adds []*MazeEquipCache.MazeEquipInfoDb) {
	adds = make([]*MazeEquipCache.MazeEquipInfoDb, 0)
	rems = make([]int64, 0)
	for _, equip := range b.changeEquips {
		if equip.ChgType == DB_OP_TYPE_REM {
			rems = append(rems, equip.BagEquip.GetEquipGuid())
		} else {
			adds = append(adds, equip.BagEquip)
		}
	}
	return rems, adds
}
