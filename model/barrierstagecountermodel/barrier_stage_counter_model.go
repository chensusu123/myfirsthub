package barrierstagecountermodel

import (
	"context"
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"maze_game_server/io"
)

func getRedisKey(userId uint64, barrier int32) string {
	return fmt.Sprintf("maze:u:%d:barrier:%d:stage:count:record", userId, barrier)
}

// 关卡计数器
type BarrierStageCounterModel struct {
	ExpMap               map[int32]int64   `json:"exp_map,omitempty"`                 //杀怪获得的经验
	DamageRecordMap      map[int32]int64   `json:"damage_record_map,omitempty"`       //伤害值
	KillMonsterRecordMap map[int32]int32   `json:"kill_monster_record_map,omitempty"` //杀怪数
	KillMonsterGuidMap   map[int32][]int64 `json:"kill_monster_guid_map,omitempty"`   //杀怪的guid列表
}

func NewBarrierStageCounterModel(logger fklog.FKLogI, userID uint64, barrierId int32) (*BarrierStageCounterModel, error) {
	passArea := &BarrierStageCounterModel{
		ExpMap:               make(map[int32]int64),
		DamageRecordMap:      make(map[int32]int64),
		KillMonsterRecordMap: make(map[int32]int32),
		KillMonsterGuidMap:   make(map[int32][]int64),
	}

	if err := passArea.load(logger, userID, barrierId); err != nil {
		return nil, err
	}
	return passArea, nil
}

func (p *BarrierStageCounterModel) load(logger fklog.FKLogI, userID uint64, barrierId int32) (err error) {
	return io.LoadSvrData(context.TODO(), getRedisKey(userID, barrierId), p)
}

func (p *BarrierStageCounterModel) Save(logger fklog.FKLogI, userID uint64, barrierId int32) (err error) {
	return io.SaveSvrData(context.TODO(), getRedisKey(userID, barrierId), p)
}

func (p *BarrierStageCounterModel) Del(logger fklog.FKLogI, userID uint64, barrierId int32) (err error) {
	return io.DeleteSvrData(context.TODO(), getRedisKey(userID, barrierId))
}
