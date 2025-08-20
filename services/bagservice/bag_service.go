/*
@Author: xiaobo
@Date: 2025/3/21 11:42
@Description:
*/

package bagservice

import (
	"context"
	"maze_game_server/services/itemservice"
)

const EnterTypeMazeBag = 3 // 进入背包类型 迷宫背包

func Register(reg *itemservice.RegisterInfo) {
	reg.RegisterByBagType(EnterTypeMazeBag, gBagService)
}

type BagService interface {
	// 获取全部背包物品
	GetAllBagItem(ctx context.Context, userID uint64) (map[int32]int64, error)
	// 删除全部背包物品
	DelAllBagItem(ctx context.Context, userID uint64) error
}

var GlobalBagService BagService
var gBagService *service

func init() {
	GlobalBagService = newBagService()
}

type service struct {
}

func newBagService() BagService {
	gBagService = &service{}
	return gBagService
}
