package bagmodel

import (
	"context"
	"maze_game_server/common/function/maputil"
	"maze_game_server/io/redis/bagredis"
)

type BagModel struct {
}

func NewBagModel(ctx context.Context, userID uint64) (*BagModel, error) {
	bagModel := &BagModel{}
	return bagModel, nil
}

//func (t *BagModel) Load(ctx context.Context, userID uint64, itemId int32) (int64, error) {
//	count, err := bagredis.GetAllBagItem(ctx, userID, itemId)
//	if err != nil {
//		return 0, err
//	}
//	return count, nil
//}

func (t *BagModel) BatchLoad(ctx context.Context, userID uint64, itemIds []int32) (map[int32]int64, error) {
	if itemIds == nil || len(itemIds) == 0 {
		return make(map[int32]int64), nil
	}
	itemIdsStr := maputil.SliceInt32ToString(itemIds)
	countList, err := bagredis.BatchGet(ctx, userID, itemIdsStr)
	if err != nil {
		return nil, err
	}
	res := maputil.SliceToMap(itemIds, countList)
	return res, nil
}

func (t *BagModel) BatchDel(ctx context.Context, userID uint64, itemIds []int32) error {
	if itemIds == nil || len(itemIds) == 0 {
		return nil
	}
	itemIdsStr := maputil.SliceInt32ToString(itemIds)
	err := bagredis.BatchDel(ctx, userID, itemIdsStr)
	if err != nil {
		return err
	}
	return nil
}

func (t *BagModel) LoadAll(ctx context.Context, userID uint64) (map[int32]int64, error) {
	mpStr, err := bagredis.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	mp, err := maputil.MapStrStrToI32I64(mpStr)
	if err != nil {
		return nil, err
	}
	return mp, nil
}

//func (t *BagModel) SetValue(ctx context.Context, userID uint64, itemId int32, value int64) (err error) {
//	return bagredis.IncrBagItem(ctx, userID, itemId, value)
//}

func (t *BagModel) IncrValue(ctx context.Context, userID uint64, itemId int32, value int64) (int64, error) {
	return bagredis.Incr(ctx, userID, itemId, value)
}

func (t *BagModel) Del(ctx context.Context, userID uint64) (err error) {
	return bagredis.Del(ctx, userID)
}
