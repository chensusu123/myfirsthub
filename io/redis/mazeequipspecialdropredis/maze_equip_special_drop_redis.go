package mazeequipspecialdropredis

import (
	"fmt"
)

func getRedisKey(userId uint64) string {
	return fmt.Sprintf("maze:equip:special:drop:u:%d", userId)
}

//func GetMazeEquipSpecialDropInfo(ctx context.Context, userId uint64) ([]byte, error) {
//	logger := fklog.ContextAppLogger(ctx)
//	db, err := globalredis.GCli.GetDB()
//	if err != nil {
//		return nil, err
//	}
//	key := db.MakeSectionKey(getRedisKey(userId))
//	bytes, err := db.Get(ctx, key).Bytes()
//	if err != nil {
//		if err == redis.Nil {
//			return nil, nil
//		}
//		logger.CtxError(ctx, "GetMazeEquipSpecialDropInfo GET", zap.String("key", key), zap.Error(err))
//		return nil, err
//	}
//	logger.CtxInfo(ctx, "GetMazeEquipSpecialDropInfo success", zap.String("key", key))
//	return bytes, nil
//}

//func SetMazeEquipSpecialDropInfo(ctx context.Context, userId uint64, bytes []byte) (err error) {
//	logger := fklog.ContextAppLogger(ctx)
//	db, err := globalredis.GCli.GetDB()
//	if err != nil {
//		return err
//	}
//	key := db.MakeSectionKey(getRedisKey(userId))
//	_, err = db.Set(ctx, key, bytes, 0).Result()
//	if err != nil {
//		if err == redis.Nil {
//			return nil
//		}
//		logger.CtxError(ctx, "SetMazeEquipSpecialDropInfo GET", zap.String("key", key), zap.Error(err))
//		return err
//	}
//	logger.CtxInfo(ctx, "SetMazeEquipSpecialDropInfo success", zap.Any("bytes", bytes), zap.String("key", key))
//	return
//}

// gm删除
//func GMDel(ctx context.Context, userId uint64) (err error) {
//	logger := fklog.ContextAppLogger(ctx)
//	db, err := globalredis.GCli.GetDB()
//	if err != nil {
//		return err
//	}
//	key := db.MakeSectionKey(getRedisKey(userId))
//	err = db.Del(ctx, key).Err()
//	if err != nil {
//		logger.CtxError(ctx, "SetMazeEquipSpecialDropInfo GMDel fail", zap.String("key", key), zap.Error(err))
//		return
//	}
//	return
//}
