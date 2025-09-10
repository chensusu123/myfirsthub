package barrierstagecountermodel

import (
	"context"
	"fmt"
	"maze_game_server/io"
	globalredis "maze_game_server/io/redis"
)

func getRedisKey(userId uint64, barrier int32) string {
	return fmt.Sprintf("maze:u:%d:barrier:%d:stage:count:record", userId, barrier)
}

func getSectionKey(s string) (string, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return "", err
	}
	key := db.MakeSectionKey(s)
	return key, err
}

// 关卡计数器
type BarrierStageCounterModel struct {
	ExpMap               map[int32]int64              `json:"exp_map,omitempty"`                 //杀怪获得的经验
	DamageRecordMap      map[int32]int64              `json:"damage_record_map,omitempty"`       //伤害值
	KillMonsterRecordMap map[int32]int32              `json:"kill_monster_record_map,omitempty"` //杀怪数
	KillMonsterGuidMap   map[int32]map[int64]struct{} `json:"kill_monster_guid_map,omitempty"`   //杀怪的guid列表
}

func NewBarrierStageCounterModel(ctx context.Context, userID uint64, barrierId int32) (*BarrierStageCounterModel, error) {
	passArea := &BarrierStageCounterModel{
		ExpMap:               make(map[int32]int64),
		DamageRecordMap:      make(map[int32]int64),
		KillMonsterRecordMap: make(map[int32]int32),
		KillMonsterGuidMap:   make(map[int32]map[int64]struct{}),
	}

	if err := passArea.load(ctx, userID, barrierId); err != nil {
		return nil, err
	}
	return passArea, nil
}

func (p *BarrierStageCounterModel) load(ctx context.Context, userID uint64, barrierId int32) (err error) {
	return io.LoadSvrData(ctx, getRedisKey(userID, barrierId), p)
}

func (p *BarrierStageCounterModel) Save(ctx context.Context, userID uint64, barrierId int32) (err error) {
	return io.SaveSvrData(ctx, getRedisKey(userID, barrierId), p)
}

func (p *BarrierStageCounterModel) Del(ctx context.Context, userID uint64, barrierId int32) (err error) {
	return io.DeleteSvrData(ctx, getRedisKey(userID, barrierId))
}

func GMDel(ctx context.Context, userID uint64, barrierID int32) (err error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return err
	}

	key, err := getSectionKey(getRedisKey(userID, barrierID))
	if err != nil {
		return err
	}

	return db.Del(ctx, key).Err()
}
