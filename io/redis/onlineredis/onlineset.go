package onlineredis

import (
	"context"
	"fmt"
	"strconv"

	globalredis "maze_game_server/io/redis"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func getKey(ctx context.Context) string {
	return fmt.Sprintf("online:user")
}

func AddOnline(ctx context.Context, userID uint64) error {
	logger := fklog.ContextAppLogger(ctx)
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("AddOnline get redis err", zap.Error(err), zap.Any("userID", userID))
		return err
	}

	err = db.SAdd(ctx, db.MakeSectionKey(getKey(ctx)), userID).Err()
	if err != nil {
		logger.ErrorWF("AddOnline SAdd err", zap.Error(err), zap.Any("userID", userID))
		return err
	}

	return nil
}

func DelOnline(ctx context.Context, userID uint64) error {
	logger := fklog.ContextAppLogger(ctx)
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("DelOnline get redis err", zap.Error(err), zap.Any("userID", userID))
		return err
	}

	err = db.SRem(ctx, db.MakeSectionKey(getKey(ctx)), userID).Err()
	if err != nil {
		logger.ErrorWF("DelOnline SRem err", zap.Error(err), zap.Any("userID", userID))
		return err
	}

	return nil
}

func GetAllOnline(ctx context.Context) ([]uint64, error) {
	logger := fklog.ContextAppLogger(ctx)
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("GetAllOnline get redis err", zap.Error(err))
		return nil, err
	}

	members, err := db.SMembers(ctx, db.MakeSectionKey(getKey(ctx))).Result()
	if err != nil {
		logger.ErrorWF("GetAllOnline SMembers err", zap.Error(err))
		return nil, err
	}

	var userIDs []uint64
	for _, member := range members {
		uid, err := strconv.ParseUint(member, 10, 64)
		if err != nil {
			logger.ErrorWF("GetAllOnline parse uint64 err", zap.Error(err), zap.String("member", member))
			return nil, err
		}
		userIDs = append(userIDs, uid)
	}

	return userIDs, nil
}
