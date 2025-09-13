package friendmodel

import (
	"context"
	"maze_game_server/io/redis/friendredis"
	"maze_game_server/lib/serialize"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
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
	logger := fklog.ContextAppLogger(ctx)
	bytes, err := friendredis.GetReceiveFriendRequest(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "load receive redis err", zap.Error(err))
		return nil
	}
	if len(bytes) == 0 {
		f.ReceiveList = make([]*ReceiveFriendRequestInfo, 0)
		return nil
	}
	err = serialize.Unmarshal(bytes, f)
	if err != nil {
		logger.CtxError(ctx, "load receive unmarshal err", zap.Error(err), zap.Any("bytes", string(bytes)))
		return err
	}

	return nil
}

func (f *ReceiveFriendRequestModel) Save(logger fklog.FKLogI, userId uint64) (err error) {
	bytes, err := serialize.Marshal(f)
	if err != nil {
		logger.ErrorWF("save receive marshal failed", zap.Error(err), zap.Any("friends", f))
		return err
	}
	if err = friendredis.SetReceiveFriendRequest(logger, userId, bytes); err != nil {
		logger.ErrorWF("save receive marshal failed", zap.Error(err), zap.Any("friends", f))
		return err
	}
	logger.InfoWF("save receive success", zap.Any("friends", f))
	return nil
}

// 删除所有收到的好友申请
func (f *ReceiveFriendRequestModel) Del(logger fklog.FKLogI, userId uint64) (err error) {
	if err = friendredis.DelReceiveFriendRequest(logger, userId); err != nil {
		logger.ErrorWF("del err", zap.Error(err))
		return err
	}
	logger.InfoWF("del receives success")
	return nil
}
