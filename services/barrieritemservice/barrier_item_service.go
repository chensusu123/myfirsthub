package barrieritemservice

import (
	"context"
	"maze_game_server/services/itemservice"
)

// 关卡内物品统一掉落存储
type barrierItemsService interface {
	// 清除当前关卡存储
	ClearBarrierItems(ctx context.Context, userID uint64, barrierID int32) (int64, int64, error)
	// 清除装备以外的物品
	DelInAdditionToEquips(ctx context.Context, userID uint64, barrierID int32) error
	// 尝试扣除道具
	TrySubBarrierItems(ctx context.Context, userID uint64, barrierID int32, items []*itemservice.ItemInfo, equips []*itemservice.ItemInfo) (bool, error)
	// 掉落装备信息生成
	FallOffEquip(ctx context.Context, userID uint64, barrierID int32, equipNum int32) ([]*itemservice.ItemInfo, error)
	// 增加装备掉落分数 如果可以增加装备 将会增加装备
	AddEquipScore(ctx context.Context, userID uint64, barrierID, score, killMonsterNum int32, guid int64, pos string) (res int32, dropItems []*itemservice.ItemInfo, err error)
	// 增加 带掉落分数的物品
	AddScoreItem(ctx context.Context, userID uint64, barrierID int32, itemType int32, score int32, guid int64, pos string) (res int32, dropItems []*itemservice.ItemInfo, err error)
	// 增加普通物品
	AddItems(ctx context.Context, userID uint64, barrierID int32, items []*itemservice.ItemInfo, guid int64, pos string) error
	// 技能道具掉落
	FallOffSkillItems(ctx context.Context, userID uint64, barrierID int32, killMonsterNum int32, nowBloodVolume int64, allBloodVolume int64, guid int64, pos string) (dropItems []*itemservice.ItemInfo, err error)
	// 增加血瓶
	SpecialAddBloodBottles(ctx context.Context, userID uint64, barrierID int32) error
	// 检测血瓶属性
	// -- 血瓶当前数量 同时拥有血瓶上限 血瓶使用cd
	CheckBloodAttr(ctx context.Context, userID uint64, barrierID int32) (bloodBottleLimit, bloodBottleCd int64, err error)
}

var GbarrierItemsService barrierItemsService

type service struct {
}

func NewBarrierItemsService() barrierItemsService {
	return &service{}
}

func init() {
	GbarrierItemsService = NewBarrierItemsService()
}
