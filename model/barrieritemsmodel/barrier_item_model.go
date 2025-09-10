package barrieritemsmodel

import (
	"context"
	"fmt"
	"maze_game_server/io"
	"maze_game_server/services/itemservice"
)

var (
	SpecialType = 483 // 技能道具类型
)

func getKey(userID uint64, barrierID int32) string {
	return fmt.Sprintf("barrier:items:u:%d:%d", userID, barrierID)
}

type barrierItems struct {
	Items           map[int64]*itemservice.ItemInfo `json:"items,omitempty"`              // 物品
	Equips          map[int64]*itemservice.ItemInfo `json:"equips,omitempty"`             // 装备
	EquipScore      int32                           `json:"equip_score,omitempty"`        // 杀怪获得的装备分数
	ItemsScore      map[int32]int32                 `json:"items_score,omitempty"`        // 杀怪获得的物品分数
	SkillsCount     map[int32]int32                 `json:"blood_bottle_count,omitempty"` // 技能道具使用次数存储
	SkillDropTime   map[int32]int64                 `json:"skill_drop_time,omitempty"`    // 技能道具掉落间隔
	BloodBottleAttr map[int32]int64                 `json:"blood_bottle_attr,omitempty"`  // 血瓶属性存储 感知变化时使用 跟着关卡走
}

func NewBarrierItems(ctx context.Context, userID uint64, barrierID int32) (*barrierItems, error) {
	res := &barrierItems{
		Items:           make(map[int64]*itemservice.ItemInfo),
		Equips:          make(map[int64]*itemservice.ItemInfo),
		EquipScore:      0,
		ItemsScore:      make(map[int32]int32),
		SkillsCount:     make(map[int32]int32),
		SkillDropTime:   make(map[int32]int64),
		BloodBottleAttr: make(map[int32]int64),
	}
	if err := res.load(ctx, userID, barrierID); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *barrierItems) load(ctx context.Context, userID uint64, barrierID int32) (err error) {
	return io.LoadSvrData(ctx, getKey(userID, barrierID), s)
}

func (s *barrierItems) Save(ctx context.Context, userID uint64, barrierID int32) (err error) {
	return io.SaveSvrData(ctx, getKey(userID, barrierID), s)
}

func (s *barrierItems) Del(ctx context.Context, userID uint64, barrierID int32) (err error) {
	return io.DeleteSvrData(ctx, getKey(userID, barrierID))
}
