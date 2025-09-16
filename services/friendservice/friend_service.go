package friendservice

import (
	"context"
	"maze_game_server/model/friendmodel"
)

// 外部系统可调用方法
type FriendService interface {
	// 好友请求
	AddFriendRequest(ctx context.Context, userId, toId uint64, from int32) (sendrq *friendmodel.SendFriendRequestInfo, err error)
	// // 好友列表
	FriendList(ctx context.Context, userId uint64, page, pageSize int32) ([]*friendmodel.FriendInfo, bool, error)
	// // 添加到黑名单   会把好友和好友申请记录都删除
	AddBlacklist(ctx context.Context, userId, toId uint64) error
	// 移除黑名单
	RemoveBlacklist(ctx context.Context, userID, toId uint64) error
	// 移除好友
	RemoveFriend(ctx context.Context, userId, toId uint64) (to *friendmodel.FriendInfo, in *friendmodel.FriendInfo, err error)
	// // 收到的好友请求列表  申请列表30天清除
	ReceiveFriendRequestList(ctx context.Context, userId uint64, page, pageSize int32) ([]*friendmodel.ReceiveFriendRequestInfo, bool, error)
	// 发送的好友请求列表  发送列表30天清除
	SendFriendRequestList(ctx context.Context, userId uint64, page, pageSize int32) ([]*friendmodel.SendFriendRequestInfo, bool, error)
	// 黑名单列表
	Blacklist(ctx context.Context, userId uint64, page, pageSize int32) ([]*friendmodel.BlacklistInfo, bool, error)
	// 批量同意好友请求
	AgreeFriendApply(ctx context.Context, userID uint64, toID []int64) (rs []*friendmodel.ReceiveFriendRequestInfo, change []*friendmodel.ReceiveFriendRequestInfo, isSkip bool, err error)
	// 批量拒绝好友请求
	RefuseFriendApply(ctx context.Context, userID uint64, toID []int64) (rs []*friendmodel.ReceiveFriendRequestInfo, change []*friendmodel.ReceiveFriendRequestInfo, isSkip bool, err error)
	// 好友推荐
	FriendRecommend(ctx context.Context, userID uint64, pageSize int32) ([]uint64, error)
	// 检测是否是好友
	CheckFriend(ctx context.Context, userID uint64, toID uint64) (bool, error)
	// 检测是否是黑名单
	IsBlacklist(ctx context.Context, userID uint64, toID uint64) (bool, error)
}

// GlobalFriendService 好友服务可用全局唯一对象
var GlobalFriendService FriendService

func init() {
	GlobalFriendService = newFriendService()
}

type service struct {
}

func newFriendService() FriendService {
	return &service{}
}
