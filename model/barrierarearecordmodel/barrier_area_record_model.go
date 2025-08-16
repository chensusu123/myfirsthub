package barrierarearecordmodel

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/redis/barrierarearecordredis"
	"maze_game_server/lib/serialize"
)

//type BarrierAreaNumRecordInfo struct {
//	//Damage         int64 `json:"damage,omitempty"`              //伤害值
//	//KillMonsterNum int32 `json:"kill_monster_number,omitempty"` //杀怪数
//	//AreaIndex      int32 `json:"area_index,omitempty"`
//	//AreaId         int32 `json:"area_id,omitempty"`
//
//	DamageRecordMap      map[int32]int64 `json:"damage_record_map,omitempty"`       //伤害值
//	KillMonsterRecordMap map[int32]int64 `json:"kill_monster_record_map,omitempty"` //杀怪数
//
//}

type BarrierAreaNumRecordModel struct {
	//PassAreaList []*BarrierAreaNumRecordInfo `json:"pass_area_list,omitempty"`
	DamageRecordMap      map[int32]int64   `json:"damage_record_map,omitempty"`       //伤害值
	KillMonsterRecordMap map[int32]int32   `json:"kill_monster_record_map,omitempty"` //杀怪数
	KillMonsterGuidMap   map[int32][]int64 `json:"kill_monster_guid_map,omitempty"`   //杀怪的guid列表
}

func NewBarrierAreaNumRecordModel(logger fklog.FKLogI, userID uint64, stageId int32) (*BarrierAreaNumRecordModel, error) {
	passArea := &BarrierAreaNumRecordModel{
		DamageRecordMap:      make(map[int32]int64),
		KillMonsterRecordMap: make(map[int32]int32),
		KillMonsterGuidMap:   make(map[int32][]int64),
	}
	if err := passArea.load(logger, userID, stageId); err != nil {
		return nil, err
	}
	return passArea, nil
}

func (p *BarrierAreaNumRecordModel) load(logger fklog.FKLogI, userID uint64, stageId int32) (err error) {
	bytes, err := barrierarearecordredis.GetBarrierAreaNumRecord(logger, userID, stageId)
	if err != nil {
		return err
	}
	if bytes == nil {
		return nil
	}
	err = serialize.Unmarshal(bytes, p)
	if err != nil {
		logger.ErrorWF("BarrierAreaNumRecordInfo load Unmarshal failed", zap.Error(err), zap.Uint64("userID", userID), zap.Int32("stageId", stageId))
		return err
	}
	return
}

func (p *BarrierAreaNumRecordModel) Save(logger fklog.FKLogI, userID uint64, stageId int32) (err error) {
	bytes, err := serialize.Marshal(p)
	if err != nil {
		logger.ErrorWF("BarrierAreaNumRecordInfo save Marshal failed", zap.Error(err), zap.Uint64("userID", userID), zap.Int32("stageId", stageId))
		return err
	}
	return barrierarearecordredis.SetBarrierAreaNumRecord(logger, userID, stageId, bytes)
}

func (p *BarrierAreaNumRecordModel) Del(logger fklog.FKLogI, userID uint64, stageId int32) (err error) {
	return barrierarearecordredis.DelBarrierAreaNumRecord(logger, userID, stageId)
}
