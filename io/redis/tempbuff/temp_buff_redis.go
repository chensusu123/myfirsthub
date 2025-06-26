package tempbuffredis

import (
	"context"
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/database/nanoredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/serverdepend"
	"go.uber.org/zap"
)

type TempBuffRedis struct {
	*nanoredis.NanoRedis
}

func New(serviceName string, name string) *TempBuffRedis {
	rt := &TempBuffRedis{}
	rt.NanoRedis = nanoredis.NewNanoRedis(serviceName, name)
	return rt
}

var gCli *TempBuffRedis

func init() {
	gCli = New("temp.buff.redis", "temp.buff.redis.maze_main")
	serverdepend.RegisterDepend(gCli)
}

func getKey(userId uint64, stageId int32) string {
	return fmt.Sprintf("tempbuff:u:%d:stage:%d", userId, stageId)
}

func SetMazeTempBuff(logger fklog.FKLogI, userId uint64, stageId int32, bytes []byte) error {
	db, err := gCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(getKey(userId, stageId))

	err = db.Set(context.TODO(), key, bytes, 0).Err()
	if err != nil {
		logger.ErrorWF("SetMazeTempBuff", zap.String("key", key), zap.Int32("stageId", stageId), zap.Any("bytes", bytes),
			zap.Error(err))
		return err
	}
	logger.InfoWF("SetMazeTempBuff success", zap.String("key", key), zap.Int32("stageId", stageId))
	return nil
}

func GetMazeTempBuff(logger fklog.FKLogI, userId uint64, stageId int32) ([]byte, error) {
	db, err := gCli.GetDB()
	if err != nil {
		return nil, err
	}
	key := db.MakeSectionKey(getKey(userId, stageId))

	bytes, err := db.Get(context.TODO(), key).Bytes()

	if err != nil {
		logger.ErrorWF("GetMazeTempBuff GET", zap.String("key", key), zap.Int32("stageId", stageId), zap.Error(err))
		return nil, err
	}

	logger.InfoWF("GetMazeTempBuff success", zap.String("key", key), zap.Int32("stageId", stageId))
	return bytes, nil
}

func DelMazeTempBuff(logger fklog.FKLogI, userId uint64, stageId int32) error {
	db, err := gCli.GetDB()
	if err != nil {
		return err
	}
	key := db.MakeSectionKey(getKey(userId, stageId))

	err = db.Del(context.TODO(), key).Err()
	if err != nil {
		logger.ErrorWF("DelMazeTempBuff DEL", zap.String("key", key), zap.Error(err), zap.Int32("stageId", stageId))
		return err
	}
	logger.InfoWF("DelMazeTempBuff success", zap.String("key", key), zap.Int32("stageId", stageId))
	return nil
}
