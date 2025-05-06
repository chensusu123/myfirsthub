/*
 * @Author: majian
 * @Date: 2024-08-21 16:15:50
 * @Last Modified by: majian
 * @Last Modified time: 2024-09-02 17:51:40
 */
package calcassembleattr

type SuitPos struct {
	Pos  int32
	Guid int64
}
type EquipSuitInfo struct {
	SuitId     int32
	MaxNum     int32
	SuitPosMap map[int32]*SuitPos
}

type EquipSuitMgr struct {
	EquipSuitMap map[int32]*EquipSuitInfo
}

func NewEquipSuitMgr() *EquipSuitMgr {
	r := new(EquipSuitMgr)
	r.EquipSuitMap = map[int32]*EquipSuitInfo{}
	return r
}

func NewEquipSuitInfo(suitId, maxCnt int32) *EquipSuitInfo {
	r := new(EquipSuitInfo)
	r.SuitPosMap = make(map[int32]*SuitPos)
	r.SuitId = suitId
	r.MaxNum = maxCnt
	return r
}

func (m *EquipSuitMgr) InitSuit(suitId int32, maxNum int32) *EquipSuitInfo {
	if v, ok := m.EquipSuitMap[suitId]; ok {
		return v
	}
	m.EquipSuitMap[suitId] = NewEquipSuitInfo(suitId, maxNum)
	return m.EquipSuitMap[suitId]
}

func (m *EquipSuitMgr) GetSuitInfo(suitId int32) *EquipSuitInfo {
	return m.EquipSuitMap[suitId]
}

func (m *EquipSuitInfo) AddActivePos(pos int32, guid int64) {
	m.SuitPosMap[pos] = &SuitPos{Pos: pos, Guid: guid}
}

// 查询指定套装的部位激活数量
func (m *EquipSuitInfo) GetSuitPosNum() int32 {
	var cnt int32
	for _, pos := range m.SuitPosMap {
		if pos.Guid > 0 {
			cnt++
		}
	}
	return cnt
}

// 获取激活的套装Id
func (m *EquipSuitMgr) GetActiveSuitId() int32 {
	for _, suitInfo := range m.EquipSuitMap {
		if suitInfo.GetSuitPosNum() >= suitInfo.MaxNum {
			return LegendSuitKey(suitInfo.SuitId, suitInfo.MaxNum)
		}
	}
	return 0
}

// 获取激活的套装以及对应的装备数量
func (m *EquipSuitMgr) GetSuitNumMap() map[int32]int32 {
	if m == nil {
		return nil
	}
	r := make(map[int32]int32)
	for _, suitInfo := range m.EquipSuitMap {
		r[suitInfo.SuitId] = suitInfo.GetSuitPosNum()
	}
	return r
}

// 获取指定套装的激活部位数量
func (m *EquipSuitMgr) GetSuitPosActiveNum(suitId int32) int32 {
	suitInfo := m.GetSuitInfo(suitId)
	if suitInfo != nil {
		return suitInfo.GetSuitPosNum()
	}
	return 0
}
