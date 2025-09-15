package friendmodel

import (
	"context"
	"fmt"
	"maze_game_server/io"
)

// 已发送好友请求key
var KeySendFriendRequest = "friend:send:%d"

func getKeySendFriendRequest(userId uint64) string {
	return fmt.Sprintf(KeySendFriendRequest, userId)
}

// 已发送的好友请求
type SendFriendRequestInfo struct {
	ToUserId uint64 `json:"to_user_id,omitempty"`
	CreateAt int64  `json:"create_at,omitempty"`
	Status   int32  `json:"status,omitempty"`
	From     int32  `json:"from,omitempty"`
}

type SendFriendRequestModel struct {
	SendList []*SendFriendRequestInfo `json:"send_list,omitempty"`
}

func NewSendFriendRequestModel(ctx context.Context, userID uint64) (*SendFriendRequestModel, error) {
	sendModel := &SendFriendRequestModel{}
	if err := sendModel.load(ctx, userID); err != nil {
		return nil, err
	}
	return sendModel, nil
}

func (f *SendFriendRequestModel) load(ctx context.Context, userId uint64) (err error) {
	return io.LoadSvrData(ctx, getKeySendFriendRequest(userId), f)
}

func (f *SendFriendRequestModel) Save(ctx context.Context, userId uint64) (err error) {
	return io.SaveSvrData(ctx, getKeySendFriendRequest(userId), f)
}

// 删除所有发送的好友申请
func (f *SendFriendRequestModel) Del(ctx context.Context, userId uint64) (err error) {
	return io.DeleteSvrData(ctx, getKeySendFriendRequest(userId))
}
