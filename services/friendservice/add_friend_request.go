package friendservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/model/friendmodel"
	"time"
)

func (s *service) AddFriendRequest(logger fklog.FKLogI, userId, toId uint64) *errors.CodeError {
	// 好友列表
	friends, err := friendmodel.NewFriendModel(logger, userId)
	if err != nil {
		logger.ErrorWF("AddFriendRequest NewFriendModel err", zap.Error(err))
		return errors.MODULE_ERROR
	}
	// 查一下发出去的申请数量
	sends, err := friendmodel.NewSendFriendRequestModel(logger, userId)
	if err != nil {
		logger.ErrorWF("AddFriendRequest NewSendFriendRequestModel err", zap.Error(err))
		return errors.MODULE_ERROR
	}
	// 1.判断好友数量
	friendCount := int64(len(friends.FriendList)) + int64(len(sends.SendList))
	// todo 读取好友数量配置
	if friendCount >= 100 {
		logger.InfoWF("AddFriendRequest 好友数量达到上限了", zap.Int64("friendCount", friendCount))
		return errors.COMMON_ERROR_TIPS.WrapMsg("好友数量达到上限了")
	}
	// 2.是否在我的黑名单中
	inBlk, err := s.isBlacklist(logger, userId, toId)
	if err != nil {
		logger.ErrorWF("OnFriendRequest IsBlacklist err", zap.Error(err))
		return errors.MODULE_ERROR
	}
	if inBlk {
		return errors.COMMON_ERROR_TIPS.WrapMsg("对方在你的黑名单中")
	}
	// 3.判断是否已经是好友了
	if isFriend := s.IsFriend(friends, toId); isFriend {
		return errors.COMMON_ERROR_TIPS.WrapMsg("对方已经是你的好友了")
	}
	// 4.检查重复发送
	if canSend := s.checkRepeatSendFriendRequest(sends, toId); !canSend {
		return errors.COMMON_ERROR_TIPS.WrapMsg("已经发送过好友请求了")
	}
	// 5.发送给对方 对方需判断还可不可以加好友 todo
	if codeErr := s.AddFriendRequestEvent(logger, toId, userId); codeErr != nil {
		logger.ErrorWF("AddFriendRequestEvent err", zap.Error(codeErr))
		return codeErr
	}

	// 6.设置已发送好友请求
	if err = s.addSendFriendRequest(logger, sends, userId, toId); err != nil {
		logger.ErrorWF("AddFriendRequest addSendFriendRequest err", zap.Error(err))
		return errors.MODULE_ERROR
	}

	return nil
}

// 获取好友数量，好友数量包括已经成为好友的+发出去的申请
func (s *service) GetUserFriendCount(logger fklog.FKLogI, userId uint64) (count int64, err error) {
	friendModel, err := friendmodel.NewFriendModel(logger, userId)
	if err != nil {
		logger.ErrorWF("GetUserFriendCount NewFriendModel err", zap.Error(err))
		return 0, err
	}
	// 查一下发出去的申请数量
	sendfriendRequestModel, err := friendmodel.NewSendFriendRequestModel(logger, userId)
	if err != nil {
		logger.ErrorWF("GetUserFriendCount NewSendFriendRequestModel err", zap.Error(err))
		return 0, err
	}
	return int64(len(friendModel.FriendList)) + int64(len(sendfriendRequestModel.SendList)), nil
}

// 等待actor逻辑完成处理
// 收到好友请求事件
func (s *service) AddFriendRequestEvent(logger fklog.FKLogI, userId, fromId uint64) *errors.CodeError {
	// 是否在我的黑名单中
	inBlk, err := s.isBlacklist(logger, userId, fromId)
	if err != nil {
		logger.ErrorWF("AddFriendRequestEvent isBlacklist err", zap.Error(err))
		return errors.MODULE_ERROR
	}
	if inBlk {
		return errors.COMMON_ERROR_TIPS.WrapMsg("已被对方拉黑")
	}

	// 设置待处理好友请求
	if err = s.addReceiveFriendRequest(logger, userId, fromId); err != nil {
		return errors.MODULE_ERROR
	}

	return nil
}

// 设置好友请求待处理
func (s *service) addReceiveFriendRequest(logger fklog.FKLogI, userId, fromId uint64) error {
	receiveModel, err := friendmodel.NewReceiveFriendRequestModel(logger, userId)
	if err != nil {
		logger.ErrorWF("addReceiveFriendRequest NewReceiveFriendRequestModel err", zap.Error(err), zap.Uint64("fromId", fromId))
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
		logger.WarnWF("addReceiveFriendRequest 已经申请过了", zap.Uint64("fromId", fromId))
		return nil
	}
	receiveModel.ReceiveList = append(receiveModel.ReceiveList, &friendmodel.ReceiveFriendRequestInfo{
		FromUserId: fromId,
		CreateAt:   time.Now().UnixMilli(),
		Status:     friendmodel.FriendRequestStatusPending,
	})

	err = receiveModel.Save(logger, userId)
	if err != nil {
		logger.ErrorWF("addReceiveFriendRequest SetReceiveFriendRequest err", zap.Error(err), zap.Uint64("fromId", fromId))
		return errors.MODULE_ERROR
	}

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
	if time.Since(sendTime) > 30*time.Minute {
		sendModel.SendList = append(sendModel.SendList[:index], sendModel.SendList[index+1:]...)
		return true
	}
	return false
}

// 设置好友请求
func (s *service) addSendFriendRequest(logger fklog.FKLogI, sendModel *friendmodel.SendFriendRequestModel, userId, toUserId uint64) (err error) {
	sendModel.SendList = append(sendModel.SendList, &friendmodel.SendFriendRequestInfo{
		ToUserId: toUserId,
		CreateAt: time.Now().UnixMilli(),
	})

	err = sendModel.Save(logger, userId)
	if err != nil {
		logger.ErrorWF("addSendFriendRequest err", zap.Error(err))
		return err
	}
	return nil
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
