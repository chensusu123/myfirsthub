package friend

// 处理好友请求 同意/拒绝
// func (f *FriendComponent) OnReplyFriendApply_10697_10698(s *session.Session, req *Friend.ReplyFriendApplyRQ) (err error) {
// 	userId := uint64(s.UID())
// 	ctx := s.Context()
// 	logger := fklog.ContextAppLogger(ctx)
// 	res := &Friend.ReplyFriendApplyRS{}

// 	logger.InfoWF("OnReplyFriendApply start", zap.Any("req", req))
// 	defer func() {
// 		err = s.Response(res)
// 		logger.InfoWF("OnReplyFriendApply end", zap.Any("res", res))
// 	}()

// 	toID := req.GetUserId()

// 	switch req.GetReplyResult() {
// 	case int32(Friend.REPLY_FRIEND_APPLY_RESULT_AGREE):
// 		err = friendservice.GlobalFriendService.AgreeFriendApply(ctx, userId, toID)
// 		if err != nil {
// 			logger.ErrorWF("OnAcceptFriendRequest AcceptFriendRequest failed", zap.Error(err), zap.Any("toID", toID))
// 			res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
// 			return err
// 		}
// 	case int32(Friend.REPLY_FRIEND_APPLY_RESULT_REFUSE):

// 	}

// 	// 通知同意加好友
// 	pushMsg := &Friend.AcceptFriendRequestID{
// 		UserId: proto.Uint64(userId),
// 	}
// 	// 通知对方同意加好友
// 	online.Push(logger, toID, 10557, pushMsg)

// 	return nil
// }
