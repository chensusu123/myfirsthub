package passarearedis

import (
	"context"
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"go.uber.org/zap"
	"maze_game_server/lib/serialize"
)

var gRedis = &fkredis.FkRedis{}

// 通过的area
type PassArea struct {
	AreaId    int32 `json:"areaId,omitempty"`
	AreaIndex int32 `json:"areaIndex,omitempty"`
}

func init() {
	fkconfig.RegisterNameNode("mazebarrierpassarearedis", 21644, gRedis)
}

// getKey 获取缓存操作key。
func getKey(userId uint64, barrierId int32) string {
	return fmt.Sprintf("u:%d:barrier:%d:area", userId, barrierId)
}

// GetBarrierPassArea
func GetBarrierPassArea(logger fklog.FKLogI, userId uint64, barrierId int32) ([]*PassArea, error) {
	key := getKey(userId, barrierId)
	res, err := redis.Bytes(gRedis.Do(context.Background(), "GET", key))
	if err == redis.ErrNil {
		return []*PassArea{}, nil
	}

	if err != nil {
		logger.ErrorWF("GetMazeTempBuff GET", zap.String("key", key), zap.Error(err))
		return nil, err
	}

	var passAreas []*PassArea
	err = serialize.Unmarshal(res, passAreas)
	if err != nil {
		logger.ErrorWF("GetBarrierPassArea Unmarshal", zap.String("key", key), zap.Error(err))
		return nil, err
	}
	logger.DebugWF("GetBarrierPassArea end", zap.String("key", key), zap.Any("passAreas", passAreas))

	return passAreas, nil
}

// SetBarrierPassArea
func SetBarrierPassArea(logger fklog.FKLogI, userId uint64, barrierId int32, passArea []*PassArea) (err error) {
	key := getKey(userId, barrierId)
	data, err := serialize.Marshal(passArea)
	if err != nil {
		return err
	}
	_, err = gRedis.Do(context.Background(), "SET", key, data)
	if err != nil {
		logger.ErrorWF("SetBarrierPassArea", zap.String("key", key), zap.Any("passArea", passArea),
			zap.Error(err))
		return err
	}
	logger.InfoWF("SetBarrierPassArea end", zap.String("key", key), zap.Any("passArea", passArea))
	return nil
}
