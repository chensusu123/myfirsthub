package friendservice

import (
	"context"
	"maze_game_server/model/friendmodel"
)

// 外部系统可调用方法
type FriendService interface {
	// 好友请求
	AddFriendRequest(ctx context.Context, userId, toId uint64, from int32) (int64, error)
	// // // 同意好友请求   同意后会把双方的申请记录删除
	// // AcceptFriendRequest(logger fklog.FKLogI, userId, toID uint64) *errors.CodeError
	// // 拒绝好友请求   拒绝后能看见拒绝信息
	// RejectFriendRequest(logger fklog.FKLogI, userId, toID uint64) *errors.CodeError
	// // 好友列表
	FriendList(ctx context.Context, userId uint64, page, pageSize int32) ([]*friendmodel.FriendInfo, bool, error)
	// // 添加到黑名单   会把好友和好友申请记录都删除
	// AddBlacklist(logger fklog.FKLogI, userID, toUserID uint64) *errors.CodeError
	// // 移除黑名单
	// RemoveBlacklist(logger fklog.FKLogI, userID, toId uint64) *errors.CodeError
	// 移除好友
	RemoveFriend(ctx context.Context, userId, toId uint64) (to *friendmodel.FriendInfo, in *friendmodel.FriendInfo, err error)
	// // 收到的好友请求列表  申请列表30天清除
	ReceiveFriendRequestList(ctx context.Context, userId uint64, page, pageSize int32) ([]*friendmodel.ReceiveFriendRequestInfo, bool, error)
	// // 发送的好友请求列表  发送列表30天清除
	// SendFriendRequestList(logger fklog.FKLogI, userId uint64, page, pageSize int32) ([]*friendmodel.SendFriendRequestInfo, *errors.CodeError)
	// // 黑名单列表
	// Blacklist(logger fklog.FKLogI, userId uint64, page, pageSize int32) ([]*friendmodel.BlacklistInfo, *errors.CodeError)
	// // 删除用户全部的好友, 发送的申请, 收到的申请, 黑名单
	// DeleteUserAll(logger fklog.FKLogI, userId uint64) error

	// 批量同意好友请求
	AgreeFriendApply(ctx context.Context, userID uint64, toID []int64) (rs []*friendmodel.ReceiveFriendRequestInfo, err error)
	// 批量拒绝好友请求
	RefuseFriendApply(ctx context.Context, userID uint64, toID []int64) (rs []*friendmodel.ReceiveFriendRequestInfo, err error)

	// 好友推荐
	FriendRecommend(ctx context.Context, userID uint64, pageSize int32) ([]uint64, error)
	// 检测是否是好友
	CheckFriend(ctx context.Context, userID uint64, toID uint64) (bool, error)
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
