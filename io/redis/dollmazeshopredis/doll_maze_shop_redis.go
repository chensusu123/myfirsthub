package dollmazeshopredis

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

var gRedis = &fkredis.FkRedis{}

func init() {
	// doll:equip:get:num:%d
	fkconfig.RegisterNameNode("dollmazeshopredis", 21560, gRedis)
}

func GetMazeShopNum(ctx context.Context, userId uint64, barrierId int32, itemId int32) (int64, error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("doll:maze:shop:%d:%d", userId, barrierId)
	ret, err := redis.Int64(gRedis.Do(ctx, "hget", key, itemId))
	if err == redis.ErrNil {
		err = nil
		return 0, err
	}
	if err != nil {
		logger.CtxError(ctx, "GetMazeShopNum get count failed with", zap.Error(err), zap.Int32("itemId", itemId))
		return 0, err
	}
	logger.CtxInfo(ctx, "GetMazeShopNum get score with", zap.Int32("itemId", itemId), zap.Int64("score", ret))
	return ret, nil
}

// 保存商店数量
func SetMazeShopNum(ctx context.Context, userId uint64, barrierId int32, itemId, count int32) error {
	key := fmt.Sprintf("doll:maze:shop:%d:%d", userId, barrierId)
	logger := fklog.ContextAppLogger(ctx)
	_, err := gRedis.Do(ctx, "hset", key, itemId, count)
	if err != nil {
		logger.CtxError(ctx, "SetMazeShopNum set score failed with", zap.Error(err),
			zap.Int32("itemId", itemId),
			zap.Int32("count", count))
		return err
	}
	logger.CtxInfo(ctx, "SetMazeShopNum set score with",
		zap.Int32("itemId", itemId),
		zap.Int32("count", count))
	return err
}

func BatchSetMazeShopNum(ctx context.Context, userId uint64, barrierId int32, mazeShopMap map[int32]int32) (err error) {
	key := fmt.Sprintf("doll:maze:shop:%d:%d", userId, barrierId)
	logger := fklog.ContextAppLogger(ctx)
	args := make([]interface{}, 0, len(mazeShopMap)*2+1)
	args = append(args, key)
	for k, v := range mazeShopMap {
		args = append(args, k)
		args = append(args, v)
	}
	if len(args) <= 1 {
		return nil
	}
	_, err = gRedis.Do(ctx, "HMSET", args...)
	if err != nil {
		logger.CtxError(ctx, "BatchSetMazeShopNum fail", zap.Error(err), zap.Any("mazeShopMap", mazeShopMap), zap.String("key", key))
		return err
	}
	logger.CtxInfo(ctx, "BatchSetMazeShopNum succ", zap.Any("mazeShopMap", mazeShopMap), zap.String("key", key))
	return
}

func GetBatchMazeShopNum(ctx context.Context, userId uint64, barrierId int32, itemIds []int32) (mazeShopMap map[int32]int32, err error) {
	logger := fklog.ContextAppLogger(ctx)
	mazeShopMap = make(map[int32]int32)
	key := fmt.Sprintf("doll:maze:shop:%d:%d", userId, barrierId)
	args := make([]interface{}, 0, len(itemIds)+1)
	args = append(args, key)
	for _, field := range itemIds {
		args = append(args, field)
	}
	resp, err := redis.Int64s(gRedis.Do(ctx, "HMGET", args...))
	if err != nil {
		if err == redis.ErrNil {
			logger.CtxInfo(ctx, "GetBatchMazeShopNum call HMGET empty", zap.Int32s("itemIds", itemIds))
			return mazeShopMap, nil
		}
		logger.CtxError(ctx, "GetBatchMazeShopNum call HMGET failed", zap.Int32s("itemIds", itemIds), zap.Any("err", err))
		return mazeShopMap, err
	}
	if len(resp) != len(itemIds) {
		logger.CtxWarn(ctx, "GetBatchMazeShopNum count not match", zap.String("key", key), zap.Int32s("itemIds", itemIds), zap.Any("resp", resp))
		return
	}
	for i, field := range itemIds {
		mazeShopMap[field] = int32(resp[i])
	}
	logger.CtxInfo(ctx, "GetBatchMazeShopNum success", zap.String("key", key), zap.Any("mazeShopMap", mazeShopMap))
	return mazeShopMap, nil
}

func GetMazeShopAllNum(ctx context.Context, userId uint64, barrierId int32) (map[int32]int32, error) {
	logger := fklog.ContextAppLogger(ctx)
	mazeShopMap := make(map[int32]int32)
	key := fmt.Sprintf("doll:maze:shop:%d:%d", userId, barrierId)
	ret, err := redis.Int64Map(gRedis.Do(ctx, "HGETALL", key))
	if err != nil {
		logger.CtxError(ctx, "GetMazeShopAllNum redis op failed with ",
			zap.Error(err),
			zap.String("key", key))
		return mazeShopMap, err
	}
	for k, v := range ret {
		field := fkutil.ToInt32(k)
		degree := int32(v)
		mazeShopMap[field] = degree
	}

	logger.CtxInfo(ctx, "GetMazeShopAllNum success",
		zap.Any("mazeShopMap", mazeShopMap),
		zap.String("key", key))
	return mazeShopMap, nil
}
