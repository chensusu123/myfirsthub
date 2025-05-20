// @Desc: \\todo
// @Author: LiuChongyu 2021/11/17 4:03 下午
// @Update: LiuChongyu 2021/11/17 4:03 下午

package UnionIDBindRedis

import (
	"context"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

var gRedisCli = &fkredis.FkRedis{}

func init() {
	fkconfig.RegisterIO(gRedisCli, 16895 /*cgk数据库类型*/)
}

func AddUnionID2UserID(logger fklog.FKLogI, unionID, userID uint64) error {
	ctx := context.TODO()

	key := fkutil.K_str("paipai:unionid:%d:to:userid:set", unionID)

	_, err := gRedisCli.Do(ctx, "SADD", key, userID)
	if err != nil {
		logger.ErrorWF("failed to sadd new userid in unionid",
			zap.Error(err),
			zap.String("key", key),
			zap.Uint64("userID", userID),
		)
		return err
	}

	return nil
}

func AddUserID2UnionID(logger fklog.FKLogI, userID, unionID uint64) error {
	ctx := context.TODO()
	key := fkutil.K_str("paipai:userid:%d:to:unionid:string", userID)

	_, err := gRedisCli.Do(ctx, "SET", key, unionID)
	if err != nil {
		logger.ErrorWF("failed to set unionid on userid",
			zap.Error(err),
			zap.String("key", key),
			zap.Uint64("unionID", unionID),
		)
		return err
	}

	return nil
}

func GetUsersWithUnionID(logger fklog.FKLogI, unionID uint64) (users []uint64, err error) {
	ctx := context.TODO()
	key := fkutil.K_str("paipai:unionid:%d:to:userid:set", unionID)
	ret, err := redis.Strings(gRedisCli.Do(ctx, "SMEMBERS", key))

	if err == redis.ErrNil {
		logger.ErrorWF("GetUsersWithUnionID empty",
			zap.String("key", key),
			fklog.Any("err", err),
		)
		return users, err
	}

	if err != nil {
		logger.ErrorWF("GetUsersWithUnionID error",
			fklog.Uint64("unionID", unionID),
			fklog.String("key", key),
			fklog.Any("err", err),
		)
		return users, err
	}

	for _, v := range ret {
		user := fkutil.ToUint64(v)
		if user <= 0 {
			logger.ErrorWF("GetUsersWithUnionID ToUint64",
				zap.Uint64("unionID", unionID),
				zap.String("v", v),
			)
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUserID2UnionID(logger fklog.FKLogI, userID uint64) (unionID uint64, err error) {
	ctx := context.TODO()
	key := fkutil.K_str("paipai:userid:%d:to:unionid:string", userID)
	ret, err := redis.String(gRedisCli.Do(ctx, "GET", key))

	if err == redis.ErrNil {
		logger.WarnWF("GetUserID2UnionID empty",
			zap.String("key", key))
		return 0, nil
	}

	if err != nil {
		logger.ErrorWF("GetValue",
			zap.String("key", key),
			zap.Error(err),
		)
		return 0, err
	}
	unionID = fkutil.ToUint64(ret)
	return unionID, nil
}
