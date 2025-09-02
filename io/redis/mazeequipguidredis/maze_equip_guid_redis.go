package mazeequipguidredis

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
)

var gRedis = &fkredis.FkRedis{}

func init() {
	//21641 maze:equip:guid:%llu 迷宫游戏装备生成guid信息
	fkconfig.RegisterNameNode("mazeequipguidredis", 21641, gRedis)
}

func GetNewGuid(ctx context.Context, userId uint64, addCount int32) (int64, error) {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("doll:equip:guid:%d", userId)
	ret, err := redis.Int64(gRedis.Do(ctx, "HINCRBY", key, "equipMaxId", addCount))
	if err != nil {
		logger.CtxError(ctx, "get new guid failed with", zap.Error(err))
		return 0, err
	}
	logger.CtxInfo(ctx, "ret new guid with", zap.Int64("inc", ret))
	return ret, nil
}

func SetEquipGuidGm(ctx context.Context, userId uint64, equipMaxGuid int64) error {
	logger := fklog.ContextAppLogger(ctx)
	key := fmt.Sprintf("doll:equip:guid:%d", userId)
	_, err := redis.Int64(gRedis.Do(ctx, "hset", key, "equipMaxId", equipMaxGuid))
	if err != nil {
		logger.CtxError(ctx, "SetEquipGuidGm failed with", zap.Error(err))
		return err
	}
	logger.CtxInfo(ctx, "SetEquipGuidGm success", zap.Int64("equipMaxGuid", equipMaxGuid))
	return nil
}
