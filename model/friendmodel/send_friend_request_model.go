package friendmodel

import (
	"context"
	"maze_game_server/io/redis/friendredis"
	"maze_game_server/lib/serialize"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 已发送的好友请求
type SendFriendRequestInfo struct {
	ToUserId uint64 `json:"to_user_id,omitempty"`
	CreateAt int64  `json:"create_at,omitempty"`
	Status   int32  `json:"status,omitempty"`
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
	logger := fklog.ContextAppLogger(ctx)
	bytes, err := friendredis.GetSendFriendRequest(logger, userId)
	if err != nil {
		logger.CtxError(ctx, "load send redis err", zap.Error(err))
		return nil
	}
	if len(bytes) == 0 {
		f.SendList = make([]*SendFriendRequestInfo, 0)
		return nil
	}
	err = serialize.Unmarshal(bytes, f)
	if err != nil {
		logger.CtxError(ctx, "load send unmarshal err", zap.Error(err), zap.Any("bytes", string(bytes)))
		return err
	}

	return nil
}

func (f *SendFriendRequestModel) Save(logger fklog.FKLogI, userId uint64) (err error) {
	bytes, err := serialize.Marshal(f)
	if err != nil {
		logger.ErrorWF("save send marshal err", zap.Error(err), zap.Any("friends", f))
		return err
	}
	if err = friendredis.SetSendFriendRequest(logger, userId, bytes); err != nil {
		logger.ErrorWF("save send err", zap.Error(err), zap.Any("friends", f))
		return err
	}
	logger.InfoWF("save send success", zap.Any("friends", f))
	return nil
}

// 删除所有发送的好友申请
func (f *SendFriendRequestModel) Del(logger fklog.FKLogI, userId uint64) (err error) {
	if err = friendredis.DelSendFriendRequest(logger, userId); err != nil {
		logger.ErrorWF("del send err", zap.Error(err))
		return err
	}
	logger.InfoWF("del send success")
	return nil
}
