package friendmodel

import (
	"context"
	"fmt"
	"maze_game_server/io"
)

var (
	ExpireTime    = 7 * 24 * 60 * 60 * 1000
	MaxFriendSize = 500
	RecommendSize = int32(5)
)

// 好友key
var KeyFriends = "friend:list:%d"

// 黑名单key
func getKeyFriends(userId uint64) string {
	return fmt.Sprintf(KeyFriends, userId)
}

// 好友信息
type FriendInfo struct {
	UserId     uint64 `json:"user_id,omitempty"`
	FriendType int32  `json:"friend_type"`         // 好友类型
	CreateAt   int64  `json:"create_at,omitempty"` // 添加时间
}
type FriendModel struct {
	FriendList []*FriendInfo `json:"friend_list,omitempty"`
}

func NewFriendModel(ctx context.Context, userID uint64) (*FriendModel, error) {
	friendModel := &FriendModel{}
	if err := friendModel.load(ctx, userID); err != nil {
		return nil, err
	}
	return friendModel, nil
}

func (f *FriendModel) load(ctx context.Context, userId uint64) (err error) {
	return io.LoadSvrData(ctx, getKeyFriends(userId), f)
}

func (f *FriendModel) Save(ctx context.Context, userId uint64) (err error) {
	return io.SaveSvrData(ctx, getKeyFriends(userId), f)
}

// todo 删除接口暴露出来
func (f *FriendModel) Del(ctx context.Context, userId uint64) (err error) {
	return io.SaveSvrData(ctx, getKeyFriends(userId), f)
}
