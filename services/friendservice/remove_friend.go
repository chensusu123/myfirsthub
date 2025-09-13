package friendservice

import (
	"context"
	"maze_game_server/model/friendmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// import (
// 	"context"
// 	"maze_game_server/common/errors"
// 	"maze_game_server/model/friendmodel"

// 	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
// 	"go.uber.org/zap"
// )

// func (s *service) RemoveFriend(logger fklog.FKLogI, userId, toId uint64) *errors.CodeError {
// 	// 好友列表
// 	friendModel, err := friendmodel.NewFriendModel(logger, userId)
// 	if err != nil {
// 		logger.ErrorWF("RemoveFriend GetFriends err", zap.Error(err))
// 		return errors.MODULE_ERROR
// 	}
// 	if isFriend := s.IsFriend(friendModel, toId); !isFriend {
// 		logger.InfoWF("RemoveFriend 已经不是好友了", zap.Uint64("toId", toId))
// 		return errors.MODULE_ERROR.WrapMsg("已经不是好友了")
// 	}

// 	// 通知对方删好友 todo异步
// 	if codeErr := s.RemoveFriendEvent(logger, toId, userId); codeErr != nil {
// 		logger.ErrorWF("RemoveFriendEvent err", zap.Error(codeErr))
// 		return codeErr
// 	}

// 	if err = s.delFriendAndFriendRequest(logger, friendModel, userId, toId); err != nil {
// 		logger.ErrorWF("RemoveFriend delFriend err", zap.Error(err))
// 		return errors.MODULE_ERROR
// 	}

// 	return nil
// }

// // 删除好友及好友相关的申请
// func (s *service) delFriendAndFriendRequest(logger fklog.FKLogI, friendModel *friendmodel.FriendModel, userId, toId uint64) error {
// 	var index = -1
// 	for i, j := range friendModel.FriendList {
// 		if j.UserId == toId {
// 			index = i
// 			break
// 		}
// 	}
// 	if index == -1 {
// 		logger.WarnWF("delFriendAndFriendRequest toId not found", zap.Uint64("toId", toId))
// 		return nil
// 	}
// 	friendModel.FriendList = append(friendModel.FriendList[:index], friendModel.FriendList[index+1:]...)
// 	if err := friendModel.Save(logger, userId); err != nil {
// 		logger.ErrorWF("delFriendAndFriendRequest SetFriends err", zap.Error(err), zap.Uint64("toId", toId))
// 		return err
// 	}
// 	// 失败了也没关系，好友删除了就行
// 	s.delReceiveFriendRequest(logger, userId, toId)
// 	s.delSendFriendRequest(logger, userId, toId)
// 	logger.InfoWF("delFriendAndFriendRequest success", zap.Any("toId", toId))
// 	return nil
// }

// // 被删除好友事件
// func (s *service) RemoveFriendEvent(logger fklog.FKLogI, userId, fromId uint64) *errors.CodeError {
// 	// 好友列表
// 	friends, err := friendmodel.NewFriendModel(logger, userId)
// 	if err != nil {
// 		logger.ErrorWF("RemoveFriendEvent GetFriends error", zap.Error(err))
// 		return errors.MODULE_ERROR
// 	}
// 	if isFriend := s.IsFriend(friends, fromId); !isFriend {
// 		logger.WarnWF("RemoveFriendEvent 已经不是好友了", zap.Uint64("fromId", fromId))
// 		return nil
// 	}

// 	if err = s.delFriendAndFriendRequest(logger, friends, userId, fromId); err != nil {
// 		logger.ErrorWF("RemoveFriendEvent delFriend err", zap.Error(err), zap.Uint64("fromId", fromId))
// 		return errors.MODULE_ERROR
// 	}
// 	logger.InfoWF("RemoveFriendEvent success", zap.Uint64("userId", userId), zap.Any("fromId", fromId))
// 	return nil
// }

// func (s *service) delReceiveFriendRequest(ctx context.Context, userId, fromUserId uint64) error {
// 	logger := fklog.ContextAppLogger(ctx)
// 	receives, err := friendmodel.NewReceiveFriendRequestModel(ctx, userId)
// 	if err != nil {
// 		logger.ErrorWF("delReceiveFriendRequest GetReceiveFriendRequest failed", zap.Error(err))
// 		return err
// 	}
// 	if err = s.delReceiveFriendRequestNotGet(logger, receives, userId, fromUserId); err != nil {
// 		return err
// 	}
// 	logger.InfoWF("delReceiveFriendRequest success", zap.Any("fromUserId", fromUserId))
// 	return nil
// }

func (s *service) delReceiveFriendRequestNotGet(logger fklog.FKLogI, receiveModel *friendmodel.ReceiveFriendRequestModel, userId, fromUserId uint64) error {
	var index = -1
	for i, j := range receiveModel.ReceiveList {
		if j.FromUserId == fromUserId {
			index = i
			break
		}
	}
	if index == -1 {
		logger.InfoWF("delReceiveFriendRequestNotGet fromUserId not found", zap.Uint64("fromUserId", fromUserId))
		return nil
	}
	receiveModel.ReceiveList = append(receiveModel.ReceiveList[:index], receiveModel.ReceiveList[index+1:]...)
	if err := receiveModel.Save(logger, userId); err != nil {
		logger.ErrorWF("delReceiveFriendRequestNotGet SetReceiveFriendRequest err", zap.Error(err), zap.Uint64("fromUserId", fromUserId))
		return err
	}
	logger.InfoWF("delReceiveFriendRequestNotGet success", zap.Any("fromUserId", fromUserId))
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
	if err := sendModel.Save(logger, userId); err != nil {
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
