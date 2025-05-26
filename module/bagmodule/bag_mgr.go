package bagmodule

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/redis/dollassemblesuitredis"
	"maze_game_server/io/redis/mazebagequipredis"
	"maze_game_server/pb/server/MazeEquipCache"
)

// 背包数据管理
type BagEquipMgr struct {
	UserID uint64
	// 背包列表
	MainBagEquips *MainBagEquip
	// 已装配装备
	AssembleEquips *AssembleEquip
	// 装备变化列表
	dbChange *BagChangeInfo
	fklog.FKLogI
}

func NewBagEquipMgr(logger fklog.FKLogI, userId uint64) *BagEquipMgr {
	ret := &BagEquipMgr{UserID: userId,
		MainBagEquips:  NewMainBagEquip(),
		AssembleEquips: NewAssembleEquip(),
		dbChange:       NewBagChangeInfo(logger),
		FKLogI:         logger,
	}
	return ret
}

// 加载数据
func (m *BagEquipMgr) LoadBagFromRedis() error {
	equipMap, err := mazebagequipredis.GetAllEquipInfo(m, m.UserID)
	if err != nil {
		m.ErrorWF("LoadBagFromRedis GetAllEquipInfo error", zap.Error(err))
		return err
	}

	// 获取身上的装备信息
	assembleInfoMap, err := dollassemblesuitredis.GetAllDollAssembleSuit(m, m.UserID)
	if err != nil {
		m.ErrorWF("LoadBagFromRedis GetAllDollAssembleSuit error", zap.Error(err))
		return err
	}

	for _, equipList := range assembleInfoMap {
		for _, v := range equipList {
			if v.GetEquipGuid() > 0 {
				equip, ok := equipMap[v.GetEquipGuid()]
				if ok {
					m.AssembleEquips.LoadAssembleEquip(equip)
					delete(equipMap, v.GetEquipGuid())
				}
			}
		}
	}
	for _, equip := range equipMap {
		m.MainBagEquips.LoadMainEquip(equip)
	}
	return nil
}

func (m *BagEquipMgr) SaveBagInfoToRedis() error {
	rems, adds := m.dbChange.MergeAddAndRemoveUpdateInfo()
	if len(rems) > 0 {
		err := mazebagequipredis.BatchDelEquip(m, m.UserID, rems...)
		if err != nil {
			m.ErrorWF("SaveBagInfoToRedis BatchDelEquip error", zap.Error(err), zap.Any("rems", rems))
			return err
		}
	}
	if len(adds) > 0 {
		err := mazebagequipredis.BatchSaveEquipInfo(m, m.UserID, adds)
		if err != nil {
			m.ErrorWF("SaveBagInfoToRedis BatchSaveEquipInfo error", zap.Error(err), zap.Any("adds", adds))
			return err
		}
	}
	return nil
}

// func (m *BagEquipMgr) AddBagChange(equip *MazeEquipCache.MazeEquipInfoDb) {
//	m.dbChange.AddBagChange(equip)
//	m.MainBagEquips.LoadMainEquip(equip)
// }

func (m *BagEquipMgr) AddBagRem(equip *MazeEquipCache.MazeEquipInfoDb) {
	m.dbChange.AddBagRem(equip)
	m.MainBagEquips.RemMainEquip(equip)
}

func (m *BagEquipMgr) AddBagAdd(equip *MazeEquipCache.MazeEquipInfoDb) {
	m.dbChange.AddBagAdd(equip)
	m.MainBagEquips.LoadMainEquip(equip)
}

func (m *BagEquipMgr) CheckEquipExist(equipGuid int64) bool {
	_, ok := m.MainBagEquips.equips[equipGuid]
	if ok {
		return true
	}
	_, ok = m.AssembleEquips.equips[equipGuid]
	if ok {
		return true
	}
	return false
}
