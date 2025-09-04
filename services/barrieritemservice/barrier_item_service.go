package barrieritemservice

import (
	"context"
	"maze_game_server/services/itemservice"
)

// 关卡内物品统一掉落存储
type barrierItemsService interface {
	// 清除当前关卡存储
	ClearBarrierItems(ctx context.Context, userID uint64, barrierID int32) error
	// 清除装备以外的物品
	DelInAdditionToEquips(ctx context.Context, userID uint64, barrierID int32) error
	// 尝试扣除道具
	TrySubBarrierItems(ctx context.Context, userID uint64, barrierID int32, items []*itemservice.ItemInfo, equips []*itemservice.ItemInfo) (bool, error)
	// 掉落装备信息生成
	FallOffEquip(ctx context.Context, userID uint64, barrierID int32, equipNum int32) ([]*itemservice.ItemInfo, error)
	// 增加装备掉落分数 如果可以增加装备 将会增加装备
	AddEquipScore(ctx context.Context, userID uint64, barrierID int32, score int32, monsterGuid int64, monsterPos string) error
	// 增加物品掉落分数 如果可以增加物品 将会增加物品 如果不是特殊展示物品 则 itemType为itemID score为Count
	AddItemScore(ctx context.Context, userID uint64, barrierID int32, itemType int32, score int32, monsterGuid int64, monsterPos string) error
	// 技能道具掉落
	FallOffSkillItems(ctx context.Context, userID uint64, barrierID int32, killMonsterNum int32, nowBloodVolume int64, allBloodVolume int64, monsterGuid int64, monsterPos string) error
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
