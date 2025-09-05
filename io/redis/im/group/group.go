package group

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	globalredis "maze_game_server/io/redis"

	"github.com/redis/go-redis/v9"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/database/nanoredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

type Group struct {
	ID         int64    `json:"id,omitempty"`
	Creator    uint64   `json:"creator,omitempty"`
	CreateTime int64    `json:"create_time,omitempty"`
	GroupType  string   `json:"group_type,omitempty"`
	Members    []Member `json:"-"`
}

type Member struct {
	CreateTime int64  `json:"create_time,omitempty"`
	UserID     uint64 `json:"user_id,omitempty"`
	Nickname   string `json:"nickname,omitempty"`
}

// getKey 获取缓存操作key。
func getKey(db nanoredis.NanoRedisClient, args ...interface{}) string {
	return db.MakeSectionKey(fmt.Sprintf("im:app:%d:group:%d", args...))
}

// GetGroupInfo
func GetGroupInfo(ctx context.Context, appID int32, groupID int64) (group *Group, err error) {
	logger := fklog.ContextAppLogger(ctx)

	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "GetGroupInfo Client fail",
			zap.Error(err),
			zap.Int64("groupID", groupID),
		)
		return nil, err
	}
	key := getKey(cli, appID, groupID)
	ret, err := cli.HGetAll(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		} else {
			logger.CtxError(ctx, "GetGroupInfo HGetAll fail",
				zap.Error(err),
				zap.Any("key", key),
			)
			return nil, err
		}
	}
	group = new(Group)
	for key, value := range ret {
		if key == "info" {
			err = json.Unmarshal([]byte(value), &group)
		} else {
			var member Member
			err = json.Unmarshal([]byte(value), &member)
			if err == nil {
				group.Members = append(group.Members, member)
			}
		}
		if err != nil {
			logger.CtxError(ctx, "GetGroupInfo Unmarshal fail", zap.Error(err), zap.String("key", key), zap.String("value", value))
			return
		}

	}
	logger.CtxInfo(ctx, "GetGroupInfo success", zap.Any("key", key), zap.Any("group", group))
	return
}

// CreateGroup
func CreateGroup(ctx context.Context, appID int32, creator uint64, groupID int64, grouptype string, invitees []uint64) (group *Group, err error) {
	var (
		now    = time.Now()
		logger = fklog.ContextAppLogger(ctx)
	)
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "CreateGroup Client fail",
			zap.Error(err),
			zap.Int64("groupID", groupID),
		)
		return nil, err
	}
	key := getKey(cli, appID, groupID)
	values := make([]any, 0)
	group = &Group{
		ID:         groupID,
		Creator:    creator,
		CreateTime: now.Unix(),
		GroupType:  grouptype,
	}
	info, err := json.Marshal(group)
	if err != nil {
		logger.CtxError(ctx, "CreateGroup Marshal fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.Int64("groupID", groupID),
			zap.Any("group", group),
		)
		return nil, err
	} else {
		values = append(values, "info", info)
	}
	memberIDs := invitees
	if creator != 0 {
		memberIDs = append(memberIDs, creator)
	}
	for _, memberID := range memberIDs {
		member := Member{
			UserID:     memberID,
			CreateTime: now.Unix(),
		}
		data, err := json.Marshal(member)
		if err != nil {
			logger.CtxError(ctx, "CreateGroup Marshal fail",
				zap.Error(err),
				zap.Any("key", key),
				zap.Int64("groupID", groupID),
				zap.Uint64("memberID", memberID),
				zap.Any("member", member),
			)
			return nil, err
		} else {
			values = append(values, memberID, data)
		}
		group.Members = append(group.Members, member)
	}
	err = cli.HMSet(ctx, key, values...).Err()
	if err != nil {
		logger.CtxError(ctx, "CreateGroup HMSet fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.Int64("groupID", groupID),
			zap.Any("group", group),
		)
		return nil, err
	}
	logger.CtxInfo(ctx, "CreateGroup success", zap.Any("key", key), zap.Int64("groupID", groupID), zap.Any("group", group))
	return group, nil
}

// InviteMember
func InviteMember(ctx context.Context, appID int32, groupID int64, memberID uint64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "InviteMember Client fail",
			zap.Error(err),
			zap.Int64("groupID", groupID),
			zap.Uint64("memberID", memberID),
		)
		return err
	}
	key := getKey(cli, appID, groupID)
	member := Member{
		UserID:     memberID,
		CreateTime: time.Now().Unix(),
	}
	data, err := json.Marshal(member)
	if err != nil {
		logger.CtxError(ctx, "InviteMember Marshal fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.Int64("groupID", groupID),
			zap.Uint64("memberID", memberID),
			zap.Any("member", member),
		)
		return err
	}
	err = cli.HSet(ctx, key, memberID, data).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		} else {
			logger.CtxError(ctx, "InviteMember HDel fail",
				zap.Error(err),
				zap.Any("key", key),
				zap.Int64("groupID", groupID),
				zap.Uint64("memberID", memberID),
				zap.Any("member", member),
			)
			return err
		}
	}
	logger.CtxInfo(ctx, "InviteMember success", zap.Any("key", key), zap.Int64("groupID", groupID), zap.Uint64("memberID", memberID))
	return
}

// 是否为群成员
func IsMember(ctx context.Context, appID int32, groupID int64, memberID uint64) (is bool, err error) {
	logger := fklog.ContextAppLogger(ctx)
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.CtxError(ctx, "IsMember Client fail",
			zap.Error(err),
			zap.Int64("groupID", groupID),
			zap.Uint64("memberID", memberID),
		)
		return false, err
	}
	key := getKey(cli, appID, groupID)
	ret, err := cli.HExists(ctx, key, fmt.Sprintf("%d", memberID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		} else {
			logger.CtxError(ctx, "IsMember HExists fail",
				zap.Error(err),
				zap.Any("key", key),
				zap.Int64("groupID", groupID),
				zap.Uint64("memberID", memberID),
			)
			return false, err
		}
	}

	logger.CtxInfo(ctx, "IsMember success", zap.Any("key", key), zap.Int64("groupID", groupID), zap.Uint64("memberID", memberID), zap.Bool("is", is))
	return ret, nil
}

// // RemoveGroup
// func RemoveGroup(logger fklog.FKLogI, appID int32, groupID int32) (err error) {

// }
