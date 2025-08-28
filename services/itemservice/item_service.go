package itemservice

import (
	"context"
	"maze_game_server/pb/common/MessageType"
)

type ItemService interface {
	// 添加道具
	AddItem(ctx context.Context, userId uint64, opType ItemOpType, tradeNo uint64, items ...*ItemInfo) (errInfo *MessageType.ErrorInfo)
	// 减少道具
	SubItem(ctx context.Context, userId uint64, opType ItemOpType, tradeNo uint64, items ...*ItemInfo) (errInfo *MessageType.ErrorInfo)
	// 查询道具
	QueryItems(ctx context.Context, userId uint64, items ...*ItemInfo) (queryItems []*ItemInfo, errInfo *MessageType.ErrorInfo)
}

// GlobalItemService 道具可用全局唯一对象
var GlobalItemService ItemService

func init() {
	GlobalItemService = newItemService()
}

type service struct {
}

func newItemService() ItemService {
	return &service{}
}

// checkAndMergeItem 检查并合并物品
func (s *service) checkAndMergeItem(input []*ItemInfo) (output []*ItemInfo) {
	if len(input) == 0 {
		return
	}

	itemMap := make(map[int32]*ItemInfo, len(input))
	output = make([]*ItemInfo, 0, len(input))

	var (
		tmpItem *ItemInfo
		isExist bool
	)
	for _, item := range input {
		if item.ItemId <= 0 || item.Count <= 0 {
			continue
		}

		if tmpItem, isExist = itemMap[item.ItemId]; !isExist {
			itemMap[item.ItemId] = &ItemInfo{
				ItemId: item.ItemId,
				Count:  item.Count,
			}
			continue
		}

		itemMap[item.ItemId].Count = tmpItem.Count + item.Count
	}

	for _, item := range itemMap {
		output = append(output, item)
	}
	return
}

type ItemOpType int32

// 道具操作原因类型
const (
	ItemOpTypeSweep           ItemOpType = 1  // 扫荡
	ItemOpTypePass            ItemOpType = 2  // 通关奖励
	ItemOpTypeDeath           ItemOpType = 3  // 死亡奖励
	ItemOpTypeUseItem         ItemOpType = 4  // 使用道具
	ItemOpTypeMonsterDeath    ItemOpType = 5  // 怪物掉落
	ItemOpTypeOpenBox         ItemOpType = 6  // 开箱子
	ItemOpTypePickItem        ItemOpType = 7  // 拾取
	ItemOpTypeGM              ItemOpType = 8  // Gm添加
	ItemOpTypeDismantle       ItemOpType = 9  // 装备分解
	ItemOpTypeCollect         ItemOpType = 10 // 挂机
	ItemOpTypeMail            ItemOpType = 11 // 邮件附件
	ItemOpTypePay             ItemOpType = 12 // 支付
	ItemOpTypeEquipPosLvUp    ItemOpType = 13 // 装备位升级
	ItemOpTypeReborn          ItemOpType = 14 // 重生
	ItemOpTypeUseEnergy       ItemOpType = 15 // 使用体力道具
	ItemOpTypeEquipMix        ItemOpType = 16 // 装备合成
	ItemOpTypeRefreshTempBuff ItemOpType = 17 // 刷新三选一词条
)
