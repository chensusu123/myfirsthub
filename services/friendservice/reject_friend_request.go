package friendservice

import (
	"context"
	"maze_game_server/common/errors"
	"maze_game_server/model/friendmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// func (s *service) RejectFriendRequest(ctx , userId, toID uint64) error {
// 	// 检查有没有收到过好友请求
// 	receiveModel, err := friendmodel.NewReceiveFriendRequestModel(logger, userId)
// 	if err != nil {
// 		logger.ErrorWF("RejectFriendRequest GetReceiveFriendRequest err", zap.Error(err))
// 		return errors.MODULE_ERROR
// 	}
// 	var receive *friendmodel.ReceiveFriendRequestInfo = nil
// 	for _, i := range receiveModel.ReceiveList {
// 		if i.FromUserId == toID {
// 			receive = i
// 			break
// 		}
// 	}
// 	if receive == nil || receive.Status != friendmodel.FriendRequestStatusPending {
// 		return errors.COMMON_ERROR_TIPS.WrapMsg("已经处理过了")
// 	}

// 	// 通知对方拒绝 todo
// 	if codeErr := s.RejectFriendRequestEvent(logger, toID, userId); codeErr != nil {
// 		logger.ErrorWF("RejectFriendRequestEvent err", zap.Error(codeErr))
// 		return codeErr
// 	}

// 	receive.Status = friendmodel.FriendRequestStatusRejected
// 	err = receiveModel.Save(logger, userId)
// 	if err != nil {
// 		logger.ErrorWF("RejectFriendRequest SetReceiveFriendRequest err", zap.Error(err))
// 		return errors.MODULE_ERROR
// 	}

// 	return nil
// }

func (s *service) RefuseFriendApply(ctx context.Context, userID uint64, toID []int64) (rs []*friendmodel.ReceiveFriendRequestInfo, err error) {
	logger := fklog.ContextAppLogger(ctx)
	// 检查有没有收到过好友请求
	receiveModel, err := friendmodel.NewReceiveFriendRequestModel(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "RefuseFriendApply GetReceiveFriendRequest err", zap.Error(err))
		return
	}

	// toID 可靠性校验
	for _, realyID := range toID {
		index := -1
		for j, i := range receiveModel.ReceiveList {
			if i.FromUserId == uint64(realyID) {
				index = j
				break
			}
		}

		if index == -1 {
			logger.CtxError(ctx, "RefuseFriendApply Find friend fail",
				zap.Error(err),
				zap.Any("userID", userID),
				zap.Any("toID", realyID),
			)
			continue
		}

		// 通知对方拒绝 todo
		if codeErr := s.RejectFriendRequestEvent(ctx, uint64(realyID), userID); codeErr != nil {
			logger.CtxError(ctx, "RejectFriendRequestEvent err", zap.Error(codeErr))
			return
		}

		receiveModel.ReceiveList = append(receiveModel.ReceiveList[:index], receiveModel.ReceiveList[index+1:]...)
	}
	return
}

// 收到拒绝好友请求事件
func (s *service) RejectFriendRequestEvent(ctx context.Context, userId, fromId uint64) error {
	sendModel, err := friendmodel.NewSendFriendRequestModel(ctx, userId)
	if err != nil {
		return err
	}

	index := -1
	for j, i := range sendModel.SendList {
		if i.ToUserId == fromId {
			index = j
			break
		}
	}

	if index == -1 {
		err = errors.New("对方不存在")
		return err
	}

	return nil
}
