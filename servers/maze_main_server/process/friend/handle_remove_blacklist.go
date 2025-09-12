package friend

// func (f *FriendComponent) OnRemoveBlacklist_10566_10567(s *session.Session, req *Friend.RemoveBlacklistRQ) (err error) {
// 	userId := uint64(s.UID())
// 	logger := log.Clone("OnRemoveBlacklist", userId, 0)
// 	res := &Friend.RemoveBlacklistRS{}

// 	logger.InfoWF("OnRemoveBlacklist start", zap.Any("req", req))
// 	defer func() {
// 		err = s.Response(res)
// 		logger.InfoWF("OnRemoveBlacklist end", zap.Any("res", res))
// 	}()

// 	toID := req.GetUserId()
// 	if userId == toID {
// 		return errors.New("不能加自己为好友")
// 	}

// 	codeErr := friendservice.GlobalFriendService.RemoveBlacklist(logger, userId, toID)
// 	if codeErr != nil {
// 		logger.ErrorWF("OnRemoveBlacklist RemoveBlacklist failed", zap.Error(err), zap.Uint64("userId", userId), zap.Uint64("toID", toID))
// 		res.ErrInfo = codeErr.ToInfo()
// 		return codeErr
// 	}

// 	return nil
// }
