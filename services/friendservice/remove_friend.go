package friendservice

import (
	"context"
	"fmt"
	"maze_game_server/model/friendmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) RemoveFriend(ctx context.Context, userId, toId uint64) (to *friendmodel.FriendInfo, in *friendmodel.FriendInfo, err error) {
	logger := fklog.ContextAppLogger(ctx)
	// 好友列表
	friendModel, err := friendmodel.NewFriendModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "RemoveFriend GetFriends err", zap.Error(err))
		return
	}
	if isFriend := s.IsFriend(friendModel, toId); !isFriend {
		logger.CtxInfo(ctx, "RemoveFriend 已经不是好友了", zap.Uint64("toId", toId))
		err = fmt.Errorf("已经不是好友了")
		return
	}

	to = s.getFriendInfo(ctx, friendModel, toId)

	// 通知对方删好友 todo异步
	if in, err = s.RemoveFriendEvent(ctx, toId, userId); err != nil {
		logger.CtxError(ctx, "RemoveFriendEvent err", zap.Error(err))
		return
	}

	if err = s.delFriendAndFriendRequest(ctx, friendModel, userId, toId); err != nil {
		logger.CtxError(ctx, "RemoveFriend delFriend err", zap.Error(err))
		return
	}

	return
}

func (s *service) getFriendInfo(ctx context.Context, friendModel *friendmodel.FriendModel, toId uint64) (to *friendmodel.FriendInfo) {
	for _, friendInfo := range friendModel.FriendList {
		if friendInfo.UserId == toId {
			to = friendInfo
			return
		}
	}
	return
}

// 删除好友及好友相关的申请
func (s *service) delFriendAndFriendRequest(ctx context.Context, friendModel *friendmodel.FriendModel, userId, toId uint64) error {
	logger := fklog.ContextAppLogger(ctx)
	var index = -1
	for i, j := range friendModel.FriendList {
		if j.UserId == toId {
			index = i
			break
		}
	}
	if index == -1 {
		logger.CtxWarn(ctx, "delFriendAndFriendRequest toId not found", zap.Uint64("toId", toId))
		return nil
	}
	friendModel.FriendList = append(friendModel.FriendList[:index], friendModel.FriendList[index+1:]...)
	if err := friendModel.Save(ctx, userId); err != nil {
		logger.CtxError(ctx, "delFriendAndFriendRequest SetFriends err", zap.Error(err), zap.Uint64("toId", toId))
		return err
	}
	// 失败了也没关系，好友删除了就行
	s.delReceiveFriendRequest(ctx, userId, toId)
	s.delSendFriendRequest(ctx, userId, toId)
	logger.CtxInfo(ctx, "delFriendAndFriendRequest success", zap.Any("toId", toId))
	return nil
}

// 被删除好友事件
func (s *service) RemoveFriendEvent(ctx context.Context, userId, fromId uint64) (in *friendmodel.FriendInfo, err error) {
	logger := fklog.ContextAppLogger(ctx)
	// 好友列表
	friends, err := friendmodel.NewFriendModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "RemoveFriendEvent GetFriends error", zap.Error(err))
		return
	}
	if isFriend := s.IsFriend(friends, fromId); !isFriend {
		logger.CtxWarn(ctx, "RemoveFriendEvent 已经不是好友了", zap.Uint64("fromId", fromId))
		return
	}

	in = s.getFriendInfo(ctx, friends, fromId)

	if err = s.delFriendAndFriendRequest(ctx, friends, userId, fromId); err != nil {
		logger.CtxError(ctx, "RemoveFriendEvent delFriend err", zap.Error(err), zap.Uint64("fromId", fromId))
		return
	}
	logger.CtxInfo(ctx, "RemoveFriendEvent success", zap.Uint64("userId", userId), zap.Any("fromId", fromId))
	return
}

func (s *service) delReceiveFriendRequest(ctx context.Context, userId, fromUserId uint64) error {
	logger := fklog.ContextAppLogger(ctx)
	receives, err := friendmodel.NewReceiveFriendRequestModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "delReceiveFriendRequest GetReceiveFriendRequest failed", zap.Error(err))
		return err
	}
	if err = s.delReceiveFriendRequestNotGet(ctx, receives, userId, fromUserId); err != nil {
		return err
	}
	logger.CtxInfo(ctx, "delReceiveFriendRequest success", zap.Any("fromUserId", fromUserId))
	return nil
}

func (s *service) delReceiveFriendRequestNotGet(ctx context.Context, receiveModel *friendmodel.ReceiveFriendRequestModel, userId, fromUserId uint64) error {
	logger := fklog.ContextAppLogger(ctx)
	var index = -1
	for i, j := range receiveModel.ReceiveList {
		if j.FromUserId == fromUserId {
			index = i
			break
		}
	}
	if index == -1 {
		logger.CtxInfo(ctx, "delReceiveFriendRequestNotGet fromUserId not found", zap.Uint64("fromUserId", fromUserId))
		return nil
	}
	receiveModel.ReceiveList = append(receiveModel.ReceiveList[:index], receiveModel.ReceiveList[index+1:]...)
	if err := receiveModel.Save(ctx, userId); err != nil {
		logger.CtxError(ctx, "delReceiveFriendRequestNotGet SetReceiveFriendRequest err", zap.Error(err), zap.Uint64("fromUserId", fromUserId))
		return err
	}
	logger.CtxInfo(ctx, "delReceiveFriendRequestNotGet success", zap.Any("fromUserId", fromUserId))
	return nil
}

func (s *service) delSendFriendRequest(ctx context.Context, userId, toId uint64) error {
	logger := fklog.ContextAppLogger(ctx)
	sendModel, err := friendmodel.NewSendFriendRequestModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "delSendFriendRequest GetSendFriendRequest failed", zap.Error(err))
		return err
	}
	if err = s.delSendFriendRequestNotGet(ctx, sendModel, userId, toId); err != nil {
		return err
	}

	return nil
}

func (s *service) delSendFriendRequestNotGet(ctx context.Context, sendModel *friendmodel.SendFriendRequestModel, userId, toId uint64) error {
	logger := fklog.ContextAppLogger(ctx)
	var index = -1
	for i, j := range sendModel.SendList {
		if j.ToUserId == toId {
			index = i
			break
		}
	}
	if index == -1 {
		logger.CtxInfo(ctx, "delSendFriendRequest toId not found", zap.Uint64("toId", toId))
		return nil
	}
	sendModel.SendList = append(sendModel.SendList[:index], sendModel.SendList[index+1:]...)
	if err := sendModel.Save(ctx, userId); err != nil {
		logger.CtxError(ctx, "delSendFriendRequest SetSendFriendRequest err", zap.Error(err), zap.Uint64("toId", toId))
		return err
	}
	logger.CtxInfo(ctx, "delSendFriendRequest success", zap.Any("toId", toId))
	return nil
}

// // 删除用户全部的好友, 发送的申请, 收到的申请, 黑名单
// func (s *service) DeleteUserAll(logger fklog.FKLogI, userId uint64) error {
// 	var receiveModel friendmodel.ReceiveFriendRequestModel
// 	var sendModel friendmodel.SendFriendRequestModel
// 	var blacklistModel friendmodel.BlacklistModel
// 	var friendModel friendmodel.FriendModel

// 	err := friendModel.Del(logger, userId)
// 	if err != nil {
// 		return err
// 	}
// 	if err = receiveModel.Del(logger, userId); err != nil {
// 		return err
// 	}
// 	if err = sendModel.Del(logger, userId); err != nil {
// 		return err
// 	}
// 	if err = blacklistModel.Del(logger, userId); err != nil {
// 		return err
// 	}
// 	logger.InfoWF("DeleteUserAll success", zap.Uint64("userId", userId))
// 	return nil
// }
