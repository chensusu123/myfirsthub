package tempbuffredis

import (
	"context"
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// todo 通过区域的redis直接使用了临时buff的redis， 因为通过区域只有临时buff这使用了

// getKey 获取缓存操作key。
func gePassAreaKey(userId uint64, barrierId int32) string {
	return fmt.Sprintf("area:u:%d:barrier:%d", userId, barrierId)
}

// GetBarrierPassArea
func GetBarrierPassArea(logger fklog.FKLogI, userId uint64, barrierId int32) ([]byte, error) {
	db, err := gCli.GetDB()
	if err != nil {
		return nil, err
	}
	key := db.MakeSectionKey(gePassAreaKey(userId, barrierId))

	bytes, err := db.Get(context.TODO(), key).Bytes()
	if err != nil {
		logger.ErrorWF("GetBarrierPassArea GET err", zap.String("key", key), zap.Error(err))
		return nil, err
	}

	return bytes, nil
}

// SetBarrierPassArea
func SetBarrierPassArea(logger fklog.FKLogI, userId uint64, barrierId int32, data []byte) (err error) {
	db, err := gCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(gePassAreaKey(userId, barrierId))

	err = db.Set(context.TODO(), key, data, 0).Err()
	if err != nil {
		logger.ErrorWF("SetBarrierPassArea err", zap.String("key", key), zap.Any("passArea", string(data)),
			zap.Error(err))
		return err
	}
	logger.InfoWF("SetBarrierPassArea end", zap.String("key", key), zap.Any("passArea", string(data)))
	return nil
}

// DelBarrierPassArea
func DelBarrierPassArea(logger fklog.FKLogI, userId uint64, barrierId int32) (err error) {
	db, err := gCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(gePassAreaKey(userId, barrierId))

	err = db.Del(context.TODO(), key).Err()
	if err != nil {
		logger.ErrorWF("DelBarrierPassArea err", zap.String("key", key), zap.Error(err))
		return err
	}
	logger.InfoWF("DelBarrierPassArea end", zap.String("key", key))
	return nil
}
