package bagmodule

import (
	"context"
	"maze_game_server/pb/server/MazeEquipCache"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"google.golang.org/protobuf/proto"
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

func NewBagChangeInfo(ctx context.Context) *BagChangeInfo {
	return &BagChangeInfo{FKLogI: fklog.ContextAppLogger(ctx), changeEquips: []*BagChangeEquip{}}
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
