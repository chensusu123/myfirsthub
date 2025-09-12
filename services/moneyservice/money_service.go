package moneyservice

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/services/itemservice"
)

func Register(reg *itemservice.RegisterInfo) {
	reg.RegisterByItemID(constdef.MazeCommonItemCoin, gMoneyService)
	reg.RegisterByItemID(constdef.MazeCommonItemDiamond, gMoneyService)
}

type MoneyService interface {
	// 获取用户货币
	GetUserMoney(ctx context.Context, userId uint64) (coin, diamond int64, err error)
	// 不要用，用itemservice里面的AddItem方法
	SetMoney(ctx context.Context, userId uint64, itemId int32, value int64) error
}

var GlobalMoneyService MoneyService
var gMoneyService *service

func init() {
	GlobalMoneyService = newMoneyService()
}

type service struct {
}

func newMoneyService() MoneyService {
	gMoneyService = &service{}
	return gMoneyService
}
