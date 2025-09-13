package friendservice

import (
	"context"
	"maze_game_server/common/errors"
	"maze_game_server/model/friendmodel"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// func (s *service) AcceptFriendRequest(logger fklog.FKLogI, userId, toId uint64) *errors.CodeError {
// 	// 检查有没有收到好友申请
// 	receiveModel, err := friendmodel.NewReceiveFriendRequestModel(logger, userId)
// 	if err != nil {
// 		logger.ErrorWF("AcceptFriendRequest NewReceiveFriendRequestModel err", zap.Error(err))
// 		return errors.MODULE_ERROR
// 	}
// 	index := -1
// 	var receive *friendmodel.ReceiveFriendRequestInfo
// 	for j, i := range receiveModel.ReceiveList {
// 		if i.FromUserId == toId {
// 			index = j
// 			receive = i
// 			break
// 		}
// 	}
// 	// 1. 检查申请记录
// 	if receive == nil {
// 		logger.ErrorWF("AcceptFriendRequest failed 没有好友申请记录")
// 		return errors.COMMON_ERROR_TIPS.WrapMsg("没有好友申请记录")
// 	}
// 	if receive.Status != friendmodel.FriendRequestStatusPending {
// 		logger.ErrorWF("AcceptFriendRequest failed 好友请求已经处理过了")
// 		return errors.COMMON_ERROR_TIPS.WrapMsg("好友请求已经处理过了")
// 	}
// 	// 好友列表
// 	friendModel, err := friendmodel.NewFriendModel(logger, userId)
// 	if err != nil {
// 		logger.ErrorWF("AcceptFriendRequest NewFriendModel err", zap.Error(err))
// 		return errors.MODULE_ERROR
// 	}
// 	// 2.检查好友
// 	if isFriend := s.IsFriend(friendModel, toId); isFriend {
// 		// 把对方的发送也删掉
// 		s.delSendFriendRequestEvent(logger, toId, userId)
// 		// 忽略错误 已经是好友了就把收到的申请删掉
// 		s.delReceiveFriendRequestNotGet(logger, receiveModel, userId, toId)
// 		return errors.COMMON_ERROR_TIPS.WrapMsg("对方已经是你的好友了")
// 	}
// 	// 查一下发出去的申请数量
// 	sendModel, err := friendmodel.NewSendFriendRequestModel(logger, userId)
// 	if err != nil {
// 		logger.ErrorWF("AddFriendRequest GetSendFriendRequest err", zap.Error(err))
// 		return errors.MODULE_ERROR
// 	}
// 	// 3.检查好友数量
// 	friendCount := int64(len(friendModel.FriendList)) + int64(len(sendModel.SendList))
// 	// todo 读取好友数量配置
// 	if friendCount >= 100 {
// 		logger.InfoWF("AcceptFriendRequest 好友数量达到上限了", zap.Int64("friendCount", friendCount))
// 		return errors.COMMON_ERROR_TIPS.WrapMsg("好友数量达到上限了")
// 	}
// 	// 4.通知同意好友申请 todo
// 	if codeErr := s.AcceptFriendRequestEvent(logger, toId, userId); codeErr != nil {
// 		logger.ErrorWF("AcceptFriendRequestEvent err", zap.Error(codeErr), zap.Uint64("toId", toId))
// 		return codeErr
// 	}
// 	// 5.设置好友
// 	if err = s.addFriend(logger, friendModel, userId, toId); err != nil {
// 		logger.ErrorWF("AcceptFriendRequest addFriend err", zap.Error(err))
// 		return errors.MODULE_ERROR
// 	}
// 	// 6.删除收到的申请
// 	receiveModel.ReceiveList = append(receiveModel.ReceiveList[:index], receiveModel.ReceiveList[index+1:]...)
// 	// 7.更新待处理记录
// 	err = receiveModel.Save(logger, userId)
// 	if err != nil {
// 		logger.ErrorWF("AcceptFriendRequest SetReceiveFriendRequest err", zap.Error(err))
// 		return errors.MODULE_ERROR
// 	}
// 	return nil
// }

func (s *service) AgreeFriendApply(ctx context.Context, userID uint64, toID []int64) (rs []*friendmodel.ReceiveFriendRequestInfo, err error) {
	logger := fklog.ContextAppLogger(ctx)
	// 好友列表
	friendModel, err := friendmodel.NewFriendModel(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "AgreeFriendApply NewFriendModel err",
			zap.Uint64("userID", userID),
			zap.Error(err),
		)
		return
	}

	// 收到的好友申请
	receiveModel, err := friendmodel.NewReceiveFriendRequestModel(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "AgreeFriendApply NewReceiveFriendRequestModel err",
			zap.Uint64("userID", userID),
			zap.Error(err),
		)
		return
	}

	// TODO 判断是否toID每一个同意是否合法 todo 后续黑名单校验 容量校验

	for _, realyID := range toID {
		index := -1
		var request *friendmodel.ReceiveFriendRequestInfo
		for j, i := range receiveModel.ReceiveList {
			if i.FromUserId == uint64(realyID) {
				index = j
				request = i
				break
			}
		}

		if index == -1 {
			logger.CtxError(ctx, "AgreeFriendApply Find friend fail",
				zap.Error(err),
				zap.Any("userID", userID),
				zap.Any("toID", toID),
			)
			continue
		}

		// 逐步处理好友请求
		if isFriend := s.IsFriend(friendModel, uint64(realyID)); isFriend {
			// 把对方的发送也删掉
			s.delSendFriendRequestEvent(ctx, uint64(realyID), userID)
			// 忽略错误 已经是好友了就把收到的申请删掉
			s.delReceiveFriendRequestNotGet(logger, receiveModel, userID, uint64(realyID))
		}

		// 4.通知同意好友申请 todo
		if err = s.AcceptFriendRequestEvent(ctx, uint64(realyID), userID); err != nil {
			logger.ErrorWF("AgreeFriendApply err", zap.Error(err), zap.Uint64("realyID", uint64(realyID)))
			return
		}
		// 5.设置好友
		if err = s.addFriend(ctx, friendModel, userID, uint64(realyID)); err != nil {
			logger.ErrorWF("AgreeFriendApply addFriend err", zap.Error(err))
			return
		}
		// 6.删除收到的申请
		receiveModel.ReceiveList = append(receiveModel.ReceiveList[:index], receiveModel.ReceiveList[index+1:]...)
		rs = append(rs, request)
	}

	err = receiveModel.Save(logger, userID)
	if err != nil {
		logger.ErrorWF("AgreeFriendApply SetReceiveFriendRequest err", zap.Error(err))
		return
	}
	return
}

func (s *service) addFriend(ctx context.Context, friendModel *friendmodel.FriendModel, userId, toID uint64) error {
	logger := fklog.ContextAppLogger(ctx)
	friendModel.FriendList = append(friendModel.FriendList, &friendmodel.FriendInfo{
		UserId:   toID,
		CreateAt: time.Now().UnixMilli(),
	})
	err := friendModel.Save(ctx, userId)
	if err != nil {
		logger.ErrorWF("addFriend SetFriends failed", zap.Error(err))
		return errors.MODULE_ERROR
	}
	logger.InfoWF("addFriend success", zap.Uint64("userId", userId), zap.Uint64("toId", toID))
	return nil
}

// 收到同意好友请求事件, 已经是好友了将忽略错误
func (s *service) AcceptFriendRequestEvent(ctx context.Context, userId, fromId uint64) *errors.CodeError {
	logger := fklog.ContextAppLogger(ctx)
	// 好友列表
	friendModel, err := friendmodel.NewFriendModel(ctx, userId)
	if err != nil {
		logger.ErrorWF("AcceptFriendRequestEvent GetFriends err", zap.Error(err))
		return errors.MODULE_ERROR
	}
	sendModel, err := friendmodel.NewSendFriendRequestModel(ctx, userId)
	if err != nil {
		return errors.MODULE_ERROR
	}

	index := -1
	// 更新已发送好友请求
	for j, i := range sendModel.SendList {
		if i.ToUserId == fromId {
			index = j
		}
	}
	if index != -1 {
		sendModel.SendList = append(sendModel.SendList[:index], sendModel.SendList[index+1:]...)
		if err = sendModel.Save(logger, userId); err != nil {
			logger.ErrorWF("AcceptFriendRequestEvent SetSendFriendRequest err", zap.Error(err))
			return errors.MODULE_ERROR
		}
	}

	if isFriend := s.IsFriend(friendModel, fromId); isFriend {
		logger.WarnWF("AcceptFriendRequestEvent 对方已经是你的好友了")
		return nil
	}
	// 设置好友
	if err = s.addFriend(ctx, friendModel, userId, fromId); err != nil {
		logger.ErrorWF("AcceptFriendRequestEvent addFriend err", zap.Error(err))
		return errors.MODULE_ERROR
	}
	logger.InfoWF("AcceptFriendRequestEvent success", zap.Uint64("userId", userId), zap.Uint64("fromId", fromId))
	return nil
}

// 删除好友请求事件
func (s *service) delSendFriendRequestEvent(ctx context.Context, userId, toId uint64) {
	logger := fklog.ContextAppLogger(ctx)
	// // 补偿操作可能没必要
	// friendModel, err := friendmodel.NewFriendModel(ctx, userId)
	// if err != nil {
	// 	logger.ErrorWF("AcceptFriendRequest GetFriends err", zap.Error(err))
	// 	return
	// }
	// exist := false
	// for _, i := range friendModel.FriendList {
	// 	if i.UserId == toId {
	// 		exist = true
	// 		break
	// 	}
	// }
	// if !exist {
	// 	if err = s.addFriend(logger, friendModel, userId, toId); err != nil {
	// 		logger.ErrorWF("delSendFriendRequestEvent single friend rollback err", zap.Error(err))
	// 	}
	// 	logger.WarnWF("delSendFriendRequestEvent single friend rollback success", zap.Uint64("userId", userId), zap.Uint64("toId", toId))
	// }

	// 好友列表
	err := s.delSendFriendRequest(ctx, userId, toId)
	if err != nil {
		logger.CtxError(ctx, "delSendFriendRequestEvent delSendFriendRequest err", zap.Error(err), zap.Uint64("userId", userId), zap.Uint64("toId", toId))
	}
	logger.CtxInfo(ctx, "delSendFriendRequestEvent success", zap.Uint64("userId", userId), zap.Uint64("toId", toId))
}
