package friendmodel

import (
	"context"
	"fmt"
	"maze_game_server/io"
)

// 黑名单key
var KeyBlacklist = "blacklist:%d"

func getKeyBlacklist(userId uint64) string {
	return fmt.Sprintf(KeyBlacklist, userId)
}

// 黑名单信息
type BlacklistInfo struct {
	UserId   uint64 `json:"user_id,omitempty"`
	CreateAt int64  `json:"create_at,omitempty"`
}
type BlacklistModel struct {
	Blacklist []*BlacklistInfo `json:"blacklist,omitempty"`
}

func NewBlacklistModel(ctx context.Context, userID uint64) (*BlacklistModel, error) {
	sendModel := &BlacklistModel{}
	if err := sendModel.load(ctx, userID); err != nil {
		return nil, err
	}
	return sendModel, nil
}

func (f *BlacklistModel) load(ctx context.Context, userId uint64) (err error) {
	return io.LoadSvrData(ctx, getKeyBlacklist(userId), f)
}

func (f *BlacklistModel) Save(ctx context.Context, userId uint64) (err error) {
	return io.SaveSvrData(ctx, getKeyBlacklist(userId), f)
}

// 删除所有发送的好友申请
func (f *BlacklistModel) Del(ctx context.Context, userId uint64) (err error) {
	return io.DeleteSvrData(ctx, getKeyBlacklist(userId))
}
