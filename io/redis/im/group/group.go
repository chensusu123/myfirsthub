package group

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	globalredis "maze_game_server/io/redis"

	"github.com/redis/go-redis/v9"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

type Group struct {
	ID         int32    `json:"id,omitempty"`
	Creator    uint64   `json:"creator,omitempty"`
	CreateTime int64    `json:"create_time,omitempty"`
	Members    []Member `json:"-"`
}

type Member struct {
	CreateTime int64  `json:"create_time,omitempty"`
	UserID     uint64 `json:"user_id,omitempty"`
	Nickname   string `json:"nickname,omitempty"`
}

// getKey 获取缓存操作key。
func getKey(args ...interface{}) string {
	return fmt.Sprintf("im:app:%d:group:%d:", args...)
}

// GetGroupInfo
func GetGroupInfo(logger fklog.FKLogI, appID int32, groupID int32) (group *Group, err error) {
	var (
		key = getKey(appID, groupID)
		ctx = context.Background()
	)
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("GetGroupInfo Client fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.Int32("groupID", groupID),
		)
		return nil, err
	}
	ret, err := cli.HGetAll(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = nil
		} else {
			logger.ErrorWF("GetGroupInfo HGetAll fail",
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
			logger.ErrorWF("GetGroupInfo Unmarshal fail", zap.Error(err), zap.String("key", key), zap.String("value", value))
			return
		}

	}
	logger.DebugWF("GetGroupInfo success", zap.Any("key", key), zap.Any("group", group))
	return
}

// CreateGroup
func CreateGroup(logger fklog.FKLogI, appID int32, creator uint64, groupID int32, invitees []uint64) (group *Group, err error) {
	var (
		now = time.Now()
		key = getKey(appID, groupID)
		ctx = context.Background()
	)
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("CreateGroup Client fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.Int32("groupID", groupID),
		)
		return nil, err
	}
	values := make([]any, 0)
	group = &Group{
		ID:         groupID,
		Creator:    creator,
		CreateTime: now.Unix(),
	}
	info, err := json.Marshal(group)
	if err != nil {
		logger.ErrorWF("CreateGroup Marshal fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.Int32("groupID", groupID),
			zap.Any("group", group),
		)
		return nil, err
	} else {
		values = append(values, "info", info)
	}
	memberIDs := append([]uint64{creator}, invitees...)
	for _, memberID := range memberIDs {
		member := Member{
			UserID:     memberID,
			CreateTime: now.Unix(),
		}
		data, err := json.Marshal(member)
		if err != nil {
			logger.ErrorWF("CreateGroup Marshal fail",
				zap.Error(err),
				zap.Any("key", key),
				zap.Int32("groupID", groupID),
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
		logger.ErrorWF("CreateGroup HMSet fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.Int32("groupID", groupID),
			zap.Any("group", group),
		)
		return nil, err
	}
	logger.DebugWF("CreateGroup success", zap.Any("key", key), zap.Int32("groupID", groupID), zap.Any("group", group))
	return group, nil
}

// InviteMember
func InviteMember(logger fklog.FKLogI, appID int32, groupID int32, memberID uint64) (err error) {
	var (
		key = getKey(appID, groupID)
		ctx = context.Background()
	)
	cli, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("InviteMember Client fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.Int32("groupID", groupID),
			zap.Uint64("memberID", memberID),
		)
		return err
	}
	member := Member{
		UserID:     memberID,
		CreateTime: time.Now().Unix(),
	}
	data, err := json.Marshal(member)
	if err != nil {
		logger.ErrorWF("InviteMember Marshal fail",
			zap.Error(err),
			zap.Any("key", key),
			zap.Int32("groupID", groupID),
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
			logger.ErrorWF("InviteMember HDel fail",
				zap.Error(err),
				zap.Any("key", key),
				zap.Int32("groupID", groupID),
				zap.Uint64("memberID", memberID),
				zap.Any("member", member),
			)
			return err
		}
	}
	logger.DebugWF("InviteMember success", zap.Any("key", key), zap.Int32("groupID", groupID), zap.Uint64("memberID", memberID))
	return
}

// // RemoveGroup
// func RemoveGroup(logger fklog.FKLogI, appID int32, groupID int32) (err error) {

// }
