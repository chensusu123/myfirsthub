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

func (s *service) AgreeFriendApply(ctx context.Context, userID uint64, toID []int64) (rs []*friendmodel.ReceiveFriendRequestInfo, change []*friendmodel.ReceiveFriendRequestInfo, isSkip bool, err error) {
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
			err = fmt.Errorf("请求好友不存在")
			continue
		}

		rs = append(rs, request)

		// 删除对方的发送请求
		s.delSendFriendRequestEvent(ctx, uint64(realyID), userID)
		// 删除自己的收到请求
		s.delReceiveFriendRequestNotGet(ctx, receiveModel, userID, uint64(realyID))

		// 校验是否在黑名单
		var isBlack bool
		isBlack, err = s.IsBlacklist(ctx, userID, uint64(realyID))
		if err != nil {
			logger.CtxError(ctx, "AgreeFriendApply IsBlacklist Fail", zap.Any("userID", userID), zap.Any("toID", toID), zap.Error(err))
			return
		}

		if isBlack {
			err = errors.New("对方在黑名单中")
			isSkip = true
			return
		}

		// 逐步处理好友请求
		if isFriend := s.IsFriend(friendModel, uint64(realyID)); isFriend {
			logger.CtxWarn(ctx, "AgreeFriendApply IsFriend", zap.Any("err", "已经是好友了"), zap.Any("userID", userID), zap.Any("realyID", realyID))
			err = errors.New("已经是好友了")
			isSkip = true
			continue
		}

		// 判断好友列表长度
		if len(friendModel.FriendList) >= friendmodel.MaxFriendSize {
			logger.CtxWarn(ctx, "AgreeFriendApply Equal to Limit", zap.Any("err", "当前用户好友满了"), zap.Any("userID", userID), zap.Any("realyID", realyID))
			err = errors.New("当前用户好友满了")
			isSkip = true
			continue
		}

		// 4.通知同意好友申请 todo
		if isSkip, err = s.AcceptFriendRequestEvent(ctx, uint64(realyID), userID, request.From); err != nil {
			logger.CtxError(ctx, "AgreeFriendApply err", zap.Error(err), zap.Uint64("realyID", uint64(realyID)))
			return
		}
		// 5.设置好友
		if err = s.addFriend(ctx, friendModel, userID, uint64(realyID), request.From); err != nil {
			logger.CtxError(ctx, "AgreeFriendApply addFriend err", zap.Error(err))
			return
		}
		// 6.删除收到的申请
		receiveModel.ReceiveList = append(receiveModel.ReceiveList[:index], receiveModel.ReceiveList[index+1:]...)
		change = append(change, request)
	}

	err = receiveModel.Save(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "AgreeFriendApply SetReceiveFriendRequest err", zap.Error(err))
		return
	}
	return
}

// userID校验toID是否合法
func (s *service) checkReliability(ctx context.Context, userID, toID uint64) (bool, error) {
	logger := fklog.ContextAppLogger(ctx)
	friendList, err := friendmodel.NewFriendModel(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "checkReliability NewFriendModel Fail", zap.Any("userID", userID), zap.Any("toID", toID), zap.Error(err))
		return false, err
	}
	// 校验容量
	if len(friendList.FriendList) >= friendmodel.MaxFriendSize {
		return false, nil
	}

	// 校验黑名单
	isBlack, err := s.IsBlacklist(ctx, userID, toID)
	if err != nil {
		logger.CtxError(ctx, "checkReliability IsBlacklist Fail", zap.Any("userID", userID), zap.Any("toID", toID), zap.Error(err))
		return false, err
	}
	if isBlack {
		return !isBlack, nil
	}

	// 校验好友列表是否存在对方
	isFriend := s.IsFriend(friendList, toID)
	if isFriend {
		return !isFriend, nil
	}

	return true, nil
}

func (s *service) addFriend(ctx context.Context, friendModel *friendmodel.FriendModel, userId, toID uint64, from int32) error {
	logger := fklog.ContextAppLogger(ctx)
	friendModel.FriendList = append(friendModel.FriendList, &friendmodel.FriendInfo{
		UserId:     toID,
		FriendType: from,
		CreateAt:   time.Now().UnixMilli(),
	})
	err := friendModel.Save(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "addFriend SetFriends failed", zap.Error(err))
		return errors.MODULE_ERROR
	}
	logger.CtxInfo(ctx, "addFriend success", zap.Uint64("userId", userId), zap.Uint64("toId", toID))
	return nil
}

// 收到同意好友请求事件, 已经是好友了将忽略错误
func (s *service) AcceptFriendRequestEvent(ctx context.Context, userId, fromId uint64, from int32) (bool, error) {
	logger := fklog.ContextAppLogger(ctx)
	// 好友列表
	friendModel, err := friendmodel.NewFriendModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "AcceptFriendRequestEvent GetFriends err", zap.Error(err))
		return false, err
	}
	sendModel, err := friendmodel.NewSendFriendRequestModel(ctx, userId)
	if err != nil {
		return false, err
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
		if err = sendModel.Save(ctx, userId); err != nil {
			logger.CtxError(ctx, "AcceptFriendRequestEvent SetSendFriendRequest err", zap.Error(err))
			return false, err
		}
	}

	if isFriend := s.IsFriend(friendModel, fromId); isFriend {
		logger.CtxWarn(ctx, "AcceptFriendRequestEvent 对方已经是你的好友了")
		return true, fmt.Errorf("对方已经是你的好友了")
	}

	// 检测容量
	if len(friendModel.FriendList) >= friendmodel.MaxFriendSize {
		logger.CtxWarn(ctx, "AcceptFriendRequestEvent Equal to limit", zap.String("err", "好友容量满了"), zap.Any("userId", userId), zap.Any("fromId", fromId), zap.Any("from", from))
		return true, fmt.Errorf("对方好友容量满了")
	}

	// 检测黑名单
	isBlack, err := s.IsBlacklist(ctx, userId, fromId)
	if err != nil {
		logger.CtxError(ctx, "AcceptFriendRequestEvent IsBlacklist Fail", zap.Uint64("userId", userId), zap.Uint64("fromId", fromId), zap.Error(err))
		return false, err
	}

	if isBlack {
		logger.CtxWarn(ctx, "AcceptFriendRequestEvent isBlack", zap.String("err", "对方已经将你拉黑"), zap.Any("userId", userId), zap.Any("fromId", fromId), zap.Any("from", from))
		return true, fmt.Errorf("对方已经将你拉黑")
	}

	// 设置好友
	if err = s.addFriend(ctx, friendModel, userId, fromId, from); err != nil {
		logger.CtxError(ctx, "AcceptFriendRequestEvent addFriend err", zap.Error(err))
		return false, err
	}
	logger.CtxInfo(ctx, "AcceptFriendRequestEvent success", zap.Uint64("userId", userId), zap.Uint64("fromId", fromId))
	return false, nil
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
