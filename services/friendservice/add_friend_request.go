package friendservice

import (
	"context"
	"fmt"
	"maze_game_server/common/errors"
	"maze_game_server/model/friendmodel"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) AddFriendRequest(ctx context.Context, userId, toId uint64, from int32) (sendrq *friendmodel.SendFriendRequestInfo, err error) {
	logger := fklog.ContextAppLogger(ctx)
	// 好友列表
	friends, err := friendmodel.NewFriendModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "AddFriendRequest NewFriendModel err",
			zap.Uint64("userID", userId),
			zap.Uint64("toID", toId),
			zap.Error(err))
		return
	}

	sends, err := friendmodel.NewSendFriendRequestModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "AddFriendRequest NewSendFriendRequestModel err",
			zap.Uint64("userID", userId),
			zap.Uint64("toID", toId),
			zap.Error(err))
		return
	}

	// 1.是否在我的黑名单中
	inBlk, err := s.IsBlacklist(ctx, userId, toId)
	if err != nil {
		logger.CtxError(ctx, "AddFriendRequest IsBlacklist err", zap.Uint64("userID", userId),
			zap.Uint64("toID", toId),
			zap.Error(err))
		return
	}

	if inBlk {
		logger.CtxWarn(ctx, "AddFriendRequest isBlk", zap.String("err", "对方中你的黑名单中"))
		err = fmt.Errorf("SKIP")
		return
	}

	// 3.判断是否已经是好友了
	if isFriend := s.IsFriend(friends, toId); isFriend {
		logger.CtxWarn(ctx, "AddFriendRequest isBlk", zap.String("err", "对方已经是你的好友了"))
		err = fmt.Errorf("SKIP")
		return
	}
	// 4.检查重复发送
	if canSend := s.checkRepeatSendFriendRequest(sends, toId); !canSend {
		logger.CtxWarn(ctx, "AddFriendRequest isBlk", zap.String("err", "已经发送过好友请求了, 2分钟内不能重复发送"))
		err = fmt.Errorf("SKIP")
		return
	}

	// // 5.检测长度
	// if

	// 6.发送给对方
	if err = s.AddFriendRequestEvent(ctx, toId, userId, from); err != nil {
		logger.CtxError(ctx, "AddFriendRequest AddFriendRequestEvent err",
			zap.Uint64("userID", userId),
			zap.Uint64("toID", toId),
			zap.Error(err))
		return
	}

	// 7.设置已发送好友请求
	if sendrq, err = s.addSendFriendRequest(ctx, sends, userId, toId, from); err != nil {
		logger.CtxError(ctx, "AddFriendRequest addSendFriendRequest err",
			zap.Uint64("userID", userId),
			zap.Uint64("toID", toId),
			zap.Error(err))
		return
	}

	logger.CtxInfo(ctx, "AddFriendRequest Successful", zap.Any("userId", userId), zap.Any("toId", toId), zap.Any("from", from), zap.Any("sendrq", sendrq))
	return
}

// // 获取好友数量，好友数量包括已经成为好友的+发出去的申请
// func (s *service) GetUserFriendCount(logger fklog.FKLogI, userId uint64) (count int64, err error) {
// 	friendModel, err := friendmodel.NewFriendModel(logger, userId)
// 	if err != nil {
// 		logger.ErrorWF("GetUserFriendCount NewFriendModel err", zap.Error(err))
// 		return 0, err
// 	}
// 	// 查一下发出去的申请数量
// 	sendfriendRequestModel, err := friendmodel.NewSendFriendRequestModel(logger, userId)
// 	if err != nil {
// 		logger.ErrorWF("GetUserFriendCount NewSendFriendRequestModel err", zap.Error(err))
// 		return 0, err
// 	}
// 	return int64(len(friendModel.FriendList)) + int64(len(sendfriendRequestModel.SendList)), nil
// }

// 等待actor逻辑完成处理
// 收到好友请求事件
func (s *service) AddFriendRequestEvent(ctx context.Context, userId, fromId uint64, from int32) error {
	logger := fklog.ContextAppLogger(ctx)
	// 是否在我的黑名单中
	inBlk, err := s.IsBlacklist(ctx, userId, fromId)
	if err != nil {
		logger.CtxError(ctx, "AddFriendRequestEvent isBlacklist err", zap.Error(err))
		return err
	}
	if inBlk {
		return fmt.Errorf("已被对方拉黑")
	}

	// 设置待处理好友请求
	if err := s.addReceiveFriendRequest(ctx, userId, fromId, from); err != nil {
		return err
	}

	logger.CtxInfo(ctx, "AddFriendRequestEvent Successful", zap.Any("userId", userId), zap.Any("fromId", fromId), zap.Any("from", from))
	return nil
}

// 设置好友请求待处理
func (s *service) addReceiveFriendRequest(ctx context.Context, userId, fromId uint64, from int32) error {
	logger := fklog.ContextAppLogger(ctx)
	receiveModel, err := friendmodel.NewReceiveFriendRequestModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "addReceiveFriendRequest NewReceiveFriendRequestModel err", zap.Error(err), zap.Uint64("fromId", fromId))
		return errors.MODULE_ERROR
	}
	exist := false
	for _, i := range receiveModel.ReceiveList {
		if i.FromUserId == fromId {
			exist = true
		}
	}
	// 可以不报错
	if exist {
		logger.CtxWarn(ctx, "addReceiveFriendRequest 已经申请过了", zap.Uint64("fromId", fromId))
		return nil
	}
	receiveModel.ReceiveList = append(receiveModel.ReceiveList, &friendmodel.ReceiveFriendRequestInfo{
		FromUserId: fromId,
		CreateAt:   time.Now().UnixMilli(),
		Status:     friendmodel.FriendRequestStatusPending,
		From:       from,
	})

	err = receiveModel.Save(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "addReceiveFriendRequest SetReceiveFriendRequest err", zap.Error(err), zap.Uint64("fromId", fromId))
		return errors.MODULE_ERROR
	}

	logger.CtxInfo(ctx, "addReceiveFriendRequest Successful", zap.Any("userId", userId), zap.Any("fromId", fromId), zap.Any("from", from))
	return nil
}

// 检查重复发送 true可以发送， false 不可以发送
func (s *service) checkRepeatSendFriendRequest(sendModel *friendmodel.SendFriendRequestModel, toID uint64) bool {
	var send *friendmodel.SendFriendRequestInfo = nil
	index := -1
	for j, i := range sendModel.SendList {
		if i.ToUserId == toID {
			send = i
			index = j
		}
	}
	// 没有发送记录
	if send == nil {
		return true
	}
	// 30分钟内不可以重复发送
	sendTime := time.UnixMilli(send.CreateAt)
	if time.Since(sendTime) > 2*time.Minute {
		sendModel.SendList = append(sendModel.SendList[:index], sendModel.SendList[index+1:]...)
		return true
	}
	return false
}

// 设置好友请求
func (s *service) addSendFriendRequest(ctx context.Context, sendModel *friendmodel.SendFriendRequestModel, userId, toUserId uint64, from int32) (sendFriend *friendmodel.SendFriendRequestInfo, err error) {
	logger := fklog.ContextAppLogger(ctx)
	nowTime := time.Now()
	sendFriend = &friendmodel.SendFriendRequestInfo{
		ToUserId: toUserId,
		CreateAt: nowTime.UnixMilli(),
		From:     from,
	}
	sendModel.SendList = append(sendModel.SendList, sendFriend)

	err = sendModel.Save(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "addSendFriendRequest err", zap.Error(err))
		return
	}

	logger.CtxInfo(ctx, "addSendFriendRequest Successful", zap.Any("sendModel", sendModel), zap.Any("userId", userId), zap.Any("toUserId", toUserId), zap.Any("from", from), zap.Any("sendFriend", sendFriend))
	return
}

// 是否是好友
func (s *service) IsFriend(friends *friendmodel.FriendModel, toId uint64) bool {
	for _, b := range friends.FriendList {
		if b.UserId == toId {
			return true
		}
	}
	return false
}

func (s *service) CheckFriend(ctx context.Context, userID uint64, toID uint64) (bool, error) {
	logger := fklog.ContextAppLogger(ctx)
	// 好友列表
	friends, err := friendmodel.NewFriendModel(ctx, userID)
	logger.CtxDebug(ctx, "CheckFriend GetSuccesful", zap.Any("friends", friends))
	if err != nil {
		logger.CtxError(ctx, "CheckFriend NewFriendModel err",
			zap.Uint64("userID", userID),
			zap.Uint64("toID", toID),
			zap.Error(err))
		return false, err
	}

	for _, friend := range friends.FriendList {
		if friend.UserId == toID {
			return true, nil
		}
	}
	return false, nil
}
