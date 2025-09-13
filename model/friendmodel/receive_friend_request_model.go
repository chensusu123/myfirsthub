package friendmodel

import (
	"context"
	"fmt"
	"maze_game_server/io"
)

// 待处理的好友申请
type ReceiveFriendRequestInfo struct {
	FromUserId uint64 `json:"from_user_id,omitempty"`
	CreateAt   int64  `json:"create_at,omitempty"`
	Status     int32  `json:"status,omitempty"`
	From       int32  `json:"from,omitempty"`
}

type ReceiveFriendRequestModel struct {
	ReceiveList []*ReceiveFriendRequestInfo `json:"receive_list,omitempty"`
}

// 待处理好友请求key
var (
	KeyReceiveFriendRequest = "friend:receive:%d"
	ReceivePage             = 10
)

func getKeyReceiveFriendRequest(userId uint64) string {
	return fmt.Sprintf(KeyReceiveFriendRequest, userId)
}

const (
	FriendRequestStatusPending  int32 = 0 // 待处理好友请求
	FriendRequestStatusAccepted int32 = 1 // 已经同意
	FriendRequestStatusRejected int32 = 2 // 已经拒绝
)

func NewReceiveFriendRequestModel(ctx context.Context, userID uint64) (*ReceiveFriendRequestModel, error) {
	receiveModel := &ReceiveFriendRequestModel{}
	if err := receiveModel.load(ctx, userID); err != nil {
		return nil, err
	}
	return receiveModel, nil
}

func (f *ReceiveFriendRequestModel) load(ctx context.Context, userId uint64) (err error) {
	return io.LoadSvrData(ctx, getKeyReceiveFriendRequest(userId), f)
}

func (f *ReceiveFriendRequestModel) Save(ctx context.Context, userId uint64) (err error) {
	return io.SaveSvrData(ctx, getKeyReceiveFriendRequest(userId), f)
}

// 删除所有收到的好友申请
func (f *ReceiveFriendRequestModel) Del(ctx context.Context, userId uint64) (err error) {
	return io.DeleteSvrData(ctx, getKeyReceiveFriendRequest(userId))
}
