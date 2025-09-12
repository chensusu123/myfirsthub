package moneymodel

import (
	"context"
	"maze_game_server/common/function/maputil"
	"maze_game_server/io/redis/moneyredis"
)

type MoneyModel struct {
}

func NewMoneyModel(ctx context.Context, userID uint64) (*MoneyModel, error) {
	moneyModel := &MoneyModel{}
	return moneyModel, nil
}

func (m *MoneyModel) Load(ctx context.Context, userID uint64, itemId int32) (int64, error) {
	count, err := moneyredis.Get(ctx, userID, itemId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *MoneyModel) BatchLoad(ctx context.Context, userID uint64, itemIds []int32) (map[int32]int64, error) {
	if itemIds == nil || len(itemIds) == 0 {
		return make(map[int32]int64), nil
	}
	itemIdsStr := maputil.SliceInt32ToString(itemIds)
	countList, err := moneyredis.BatchGet(ctx, userID, itemIdsStr)
	if err != nil {
		return nil, err
	}
	res := maputil.SliceToMap(itemIds, countList)
	return res, nil
}

func (m *MoneyModel) LoadAll(ctx context.Context, userID uint64) (map[int32]int64, error) {
	mpStr, err := moneyredis.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	mp, err := maputil.MapStrStrToI32I64(mpStr)
	if err != nil {
		return nil, err
	}
	return mp, nil
}

func (m *MoneyModel) SetValue(ctx context.Context, userID uint64, itemId int32, value int64) (err error) {
	return moneyredis.Set(ctx, userID, itemId, value)
}

func (m *MoneyModel) IncrValue(ctx context.Context, userID uint64, itemId int32, value int64) (int64, error) {
	return moneyredis.Incr(ctx, userID, itemId, value)
}

func (m *MoneyModel) Del(ctx context.Context, userID uint64) (err error) {
	return moneyredis.Del(ctx, userID)
}

func (m *MoneyModel) BatchSet(ctx context.Context, userID uint64, itemMap map[int32]int64) (err error) {
	return moneyredis.BatchSet(ctx, userID, itemMap)
}
