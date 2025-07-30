package friendservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/model/friendmodel"
)

func (s *service) RejectFriendRequest(logger fklog.FKLogI, userId, toID uint64) *errors.CodeError {
	// 检查有没有收到过好友请求
	receiveModel, err := friendmodel.NewReceiveFriendRequestModel(logger, userId)
	if err != nil {
		logger.ErrorWF("RejectFriendRequest GetReceiveFriendRequest err", zap.Error(err))
		return errors.MODULE_ERROR
	}
	var receive *friendmodel.ReceiveFriendRequestInfo = nil
	for _, i := range receiveModel.ReceiveList {
		if i.FromUserId == toID {
			receive = i
			break
		}
	}
	if receive == nil || receive.Status != friendmodel.FriendRequestStatusPending {
		return errors.COMMON_ERROR_TIPS.WrapMsg("已经处理过了")
	}

	// 通知对方拒绝 todo
	if codeErr := s.RejectFriendRequestEvent(logger, toID, userId); codeErr != nil {
		logger.ErrorWF("RejectFriendRequestEvent err", zap.Error(codeErr))
		return codeErr
	}

	receive.Status = friendmodel.FriendRequestStatusRejected
	err = receiveModel.Save(logger, userId)
	if err != nil {
		logger.ErrorWF("RejectFriendRequest SetReceiveFriendRequest err", zap.Error(err))
		return errors.MODULE_ERROR
	}

	return nil
}

// 收到拒绝好友请求事件
func (s *service) RejectFriendRequestEvent(logger fklog.FKLogI, userId, fromId uint64) *errors.CodeError {
	sendModel, err := friendmodel.NewSendFriendRequestModel(logger, userId)
	if err != nil {
		return errors.MODULE_ERROR
	}
	for _, i := range sendModel.SendList {
		if i.ToUserId == fromId {
			i.Status = friendmodel.FriendRequestStatusRejected
			break
		}
	}
	if err = sendModel.Save(logger, userId); err != nil {
		logger.ErrorWF("RejectFriendRequestEvent Save err", zap.Error(err))
		return errors.MODULE_ERROR
	}

	return nil
}
